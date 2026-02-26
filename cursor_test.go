package entdomain

import (
	"testing"
)

func TestEncodeDecode_IDOnly(t *testing.T) {
	original := &Cursor{ID: int64(12345)}
	encoded := EncodeCursor(original)
	if encoded == "" {
		t.Fatal("EncodeCursor returned empty string")
	}

	decoded, err := DecodeCursor(encoded)
	if err != nil {
		t.Fatalf("DecodeCursor failed: %v", err)
	}

	// After normalization, int64 should be preserved (not float64)
	if decoded.ID != int64(12345) {
		t.Errorf("ID = %v (%T), want int64(12345)", decoded.ID, decoded.ID)
	}
	if decoded.Value != nil {
		t.Errorf("Value = %v, want nil", decoded.Value)
	}
}

func TestEncodeDecode_WithStringValue(t *testing.T) {
	original := &Cursor{ID: int64(42), Value: "Alice"}
	encoded := EncodeCursor(original)

	decoded, err := DecodeCursor(encoded)
	if err != nil {
		t.Fatalf("DecodeCursor failed: %v", err)
	}

	if decoded.ID != int64(42) {
		t.Errorf("ID = %v (%T), want int64(42)", decoded.ID, decoded.ID)
	}
	if decoded.Value != "Alice" {
		t.Errorf("Value = %v, want Alice", decoded.Value)
	}
}

func TestNormalizeJSONNumber(t *testing.T) {
	// float64 representing whole number → int64
	if v := normalizeJSONNumber(float64(42)); v != int64(42) {
		t.Errorf("normalizeJSONNumber(42.0) = %v (%T), want int64(42)", v, v)
	}
	// float64 with decimal → stays float64
	if v := normalizeJSONNumber(float64(3.14)); v != float64(3.14) {
		t.Errorf("normalizeJSONNumber(3.14) = %v (%T), want float64(3.14)", v, v)
	}
	// string → unchanged
	if v := normalizeJSONNumber("hello"); v != "hello" {
		t.Errorf("normalizeJSONNumber(\"hello\") = %v, want hello", v)
	}
	// nil → unchanged
	if v := normalizeJSONNumber(nil); v != nil {
		t.Errorf("normalizeJSONNumber(nil) = %v, want nil", v)
	}
}

func TestEncodeDecode_WithStringID(t *testing.T) {
	original := &Cursor{ID: "uuid-abc-123", Value: "test@example.com"}
	encoded := EncodeCursor(original)

	decoded, err := DecodeCursor(encoded)
	if err != nil {
		t.Fatalf("DecodeCursor failed: %v", err)
	}

	if decoded.ID != "uuid-abc-123" {
		t.Errorf("ID = %v, want uuid-abc-123", decoded.ID)
	}
}

func TestEncodeCursor_Nil(t *testing.T) {
	if got := EncodeCursor(nil); got != "" {
		t.Errorf("EncodeCursor(nil) = %q, want empty", got)
	}
}

func TestDecodeCursor_Empty(t *testing.T) {
	_, err := DecodeCursor("")
	if err == nil {
		t.Error("DecodeCursor(\"\") should return error")
	}
}

func TestDecodeCursor_InvalidBase64(t *testing.T) {
	_, err := DecodeCursor("!!!not-base64!!!")
	if err == nil {
		t.Error("DecodeCursor with invalid base64 should return error")
	}
}

func TestDecodeCursor_InvalidJSON(t *testing.T) {
	// Valid base64 but invalid JSON
	_, err := DecodeCursor("bm90LWpzb24")
	if err == nil {
		t.Error("DecodeCursor with invalid JSON should return error")
	}
}

func TestDecodeCursor_MissingID(t *testing.T) {
	// JSON with no id field: {"value": "test"}
	_, err := DecodeCursor("eyJ2YWx1ZSI6InRlc3QifQ")
	if err == nil {
		t.Error("DecodeCursor with missing ID should return error")
	}
}

func TestPageInfo_Defaults(t *testing.T) {
	p := PageInfo{}
	if p.HasNextPage {
		t.Error("default HasNextPage should be false")
	}
	if p.EndCursor != "" {
		t.Error("default EndCursor should be empty")
	}
}
