package server

import (
	"github.com/spf13/cobra"
)

func main() error {
	return nil
}

// Register registers server command to the root gokkan command
//nolint:gomnd
func Register(root *cobra.Command) {
	root.AddCommand(
		&cobra.Command{
			Use:   "server {--port port number}",
			Short: "start the server",
			RunE: func(cmd *cobra.Command, args []string) error {
				return main()
			},
		},
	)
}
