// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package ibeacon

import "github.com/hashicorp/terraform-plugin-framework/types"

type IBeaconResourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	UUID        types.String `tfsdk:"uuid"`
	Major       types.Int64  `tfsdk:"major"`
	Minor       types.Int64  `tfsdk:"minor"`
}

type IBeaconDataSourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	UUID        types.String `tfsdk:"uuid"`
	Major       types.Int64  `tfsdk:"major"`
	Minor       types.Int64  `tfsdk:"minor"`
}
