// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package user_group

import "github.com/hashicorp/terraform-plugin-framework/types"

// UserGroupResourceModel describes the resource data model.
type UserGroupResourceModel struct {
	ID                types.Int64  `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	LocationID        types.Int64  `tfsdk:"location_id"`
	UserCount         types.Int64  `tfsdk:"user_count"`
	JamfSchoolTeacher types.String `tfsdk:"jamf_school_teacher"`
	JamfParent        types.String `tfsdk:"jamf_parent"`
}

// UserGroupDataSourceModel describes the data source data model.
type UserGroupDataSourceModel struct {
	ID                types.Int64  `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Description       types.String `tfsdk:"description"`
	LocationID        types.Int64  `tfsdk:"location_id"`
	UserCount         types.Int64  `tfsdk:"user_count"`
	JamfSchoolTeacher types.String `tfsdk:"jamf_school_teacher"`
	JamfParent        types.String `tfsdk:"jamf_parent"`
}
