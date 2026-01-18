package templates

import (
	"embed"
	"html/template"
	"io/fs"
	"path/filepath"
)

//go:embed assets/*
var templatesFS embed.FS

func ParseTemplates() (*template.Template, error) {
	funcs := template.FuncMap{}

	root := template.New("root").Funcs(funcs)

	files, err := templatesFS.ReadDir("assets")
	if err != nil {
		return nil, err
	}

	for _, entry := range files {
		if entry.IsDir() {
			continue
		}

		if filepath.Ext(entry.Name()) != ".gohtml" {
			continue
		}

		content, err := templatesFS.ReadFile(filepath.Join("assets", entry.Name()))
		if err != nil {
			return nil, err
		}

		if _, err := root.New(entry.Name()).Parse(string(content)); err != nil {
			return nil, err
		}
	}

	return root, nil
}

func StaticFS() (fs.FS, error) {
	return fs.Sub(templatesFS, "assets")
}
