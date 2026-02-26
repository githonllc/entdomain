package entdomain

import (
	"testing"
)

func TestNewIDFromInt64(t *testing.T) {
	tests := []struct {
		name     string
		value    int64
		expected string
		isZero   bool
	}{
		{
			name:     "positive value",
			value:    123,
			expected: "123",
			isZero:   false,
		},
		{
			name:     "zero value",
			value:    0,
			expected: "0",
			isZero:   true,
		},
		{
			name:     "negative value",
			value:    -456,
			expected: "-456",
			isZero:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := NewIDFromInt64(tt.value)

			if id.String() != tt.expected {
				t.Errorf("String() = %v, want %v", id.String(), tt.expected)
			}

			if id.IsZero() != tt.isZero {
				t.Errorf("IsZero() = %v, want %v", id.IsZero(), tt.isZero)
			}

			val, err := id.Int64()
			if err != nil {
				t.Errorf("Int64() error = %v", err)
			}
			if val != tt.value {
				t.Errorf("Int64() = %v, want %v", val, tt.value)
			}
		})
	}
}

func TestNewIDFromString(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
		isZero   bool
	}{
		{
			name:     "non-empty string",
			value:    "abc-123",
			expected: "abc-123",
			isZero:   false,
		},
		{
			name:     "empty string",
			value:    "",
			expected: "",
			isZero:   true,
		},
		{
			name:     "uuid string",
			value:    "550e8400-e29b-41d4-a716-446655440000",
			expected: "550e8400-e29b-41d4-a716-446655440000",
			isZero:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id := NewIDFromString(tt.value)

			if id.String() != tt.expected {
				t.Errorf("String() = %v, want %v", id.String(), tt.expected)
			}

			if id.IsZero() != tt.isZero {
				t.Errorf("IsZero() = %v, want %v", id.IsZero(), tt.isZero)
			}
		})
	}
}

func TestStringIDInt64Conversion(t *testing.T) {
	// String ID should return error when converting to int64
	id := NewIDFromString("abc-123")
	_, err := id.Int64()
	if err == nil {
		t.Error("Expected error when converting string ID to int64")
	}
}

func TestInt64IDStringConversion(t *testing.T) {
	// Int64 ID should work fine with string conversion
	id := NewIDFromInt64(123)
	str := id.String()
	if str != "123" {
		t.Errorf("String() = %v, want %v", str, "123")
	}
}

func TestIDEquality(t *testing.T) {
	id1 := NewIDFromInt64(123)
	id2 := NewIDFromInt64(123)
	id3 := NewIDFromInt64(456)

	if id1.String() != id2.String() {
		t.Error("Same int64 IDs should have same string representation")
	}

	if id1.String() == id3.String() {
		t.Error("Different int64 IDs should have different string representation")
	}

	str1 := NewIDFromString("test")
	str2 := NewIDFromString("test")
	str3 := NewIDFromString("other")

	if str1.String() != str2.String() {
		t.Error("Same string IDs should have same string representation")
	}

	if str1.String() == str3.String() {
		t.Error("Different string IDs should have different string representation")
	}
}
