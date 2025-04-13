package email

import (
	"bytes"
	"html/template"
	"path/filepath"
)

type HTMLTemplateLoader struct {
	basePath string
}

func NewHTMLTemplateLoader(basePath string) *HTMLTemplateLoader {
	return &HTMLTemplateLoader{basePath: basePath}
}

func (f *HTMLTemplateLoader) Render(name string, data any) (string, error) {
	templatePath := filepath.Join(f.basePath, name)
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
