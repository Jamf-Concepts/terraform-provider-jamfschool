// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package user

import "github.com/hashicorp/terraform-plugin-framework/types"

// UserResourceModel describes the resource data model.
type UserResourceModel struct {
	ID                          types.Int64  `tfsdk:"id"`
	Username                    types.String `tfsdk:"username"`
	Password                    types.String `tfsdk:"password"`
	Email                       types.String `tfsdk:"email"`
	FirstName                   types.String `tfsdk:"first_name"`
	LastName                    types.String `tfsdk:"last_name"`
	Notes                       types.String `tfsdk:"notes"`
	Exclude                     types.Bool   `tfsdk:"exclude"`
	StoreMailContactsCalendars  types.Bool   `tfsdk:"store_mail_contacts_calendars"`
	MailContactsCalendarsDomain types.String `tfsdk:"mail_contacts_calendars_domain"`
	LocationID                  types.Int64  `tfsdk:"location_id"`
	MoveDevicesOnLocationChange types.Bool   `tfsdk:"move_devices_on_location_change"`
	MemberOf                    types.Set    `tfsdk:"member_of"`
	Status                      types.String `tfsdk:"status"`
	DeviceCount                 types.Int64  `tfsdk:"device_count"`
}

// UserDataSourceModel describes the data source data model.
type UserDataSourceModel struct {
	ID          types.Int64  `tfsdk:"id"`
	Username    types.String `tfsdk:"username"`
	Email       types.String `tfsdk:"email"`
	FirstName   types.String `tfsdk:"first_name"`
	LastName    types.String `tfsdk:"last_name"`
	Notes       types.String `tfsdk:"notes"`
	Exclude     types.Bool   `tfsdk:"exclude"`
	LocationID  types.Int64  `tfsdk:"location_id"`
	MemberOf    types.Set    `tfsdk:"member_of"`
	GroupNames  types.List   `tfsdk:"group_names"`
	Status      types.String `tfsdk:"status"`
	DeviceCount types.Int64  `tfsdk:"device_count"`
}
