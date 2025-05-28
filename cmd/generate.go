package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/Metadandy/GoCelerator/helper"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

/*
! TODO:
- Una vez creado los templates, añadirlo en container.
- Lo que se ha hecho en container, añádir a api para registrar las rutas
*/

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().BoolVarP(&useFiber, "fiber", "f", false,
		"use Fiber in the templates (default net/http)")
}

var generateCmd = &cobra.Command{
	Use:   "generate [name]",
	Short: "Generates scaffold for a module (model, repo, service, handler, dto, container)",
	Args:  cobra.ExactArgs(1),
	RunE:  runGenerate,
}

func runGenerate(cmd *cobra.Command, args []string) error {
	raw := args[0]
	parts := strings.Split(raw, "/")
	base := parts[len(parts)-1]
	sub := filepath.Join(parts...)
	framework := "nethttp"
	if useFiber {
		framework = "fiber"
	}

	basePath := filepath.Join("templates", "generate", framework)
	tplFS, err := fs.Sub(templatesFS, basePath)
	if err != nil {
		return fmt.Errorf("could not find templates in %s: %w", basePath, err)
	}

	modelDir := filepath.Join("src", "models")
	moduleDir := filepath.Join("src", sub)
	os.MkdirAll(modelDir, 0755)
	os.MkdirAll(moduleDir, 0755)

	modulePath, err := helper.DetectModulePath()
	if err != nil {
		return err
	}

	data := map[string]string{
		"Name":       cases.Title(language.Und, cases.NoLower).String(base),
		"Package":    base,
		"ModulePath": modulePath,
	}

	files := []struct{ tmpl, dstDir, dstName string }{
		{"model.go.tmpl", modelDir, base + ".go"},
		{"repository.go.tmpl", moduleDir, base + "_repository.go"},
		{"service.go.tmpl", moduleDir, base + "_service.go"},
		{"handler.go.tmpl", moduleDir, base + "_handler.go"},
		{"dto.go.tmpl", moduleDir, base + "_dto.go"},
	}

	for _, f := range files {
		if err := helper.CopyTemplatesFS(tplFS, f.tmpl, f.dstDir, data); err != nil {
			return fmt.Errorf("generating %s: %w", f.dstName, err)
		}
	}

	fmt.Printf("Module %q generated under src/%s (framework: %s)\n", base, sub, framework)
	return nil
}
