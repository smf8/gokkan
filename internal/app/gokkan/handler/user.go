package handler

import (
	"net/http"

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

// Login handles user/admin login.
func (u UserHandler) Login(c echo.Context) error {
	req := &request.Login{}

	if err := c.Bind(req); err != nil {
		logrus.Errorf("failed to bind request: %s", err.Error())

		return echo.NewHTTPError(http.StatusBadRequest, "bind request failed")
	}

	if req.IsAdmin {
		return u.loginAdmin(c, req)
	}

	return u.loginUser(c, req)
}

func (u UserHandler) loginUser(c echo.Context, req *request.Login) error {
	user, err := u.UserRepo.Find(req.Username)
	if err != nil {
		logrus.Debugf("failed to login user %s: %s", req.Username, err.Error())

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
