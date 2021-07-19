package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Echo creates a new echo router with given config
func Echo() *echo.Echo {
	e := echo.New()

	// register default middlewares
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RemoveTrailingSlash())

	return e
}
