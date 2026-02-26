package entdomain

import (
	"entgo.io/ent/entc/gen"
)

// domainFields returns all fields with DomainField annotation
func domainFields(node *gen.Type) []*gen.Field {
	var fields []*gen.Field
	for _, field := range node.Fields {
		if annotation := getDomainFieldAnnotation(field); annotation != nil {
			fields = append(fields, field)
		}
	}
	return fields
}

// createFields returns fields that can be used in create requests
func createFields(node *gen.Type) []*gen.Field {
	var fields []*gen.Field
	for _, field := range node.Fields {
		if annotation := getDomainFieldAnnotation(field); annotation != nil {
			if hasDomainScope(field, ScopeCreate) {
				fields = append(fields, field)
			}
		}
	}
	return fields
}

// updateFields returns fields that can be used in update requests
func updateFields(node *gen.Type) []*gen.Field {
	var fields []*gen.Field
	for _, field := range node.Fields {
		if annotation := getDomainFieldAnnotation(field); annotation != nil {
			if hasDomainScope(field, ScopeUpdate) {
				fields = append(fields, field)
			}
		}
	}
	return fields
}

// responseFields returns fields that can be used in responses
func responseFields(node *gen.Type) []*gen.Field {
	var fields []*gen.Field
	for _, field := range node.Fields {
		if annotation := getDomainFieldAnnotation(field); annotation != nil {
			if hasDomainScope(field, ScopeResponse) {
				fields = append(fields, field)
			}
		}
	}
	return fields
}

// queryFields returns fields that can be used for searching
func queryFields(node *gen.Type) []*gen.Field {
	var fields []*gen.Field
	for _, field := range node.Fields {
		annotation := getDomainFieldAnnotation(field)
		if annotation != nil {
			// Check if field is searchable OR has ScopeQuery
			if annotation.Searchable || hasDomainScope(field, ScopeQuery) {
				fields = append(fields, field)
			}
		}
	}
	return fields
}

// searchableFields returns fields that can be searched
func searchableFields(node *gen.Type) []*gen.Field {
	var fields []*gen.Field
	for _, field := range node.Fields {
		annotation := getDomainFieldAnnotation(field)
		if annotation != nil && annotation.Searchable {
			fields = append(fields, field)
		}
	}
	return fields
}

// sortableFields returns fields that can be sorted
func sortableFields(node *gen.Type) []*gen.Field {
	var fields []*gen.Field
	for _, field := range node.Fields {
		annotation := getDomainFieldAnnotation(field)
		if annotation != nil && annotation.Sortable {
			// Filter out complex field types that do not support sorting
			if !isComplexFieldType(field.Type.String()) {
				fields = append(fields, field)
			}
		}
	}
	return fields
}

// updateableFields returns all fields that can be updated in Repository layer operations.
// Excludes: ID field and immutable fields.
func updateableFields(node *gen.Type) []*gen.Field {
	var fields []*gen.Field
	for _, field := range node.Fields {
		// Skip the ID field
		if field.Name == "id" {
			continue
		}

		// Skip immutable fields
		if field.Immutable {
			continue
		}

		if annotation := getDomainFieldAnnotation(field); annotation != nil {
			// The Repository layer can update all non-immutable annotated fields
			fields = append(fields, field)
		}
	}
	return fields
}

// nonDefaultDomainFields returns all annotated fields that do NOT have a
// schema-level Default value (e.g., Default(time.Now)). Used by Create
// templates to skip fields that should use Ent's default value hook.
func nonDefaultDomainFields(node *gen.Type) []*gen.Field {
	var fields []*gen.Field
	for _, field := range domainFields(node) {
		if !field.Default {
			fields = append(fields, field)
		}
	}
	return fields
}

// uniqueLookupFields returns all fields with UniqueLookup annotation
func uniqueLookupFields(node *gen.Type) []*gen.Field {
	var fields []*gen.Field
	for _, field := range domainFields(node) {
		if isUniqueLookupField(field) {
			fields = append(fields, field)
		}
	}
	return fields
}

// rangeLookupFields returns all fields with RangeLookup annotation
func rangeLookupFields(node *gen.Type) []*gen.Field {
	var fields []*gen.Field
	for _, field := range domainFields(node) {
		annotation := getDomainFieldAnnotation(field)
		if annotation != nil && annotation.RangeLookup {
			fields = append(fields, field)
		}
	}
	return fields
}
