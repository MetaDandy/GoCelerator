package helper

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func CopyTemplatesFS(templatesFS fs.FS, baseDir, dest string, data map[string]string) error {
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
		content, err := fs.ReadFile(templatesFS, path)
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
