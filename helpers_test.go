package entdomain

import (
	"testing"
	"time"
)

func TestPtr(t *testing.T) {
	// int
	i := 42
	p := Ptr(i)
	if *p != 42 {
		t.Errorf("Ptr(42) = %d, want 42", *p)
	}

	// string
	s := "hello"
	ps := Ptr(s)
	if *ps != "hello" {
		t.Errorf("Ptr(\"hello\") = %q, want \"hello\"", *ps)
	}

	// bool
	b := true
	pb := Ptr(b)
	if *pb != true {
		t.Errorf("Ptr(true) = %v, want true", *pb)
	}
}

func TestPtrOrNil(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		isNil  bool
		expect string
	}{
		{"non-empty string", "hello", false, "hello"},
		{"empty string", "", true, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PtrOrNil(tt.input)
			if tt.isNil {
				if got != nil {
					t.Errorf("PtrOrNil(%q) = %v, want nil", tt.input, *got)
				}
			} else {
				if got == nil {
					t.Errorf("PtrOrNil(%q) = nil, want %q", tt.input, tt.expect)
				} else if *got != tt.expect {
					t.Errorf("PtrOrNil(%q) = %q, want %q", tt.input, *got, tt.expect)
				}
			}
		})
	}
}

func TestPtrTimeOrNil(t *testing.T) {
	now := time.Now()
	zero := time.Time{}

	t.Run("non-zero time", func(t *testing.T) {
		got := PtrTimeOrNil(now)
		if got == nil {
			t.Fatal("PtrTimeOrNil(now) = nil, want non-nil")
		}
		if !got.Equal(now) {
			t.Errorf("PtrTimeOrNil(now) = %v, want %v", *got, now)
		}
	})

	t.Run("zero time", func(t *testing.T) {
		got := PtrTimeOrNil(zero)
		if got != nil {
			t.Errorf("PtrTimeOrNil(zero) = %v, want nil", *got)
		}
	})
}
