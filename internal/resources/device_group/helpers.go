// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package device_group

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// extractStringSet extracts string values from a types.Set.
func extractStringSet(s types.Set) []string {
	if s.IsNull() || s.IsUnknown() {
		return nil
	}
	elements := s.Elements()
	vals := make([]string, 0, len(elements))
	for _, e := range elements {
		if v, ok := e.(types.String); ok {
			vals = append(vals, v.ValueString())
		}
	}
	return vals
}

// stringSliceToSet converts a string slice to a types.Set.
func stringSliceToSet(vals []string) types.Set {
	elements := make([]attr.Value, len(vals))
	for i, v := range vals {
		elements[i] = types.StringValue(v)
	}
	s, _ := types.SetValue(types.StringType, elements)
	return s
}

// stringSetDiff returns elements in a that are not in b.
func stringSetDiff(a, b []string) []string {
	bSet := make(map[string]struct{}, len(b))
	for _, v := range b {
		bSet[v] = struct{}{}
	}
	var diff []string
	for _, v := range a {
		if _, ok := bSet[v]; !ok {
			diff = append(diff, v)
		}
	}
	return diff
}
