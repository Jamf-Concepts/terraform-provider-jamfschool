// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import "github.com/hashicorp/terraform-plugin-framework/types"

// Ptr returns a pointer to any value.
//
//go:fix inline
func Ptr[T any](v T) *T {
	return new(v)
}

// Int64PtrIfKnown returns a pointer to the int64 value if set, nil otherwise.
func Int64PtrIfKnown(v types.Int64) *int64 {
	if v.IsNull() || v.IsUnknown() {
		return nil
	}
	val := v.ValueInt64()
	return &val
}

// StringValueOrNull returns a types.String from a Go string, null if empty.
func StringValueOrNull(s string) types.String {
	if s == "" {
		return types.StringNull()
	}
	return types.StringValue(s)
}

// Int64ValueOrNull returns a types.Int64 from a Go int64, null if zero.
func Int64ValueOrNull(v int64) types.Int64 {
	if v == 0 {
		return types.Int64Null()
	}
	return types.Int64Value(v)
}

// StringPtrValueOrNull returns a types.String from a *string, null if nil.
func StringPtrValueOrNull(s *string) types.String {
	if s == nil {
		return types.StringNull()
	}
	return types.StringValue(*s)
}
