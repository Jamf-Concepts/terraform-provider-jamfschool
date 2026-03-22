// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package user_group

import (
	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// modelFromGroup maps a Jamf School API group to the Terraform resource model.
func modelFromGroup(m *UserGroupResourceModel, g *jamfschool.Group) {
	m.ID = types.Int64Value(g.ID)
	m.Name = types.StringValue(g.Name)
	m.Description = types.StringValue(g.Description)
	m.LocationID = types.Int64Value(g.LocationID)
	m.UserCount = types.Int64Value(g.UserCount)
	m.JamfSchoolTeacher = types.StringValue(g.ACL.Teacher)
	m.JamfParent = types.StringValue(g.ACL.Parent)
}
