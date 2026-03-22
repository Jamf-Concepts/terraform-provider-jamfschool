// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package profile

import "github.com/hashicorp/terraform-plugin-framework/types"

// ProfileDataSourceModel describes the data source data model for a configuration profile.
type ProfileDataSourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	LocationID  types.Int64  `tfsdk:"location_id"`
	Identifier  types.String `tfsdk:"identifier"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Platform    types.String `tfsdk:"platform"`
}
