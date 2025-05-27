/*
Copyright © 2025 MetaDandy benitezarroyojoseph@gmail.com
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Metadandy/GoCelerator/helper"
	"github.com/spf13/cobra"
)

var watch bool

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the development server",
	Long: `Serve launches your API in dev mode.
	By default runs “go run ./cmd/main”. 
	Use --watch for Air hot-reload.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dest, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("could not get working dir: %w", err)
		}

		if watch {
			airConfig := filepath.Join(dest, "air.toml")
			fmt.Println("Starting Air with config:", airConfig)
			c := exec.Command("air", "-c", airConfig)
			c.Dir = dest
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			return c.Run()
		}

		return helper.RunInDir(dest, "go", "run", "./cmd/main")
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().BoolVarP(&watch, "watch", "w", false, "use Air hot-reload")
}
