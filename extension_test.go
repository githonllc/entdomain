package entdomain

import (
	"testing"

	"entgo.io/ent/entc/gen"
)

func TestExtension_NewExtension(t *testing.T) {
	config := &ExtensionConfig{
		OutputDir:          "test/output",
		PackageName:        "testpkg",
		GenerateRepository: true,
		GenerateService:    false,
	}

	ext := NewExtension(config)

	if ext.Config.OutputDir != "test/output" {
		t.Errorf("OutputDir = %v, want %v", ext.Config.OutputDir, "test/output")
	}

	if ext.Config.PackageName != "testpkg" {
		t.Errorf("PackageName = %v, want %v", ext.Config.PackageName, "testpkg")
	}

	if !ext.Config.GenerateRepository {
		t.Error("GenerateRepository should be true")
	}

	if ext.Config.GenerateService {
		t.Error("GenerateService should be false")
	}
}

func TestExtension_NewExtensionWithOptions(t *testing.T) {
	ext := NewExtensionWithOptions(
		WithOutputDir("test/output"),
		WithPackageName("testpkg"),
		WithService(true),
		WithRepository(false),
	)

	if ext.Config.OutputDir != "test/output" {
		t.Errorf("OutputDir = %v, want %v", ext.Config.OutputDir, "test/output")
	}

	if ext.Config.PackageName != "testpkg" {
		t.Errorf("PackageName = %v, want %v", ext.Config.PackageName, "testpkg")
	}

	if !ext.Config.GenerateService {
		t.Error("GenerateService should be true")
	}

	if ext.Config.GenerateRepository {
		t.Error("GenerateRepository should be false")
	}
}

func TestExtension_Templates(t *testing.T) {
	ext := NewExtension(&ExtensionConfig{
		GenerateRepository: true,
		GenerateService:    false,
	})

	templates := ext.Templates()

	// Extension uses Hook-based generation, Templates() should return empty slice
	if len(templates) != 0 {
		t.Errorf("Expected 0 templates (Hook-based generation), got %d", len(templates))
	}
}

func TestExtension_Options(t *testing.T) {
	ext := NewExtension(&ExtensionConfig{})

	options := ext.Options()
	if options == nil {
		t.Error("Options should not be nil")
	}
}

func TestExtension_Annotations(t *testing.T) {
	config := &ExtensionConfig{
		OutputDir:   "test/domain",
		PackageName: "testdomain",
	}
	ext := NewExtension(config)
	annotations := ext.Annotations()

	if len(annotations) != 1 {
		t.Errorf("Expected 1 annotation, got %d", len(annotations))
	}

	configAnnotation, ok := annotations[0].(*ConfigAnnotation)
	if !ok {
		t.Errorf("Expected *ConfigAnnotation, got %T", annotations[0])
	}

	if configAnnotation.Config.OutputDir != config.OutputDir {
		t.Errorf("OutputDir = %v, want %v", configAnnotation.Config.OutputDir, config.OutputDir)
	}
}

func TestConfigAnnotation_Name(t *testing.T) {
	annotation := &ConfigAnnotation{}
	if annotation.Name() != "ExtensionConfig" {
		t.Errorf("Name() = %v, want %v", annotation.Name(), "ExtensionConfig")
	}
}

func TestExtensionOptions(t *testing.T) {
	t.Run("WithOutputDir", func(t *testing.T) {
		config := &ExtensionConfig{}
		opt := WithOutputDir("test/dir")
		opt(config)

		if config.OutputDir != "test/dir" {
			t.Errorf("OutputDir = %v, want %v", config.OutputDir, "test/dir")
		}
	})

	t.Run("WithPackageName", func(t *testing.T) {
		config := &ExtensionConfig{}
		opt := WithPackageName("testpkg")
		opt(config)

		if config.PackageName != "testpkg" {
			t.Errorf("PackageName = %v, want %v", config.PackageName, "testpkg")
		}
	})

	t.Run("WithRepository", func(t *testing.T) {
		config := &ExtensionConfig{}
		opt := WithRepository(true)
		opt(config)

		if !config.GenerateRepository {
			t.Error("GenerateRepository should be true")
		}
	})

	t.Run("WithService", func(t *testing.T) {
		config := &ExtensionConfig{}
		opt := WithService(true)
		opt(config)

		if !config.GenerateService {
			t.Error("GenerateService should be true")
		}
	})
}

func TestWithEntDomainPackage(t *testing.T) {
	t.Run("default value", func(t *testing.T) {
		ext := NewExtension(nil)
		want := "github.com/githonllc/entdomain"
		if ext.Config.EntDomainPackage != want {
			t.Errorf("default EntDomainPackage = %q, want %q", ext.Config.EntDomainPackage, want)
		}
	})

	t.Run("custom value via WithEntDomainPackage", func(t *testing.T) {
		config := &ExtensionConfig{}
		opt := WithEntDomainPackage("custom/path")
		opt(config)

		if config.EntDomainPackage != "custom/path" {
			t.Errorf("EntDomainPackage = %q, want %q", config.EntDomainPackage, "custom/path")
		}
	})
}

func TestNewExtension_Defaults(t *testing.T) {
	ext := NewExtension(nil)

	if ext.Config.OutputDir != "." {
		t.Errorf("OutputDir = %q, want %q", ext.Config.OutputDir, ".")
	}

	if ext.Config.PackageName != "domain" {
		t.Errorf("PackageName = %q, want %q", ext.Config.PackageName, "domain")
	}

	if !ext.Config.GenerateRepository {
		t.Error("GenerateRepository should default to true")
	}

	if !ext.Config.GenerateService {
		t.Error("GenerateService should default to true")
	}

	if ext.Config.EntDomainPackage != "github.com/githonllc/entdomain" {
		t.Errorf("EntDomainPackage = %q, want %q", ext.Config.EntDomainPackage, "github.com/githonllc/entdomain")
	}
}

func TestExtension_Hooks(t *testing.T) {
	ext := NewExtension(nil)
	hooks := ext.Hooks()

	if len(hooks) != 1 {
		t.Errorf("Hooks() returned %d hooks, want exactly 1", len(hooks))
	}
}

func TestExtension_TemplateFuncMap(t *testing.T) {
	customPkg := "my/custom/entdomain"
	ext := NewExtension(&ExtensionConfig{
		EntDomainPackage: customPkg,
	})

	funcMap := ext.templateFuncMap()

	// Verify that gen.Funcs entries are included
	for key := range gen.Funcs {
		if _, ok := funcMap[key]; !ok {
			t.Errorf("templateFuncMap() is missing gen.Funcs key %q", key)
		}
	}

	// Verify that custom templateFuncs entries are included
	customKeys := templateFuncs()
	for key := range customKeys {
		if _, ok := funcMap[key]; !ok {
			t.Errorf("templateFuncMap() is missing custom templateFuncs key %q", key)
		}
	}

	// Verify "entdomainPkg" function exists
	entdomainPkgFn, ok := funcMap["entdomainPkg"]
	if !ok {
		t.Fatal("templateFuncMap() is missing the 'entdomainPkg' function")
	}

	// Verify calling entdomainPkg returns the configured package path
	fn, ok := entdomainPkgFn.(func() string)
	if !ok {
		t.Fatalf("entdomainPkg has unexpected type %T, want func() string", entdomainPkgFn)
	}
	got := fn()
	if got != customPkg {
		t.Errorf("entdomainPkg() = %q, want %q", got, customPkg)
	}

	// Verify it does not mutate the global gen.Funcs map
	// Snapshot gen.Funcs keys before calling templateFuncMap
	genFuncsBefore := make(map[string]bool, len(gen.Funcs))
	for k := range gen.Funcs {
		genFuncsBefore[k] = true
	}

	// Call templateFuncMap multiple times
	_ = ext.templateFuncMap()
	_ = ext.templateFuncMap()

	// "entdomainPkg" is injected by templateFuncMap and must NOT leak into gen.Funcs
	if _, exists := gen.Funcs["entdomainPkg"]; exists {
		t.Error("templateFuncMap() mutated global gen.Funcs: found 'entdomainPkg' key")
	}

	// gen.Funcs should have exactly the same keys as before
	for k := range gen.Funcs {
		if !genFuncsBefore[k] {
			t.Errorf("templateFuncMap() added unexpected key %q to global gen.Funcs", k)
		}
	}
}

func TestConfigAnnotation_NameRenamed(t *testing.T) {
	annotation := ConfigAnnotation{}
	got := annotation.Name()
	want := "ExtensionConfig"
	if got != want {
		t.Errorf("ConfigAnnotation.Name() = %q, want %q (was renamed from DomainConfig)", got, want)
	}
}

func TestNewExtensionWithOptions_EntDomainPackage(t *testing.T) {
	customPkg := "github.com/myorg/myentdomain"
	ext := NewExtensionWithOptions(
		WithEntDomainPackage(customPkg),
	)

	if ext.Config.EntDomainPackage != customPkg {
		t.Errorf("EntDomainPackage = %q, want %q", ext.Config.EntDomainPackage, customPkg)
	}

	// Other defaults should still apply
	if ext.Config.OutputDir != "." {
		t.Errorf("OutputDir = %q, want %q", ext.Config.OutputDir, ".")
	}
	if ext.Config.PackageName != "domain" {
		t.Errorf("PackageName = %q, want %q", ext.Config.PackageName, "domain")
	}
	if !ext.Config.GenerateRepository {
		t.Error("GenerateRepository should default to true")
	}
	if !ext.Config.GenerateService {
		t.Error("GenerateService should default to true")
	}
}
