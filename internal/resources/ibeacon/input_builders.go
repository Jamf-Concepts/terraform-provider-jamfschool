// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package ibeacon

import (
	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
	"github.com/Jamf-Concepts/terraform-provider-jamfschool/internal/common/helpers"
)

// buildCreateInput constructs an IBeaconCreateInput from the Terraform plan.
func buildCreateInput(plan *IBeaconResourceModel) jamfschool.IBeaconCreateInput {
	return jamfschool.IBeaconCreateInput{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		UUID:        plan.UUID.ValueString(),
		Major:       helpers.Int64PtrIfKnown(plan.Major),
		Minor:       helpers.Int64PtrIfKnown(plan.Minor),
	}
}

// buildUpdateInput constructs an IBeaconUpdateInput from the Terraform plan.
func buildUpdateInput(plan *IBeaconResourceModel) jamfschool.IBeaconUpdateInput {
	return jamfschool.IBeaconUpdateInput{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		UUID:        plan.UUID.ValueString(),
		Major:       helpers.Int64PtrIfKnown(plan.Major),
		Minor:       helpers.Int64PtrIfKnown(plan.Minor),
	}
}
