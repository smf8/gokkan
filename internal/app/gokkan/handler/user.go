package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/smf8/gokkan/internal/app/gokkan/auth"
	"github.com/smf8/gokkan/internal/app/gokkan/model"
	"github.com/smf8/gokkan/internal/app/gokkan/request"
	"github.com/smf8/gokkan/internal/app/gokkan/response"
)

// UserHandler handles user and admin related endpoints.
type UserHandler struct {
	UserRepo  model.UserRepo
	AdminRepo model.AdminRepo
	jwtSecret string
}

// NewUserHandler creates a new user handler.
func NewUserHandler(userRepo model.UserRepo,
	adminRepo model.AdminRepo, jwtSecret string) UserHandler {
	return UserHandler{
		UserRepo:  userRepo,
		AdminRepo: adminRepo,
		jwtSecret: jwtSecret,
	}
}

// Signup handles user signup.
func (u UserHandler) Signup(c echo.Context) error {
	req := &request.Signup{}

	if err := c.Bind(req); err != nil {
		logrus.Errorf("failed to bind signup request: %s", err.Error())

		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("bind request failed: %s", err))
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("bad login request: %s", err.Error()))
	}

	user := &model.User{
		Username:       strings.ToLower(req.Username),
		Password:       auth.Hash(req.Password),
		FullName:       req.FullName,
		BillingAddress: req.BillingAddress,
	}

	// it's better to handle duplicate user signup error differently
	if err := u.UserRepo.Save(user); err != nil {
		logrus.Errorf("failed to create user %+v : %s", *user, err)

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create user, user exists")
	}

	jwtToken, err := auth.Generate(u.jwtSecret, user.Username, false)
	if err != nil {
		logrus.Errorf("failed to create jwt token: %s", err.Error())

		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	response := response.User{
		ID:             user.ID,
		Username:       user.Username,
		FullName:       user.FullName,
		BillingAddress: user.BillingAddress,
		Token:          jwtToken,
	}

	return c.JSON(http.StatusCreated, response)
}

// Login handles user/admin login.
func (u UserHandler) Login(c echo.Context) error {
	req := &request.Login{}

	if err := c.Bind(req); err != nil {
		logrus.Errorf("failed to bind login request: %s", err.Error())

		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("bind request failed: %s", err))
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("bad login request: %s", err.Error()))
	}

	if req.IsAdmin {
		return u.loginAdmin(c, req)
	}

	// convert username to lowercase
	req.Username = strings.ToLower(req.Username)

	return u.loginUser(c, req)
}

func (u UserHandler) loginUser(c echo.Context, req *request.Login) error {
	user, err := u.UserRepo.Find(req.Username)
	if err != nil {
		logrus.Errorf("failed to login user %s: %s", req.Username, err.Error())

		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	if !auth.CheckPassword(req.Password, user.Password) {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	jwtToken, err := auth.Generate(u.jwtSecret, user.Username, false)
	if err != nil {
		logrus.Errorf("failed to create jwt token: %s", err.Error())

		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	userResponse := &response.User{
		ID:             user.ID,
		Username:       user.Username,
		FullName:       user.FullName,
		BillingAddress: user.BillingAddress,
		Balance:        user.Balance,
		Token:          jwtToken,
		IsAdmin:        false,
	}

	return c.JSON(http.StatusOK, userResponse)
}

func (u UserHandler) loginAdmin(c echo.Context, req *request.Login) error {
	admin, err := u.AdminRepo.Find(req.Username)
	if err != nil {
		logrus.Debugf("failed to find admin %s: %s", req.Username, err.Error())
	}

	if !auth.CheckPassword(req.Password, admin.Password) {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	jwtToken, err := auth.Generate(u.jwtSecret, req.Username, true)
	if err != nil {
		logrus.Errorf("failed to create jwt token: %s", err.Error())

		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	adminResponse := &response.User{
		ID:       admin.ID,
		Username: admin.Username,
		Token:    jwtToken,
		IsAdmin:  true,
	}

	return c.JSON(http.StatusOK, adminResponse)
}
