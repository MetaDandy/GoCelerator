/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

//go:embed templates/*
var templatesFS embed.FS

var (
	useFiber bool
)

// initCmd represents the init command
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
		fmt.Printf("Generating proyect '%s' using %s...\n", name,
			map[bool]string{true: "Fiber", false: "net/http"}[useFiber],
		)
		return fs.WalkDir(templatesFS, baseDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			rel, err := filepath.Rel(baseDir, path)
			if err != nil {
				return err
			}
			if rel == "." {
				return os.MkdirAll(dest, 0755)
			}
			var targetRel string
			if d.IsDir() {
				targetRel = rel
			} else {
				targetRel = strings.TrimSuffix(rel, ".tmpl")
			}
			target := filepath.Join(dest, targetRel)
			if d.IsDir() {
				return os.MkdirAll(target, 0755)
			}
			if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
				return err
			}
			data, err := templatesFS.ReadFile(path)
			if err != nil {
				return err
			}
			t := template.Must(template.New(rel).Parse(string(data)))
			f, err := os.Create(target)
			if err != nil {
				return err
			}
			defer f.Close()
			return t.Execute(f, map[string]string{"ProjectName": name})
		})
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVarP(&useFiber, "fiber", "f", false, "use Fiber in the template")
}
