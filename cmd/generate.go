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

	baseTitle := cases.Title(language.Und, cases.NoLower).String(base)

	data := map[string]string{
		"Name":       baseTitle,
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

	containerPath := filepath.Join("src", "container.go")
	content, err := os.ReadFile(containerPath)
	if err != nil {
		return err
	}
	text := string(content)

	imp := fmt.Sprintf("\t\"%s/src/%s\"\n", modulePath, sub)
	field := fmt.Sprintf("\t// %s\n\t%[1]sRepo *%[1]s.Repository\n\t%[1]sSvc  *%[1]s.Service\n\t%[1]sHdl  *%[1]s.Handler\n", baseTitle)
	init := fmt.Sprintf("\t// %s\n\t%[1]sRepo := %s.NewRepository(config.DB)\n\t%[1]sSvc := %s.NewService(%[1]sRepo)\n\t%[1]sHdl := %s.NewHandler(%[1]sSvc)\n", baseTitle,
		base, base, base, base, base,
	)
	ctor := fmt.Sprintf("\t// %s\n\t%[1]sRepo: %[1]sRepo,\n\t%[1]sSvc:  %[1]sSvc,\n\t%[1]sHdl:  %[1]sHdl,\n", baseTitle)

	text = strings.Replace(text,
		"// ##module_imports##",
		"// ##module_imports##\n"+imp,
		1,
	)
	text = strings.Replace(text,
		"// ##module_fields##",
		"// ##module_fields##\n"+field,
		1,
	)
	text = strings.Replace(text,
		"// ##module_inits##",
		"// ##module_inits##\n"+init,
		1,
	)
	text = strings.Replace(text,
		"// ##module_ctor##",
		"// ##module_ctor##\n"+ctor,
		1,
	)

	if err := os.WriteFile(containerPath, []byte(text), 0644); err != nil {
		return err
	}

	fmt.Printf("Module %q setup in container\n", base)

	apiPath := filepath.Join("cmd", "api", "api.go")
	content, err = os.ReadFile(apiPath)
	if err != nil {
		return err
	}
	text = string(content)

	handlerLine := fmt.Sprintf("\tc.%sHdl.RegisterRoutes,\n", baseTitle)

	text = strings.Replace(text,
		"// ##module_api_handlers##",
		"// ##module_api_handlers##\n"+handlerLine,
		1,
	)

	if err := os.WriteFile(apiPath, []byte(text), 0644); err != nil {
		return err
	}

	fmt.Printf("Module %q setup in api\n", base)

	return nil
}
