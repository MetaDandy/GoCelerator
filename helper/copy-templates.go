package helper

import (
	"embed"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Tour templates in templatesFS  under baseDir
// process each file as a template with 'data' and dumps it into dest,
// removing the ".tmpl" extension.
func CopyTemplates(baseDir, dest string, templatesFS embed.FS, data map[string]string) error {
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
		content, err := templatesFS.ReadFile(path)
		if err != nil {
			return err
		}
		tmpl := template.Must(template.New(rel).Parse(string(content)))
		f, err := os.Create(target)
		if err != nil {
			return err
		}
		defer f.Close()
		return tmpl.Execute(f, data)
	})
}
