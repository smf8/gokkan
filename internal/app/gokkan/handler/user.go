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
	UserRepo      model.UserRepo
	BlacklistRepo model.TokenBlacklistRepo
	jwtSecret     string
}

// NewUserHandler creates a new user handler.
func NewUserHandler(userRepo model.UserRepo,
	tokenRepo model.TokenBlacklistRepo, jwtSecret string) UserHandler {
	return UserHandler{
		UserRepo:      userRepo,
		BlacklistRepo: tokenRepo,
		jwtSecret:     jwtSecret,
	}
}

// GetInfo handles user info retrieval.
func (u UserHandler) GetInfo(c echo.Context) error {
	claims, err := auth.ExtractClaims(c)
	if err != nil {
		logrus.Errorf("get info: failed to extract claims: %s", err)

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to extract jwt claims")
	}

	user, err := u.UserRepo.Find(claims.Sub)
	if err != nil {
		logrus.Errorf("get info: failed to find user %s: %s", claims.Sub, err.Error())

		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	response := response.User{
		ID:             user.ID,
		Username:       user.Username,
		FullName:       user.FullName,
		BillingAddress: user.BillingAddress,
		Balance:        user.Balance,
		IsAdmin:        user.IsAdmin,
	}

	return c.JSON(http.StatusOK, response)
}

// ChargeBalance handles balance increase.
func (u UserHandler) ChargeBalance(c echo.Context) error {
	req := &request.ChargeBalance{}

	if err := c.Bind(req); err != nil {
		logrus.Errorf("charge balance: failed to bind signup request: %s", err.Error())

		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("bind request failed: %s", err))
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("bad login request: %s", err.Error()))
	}

	claims, err := auth.ExtractClaims(c)
	if err != nil {
		logrus.Errorf("charge balance: failed to extract claims: %s", err)

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to extract jwt claims")
	}

	user, err := u.UserRepo.Find(claims.Sub)
	if err != nil {
		logrus.Errorf("charge balance: failed to find user %s: %s", claims.Sub, err.Error())

		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	user.Balance += req.Amount

	if err := u.UserRepo.Save(user); err != nil {
		logrus.Errorf("failed to update user %s: %s", user.Username, err.Error())

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update user balance: internal error")
	}

	return c.NoContent(http.StatusOK)
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
		logrus.Errorf("signup: failed to create user %+v : %s", *user, err)

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create user, user exists")
	}

	jwtToken, err := auth.Generate(u.jwtSecret, user.Username, false)
	if err != nil {
		logrus.Errorf("signup: failed to create jwt token: %s", err.Error())

		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	response := response.User{
		ID:             user.ID,
		Username:       user.Username,
		FullName:       user.FullName,
		BillingAddress: user.BillingAddress,
		Balance:        user.Balance,
		Token:          jwtToken,
	}

	return c.JSON(http.StatusCreated, response)
}

// Login handles user/admin login.
func (u UserHandler) Login(c echo.Context) error {
	req := &request.Login{}

	if err := c.Bind(req); err != nil {
		logrus.Errorf("login: failed to bind login request: %s", err.Error())

		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("bind request failed: %s", err))
	}

	if err := c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("bad login request: %s", err.Error()))
	}

	req.Username = strings.ToLower(req.Username)

	user, err := u.UserRepo.Find(req.Username)
	if err != nil {
		logrus.Errorf("user login: failed to login user %s: %s", req.Username, err.Error())

		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	if !auth.CheckPassword(req.Password, user.Password) {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	jwtToken, err := auth.Generate(u.jwtSecret, user.Username, user.IsAdmin)
	if err != nil {
		logrus.Errorf("user login: failed to create jwt token: %s", err.Error())

		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	userResponse := &response.User{
		ID:             user.ID,
		Username:       user.Username,
		FullName:       user.FullName,
		BillingAddress: user.BillingAddress,
		Balance:        user.Balance,
		Token:          jwtToken,
		IsAdmin:        user.IsAdmin,
	}

	return c.JSON(http.StatusOK, userResponse)
}

// Logout handles user logout.
func (u UserHandler) Logout(c echo.Context) error {
	claims, err := auth.ExtractClaims(c)
	if err != nil {
		logrus.Errorf("logout: failed to extract claims: %s", err)

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to extract jwt claims")
	}

	if err := u.BlacklistRepo.Save(claims.JTI); err != nil {
		logrus.Errorf("logout: failed to block token: %s", err.Error())

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to block jwt token")
	}

	return c.NoContent(http.StatusOK)
}

// Update updates user info.
func (u UserHandler) Update(c echo.Context) error {
	claims, err := auth.ExtractClaims(c)
	if err != nil {
		logrus.Errorf("update info: failed to extract claims: %s", err)

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to extract jwt claims")
	}

	user, err := u.UserRepo.Find(claims.Sub)
	if err != nil {
		logrus.Errorf("update info: failed to find user %s: %s", claims.Sub, err.Error())

		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	req := &request.UpdateUser{}

	if err = c.Bind(req); err != nil {
		logrus.Errorf("update info: failed to bind update info request: %s", err.Error())

		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("bind request failed: %s", err))
	}

	if err = c.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("bad get info request: %s", err.Error()))
	}

	if req.Password != "" {
		user.Password = auth.Hash(req.Password)
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}

	if req.BillingAddress != "" {
		user.BillingAddress = req.BillingAddress
	}

	if err = u.UserRepo.Save(user); err != nil {
		logrus.Errorf("failed to update user %+v: %s", req, err.Error())

		return echo.NewHTTPError(http.StatusInternalServerError, "failed to update user")
	}

	return c.JSON(http.StatusOK, user)
}
