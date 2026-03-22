// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package class

import (
	"slices"

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

// equalInt64Sets compares two int64 slices as sets (order-independent).
func equalInt64Sets(a, b []int64) bool {
	if len(a) != len(b) {
		return false
	}
	aCopy := make([]int64, len(a))
	bCopy := make([]int64, len(b))
	copy(aCopy, a)
	copy(bCopy, b)
	slices.Sort(aCopy)
	slices.Sort(bCopy)
	for i := range aCopy {
		if aCopy[i] != bCopy[i] {
			return false
		}
	}
	return true
}
