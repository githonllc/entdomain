package entdomain

import (
	"testing"

	"entgo.io/ent/entc/gen"
)

func TestExtension_NewExtension(t *testing.T) {
	config := &ExtensionConfig{
		GenerateBaseService: true,
	}

	ext := NewExtension(config)

	if !ext.Config.GenerateBaseService {
		t.Error("GenerateBaseService should be true")
	}
}

func TestExtension_NewExtensionWithOptions(t *testing.T) {
	ext := NewExtensionWithOptions(
		WithBaseService(true),
		WithBaseHandler(true),
	)

	if !ext.Config.GenerateBaseService {
		t.Error("GenerateBaseService should be true")
	}

	if !ext.Config.GenerateBaseHandler {
		t.Error("GenerateBaseHandler should be true")
	}
}

func TestExtension_Templates(t *testing.T) {
	ext := NewExtension(&ExtensionConfig{
		GenerateBaseService: true,
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
		GenerateBaseService: true,
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

	if !configAnnotation.Config.GenerateBaseService {
		t.Error("GenerateBaseService should be true")
	}
}

func TestConfigAnnotation_Name(t *testing.T) {
	annotation := &ConfigAnnotation{}
	if annotation.Name() != "ExtensionConfig" {
		t.Errorf("Name() = %v, want %v", annotation.Name(), "ExtensionConfig")
	}
}

func TestExtensionOptions(t *testing.T) {
	t.Run("WithBaseService", func(t *testing.T) {
		config := &ExtensionConfig{}
		opt := WithBaseService(true)
		opt(config)

		if !config.GenerateBaseService {
			t.Error("GenerateBaseService should be true")
		}
	})

	t.Run("WithBaseHandler", func(t *testing.T) {
		config := &ExtensionConfig{}
		opt := WithBaseHandler(true)
		opt(config)

		if !config.GenerateBaseHandler {
			t.Error("GenerateBaseHandler should be true")
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

	// Verify "entdomainPkg" function exists and returns correct value
	entdomainPkgFn, ok := funcMap["entdomainPkg"]
	if !ok {
		t.Fatal("templateFuncMap() is missing the 'entdomainPkg' function")
	}

	fn, ok := entdomainPkgFn.(func() string)
	if !ok {
		t.Fatalf("entdomainPkg has unexpected type %T, want func() string", entdomainPkgFn)
	}
	got := fn()
	if got != customPkg {
		t.Errorf("entdomainPkg() = %q, want %q", got, customPkg)
	}

	// Verify it does not mutate the global gen.Funcs map
	genFuncsBefore := make(map[string]bool, len(gen.Funcs))
	for k := range gen.Funcs {
		genFuncsBefore[k] = true
	}

	_ = ext.templateFuncMap()
	_ = ext.templateFuncMap()

	if _, exists := gen.Funcs["entdomainPkg"]; exists {
		t.Error("templateFuncMap() mutated global gen.Funcs: found 'entdomainPkg' key")
	}

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
}
