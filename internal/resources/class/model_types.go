// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package class

import "github.com/hashicorp/terraform-plugin-framework/types"

type ClassResourceModel struct {
	UUID          types.String `tfsdk:"uuid"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	LocationID    types.Int64  `tfsdk:"location_id"`
	Source        types.String `tfsdk:"source"`
	StudentCount  types.Int64  `tfsdk:"student_count"`
	TeacherCount  types.Int64  `tfsdk:"teacher_count"`
	DeviceGroupID types.Int64  `tfsdk:"device_group_id"`
	Students      types.Set    `tfsdk:"students"`
	Teachers      types.Set    `tfsdk:"teachers"`
}

type ClassDataSourceModel struct {
	UUID          types.String `tfsdk:"uuid"`
	Name          types.String `tfsdk:"name"`
	Description   types.String `tfsdk:"description"`
	LocationID    types.Int64  `tfsdk:"location_id"`
	Source        types.String `tfsdk:"source"`
	StudentCount  types.Int64  `tfsdk:"student_count"`
	TeacherCount  types.Int64  `tfsdk:"teacher_count"`
	DeviceGroupID types.Int64  `tfsdk:"device_group_id"`
	DeviceCount   types.Int64  `tfsdk:"device_count"`
	DeviceUDIDs   types.List   `tfsdk:"device_udids"`
	Students      types.Set    `tfsdk:"students"`
	Teachers      types.Set    `tfsdk:"teachers"`
}
