// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package device_group

import (
	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
	"github.com/Jamf-Concepts/terraform-provider-jamfschool/internal/common/helpers"
)

// buildCreateInput constructs a DeviceGroupCreateInput from the Terraform plan.
func buildCreateInput(plan *DeviceGroupResourceModel) jamfschool.DeviceGroupCreateInput {
	return jamfschool.DeviceGroupCreateInput{
		Name:           plan.Name.ValueString(),
		Description:    plan.Description.ValueString(),
		Information:    plan.Information.ValueString(),
		CollectionType: showInIOSAppToAPI(plan.ShowInIOSApp.ValueString()),
		Shared:         plan.Shared.ValueBool(),
		LocationID:     helpers.Int64PtrIfKnown(plan.LocationID),
	}
}

// showInIOSAppToAPI converts user-facing show_in_ios_app values to API collectionType values.
func showInIOSAppToAPI(v string) string {
	if v == "animated_icons" {
		return "runningTiles"
	}
	return v
}

// buildUpdateInput constructs a DeviceGroupUpdateInput from the Terraform plan.
func buildUpdateInput(plan *DeviceGroupResourceModel) jamfschool.DeviceGroupUpdateInput {
	return jamfschool.DeviceGroupUpdateInput{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		Shared:      new(plan.Shared.ValueBool()),
	}
}
