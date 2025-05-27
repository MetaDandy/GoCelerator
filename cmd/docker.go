// cmd/docker.go
package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Metadandy/GoCelerator/helper"
	"github.com/spf13/cobra"
)

var devDocker bool

var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Run Docker Compose stack",
	Long: `Starts your Docker environment.
	By default uses docker-compose.yaml (production).
	Use --dev to use docker-compose.dev.yml (development).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		dest, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("could not get working dir: %w", err)
		}

		var composeFile string
		if devDocker {
			composeFile = "docker-compose.dev.yaml"
			fmt.Println("Starting development stack with", composeFile)
			return helper.RunInDir(dest,
				"docker", "compose", "-f", filepath.Join(dest, composeFile), "up", "--build",
			)
		} else {
			composeFile = "docker-compose.yaml"
			fmt.Println("Starting production stack with", composeFile)
			return helper.RunInDir(dest,
				"docker", "compose", "up", "--build",
			)
		}
	},
}

func init() {
	rootCmd.AddCommand(dockerCmd)
	dockerCmd.Flags().BoolVarP(&devDocker, "dev", "d", false,
		"use development Docker Compose (docker-compose.dev.yaml)")
}
