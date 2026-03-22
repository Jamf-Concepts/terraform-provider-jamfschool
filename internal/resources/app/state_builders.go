// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package app

import (
	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// modelFromApp maps a Jamf School API app to the Terraform resource model.
// Write-only fields (CountryCode) are not set here — preserved from plan.
func modelFromApp(m *AppResourceModel, a *jamfschool.App) {
	m.ID = types.Int64Value(a.ID)
	m.AdamID = types.Int64Value(a.AdamID)
	m.BundleID = types.StringValue(a.BundleID)
	m.Name = types.StringValue(a.Name)
	m.Vendor = types.StringValue(a.Vendor)
	m.Version = types.StringValue(a.Version)
	m.Platform = types.StringValue(a.Platform)
	m.LocationID = types.Int64Value(a.LocationID)
}
