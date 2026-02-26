package entdomain

import (
	"strings"
	"testing"

	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
)

func TestGetEntityPackageName(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect string
	}{
		{"simple name", "User", "user"},
		{"multi-word", "UserProfile", "userprofile"},
		{"already lower", "item", "item"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := newTestType(tt.input)
			got := getEntityPackageName(node)
			if got != tt.expect {
				t.Errorf("getEntityPackageName(%q) = %q, want %q", tt.input, got, tt.expect)
			}
		})
	}
}

func TestSetFieldCall(t *testing.T) {
	f := newStringField("first_name", nil)
	node := newTestType("User")

	got := setFieldCall(f, node)
	// StructField() for "first_name" is "FirstName"
	expected := "SetFirstName(model.FirstName)"
	if got != expected {
		t.Errorf("setFieldCall() = %q, want %q", got, expected)
	}
}

func TestLast(t *testing.T) {
	f1 := newStringField("a", nil)
	f2 := newStringField("b", nil)

	if last(nil) != nil {
		t.Error("last(nil) should return nil")
	}
	if last([]*gen.Field{}) != nil {
		t.Error("last([]) should return nil")
	}
	if last([]*gen.Field{f1}) != f1 {
		t.Error("last([f1]) should return f1")
	}
	if last([]*gen.Field{f1, f2}) != f2 {
		t.Error("last([f1,f2]) should return f2")
	}
}

func TestFieldPredicate_String(t *testing.T) {
	f := newStringField("name", nil)
	node := newTestType("User")

	got := fieldPredicate(f, node, "\t", false)
	assertContains(t, got, `value.(string)`)
	assertContains(t, got, `user.NameEQ(v)`)
	assertNotContains(t, got, `v != ""`)
}

func TestFieldPredicate_StringSkipEmpty(t *testing.T) {
	f := newStringField("name", nil)
	node := newTestType("User")

	got := fieldPredicate(f, node, "\t", true)
	assertContains(t, got, `v != ""`)
}

func TestFieldPredicate_Int(t *testing.T) {
	f := newIntField("age", nil)
	node := newTestType("User")

	got := fieldPredicate(f, node, "\t", false)
	assertContains(t, got, `value.(int)`)
	assertContains(t, got, `user.AgeEQ(v)`)
	// int fields also have int64 fallback
	assertContains(t, got, `value.(int64)`)
	assertContains(t, got, `int(v)`)
}

func TestFieldPredicate_Int32(t *testing.T) {
	f := newInt32Field("priority", nil)
	node := newTestType("Task")

	got := fieldPredicate(f, node, "\t", false)
	assertContains(t, got, `value.(int32)`)
	assertContains(t, got, `task.PriorityEQ(v)`)
	assertContains(t, got, `value.(int64)`)
	assertContains(t, got, `int32(v)`)
}

func TestFieldPredicate_Int64(t *testing.T) {
	f := newInt64Field("count", nil)
	node := newTestType("Item")

	got := fieldPredicate(f, node, "\t", false)
	assertContains(t, got, `value.(int64)`)
	assertContains(t, got, `item.CountEQ(v)`)
}

func TestFieldPredicate_Bool(t *testing.T) {
	f := newBoolField("active", nil)
	node := newTestType("User")

	got := fieldPredicate(f, node, "\t", false)
	assertContains(t, got, `value.(bool)`)
	assertContains(t, got, `user.ActiveEQ(v)`)
}

func TestFieldPredicate_Time(t *testing.T) {
	f := newTimeField("created_at", nil)
	node := newTestType("User")

	got := fieldPredicate(f, node, "\t", false)
	assertContains(t, got, `value.(time.Time)`)
	assertContains(t, got, `user.CreatedAtEQ(v)`)
}

func TestFieldPredicate_Enum(t *testing.T) {
	f := newEnumField("status", nil)
	node := newTestType("Order")

	got := fieldPredicate(f, node, "\t", false)
	// Enum: concrete type first, then string fallback
	assertContains(t, got, `value.(order.Status)`)
	assertContains(t, got, `value.(string)`)
	assertContains(t, got, `order.StatusEQ(v)`)
	assertContains(t, got, `order.Status(v)`)
}

func TestFieldPredicate_EnumSkipEmpty(t *testing.T) {
	f := newEnumField("status", nil)
	node := newTestType("Order")

	got := fieldPredicate(f, node, "\t", true)
	assertContains(t, got, `v != ""`)
}

func TestFieldPredicate_UnsupportedType(t *testing.T) {
	f := newField("data", &field.TypeInfo{Type: field.TypeJSON, Ident: "json.RawMessage"}, nil)
	node := newTestType("Item")

	got := fieldPredicate(f, node, "\t", false)
	assertContains(t, got, "unsupported field type")
}

func TestGenerateSearchCondition_StringField(t *testing.T) {
	f := newStringField("name", nil)
	node := newTestType("User")

	got := generateSearchCondition(f, node)
	assertContains(t, got, "user.NameContains(req.Query)")
}

func TestGenerateSearchCondition_NonString(t *testing.T) {
	f := newIntField("age", nil)
	node := newTestType("User")

	got := generateSearchCondition(f, node)
	if got != "" {
		t.Errorf("expected empty string for non-string field, got %q", got)
	}
}

func TestGenerateEntToDomainFieldAssignment_Regular(t *testing.T) {
	f := newStringField("name", nil)
	got := generateEntToDomainFieldAssignment(f)
	assertContains(t, got, "Name: entity.Name,")
}

func TestGenerateEntToDomainFieldAssignment_IDSkipped(t *testing.T) {
	f := newStringField("id", nil)
	got := generateEntToDomainFieldAssignment(f)
	if got != "" {
		t.Errorf("expected empty string for id field, got %q", got)
	}
}

func TestGenerateEntToDomainFieldAssignment_NillableTime(t *testing.T) {
	f := newTimeField("deleted_at", nil)
	f.Nillable = true
	got := generateEntToDomainFieldAssignment(f)
	assertContains(t, got, "entity.DeletedAt != nil")
	assertContains(t, got, "*entity.DeletedAt")
	assertContains(t, got, "time.Time{}")
}

// --- generateIdOperation tests ---

func TestGenerateIdOperation_Get(t *testing.T) {
	tests := []struct {
		name      string
		idType    field.Type
		idIdent   string
		fragments []string
	}{
		{
			"string id", field.TypeString, "string",
			[]string{"User.Get(ctx,", "id.String()"},
		},
		{
			"int64 id", field.TypeInt64, "int64",
			[]string{"User.Get(ctx,", "id.Int64()"},
		},
		{
			"default id", field.TypeInt, "int",
			[]string{"User.Get(ctx, id)"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &gen.Type{
				Name: "User",
				ID: &gen.Field{
					Name: "id",
					Type: &field.TypeInfo{Type: tt.idType, Ident: tt.idIdent},
				},
			}

			got := generateIdOperation(node, "get", "id")
			for _, frag := range tt.fragments {
				assertContains(t, got, frag)
			}
		})
	}
}

func TestGenerateIdOperation_Delete(t *testing.T) {
	tests := []struct {
		name      string
		idType    field.Type
		idIdent   string
		fragments []string
	}{
		{
			"string id", field.TypeString, "string",
			[]string{"DeleteOneID(id.String())", ".Exec(ctx)"},
		},
		{
			"int64 id", field.TypeInt64, "int64",
			[]string{"DeleteOneID(func()", "id.Int64()"},
		},
		{
			"default id", field.TypeInt, "int",
			[]string{"DeleteOneID(id).Exec(ctx)"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &gen.Type{
				Name: "User",
				ID: &gen.Field{
					Name: "id",
					Type: &field.TypeInfo{Type: tt.idType, Ident: tt.idIdent},
				},
			}

			got := generateIdOperation(node, "delete", "id")
			for _, frag := range tt.fragments {
				assertContains(t, got, frag)
			}
		})
	}
}

func TestGenerateIdOperation_Exists(t *testing.T) {
	tests := []struct {
		name      string
		idType    field.Type
		idIdent   string
		fragments []string
	}{
		{
			"string id", field.TypeString, "string",
			[]string{"IDEQ(id.String())"},
		},
		{
			"int64 id", field.TypeInt64, "int64",
			[]string{"IDEQ(func()", "id.Int64()"},
		},
		{
			"default id", field.TypeInt, "int",
			[]string{"IDEQ(id)"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &gen.Type{
				Name: "User",
				ID: &gen.Field{
					Name: "id",
					Type: &field.TypeInfo{Type: tt.idType, Ident: tt.idIdent},
				},
			}

			got := generateIdOperation(node, "exists", "id")
			for _, frag := range tt.fragments {
				assertContains(t, got, frag)
			}
		})
	}
}

func TestGenerateIdOperation_BatchDelete(t *testing.T) {
	tests := []struct {
		name      string
		idType    field.Type
		idIdent   string
		fragments []string
	}{
		{
			"string id", field.TypeString, "string",
			[]string{"stringIds", "id.String()", "IDIn(stringIds...)"},
		},
		{
			"int64 id", field.TypeInt64, "int64",
			[]string{"int64Ids", "id.Int64()", "IDIn(int64Ids...)"},
		},
		{
			"default id", field.TypeInt, "int",
			[]string{"IDIn(ids...)"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &gen.Type{
				Name: "User",
				ID: &gen.Field{
					Name: "id",
					Type: &field.TypeInfo{Type: tt.idType, Ident: tt.idIdent},
				},
			}

			got := generateIdOperation(node, "batchDelete", "id")
			for _, frag := range tt.fragments {
				assertContains(t, got, frag)
			}
		})
	}
}

func TestGenerateIdOperation_UpdateOneID(t *testing.T) {
	node := newTestType("User")
	got := generateIdOperation(node, "updateOneID", "model")
	assertContains(t, got, "User.UpdateOneID(model.ID)")
}

func TestGenerateIdOperation_Unknown(t *testing.T) {
	node := newTestType("User")
	got := generateIdOperation(node, "unknown_op", "id")
	assertContains(t, got, "Unknown operation: unknown_op")
}

// --- specificMethods tests ---

func TestSpecificMethods(t *testing.T) {
	// Only UniqueLookup and RangeLookup fields should generate methods.
	// Searchable string, enum, and bool fields should NOT generate methods.
	uniqueLookupField := newStringField("email", ptr(DomainField{
		UniqueLookup: true,
		Scopes:       AllFieldScopes,
	}))

	searchableField := newStringField("name", ptr(DomainField{Searchable: true, Scopes: AllFieldScopes}))

	enumField := newEnumField("status", ptr(DefaultField()))

	boolField := newBoolField("active", ptr(DefaultField()))

	rangeField := newTimeField("created_at", ptr(DomainField{
		RangeLookup: true,
		Scopes:      AllFieldScopes,
	}))

	node := newTestType("User",
		uniqueLookupField, searchableField, enumField, boolField, rangeField,
	)

	methods := specificMethods(node)

	// Only annotated lookup methods should be generated
	expectedParts := []string{
		"FindByEmail(ctx context.Context, email string)",
		"FindByCreatedAtRange(ctx context.Context, start, end time.Time)",
	}

	for _, part := range expectedParts {
		found := false
		for _, m := range methods {
			if strings.Contains(m, part) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected method containing %q not found in %v", part, methods)
		}
	}

	// Verify that auto-generated methods are NOT present
	unexpectedParts := []string{
		"SearchByName",
		"FindByStatus",
		"FindByActive",
	}

	for _, part := range unexpectedParts {
		for _, m := range methods {
			if strings.Contains(m, part) {
				t.Errorf("unexpected method %q found in %v (should not be auto-generated)", part, methods)
			}
		}
	}

	// Exactly 2 methods should be generated
	if len(methods) != 2 {
		t.Errorf("expected 2 methods, got %d: %v", len(methods), methods)
	}
}

func TestSpecificMethods_NoDuplicates(t *testing.T) {
	// A field with both UniqueLookup annotation should generate FindByX only once
	f := newStringField("email", ptr(DomainField{
		UniqueLookup: true,
		Scopes:       AllFieldScopes,
	}))

	node := newTestType("Order", f)

	methods := specificMethods(node)

	seen := map[string]int{}
	for _, m := range methods {
		key := strings.Split(m, "(")[0] // extract method name before params
		seen[key]++
	}

	for key, count := range seen {
		if count > 1 {
			t.Errorf("duplicate method generated: %s (count=%d)", key, count)
		}
	}
}

func TestSpecificMethods_Empty(t *testing.T) {
	node := newTestType("Empty")
	methods := specificMethods(node)
	if len(methods) != 0 {
		t.Errorf("expected 0 methods for empty type, got %d: %v", len(methods), methods)
	}
}

func TestSearchMethod(t *testing.T) {
	f := newStringField("name", nil)
	node := newTestType("User")

	got := searchMethod(f, node)
	// searchMethod uses indent "\t\t\t" and skipEmpty=true
	assertContains(t, got, `user.NameEQ(v)`)
	assertContains(t, got, `v != ""`)
}

func TestFindByMethod(t *testing.T) {
	f := newStringField("name", nil)
	node := newTestType("User")

	got := findByMethod(f, node)
	// findByMethod uses indent "\t\t" and skipEmpty=false
	assertContains(t, got, `user.NameEQ(v)`)
	assertNotContains(t, got, `v != ""`)
}
