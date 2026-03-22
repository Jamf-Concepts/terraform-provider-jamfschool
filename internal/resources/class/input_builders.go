// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package class

import (
	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
	"github.com/Jamf-Concepts/terraform-provider-jamfschool/internal/common/helpers"
)

// buildCreateInput constructs a ClassCreateInput from the Terraform plan.
func buildCreateInput(plan *ClassResourceModel) jamfschool.ClassCreateInput {
	return jamfschool.ClassCreateInput{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		LocationID:  helpers.Int64PtrIfKnown(plan.LocationID),
		Students:    extractInt64Set(plan.Students),
		Teachers:    extractInt64Set(plan.Teachers),
	}
}

// buildUpdateInput constructs a ClassUpdateInput from the Terraform plan.
func buildUpdateInput(plan *ClassResourceModel) jamfschool.ClassUpdateInput {
	return jamfschool.ClassUpdateInput{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
	}
}
