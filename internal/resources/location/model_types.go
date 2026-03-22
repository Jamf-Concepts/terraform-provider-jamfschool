// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package location

import "github.com/hashicorp/terraform-plugin-framework/types"

// LocationDataSourceModel describes the data source data model for a location.
type LocationDataSourceModel struct {
	ID            types.Int64  `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	IsDistrict    types.Bool   `tfsdk:"is_district"`
	Street        types.String `tfsdk:"street"`
	StreetNumber  types.String `tfsdk:"street_number"`
	PostalCode    types.String `tfsdk:"postal_code"`
	City          types.String `tfsdk:"city"`
	Source        types.String `tfsdk:"source"`
	ASMIdentifier types.String `tfsdk:"asm_identifier"`
	SchoolNumber  types.String `tfsdk:"school_number"`
}
