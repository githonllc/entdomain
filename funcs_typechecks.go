package entdomain

import (
	"strings"

	"entgo.io/ent/entc/gen"
)

// isUniqueField checks if a field has a unique constraint via ent's Unique() builder.
func isUniqueField(field *gen.Field) bool {
	return field.Unique
}

// isTimeField checks if a field is a time field.
func isTimeField(field *gen.Field) bool {
	return strings.Contains(field.Type.String(), "time.Time")
}

// hasTimeFields checks if the entity has any time fields.
func hasTimeFields(node *gen.Type) bool {
	for _, field := range domainFields(node) {
		if isTimeField(field) {
			return true
		}
	}
	return false
}

// hasTimeField checks if the entity has a specific named time field.
func hasTimeField(node *gen.Type, fieldName string) bool {
	for _, field := range domainFields(node) {
		if strings.ToLower(field.Name) == fieldName && isTimeField(field) {
			return true
		}
	}
	return false
}

// isUniqueLookupField checks if a field is annotated with UniqueLookup.
func isUniqueLookupField(field *gen.Field) bool {
	annotation := getDomainFieldAnnotation(field)
	if annotation == nil {
		return false
	}
	return annotation.UniqueLookup
}

// isUUIDType checks if the given type string represents a UUID type.
func isUUIDType(typeStr string) bool {
	return strings.Contains(strings.ToLower(typeStr), "uuid")
}

// hasSoftDelete checks if an entity has a deleted_at field (convention-based soft-delete detection).
// Returns true if the entity has a Nillable, Optional time.Time field named "deleted_at".
func hasSoftDelete(node *gen.Type) bool {
	for _, field := range node.Fields {
		if field.Name == "deleted_at" && isTimeField(field) && field.Nillable {
			return true
		}
	}
	return false
}

// isComplexFieldType checks if a field type is too complex for basic
// operations like sorting (slices, maps, JSON types).
func isComplexFieldType(fieldType string) bool {
	return strings.HasPrefix(fieldType, "[]") ||
		strings.HasPrefix(fieldType, "map[") ||
		strings.Contains(fieldType, "json.")
}
