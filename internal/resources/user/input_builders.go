// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package user

import (
	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
	"github.com/Jamf-Concepts/terraform-provider-jamfschool/internal/common/helpers"
)

// buildCreateInput constructs a UserCreateInput from the Terraform plan.
func buildCreateInput(plan *UserResourceModel) jamfschool.UserCreateInput {
	input := jamfschool.UserCreateInput{
		Username:      plan.Username.ValueString(),
		Password:      plan.Password.ValueString(),
		Email:         plan.Email.ValueString(),
		FirstName:     plan.FirstName.ValueString(),
		LastName:      plan.LastName.ValueString(),
		Domain:        plan.MailContactsCalendarsDomain.ValueString(),
		Notes:         plan.Notes.ValueString(),
		Exclude:       plan.Exclude.ValueBool(),
		StorePassword: plan.StoreMailContactsCalendars.ValueBool(),
		LocationID:    helpers.Int64PtrIfKnown(plan.LocationID),
	}
	if groupIDs := extractInt64Set(plan.MemberOf); groupIDs != nil {
		input.MemberOf = int64sToAny(groupIDs)
	}
	return input
}

// buildUpdateInput constructs a UserUpdateInput from the Terraform plan and current state.
func buildUpdateInput(plan, state *UserResourceModel) jamfschool.UserUpdateInput {
	input := jamfschool.UserUpdateInput{
		Username:      plan.Username.ValueString(),
		Email:         plan.Email.ValueString(),
		FirstName:     plan.FirstName.ValueString(),
		LastName:      plan.LastName.ValueString(),
		Domain:        plan.MailContactsCalendarsDomain.ValueString(),
		Notes:         plan.Notes.ValueString(),
		Exclude:       helpers.Ptr(plan.Exclude.ValueBool()),
		StorePassword: helpers.Ptr(plan.StoreMailContactsCalendars.ValueBool()),
	}
	if !plan.Password.Equal(state.Password) {
		input.Password = plan.Password.ValueString()
	}
	memberOfAny := int64sToAny(extractInt64Set(plan.MemberOf))
	input.MemberOf = &memberOfAny
	return input
}
