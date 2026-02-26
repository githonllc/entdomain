package entdomain

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// Cursor holds the keyset pagination position. It encodes the sort field
// value and entity ID so the next query can seek directly to the right
// position via a WHERE clause instead of counting offset rows.
type Cursor struct {
	// ID is the entity's primary key, always included for tie-breaking
	// when multiple rows share the same sort field value.
	ID any `json:"id"`

	// Value is the sort field value of the last row. Nil when sorting
	// by ID only (no secondary sort field).
	Value any `json:"value,omitempty"`
}

// PageInfo holds cursor-based pagination metadata returned alongside
// query results.
type PageInfo struct {
	// HasNextPage indicates whether more results exist beyond this page.
	HasNextPage bool `json:"hasNextPage"`

	// EndCursor is the opaque cursor string pointing to the last item
	// in the current page. Pass this as ListRequest.Cursor to fetch
	// the next page.
	EndCursor string `json:"endCursor,omitempty"`
}

// EncodeCursor serializes a Cursor to a URL-safe opaque string.
// The encoding is base64(json(cursor)).
func EncodeCursor(c *Cursor) string {
	if c == nil {
		return ""
	}
	b, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return base64.RawURLEncoding.EncodeToString(b)
}

// DecodeCursor deserializes an opaque cursor string back to a Cursor.
// JSON unmarshals numbers as float64, so this function normalizes
// float64 values that represent whole numbers back to int64.
func DecodeCursor(s string) (*Cursor, error) {
	if s == "" {
		return nil, fmt.Errorf("cursor cannot be empty")
	}
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor encoding: %w", err)
	}
	var c Cursor
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("invalid cursor data: %w", err)
	}
	if c.ID == nil {
		return nil, fmt.Errorf("cursor missing required ID field")
	}
	// Normalize float64 → int64 for JSON-unmarshaled numbers
	c.ID = normalizeJSONNumber(c.ID)
	c.Value = normalizeJSONNumber(c.Value)
	return &c, nil
}

// normalizeJSONNumber converts float64 values that represent whole
// numbers back to int64 to match the original type before JSON encoding.
func normalizeJSONNumber(v any) any {
	if f, ok := v.(float64); ok && f == float64(int64(f)) {
		return int64(f)
	}
	return v
}
