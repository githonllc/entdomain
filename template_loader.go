package entdomain

import (
	"embed"
	"fmt"
	"path/filepath"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

// loadTemplate reads a named template from the embedded filesystem.
// The name should not include the "templates/" prefix or ".tmpl" suffix.
func loadTemplate(name string) (string, error) {
	filename := filepath.Join("templates", name+".tmpl")
	content, err := templateFS.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to load template %s: %w", filename, err)
	}
	return string(content), nil
}

// mustLoadTemplate loads a named template and panics on failure.
// Use this for templates that are required at package init time.
func mustLoadTemplate(name string) string {
	content, err := loadTemplate(name)
	if err != nil {
		panic(err)
	}
	return content
}
