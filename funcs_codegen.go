package entdomain

import (
	"fmt"
	"strings"

	"entgo.io/ent/entc/gen"
)

// generateEntToDomainFieldAssignment generates field assignment for entToDomain conversion
func generateEntToDomainFieldAssignment(field *gen.Field) string {
	fieldName := field.StructField()

	// Skip ID field as it's handled separately
	if field.Name == "id" {
		return ""
	}

	// Handle nullable time fields specially
	if field.Nillable && strings.Contains(field.Type.Type.String(), "time.Time") {
		return fmt.Sprintf(`		%s: func() time.Time {
			if entity.%s != nil {
				return *entity.%s
			}
			return time.Time{}
		}(),`, fieldName, fieldName, fieldName)
	}

	// Regular field assignment
	return fmt.Sprintf("		%s: entity.%s,", fieldName, fieldName)
}

// getEntityPackageName returns the correct ent package name for an entity
func getEntityPackageName(node *gen.Type) string {
	// The ent-generated package name is the lowercase form of the entity name
	return strings.ToLower(node.Name)
}

// generateSearchCondition generates search condition for searchable string fields
func generateSearchCondition(field *gen.Field, node *gen.Type) string {
	if field.Type.String() == "string" {
		packageName := getEntityPackageName(node)
		return fmt.Sprintf("		predicates = append(predicates, %s.%sContains(req.Query))", packageName, field.StructField())
	}
	return ""
}

// generateIdOperation generates ID-related operations for the given type
func generateIdOperation(node *gen.Type, operation string, idVar string) string {
	idType := node.ID.Type.String()
	entityName := node.Name
	packageName := getEntityPackageName(node)

	switch operation {
	case "get":
		if idType == "string" {
			return fmt.Sprintf(`entity, err := r.client.%s.Get(ctx, %s.String())`, entityName, idVar)
		} else if idType == "int64" {
			return fmt.Sprintf(`entity, err := r.client.%s.Get(ctx, func() int64 {
		if i, err := %s.Int64(); err == nil {
			return i
		}
		return 0
	}())`, entityName, idVar)
		} else {
			return fmt.Sprintf(`entity, err := r.client.%s.Get(ctx, %s)`, entityName, idVar)
		}
	case "updateOneID":
		return fmt.Sprintf(`entity, err := r.client.%s.UpdateOneID(model.%s).`, entityName, node.ID.StructField())
	case "delete":
		if idType == "string" {
			return fmt.Sprintf(`return r.client.%s.DeleteOneID(%s.String()).Exec(ctx)`, entityName, idVar)
		} else if idType == "int64" {
			return fmt.Sprintf(`return r.client.%s.DeleteOneID(func() int64 {
		if i, err := %s.Int64(); err == nil {
			return i
		}
		return 0
	}()).Exec(ctx)`, entityName, idVar)
		} else {
			return fmt.Sprintf(`return r.client.%s.DeleteOneID(%s).Exec(ctx)`, entityName, idVar)
		}
	case "exists":
		if idType == "string" {
			return fmt.Sprintf(`count, err := r.client.%s.Query().Where(%s.IDEQ(%s.String())).Count(ctx)`, entityName, packageName, idVar)
		} else if idType == "int64" {
			return fmt.Sprintf(`count, err := r.client.%s.Query().Where(%s.IDEQ(func() int64 {
		if i, err := %s.Int64(); err == nil {
			return i
		}
		return 0
	}())).Count(ctx)`, entityName, packageName, idVar)
		} else {
			return fmt.Sprintf(`count, err := r.client.%s.Query().Where(%s.IDEQ(%s)).Count(ctx)`, entityName, packageName, idVar)
		}
	case "batchDelete":
		if idType == "string" {
			return fmt.Sprintf(`stringIds := make([]string, len(ids))
	for i, id := range ids {
		stringIds[i] = id.String()
	}
	_, err := r.client.%s.Delete().Where(%s.IDIn(stringIds...)).Exec(ctx)`, entityName, packageName)
		} else if idType == "int64" {
			return fmt.Sprintf(`int64Ids := make([]int64, len(ids))
	for i, id := range ids {
		if intId, err := id.Int64(); err == nil {
			int64Ids[i] = intId
		}
	}
	_, err := r.client.%s.Delete().Where(%s.IDIn(int64Ids...)).Exec(ctx)`, entityName, packageName)
		} else {
			return fmt.Sprintf(`_, err := r.client.%s.Delete().Where(%s.IDIn(ids...)).Exec(ctx)`, entityName, packageName)
		}
	default:
		return fmt.Sprintf("// Unknown operation: %s", operation)
	}
}

// setFieldCall generates a setter method call for a field (e.g., "SetName(model.Name)").
// Used by both create and update template operations.
func setFieldCall(field *gen.Field, _ *gen.Type) string {
	return fmt.Sprintf("Set%s(model.%s)", field.StructField(), field.StructField())
}

// fieldPredicate generates a type-assertion + Where predicate for a field.
// indent controls the indentation level of the generated code block.
// When skipEmpty is true, string checks include `&& v != ""`.
func fieldPredicate(field *gen.Field, node *gen.Type, indent string, skipEmpty bool) string {
	pkg := getEntityPackageName(node)
	name := field.StructField()
	ft := field.Type.String()

	where := func(cast, goType string) string {
		return fmt.Sprintf(`%sif v, ok := value.(%s); ok {
%s	query = query.Where(%s.%sEQ(v))
%s}`, indent, cast, indent, pkg, name, indent)
	}

	switch {
	case field.IsEnum():
		enumType := fmt.Sprintf("%s.%s", pkg, name)
		extra := ""
		if skipEmpty {
			extra = ` && v != ""`
		}
		// Try concrete enum type first (e.g., person.Gender), then fall back to string.
		// Go type assertions don't match underlying types, so both branches are needed.
		return fmt.Sprintf(`%sif v, ok := value.(%s); ok {
%s	query = query.Where(%s.%sEQ(v))
%s} else if v, ok := value.(string); ok%s {
%s	query = query.Where(%s.%sEQ(%s(v)))
%s}`, indent, enumType, indent, pkg, name, indent, extra, indent, pkg, name, enumType, indent)
	case ft == "string":
		extra := ""
		if skipEmpty {
			extra = ` && v != ""`
		}
		return fmt.Sprintf(`%sif v, ok := value.(string); ok%s {
%s	query = query.Where(%s.%sEQ(v))
%s}`, indent, extra, indent, pkg, name, indent)
	case ft == "int":
		return fmt.Sprintf(`%sif v, ok := value.(int); ok {
%s	query = query.Where(%s.%sEQ(v))
%s} else if v, ok := value.(int64); ok {
%s	query = query.Where(%s.%sEQ(int(v)))
%s}`, indent, indent, pkg, name, indent, indent, pkg, name, indent)
	case ft == "int32":
		return fmt.Sprintf(`%sif v, ok := value.(int32); ok {
%s	query = query.Where(%s.%sEQ(v))
%s} else if v, ok := value.(int64); ok {
%s	query = query.Where(%s.%sEQ(int32(v)))
%s}`, indent, indent, pkg, name, indent, indent, pkg, name, indent)
	case ft == "int64":
		return where("int64", ft)
	case ft == "bool":
		return where("bool", ft)
	case ft == "time.Time":
		return where("time.Time", ft)
	default:
		return fmt.Sprintf("%s// unsupported field type: %s", indent, ft)
	}
}

// searchMethod generates a filter predicate for Search/Count methods (nested indentation, skips empty strings).
func searchMethod(field *gen.Field, node *gen.Type) string {
	return fieldPredicate(field, node, "\t\t\t", true)
}

// findByMethod generates a filter predicate for FindBy methods (standard indentation).
func findByMethod(field *gen.Field, node *gen.Type) string {
	return fieldPredicate(field, node, "\t\t", false)
}

// last checks if this is the last element in slice
func last(slice []*gen.Field) *gen.Field {
	if len(slice) == 0 {
		return nil
	}
	return slice[len(slice)-1]
}

// specificMethods generates specific service/repository methods based on explicit
// field annotations only. Methods are generated for:
//   - AsUniqueLookup() → FindByX(ctx, value) (*Model, error)
//   - AsRangeLookup()  → FindByXRange(ctx, start, end time.Time) ([]*Model, error)
//
// No auto-generation for searchable string, enum, or bool fields.
func specificMethods(node *gen.Type) []string {
	var methods []string
	generated := make(map[string]bool)

	// Generate FindByX for fields annotated with AsUniqueLookup
	for _, field := range uniqueLookupFields(node) {
		fieldName := field.StructField()
		methodKey := fmt.Sprintf("FindBy%s", fieldName)
		if !generated[methodKey] {
			methodName := fmt.Sprintf("FindBy%s(ctx context.Context, %s %s) (*%sDomainModel, error)",
				fieldName,
				strings.ToLower(fieldName),
				field.Type.String(),
				node.Name)
			methods = append(methods, methodName)
			generated[methodKey] = true
		}
	}

	// Generate FindByXRange for fields annotated with AsRangeLookup
	for _, field := range rangeLookupFields(node) {
		methodKey := fmt.Sprintf("FindBy%sRange", field.StructField())
		if !generated[methodKey] {
			methodName := fmt.Sprintf("FindBy%sRange(ctx context.Context, start, end time.Time) ([]*%sDomainModel, error)",
				field.StructField(),
				node.Name)
			methods = append(methods, methodName)
			generated[methodKey] = true
		}
	}

	return methods
}
