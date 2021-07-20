package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/smf8/gokkan/internal/app/gokkan/request"
)

// Echo creates a new echo router with given config.
func Echo() *echo.Echo {
	e := echo.New()

	// register request validator
	e.Validator = request.NewValidator()

	// register default middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RemoveTrailingSlash())

	return e
}
