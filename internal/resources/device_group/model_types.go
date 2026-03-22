// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package device_group

import "github.com/hashicorp/terraform-plugin-framework/types"

type DeviceGroupResourceModel struct {
	ID           types.Int64  `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Information  types.String `tfsdk:"information"`
	ShowInIOSApp types.String `tfsdk:"show_in_ios_app"`
	LocationID   types.Int64  `tfsdk:"location_id"`
	Shared       types.Bool   `tfsdk:"shared"`
	Members      types.Int64  `tfsdk:"members"`
	IsSmartGroup types.Bool   `tfsdk:"is_smart_group"`
	MemberUDIDs  types.Set    `tfsdk:"member_udids"`
}

type DeviceGroupDataSourceModel struct {
	ID           types.Int64  `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	Information  types.String `tfsdk:"information"`
	LocationID   types.Int64  `tfsdk:"location_id"`
	Shared       types.Bool   `tfsdk:"shared"`
	Members      types.Int64  `tfsdk:"members"`
	IsSmartGroup types.Bool   `tfsdk:"is_smart_group"`
}
