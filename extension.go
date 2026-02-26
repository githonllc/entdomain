package entdomain

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"golang.org/x/tools/imports"
)

// Extension is the entdomain Ent extension.
type Extension struct {
	// Config holds the extension configuration.
	Config *ExtensionConfig
}

// ExtensionConfig holds configuration for the extension
type ExtensionConfig struct {
	// OutputDir is the output directory for generated files
	OutputDir string

	// PackageName is the Go package name for generated files
	PackageName string

	// GenerateRepository controls whether repository interfaces are generated
	GenerateRepository bool

	// GenerateService controls whether service interfaces are generated
	GenerateService bool

	// EntDomainPackage is the import path for the entdomain package
	// Default: "github.com/githonllc/entdomain"
	EntDomainPackage string
}

const defaultEntDomainPackage = "github.com/githonllc/entdomain"

// NewExtension creates a new extension instance
func NewExtension(config *ExtensionConfig) *Extension {
	if config == nil {
		config = &ExtensionConfig{
			OutputDir:          ".",
			PackageName:        "domain",
			GenerateRepository: true,
			GenerateService:    true,
		}
	}
	if config.EntDomainPackage == "" {
		config.EntDomainPackage = defaultEntDomainPackage
	}

	return &Extension{
		Config: config,
	}
}

// Hooks returns the extension's hooks — uses a Hook to generate separate files per Type.
func (e *Extension) Hooks() []gen.Hook {
	return []gen.Hook{
		e.generatePerTypeFiles, // main generation logic
	}
}

// Templates returns an empty template list — the old GraphTemplate approach is no longer used.
func (e *Extension) Templates() []*gen.Template {
	return []*gen.Template{} // removed legacy GraphTemplate generation
}

// generatePerTypeFiles is the core Hook that generates separate files for each Type.
func (e *Extension) generatePerTypeFiles(next gen.Generator) gen.Generator {
	return gen.GenerateFunc(func(g *gen.Graph) error {
		// Run the standard generation first
		if err := next.Generate(g); err != nil {
			return err
		}

		// Generate separate domain files for each Type
		for _, node := range g.Nodes {
			// Generate domain model file
			if err := e.generateDomainModelFile(g, node); err != nil {
				return fmt.Errorf("failed to generate %s domain model: %w", node.Name, err)
			}

			// Generate repository file
			if e.Config.GenerateRepository {
				if err := e.generateRepositoryFile(g, node); err != nil {
					return fmt.Errorf("failed to generate %s repository file: %w", node.Name, err)
				}
			}

			// Generate service file
			if e.Config.GenerateService {
				if err := e.generateServiceFile(g, node); err != nil {
					return fmt.Errorf("failed to generate %s service file: %w", node.Name, err)
				}
			}
		}

		// Clean up legacy single-file outputs
		if err := e.cleanupOldFiles(g); err != nil {
			return fmt.Errorf("failed to clean up old files: %w", err)
		}

		return nil
	})
}

// generateDomainModelFile generates a domain model file for a single Type.
func (e *Extension) generateDomainModelFile(g *gen.Graph, node *gen.Type) error {
	// Parse domain model template
	tmpl, err := template.New("domain_model").
		Funcs(e.templateFuncMap()).
		Parse(domainModelTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse domain model template: %w", err)
	}

	// Render template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, node); err != nil {
		return fmt.Errorf("failed to render domain model template: %w", err)
	}

	// Write file
	filename := fmt.Sprintf("%s_domain_model.go", strings.ToLower(node.Name))
	outputPath := filepath.Join(g.Config.Target, filename)

	return writeFile(outputPath, buf.Bytes())
}

// generateRepositoryFile generates a repository file for a single Type.
func (e *Extension) generateRepositoryFile(g *gen.Graph, node *gen.Type) error {
	// Parse repository template
	tmpl, err := template.New("repository").
		Funcs(e.templateFuncMap()).
		Parse(repositoryTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse repository template: %w", err)
	}

	// Render template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, node); err != nil {
		return fmt.Errorf("failed to render repository template: %w", err)
	}

	// Write file
	filename := fmt.Sprintf("%s_domain_repository.go", strings.ToLower(node.Name))
	outputPath := filepath.Join(g.Config.Target, filename)

	return writeFile(outputPath, buf.Bytes())
}

// generateServiceFile generates a service file for a single Type
func (e *Extension) generateServiceFile(g *gen.Graph, node *gen.Type) error {
	// Parse service template
	tmpl, err := template.New("service").
		Funcs(e.templateFuncMap()).
		Parse(serviceTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse service template: %w", err)
	}

	// Render template
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, node); err != nil {
		return fmt.Errorf("failed to render service template: %w", err)
	}

	// Write file
	filename := fmt.Sprintf("%s_domain_service.go", strings.ToLower(node.Name))
	outputPath := filepath.Join(g.Config.Target, filename)

	return writeFile(outputPath, buf.Bytes())
}

// cleanupOldFiles removes legacy single-file outputs that are no longer generated.
func (e *Extension) cleanupOldFiles(g *gen.Graph) error {
	oldFiles := []string{
		"domain_model.go",
		"service.go",
		"repository.go",
	}

	for _, oldFile := range oldFiles {
		oldPath := filepath.Join(g.Config.Target, oldFile)
		if err := os.Remove(oldPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to delete old file %s: %w", oldFile, err)
		}
	}

	// Clean up HTTP types files
	for _, node := range g.Nodes {
		httpTypesFile := fmt.Sprintf("%s_http_types.go", strings.ToLower(node.Name))
		httpTypesPath := filepath.Join(g.Config.Target, httpTypesFile)
		if err := os.Remove(httpTypesPath); err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("failed to delete HTTP types file %s: %w", httpTypesFile, err)
		}
	}

	return nil
}

// writeFile formats the generated Go source with goimports and writes it to disk
func writeFile(path string, content []byte) error {
	// Format with goimports to fix imports and formatting
	formatted, err := imports.Process(path, content, nil)
	if err != nil {
		log.Printf("WARNING: goimports formatting failed for %s: %v (writing unformatted)", path, err)
		formatted = content
	}
	if err := os.WriteFile(path, formatted, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", path, err)
	}
	return nil
}

// templateFuncMap returns the combined template function map with Ent standard functions
func (e *Extension) templateFuncMap() template.FuncMap {
	// Create a copy of gen.Funcs to avoid mutating the global map
	funcs := make(template.FuncMap, len(gen.Funcs))
	for k, v := range gen.Funcs {
		funcs[k] = v
	}

	// Add custom functions
	for k, v := range templateFuncs() {
		funcs[k] = v
	}

	// Inject configurable entdomain package path
	pkg := e.Config.EntDomainPackage
	funcs["entdomainPkg"] = func() string { return pkg }

	return funcs
}

// Option is a function type for configuring the extension.
type Option func(*ExtensionConfig)

// WithOutputDir sets the output directory.
func WithOutputDir(dir string) Option {
	return func(c *ExtensionConfig) {
		c.OutputDir = dir
	}
}

// WithPackageName sets the package name.
func WithPackageName(name string) Option {
	return func(c *ExtensionConfig) {
		c.PackageName = name
	}
}

// WithRepository controls whether repository interfaces are generated.
func WithRepository(generate bool) Option {
	return func(c *ExtensionConfig) {
		c.GenerateRepository = generate
	}
}

// WithService controls whether service interfaces are generated
func WithService(generate bool) Option {
	return func(c *ExtensionConfig) {
		c.GenerateService = generate
	}
}

// WithEntDomainPackage sets the import path for the entdomain package
func WithEntDomainPackage(pkg string) Option {
	return func(c *ExtensionConfig) {
		c.EntDomainPackage = pkg
	}
}

// NewExtensionWithOptions creates a new extension using functional options.
func NewExtensionWithOptions(opts ...Option) *Extension {
	config := &ExtensionConfig{
		OutputDir:          ".",
		PackageName:        "domain",
		GenerateRepository: true,
		GenerateService:    true,
	}

	for _, opt := range opts {
		opt(config)
	}

	return NewExtension(config)
}

// Annotations returns global annotations for the extension
func (e *Extension) Annotations() []entc.Annotation {
	return []entc.Annotation{
		&ConfigAnnotation{Config: e.Config},
	}
}

// Options returns the extension options
func (e *Extension) Options() []entc.Option {
	return []entc.Option{}
}

// ConfigAnnotation implements entc.Annotation for extension configuration
type ConfigAnnotation struct {
	Config *ExtensionConfig
}

// Name returns the annotation name.
// NOTE: Renamed from "DomainConfig" to "ExtensionConfig" to avoid confusion
// with the DomainConfig schema-level annotation in annotations.go.
func (ConfigAnnotation) Name() string {
	return "ExtensionConfig"
}
