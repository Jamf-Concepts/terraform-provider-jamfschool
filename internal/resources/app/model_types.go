// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package app

import "github.com/hashicorp/terraform-plugin-framework/types"

// AppDataSourceModel describes the data source data model for an app.
type AppDataSourceModel struct {
	ID       types.Int64  `tfsdk:"id"`
	BundleID types.String `tfsdk:"bundle_id"`
	AdamID   types.Int64  `tfsdk:"adam_id"`
	Name     types.String `tfsdk:"name"`
	Vendor   types.String `tfsdk:"vendor"`
	Version  types.String `tfsdk:"version"`
	Platform types.String `tfsdk:"platform"`
}

// AppResourceModel describes the resource data model for an app.
type AppResourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	AdamID      types.Int64  `tfsdk:"adam_id"`
	CountryCode types.String `tfsdk:"country_code"`
	BundleID    types.String `tfsdk:"bundle_id"`
	Name        types.String `tfsdk:"name"`
	Vendor      types.String `tfsdk:"vendor"`
	Version     types.String `tfsdk:"version"`
	Platform    types.String `tfsdk:"platform"`
	LocationID  types.Int64  `tfsdk:"location_id"`
}
