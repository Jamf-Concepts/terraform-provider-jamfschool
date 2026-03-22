// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package user

import (
	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// modelFromUser maps a Jamf School API user to the Terraform resource model.
// Write-only fields (Password, StoreMailContactsCalendars, MailContactsCalendarsDomain)
// are not set here — they are preserved from the plan/state.
func modelFromUser(m *UserResourceModel, u *jamfschool.User) {
	m.ID = types.Int64Value(u.ID)
	m.Username = types.StringValue(u.Username)
	m.Email = types.StringValue(u.Email)
	m.FirstName = types.StringValue(u.FirstName)
	m.LastName = types.StringValue(u.LastName)
	m.Notes = types.StringValue(u.Notes)
	m.Exclude = types.BoolValue(u.Exclude)
	m.LocationID = types.Int64Value(u.LocationID)
	m.MemberOf = int64SliceToSet(u.GroupIDs)
	m.Status = types.StringValue(u.Status)
	m.DeviceCount = types.Int64Value(u.DeviceCount)
}
