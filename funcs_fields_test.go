package entdomain

import (
	"testing"

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

func TestUpdateableFields(t *testing.T) {
	df := ptr(DefaultField())
	idField := newStringField("id", df)
	immutableField := newStringField("creator", df)
	immutableField.Immutable = true
	normalField := newStringField("name", df)
	unannotatedField := newStringField("internal", nil)

	node := newTestType("User", idField, immutableField, normalField, unannotatedField)

	got := updateableFields(node)
	if len(got) != 1 {
		t.Fatalf("expected 1 updateable field, got %d", len(got))
	}
	if got[0].Name != "name" {
		t.Errorf("expected 'name', got %q", got[0].Name)
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
