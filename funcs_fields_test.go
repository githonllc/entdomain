package entdomain

import (
	"testing"

	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
)

func TestDomainFields(t *testing.T) {
	df := ptr(DefaultField())
	node := newTestType("User",
		newStringField("name", df),
		newStringField("bio", nil), // no annotation
		newStringField("email", df),
	)

	got := domainFields(node)
	if len(got) != 2 {
		t.Fatalf("expected 2 domain fields, got %d", len(got))
	}
	if got[0].Name != "name" || got[1].Name != "email" {
		t.Errorf("unexpected fields: %s, %s", got[0].Name, got[1].Name)
	}
}

func TestDomainFieldsEmpty(t *testing.T) {
	node := newTestType("Empty")
	got := domainFields(node)
	if len(got) != 0 {
		t.Fatalf("expected 0 domain fields, got %d", len(got))
	}
}

func TestCreateFields(t *testing.T) {
	withCreate := ptr(DomainFieldWithScopes(ScopeCreate))
	withResponse := ptr(DomainFieldWithScopes(ScopeResponse))

	node := newTestType("User",
		newStringField("name", withCreate),
		newStringField("status", withResponse),
	)

	got := createFields(node)
	if len(got) != 1 {
		t.Fatalf("expected 1 create field, got %d", len(got))
	}
	if got[0].Name != "name" {
		t.Errorf("expected 'name', got %q", got[0].Name)
	}
}

func TestUpdateFields(t *testing.T) {
	withUpdate := ptr(DomainFieldWithScopes(ScopeUpdate))
	withCreate := ptr(DomainFieldWithScopes(ScopeCreate))

	node := newTestType("User",
		newStringField("name", withUpdate),
		newStringField("created_by", withCreate),
	)

	got := updateFields(node)
	if len(got) != 1 {
		t.Fatalf("expected 1 update field, got %d", len(got))
	}
	if got[0].Name != "name" {
		t.Errorf("expected 'name', got %q", got[0].Name)
	}
}

func TestResponseFields(t *testing.T) {
	withResp := ptr(DomainFieldWithScopes(ScopeResponse))
	withCreate := ptr(DomainFieldWithScopes(ScopeCreate))

	node := newTestType("User",
		newStringField("name", withResp),
		newStringField("password", withCreate),
	)

	got := responseFields(node)
	if len(got) != 1 {
		t.Fatalf("expected 1 response field, got %d", len(got))
	}
	if got[0].Name != "name" {
		t.Errorf("expected 'name', got %q", got[0].Name)
	}
}

func TestQueryFields(t *testing.T) {
	withQuery := ptr(DomainFieldWithScopes(ScopeQuery))
	searchable := ptr(DomainField{Searchable: true})
	plain := ptr(DomainFieldWithScopes(ScopeCreate))

	node := newTestType("User",
		newStringField("status", withQuery),
		newStringField("name", searchable),
		newStringField("bio", plain),
	)

	got := queryFields(node)
	if len(got) != 2 {
		t.Fatalf("expected 2 query fields, got %d", len(got))
	}
	if got[0].Name != "status" || got[1].Name != "name" {
		t.Errorf("unexpected fields: %s, %s", got[0].Name, got[1].Name)
	}
}

func TestSearchableFields(t *testing.T) {
	searchable := ptr(DomainField{Searchable: true, Scopes: AllFieldScopes})
	notSearchable := ptr(DomainFieldWithScopes(ScopeCreate))

	node := newTestType("User",
		newStringField("name", searchable),
		newStringField("code", notSearchable),
	)

	got := searchableFields(node)
	if len(got) != 1 {
		t.Fatalf("expected 1 searchable field, got %d", len(got))
	}
	if got[0].Name != "name" {
		t.Errorf("expected 'name', got %q", got[0].Name)
	}
}

func TestSortableFields(t *testing.T) {
	sortable := ptr(DomainField{Sortable: true, Scopes: AllFieldScopes})
	notSortable := ptr(DomainFieldWithScopes(ScopeCreate))

	node := newTestType("User",
		newStringField("name", sortable),
		newIntField("age", sortable),
		newStringField("code", notSortable),
	)

	got := sortableFields(node)
	if len(got) != 2 {
		t.Fatalf("expected 2 sortable fields, got %d", len(got))
	}
}

func TestSortableFieldsExcludesComplex(t *testing.T) {
	sortable := ptr(DomainField{Sortable: true, Scopes: AllFieldScopes})
	// Create a field with a complex type (slice)
	f := newField("tags", &field.TypeInfo{Type: field.TypeJSON, Ident: "[]string"}, sortable)

	node := newTestType("User", f)

	got := sortableFields(node)
	if len(got) != 0 {
		t.Fatalf("expected 0 sortable fields for complex type, got %d", len(got))
	}
}

func TestUniqueLookupFields(t *testing.T) {
	withLookup := ptr(DomainField{UniqueLookup: true, Scopes: AllFieldScopes})
	withoutLookup := ptr(DefaultField())

	node := newTestType("User",
		newStringField("email", withLookup),
		newStringField("name", withoutLookup),
	)

	got := uniqueLookupFields(node)
	if len(got) != 1 {
		t.Fatalf("expected 1 unique lookup field, got %d", len(got))
	}
	if got[0].Name != "email" {
		t.Errorf("expected 'email', got %q", got[0].Name)
	}
}

func TestRangeLookupFields(t *testing.T) {
	withRange := ptr(DomainField{RangeLookup: true, Scopes: AllFieldScopes})
	withoutRange := ptr(DefaultField())

	node := newTestType("User",
		newTimeField("created_at", withRange),
		newStringField("name", withoutRange),
	)

	got := rangeLookupFields(node)
	if len(got) != 1 {
		t.Fatalf("expected 1 range lookup field, got %d", len(got))
	}
	if got[0].Name != "created_at" {
		t.Errorf("expected 'created_at', got %q", got[0].Name)
	}
}

// --- responseEdges / edgeQualifiesForResponse tests ---
//
// responseEdges iterates node.Edges and delegates to edgeQualifiesForResponse.
// edge.Field() depends on unexported ent internals (Rel.fk), making it infeasible
// to construct edges that return non-nil from Field() in unit tests. The filtering
// logic is therefore tested via edgeQualifiesForResponse (pure inputs), while
// the full integration path is validated by make generate + go build.

func TestEdgeQualifiesForResponse_NilField(t *testing.T) {
	target := newTestType("Customer", newStringField("name", ptr(DefaultField())))
	if edgeQualifiesForResponse(nil, target) {
		t.Error("expected false when fkField is nil")
	}
}

func TestEdgeQualifiesForResponse_NoScopeResponse(t *testing.T) {
	// FK field exists but only has ScopeCreate, not ScopeResponse
	fkField := newUUIDField("customer_id", ptr(DomainFieldWithScopes(ScopeCreate)))
	target := newTestType("Customer", newStringField("name", ptr(DefaultField())))
	if edgeQualifiesForResponse(fkField, target) {
		t.Error("expected false when FK field lacks ScopeResponse")
	}
}

func TestEdgeQualifiesForResponse_TargetNoDomainFields(t *testing.T) {
	// FK field has ScopeResponse, but target type has no domain fields
	fkField := newUUIDField("org_id", ptr(DefaultField()))
	target := newTestType("Organization") // no annotated fields
	if edgeQualifiesForResponse(fkField, target) {
		t.Error("expected false when target type has no domain fields")
	}
}

func TestEdgeQualifiesForResponse_AllConditionsMet(t *testing.T) {
	// FK field has ScopeResponse, target type has domain fields
	fkField := newUUIDField("customer_id", ptr(DefaultField()))
	target := newTestType("Customer", newStringField("name", ptr(DefaultField())))
	if !edgeQualifiesForResponse(fkField, target) {
		t.Error("expected true when all conditions are met")
	}
}

func TestResponseEdges_NoEdges(t *testing.T) {
	node := newTestType("User", newStringField("name", ptr(DefaultField())))
	got := responseEdges(node)
	if len(got) != 0 {
		t.Fatalf("expected 0 response edges, got %d", len(got))
	}
}

func TestResponseEdges_EdgesWithoutFK(t *testing.T) {
	// Edges where Field() returns nil (no FK on this entity) should be excluded.
	// All gen.Edge constructed without ent's internal FK resolution have Field() == nil.
	node := newTestType("User", newStringField("name", ptr(DefaultField())))
	target := newTestType("Post", newStringField("title", ptr(DefaultField())))
	node.Edges = []*gen.Edge{
		{Name: "posts", Type: target, Unique: false},
		{Name: "profile", Type: target, Unique: true},
	}
	got := responseEdges(node)
	if len(got) != 0 {
		t.Fatalf("expected 0 response edges (no FK), got %d", len(got))
	}
}
