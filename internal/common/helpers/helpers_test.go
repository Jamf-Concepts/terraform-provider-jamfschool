// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestStringValueOrNull(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    string
		wantNull bool
		wantVal  string
	}{
		{"empty string returns null", "", true, ""},
		{"non-empty string returns value", "hello", false, "hello"},
		{"whitespace is not empty", " ", false, " "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := StringValueOrNull(tt.input)
			if tt.wantNull {
				if !result.IsNull() {
					t.Errorf("expected null, got %q", result.ValueString())
				}
			} else {
				if result.IsNull() {
					t.Error("expected non-null value, got null")
				}
				if result.ValueString() != tt.wantVal {
					t.Errorf("expected %q, got %q", tt.wantVal, result.ValueString())
				}
			}
		})
	}
}

func TestInt64ValueOrNull(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		input    int64
		wantNull bool
		wantVal  int64
	}{
		{"zero returns null", 0, true, 0},
		{"positive value returns value", 42, false, 42},
		{"negative value returns value", -1, false, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := Int64ValueOrNull(tt.input)
			if tt.wantNull {
				if !result.IsNull() {
					t.Errorf("expected null, got %d", result.ValueInt64())
				}
			} else {
				if result.IsNull() {
					t.Error("expected non-null value, got null")
				}
				if result.ValueInt64() != tt.wantVal {
					t.Errorf("expected %d, got %d", tt.wantVal, result.ValueInt64())
				}
			}
		})
	}
}

func TestStringPtrValueOrNull(t *testing.T) {
	t.Parallel()

	str := "hello"
	empty := ""

	tests := []struct {
		name     string
		input    *string
		wantNull bool
		wantVal  string
	}{
		{"nil pointer returns null", nil, true, ""},
		{"non-nil pointer returns value", &str, false, "hello"},
		{"empty string pointer returns value", &empty, false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := StringPtrValueOrNull(tt.input)
			if tt.wantNull {
				if !result.IsNull() {
					t.Errorf("expected null, got %q", result.ValueString())
				}
			} else {
				if result.IsNull() {
					t.Error("expected non-null value, got null")
				}
				if result.ValueString() != tt.wantVal {
					t.Errorf("expected %q, got %q", tt.wantVal, result.ValueString())
				}
			}
		})
	}
}

func TestPtr(t *testing.T) {
	t.Parallel()

	t.Run("bool true", func(t *testing.T) {
		t.Parallel()
		result := Ptr(true)
		if result == nil || *result != true {
			t.Errorf("expected pointer to true, got %v", result)
		}
	})

	t.Run("bool false", func(t *testing.T) {
		t.Parallel()
		result := Ptr(false)
		if result == nil || *result != false {
			t.Errorf("expected pointer to false, got %v", result)
		}
	})

	t.Run("string", func(t *testing.T) {
		t.Parallel()
		result := Ptr("hello")
		if result == nil || *result != "hello" {
			t.Errorf("expected pointer to hello, got %v", result)
		}
	})

	t.Run("int64", func(t *testing.T) {
		t.Parallel()
		result := Ptr(int64(42))
		if result == nil || *result != 42 {
			t.Errorf("expected pointer to 42, got %v", result)
		}
	})

	t.Run("returns distinct pointers", func(t *testing.T) {
		t.Parallel()
		a := Ptr(true)
		b := Ptr(true)
		if a == b {
			t.Error("expected distinct pointers, got same pointer")
		}
	})
}

func TestInt64PtrIfKnown(t *testing.T) {
	t.Parallel()

	t.Run("null returns nil", func(t *testing.T) {
		t.Parallel()
		result := Int64PtrIfKnown(types.Int64Null())
		if result != nil {
			t.Errorf("expected nil, got %v", *result)
		}
	})

	t.Run("unknown returns nil", func(t *testing.T) {
		t.Parallel()
		result := Int64PtrIfKnown(types.Int64Unknown())
		if result != nil {
			t.Errorf("expected nil, got %v", *result)
		}
	})

	t.Run("known value returns pointer", func(t *testing.T) {
		t.Parallel()
		result := Int64PtrIfKnown(types.Int64Value(42))
		if result == nil {
			t.Fatal("expected non-nil pointer, got nil")
		}
		if *result != 42 {
			t.Errorf("expected 42, got %d", *result)
		}
	})

	t.Run("zero value returns pointer", func(t *testing.T) {
		t.Parallel()
		result := Int64PtrIfKnown(types.Int64Value(0))
		if result == nil {
			t.Fatal("expected non-nil pointer, got nil")
		}
		if *result != 0 {
			t.Errorf("expected 0, got %d", *result)
		}
	})
}

func TestStringValueOrNull_TypeAssertion(t *testing.T) {
	t.Parallel()

	_ = StringValueOrNull("test")

	nullResult := StringValueOrNull("")
	if !nullResult.IsNull() {
		t.Error("expected null")
	}
	if nullResult.IsUnknown() {
		t.Error("null should not be unknown")
	}
}

func TestInt64ValueOrNull_TypeAssertion(t *testing.T) {
	t.Parallel()

	_ = Int64ValueOrNull(5)

	nullResult := Int64ValueOrNull(0)
	if !nullResult.IsNull() {
		t.Error("expected null")
	}
	if nullResult.IsUnknown() {
		t.Error("null should not be unknown")
	}
}
