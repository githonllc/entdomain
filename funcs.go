package entdomain

import (
	"strings"
	"text/template"
)

// templateFuncs returns the template function map for code generation.
// Only functions actually invoked by templates/*.tmpl are registered here.
// Internal helper functions used by Go code only are NOT registered.
//
// Source files:
//   - funcs_strings.go:    string manipulation utilities
//   - funcs_fields.go:     field filtering and selection
//   - funcs_scope.go:      scope and requirement checking
//   - funcs_typechecks.go: field type checking
//   - funcs_codegen.go:    code generation helpers
func templateFuncs() template.FuncMap {
	return template.FuncMap{
		// String manipulation
		"lower":    strings.ToLower,
		"hasPrefix": hasPrefix,

		// Field selection (used in template range loops)
		"domainFields":     domainFields,
		"createFields":     createFields,
		"updateFields":     updateFields,
		"responseFields":   responseFields,
		"queryFields":      queryFields,
		"searchableFields": searchableFields,
		"sortableFields":   sortableFields,
		"updateableFields": updateableFields,
		"uniqueLookupFields":    uniqueLookupFields,
		"rangeLookupFields":     rangeLookupFields,
		"nonDefaultDomainFields": nonDefaultDomainFields,

		// Scope and requirement checking
		"isDomainRequired": isDomainRequired,

		// Field type checking
		"isUniqueField": isUniqueField,
		"hasTimeFields": hasTimeFields,
		"hasTimeField":  hasTimeField,

		// Code generation helpers
		"specificMethods":    specificMethods,
		"setFieldCall":       setFieldCall,
		"searchMethod":       searchMethod,
		"findByMethod":       findByMethod,
		"last":               last,

		// Utility functions
		"contains": contains,

		// Template code generation helpers
		"generateEntToDomainFieldAssignment": generateEntToDomainFieldAssignment,
		"generateIdOperation":                generateIdOperation,
		"generateSearchCondition":            generateSearchCondition,
	}
}
