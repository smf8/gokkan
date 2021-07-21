package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4/middleware"
	"github.com/sirupsen/logrus"
	"github.com/smf8/gokkan/internal/app/gokkan/auth"
	"github.com/smf8/gokkan/internal/app/gokkan/config"
	"github.com/smf8/gokkan/internal/app/gokkan/database"
	"github.com/smf8/gokkan/internal/app/gokkan/handler"
	"github.com/smf8/gokkan/internal/app/gokkan/model"
	"github.com/smf8/gokkan/internal/app/gokkan/router"
	"github.com/spf13/cobra"
)

const shutdownTimeout = 5 * time.Second

// nolint:funlen
func main(cfg config.Config) {
	echo := router.Echo()

	db, err := database.New(cfg.Database)
	if err != nil {
		logrus.Fatalf("failed to connect to database: %s", err.Error())
	}

	userRepo := model.SQLUserRepo{
		DB: db,
	}
	categoryRepo := model.SQLCategoryRepo{
		DB: db,
	}

	blacklistRepo := model.NewCacheTokenBlacklistRepo(auth.DefaultExpiration)

	auth.SetTokenBlacklistRepo(blacklistRepo)

	jwtConfig := auth.MiddlewareConfig(cfg.Server)

	userHandler := handler.NewUserHandler(userRepo, blacklistRepo, cfg.Server.Secret)
	categoryHandler := handler.CategoryHandler{CategoryRepo: categoryRepo}

	// unrestricted endpoints
	echo.POST("/login", userHandler.Login)
	echo.POST("/signup", userHandler.Signup)
	echo.GET("/categories", categoryHandler.GetAll)

	// restricted endpoints. requires authorization
	userArea := echo.Group("/users")
	adminArea := echo.Group("/admin")

	userArea.Use(middleware.JWTWithConfig(jwtConfig))
	adminArea.Use(middleware.JWTWithConfig(jwtConfig))

	// user area routing
	userArea.PUT("/charge", userHandler.ChargeBalance)
	userArea.GET("/me", userHandler.GetInfo)
	userArea.POST("/logout", userHandler.Logout)

	// admin area routing
	adminArea.POST("/categories/create", categoryHandler.Create)
	adminArea.DELETE("/categories/delete/:id", categoryHandler.Delete)

	go func() {
		err := echo.Start(fmt.Sprintf(":%d", cfg.Server.Port))
		if err != nil {
			logrus.Fatalf("failed to start echo server: %s", err.Error())
		}
	}()

	// Handle Ctrl-C or other signals
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	s := <-sig
	logrus.Infof("got signal %s, shutting down", s)

	ctx, c := context.WithTimeout(context.Background(), shutdownTimeout)

	defer c()

	if err := echo.Shutdown(ctx); err != nil {
		logrus.Errorf("failed to shutdown echo gracefully: %s", err.Error())
	}
}

// Register registers server command to the root gokkan command
//nolint:gomnd
func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		&cobra.Command{
			Use:   "server {--port port number}",
			Short: "start the server",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg)
			},
		},
	)
}
