// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package user_group

import (
	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
	"github.com/Jamf-Concepts/terraform-provider-jamfschool/internal/common/helpers"
)

// aclFromModel builds an ACL from the Terraform plan model.
func aclFromModel(m *UserGroupResourceModel) *jamfschool.ACL {
	if m.JamfSchoolTeacher.IsNull() && m.JamfParent.IsNull() {
		return nil
	}
	return &jamfschool.ACL{
		Teacher: m.JamfSchoolTeacher.ValueString(),
		Parent:  m.JamfParent.ValueString(),
	}
}

// buildCreateInput constructs a GroupCreateInput from the Terraform plan.
func buildCreateInput(plan *UserGroupResourceModel) jamfschool.GroupCreateInput {
	return jamfschool.GroupCreateInput{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		ACL:         aclFromModel(plan),
		LocationID:  helpers.Int64PtrIfKnown(plan.LocationID),
	}
}

// buildUpdateInput constructs a GroupUpdateInput from the Terraform plan.
func buildUpdateInput(plan *UserGroupResourceModel) jamfschool.GroupUpdateInput {
	return jamfschool.GroupUpdateInput{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		ACL:         aclFromModel(plan),
	}
}
