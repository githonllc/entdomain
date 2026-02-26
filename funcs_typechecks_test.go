package entdomain

import (
	"testing"
)

func TestIsUniqueField(t *testing.T) {
	unique := newStringField("email", nil)
	unique.Unique = true

	notUnique := newStringField("name", nil)

	if !isUniqueField(unique) {
		t.Error("expected unique field to return true")
	}
	if isUniqueField(notUnique) {
		t.Error("expected non-unique field to return false")
	}
}

func TestIsTimeField(t *testing.T) {
	timeField := newTimeField("created_at", nil)
	stringField := newStringField("name", nil)
	intField := newIntField("age", nil)

	if !isTimeField(timeField) {
		t.Error("expected time field to return true")
	}
	if isTimeField(stringField) {
		t.Error("expected string field to return false")
	}
	if isTimeField(intField) {
		t.Error("expected int field to return false")
	}
}

func TestHasTimeFields(t *testing.T) {
	df := ptr(DefaultField())

	// Type with time fields
	withTime := newTestType("User",
		newStringField("name", df),
		newTimeField("created_at", df),
	)
	if !hasTimeFields(withTime) {
		t.Error("expected hasTimeFields to return true")
	}

	// Type without time fields
	withoutTime := newTestType("User",
		newStringField("name", df),
		newIntField("age", df),
	)
	if hasTimeFields(withoutTime) {
		t.Error("expected hasTimeFields to return false")
	}

	// Empty type
	empty := newTestType("Empty")
	if hasTimeFields(empty) {
		t.Error("expected hasTimeFields to return false for empty type")
	}
}

func TestHasTimeField(t *testing.T) {
	df := ptr(DefaultField())
	node := newTestType("User",
		newTimeField("created_at", df),
		newTimeField("updated_at", df),
		newStringField("name", df),
	)

	if !hasTimeField(node, "created_at") {
		t.Error("expected hasTimeField('created_at') to return true")
	}
	if !hasTimeField(node, "updated_at") {
		t.Error("expected hasTimeField('updated_at') to return true")
	}
	if hasTimeField(node, "name") {
		t.Error("expected hasTimeField('name') to return false (not a time field)")
	}
	if hasTimeField(node, "deleted_at") {
		t.Error("expected hasTimeField('deleted_at') to return false (not present)")
	}
}

func TestIsUniqueLookupField(t *testing.T) {
	withLookup := newStringField("email", ptr(DomainField{UniqueLookup: true}))
	withoutLookup := newStringField("name", ptr(DefaultField()))
	noAnnotation := newStringField("bio", nil)

	if !isUniqueLookupField(withLookup) {
		t.Error("expected unique lookup field to return true")
	}
	if isUniqueLookupField(withoutLookup) {
		t.Error("expected non-unique-lookup field to return false")
	}
	if isUniqueLookupField(noAnnotation) {
		t.Error("expected unannotated field to return false")
	}
}
