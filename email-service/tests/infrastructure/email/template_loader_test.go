package email_test

import (
	"os"
	"path/filepath"
	"testing"

	"email-service/internal/infrastructure/email"

	"github.com/stretchr/testify/assert"
)

func TestHTMLTemplateLoader_Render(t *testing.T) {
	tmpDir := t.TempDir()
	templateName := "test_template.html"
	templatePath := filepath.Join(tmpDir, templateName)

	content := `Hello, {{.name}}!`
	err := os.WriteFile(templatePath, []byte(content), 0644)
	assert.NoError(t, err)

	loader := email.NewHTMLTemplateLoader(tmpDir)

	result, err := loader.Render(templateName, map[string]string{
		"name": "Alice",
	})

	assert.NoError(t, err)
	assert.Equal(t, "Hello, Alice!", result)
}

func TestHTMLTemplateLoader_Render_FileNotFound(t *testing.T) {
	loader := email.NewHTMLTemplateLoader("nonexistent")
	_, err := loader.Render("nope.html", nil)
	assert.Error(t, err)
}
