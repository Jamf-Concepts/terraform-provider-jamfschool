// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package class

import (
	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// modelFromClass maps a Jamf School API class to the Terraform resource model.
func modelFromClass(m *ClassResourceModel, c *jamfschool.Class) {
	m.UUID = types.StringValue(c.UUID)
	m.Name = types.StringValue(c.Name)
	m.Description = types.StringValue(c.Description)
	m.LocationID = types.Int64Value(c.LocationID)
	m.Source = types.StringValue(c.Source)
	m.StudentCount = types.Int64Value(c.StudentCount)
	m.TeacherCount = types.Int64Value(c.TeacherCount)
	m.DeviceGroupID = types.Int64Value(c.DeviceGroupID)

	studentIDs := make([]attr.Value, len(c.Students))
	for i, s := range c.Students {
		studentIDs[i] = types.Int64Value(s.ID)
	}
	m.Students, _ = types.SetValue(types.Int64Type, studentIDs)

	teacherIDs := make([]attr.Value, len(c.Teachers))
	for i, t := range c.Teachers {
		teacherIDs[i] = types.Int64Value(t.ID)
	}
	m.Teachers, _ = types.SetValue(types.Int64Type, teacherIDs)
}
