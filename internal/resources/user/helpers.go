// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package user

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// extractInt64Set extracts int64 values from a types.Set.
func extractInt64Set(s types.Set) []int64 {
	if s.IsNull() || s.IsUnknown() {
		return nil
	}
	elements := s.Elements()
	ids := make([]int64, 0, len(elements))
	for _, e := range elements {
		if v, ok := e.(types.Int64); ok {
			ids = append(ids, v.ValueInt64())
		}
	}
	return ids
}

// int64SliceToSet converts an int64 slice to a types.Set.
func int64SliceToSet(vals []int64) types.Set {
	elements := make([]attr.Value, len(vals))
	for i, v := range vals {
		elements[i] = types.Int64Value(v)
	}
	s, _ := types.SetValue(types.Int64Type, elements)
	return s
}

// int64sToAny converts an int64 slice to an any slice for the API.
func int64sToAny(ids []int64) []any {
	result := make([]any, len(ids))
	for i, id := range ids {
		result[i] = id
	}
	return result
}

// stringSliceToList converts a string slice to a types.List.
func stringSliceToList(vals []string) types.List {
	elements := make([]attr.Value, len(vals))
	for i, v := range vals {
		elements[i] = types.StringValue(v)
	}
	l, _ := types.ListValue(types.StringType, elements)
	return l
}
