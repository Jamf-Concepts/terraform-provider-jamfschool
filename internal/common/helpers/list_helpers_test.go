// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestMatchesNamePrefix(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		prefix   types.String
		input    string
		expected bool
	}{
		{"null prefix matches all", types.StringNull(), "anything", true},
		{"unknown prefix matches all", types.StringUnknown(), "anything", true},
		{"empty prefix matches all", types.StringValue(""), "anything", true},
		{"matching prefix", types.StringValue("tf-"), "tf-test-group", true},
		{"non-matching prefix", types.StringValue("tf-"), "other-group", false},
		{"exact match", types.StringValue("exact"), "exact", true},
		{"prefix longer than name", types.StringValue("very-long-prefix"), "short", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			config := ListConfigModel{NamePrefix: tc.prefix}
			got := MatchesNamePrefix(config, tc.input)
			if got != tc.expected {
				t.Errorf("MatchesNamePrefix(%q, %q) = %v, want %v", tc.prefix, tc.input, got, tc.expected)
			}
		})
	}
}

func TestValidateNamePrefix_Empty(t *testing.T) {
	t.Parallel()

	config := ListConfigModel{NamePrefix: types.StringValue("   ")}
	var diags diag.Diagnostics
	ValidateNamePrefix(config, &diags)
	if !diags.HasError() {
		t.Error("expected error for whitespace-only name_prefix")
	}
}

func TestValidateNamePrefix_Valid(t *testing.T) {
	t.Parallel()

	config := ListConfigModel{NamePrefix: types.StringValue("tf-")}
	var diags diag.Diagnostics
	ValidateNamePrefix(config, &diags)
	if diags.HasError() {
		t.Errorf("unexpected error for valid name_prefix: %v", diags)
	}
}

func TestValidateNamePrefix_Null(t *testing.T) {
	t.Parallel()

	config := ListConfigModel{NamePrefix: types.StringNull()}
	var diags diag.Diagnostics
	ValidateNamePrefix(config, &diags)
	if diags.HasError() {
		t.Errorf("unexpected error for null name_prefix: %v", diags)
	}
}
