// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package device_group

import (
	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// modelFromDeviceGroup maps a Jamf School API device group to the Terraform resource model.
func modelFromDeviceGroup(m *DeviceGroupResourceModel, dg *jamfschool.DeviceGroup) {
	m.ID = types.Int64Value(dg.ID)
	m.Name = types.StringValue(dg.Name)
	m.Description = types.StringValue(dg.Description)
	m.Information = types.StringValue(dg.Information)
	m.LocationID = types.Int64Value(dg.LocationID)
	m.Shared = types.BoolValue(dg.Shared)
	m.Members = types.Int64Value(dg.Members)
	m.IsSmartGroup = types.BoolValue(dg.IsSmartGroup)
}
