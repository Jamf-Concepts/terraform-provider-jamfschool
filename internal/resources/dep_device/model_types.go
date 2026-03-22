// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package dep_device

import "github.com/hashicorp/terraform-plugin-framework/types"

// DEPDeviceDataSourceModel describes the data source data model for a DEP device.
type DEPDeviceDataSourceModel struct {
	SerialNumber    types.String `tfsdk:"serial_number"`
	ID              types.Int64  `tfsdk:"id"`
	UserID          types.Int64  `tfsdk:"user_id"`
	LocationID      types.Int64  `tfsdk:"location_id"`
	Model           types.String `tfsdk:"model"`
	Color           types.String `tfsdk:"color"`
	Status          types.String `tfsdk:"status"`
	DeviceName      types.String `tfsdk:"device_name"`
	ProfileName     types.String `tfsdk:"profile_name"`
	PlaceholderName types.String `tfsdk:"placeholder_name"`
}
