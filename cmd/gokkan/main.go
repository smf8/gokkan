package main

import (
	"os"

	"github.com/smf8/gokkan/internal/app/gokkan/cmd"
)

const exitCodeErr = 1

func main() {
	root := cmd.NewRootCommand()

	if root != nil {
		if err := root.Execute(); err != nil {
			os.Exit(exitCodeErr)
		}
	}
}
