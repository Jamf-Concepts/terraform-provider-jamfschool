// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package ibeacon

import (
	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// modelFromIBeacon maps a Jamf School API iBeacon to the Terraform resource model.
func modelFromIBeacon(m *IBeaconResourceModel, b *jamfschool.IBeacon) {
	m.ID = types.Int64Value(b.ID)
	m.Name = types.StringValue(b.Name)
	m.Description = types.StringValue(b.Description)
	m.UUID = types.StringValue(b.UUID)
	m.Major = types.Int64Value(b.Major)
	m.Minor = types.Int64Value(b.Minor)
}
