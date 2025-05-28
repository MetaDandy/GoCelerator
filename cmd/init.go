/*
Copyright © 2025 MetaDandy benitezarroyojoseph@gmail.com
*/
package cmd

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Metadandy/GoCelerator/helper"
	"github.com/spf13/cobra"
)

/*
! TODO:
hacer que la configuracion sea cargada por main.go
implementar la logica de la migración
*/

//go:embed templates/*
var templatesFS embed.FS

var (
	useFiber  bool
	goVersion string
	noAir     bool
	noDocker  bool
)

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVarP(&useFiber, "fiber", "f", false, "use Fiber in the template")
	initCmd.Flags().StringVar(&goVersion, "goversion", "", "Go version to use (e.g. 1.24.1). If empty, will detect or ask")
	initCmd.Flags().BoolVar(&noAir, "no-air", false, "skip generating Air configuration")
	initCmd.Flags().BoolVar(&noDocker, "no-docker", false, "skip generating Docker configuration")
}

var initCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Initialize a go project.",
	Long: `Init generates the base scaffolding for a Go project, following the recommended folder structure (cmd/, config/, src/core/, src/modules/) and includes templates for either net/http or Fiber v2. The generated project comes with a preconfigured Go module, an example HTTP server, and Gorm migration support.
	Options:
	-f, --fiber Generate the project using Fiber v2 instead of net/http

	Examples:

	Create a standard project using net/http
	gocelerator init myapp

	Create a project using Fiber
	gocelerator init myapp --fiber`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		dest := filepath.Join(".", name)
		baseDir := "templates/nethttp"
		if useFiber {
			baseDir = "templates/fiber"
		}
		origDir, err := os.Getwd()
		if err != nil {
			return err
		}
		fmt.Printf("Generating proyect '%s' using %s...\n", name,
			map[bool]string{true: "Fiber", false: "net/http"}[useFiber],
		)
		if goVersion == "" {
			detected, err := helper.DetectGoVersion()
			if err != nil {
				return fmt.Errorf("could not detect Go version: %w", err)
			}
			fmt.Printf("Detected Go version: %s\nPress Enter to accept or type a different version: ", detected)
			var input string
			fmt.Scanln(&input)
			if input == "" {
				goVersion = detected
			} else {
				goVersion = input
			}
		}
		dataTempl := map[string]string{
			"ProjectName": name,
			"GoVersion":   goVersion,
		}
		if err := helper.CopyTemplates(baseDir, dest, templatesFS, dataTempl); err != nil {
			return fmt.Errorf("error al copiar plantillas: %w", err)
		}

		if !noAir {
			if _, err := exec.LookPath("air"); err != nil {
				fmt.Println("Air not found, installing github.com/cosmtrek/air@latest…")
				if out, err := exec.Command("go", "install", "github.com/cosmtrek/air@latest").CombinedOutput(); err != nil {
					return fmt.Errorf("failed to install Air: %s", string(out))
				}
			}
			if err := helper.CopyTemplates("templates/air", dest, templatesFS, dataTempl); err != nil {
				return fmt.Errorf("failed to generate .air.toml: %w", err)
			}
			fmt.Println("Generated custom .air.toml for hot-reload (cmd/main.go).")
		}

		if !noDocker {
			helper.CopyTemplates("templates/docker/prod", dest, templatesFS, dataTempl)

			if !noAir {
				helper.CopyTemplates("templates/docker/dev", dest, templatesFS, dataTempl)
			}
		}

		_ = os.Chdir(origDir)

		return nil
	},
}
