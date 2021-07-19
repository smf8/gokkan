package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/smf8/gokkan/internal/app/gokkan/config"
	"github.com/smf8/gokkan/internal/app/gokkan/router"
	"github.com/spf13/cobra"
)

const shutdownTimeout = 5 * time.Second

func main(cfg config.Config) {
	echo := router.Echo()

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
		log.Fatalf("failed to shutdown echo gracefully: %s", err.Error())
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
