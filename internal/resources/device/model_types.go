// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package device

import "github.com/hashicorp/terraform-plugin-framework/types"

// DeviceDataSourceModel describes the data source data model for a device.
type DeviceDataSourceModel struct {
	UDID            types.String  `tfsdk:"udid"`
	SerialNumber    types.String  `tfsdk:"serial_number"`
	Name            types.String  `tfsdk:"name"`
	IsManaged       types.Bool    `tfsdk:"is_managed"`
	IsSupervised    types.Bool    `tfsdk:"is_supervised"`
	BatteryLevel    types.Float64 `tfsdk:"battery_level"`
	TotalCapacity   types.Float64 `tfsdk:"total_capacity"`
	Notes           types.String  `tfsdk:"notes"`
	LastCheckin     types.String  `tfsdk:"last_checkin"`
	LocationID      types.Int64   `tfsdk:"location_id"`
	EnrollType      types.String  `tfsdk:"enroll_type"`
	ModelName       types.String  `tfsdk:"model_name"`
	ModelIdentifier types.String  `tfsdk:"model_identifier"`
	OSVersion       types.String  `tfsdk:"os_version"`
	OSPrefix        types.String  `tfsdk:"os_prefix"`
	InTrash         types.Bool    `tfsdk:"in_trash"`
}
