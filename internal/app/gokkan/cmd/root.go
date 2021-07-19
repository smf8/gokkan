package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/smf8/gokkan/internal/app/gokkan/cmd/server"
	"github.com/smf8/gokkan/internal/app/gokkan/config"
	"github.com/spf13/cobra"
)

// NewRootCommand creates a new gokkan root command.
func NewRootCommand() *cobra.Command {
	var root = &cobra.Command{
		Use: "gokkan",
	}

	cfg := config.New()

	log.SetOutput(os.Stdout)
	log.SetLevel(cfg.Logger.Level)

	server.Register(root, cfg)

	return root
}
