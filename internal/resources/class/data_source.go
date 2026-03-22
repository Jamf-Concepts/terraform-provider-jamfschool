// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package class

import (
	"context"
	"fmt"

	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ClassDataSource{}

// NewClassDataSource creates a new class data source.
func NewClassDataSource() datasource.DataSource {
	return &ClassDataSource{}
}

// ClassDataSource defines the data source implementation.
type ClassDataSource struct {
	service *jamfschool.Client
}

func (d *ClassDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_class"
}

func (d *ClassDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to look up a Jamf School class by UUID.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Description: "UUID of the class.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description.",
				Computed:    true,
			},
			"location_id": schema.Int64Attribute{
				Description: "Location ID.",
				Computed:    true,
			},
			"source": schema.StringAttribute{
				Description: "Source.",
				Computed:    true,
			},
			"student_count": schema.Int64Attribute{
				Description: "Number of students.",
				Computed:    true,
			},
			"teacher_count": schema.Int64Attribute{
				Description: "Number of teachers.",
				Computed:    true,
			},
			"device_group_id": schema.Int64Attribute{
				Description: "Associated device group ID.",
				Computed:    true,
			},
			"device_count": schema.Int64Attribute{
				Description: "Number of devices in the class.",
				Computed:    true,
			},
			"device_udids": schema.ListAttribute{
				Description: "UDIDs of devices assigned to the class.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"students": schema.SetAttribute{
				Description: "Set of user IDs assigned as students.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"teachers": schema.SetAttribute{
				Description: "Set of user IDs assigned as teachers.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
		},
	}
}

func (d *ClassDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	svc, ok := req.ProviderData.(*jamfschool.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected DataSource Configure Type", fmt.Sprintf("Expected *jamfschool.Client, got %T", req.ProviderData))
		return
	}
	d.service = svc
}

func (d *ClassDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ClassDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	c, err := d.service.GetClass(ctx, config.UUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading class", fmt.Sprintf("Could not read class with UUID %q: %s", config.UUID.ValueString(), err))
		return
	}

	config.UUID = types.StringValue(c.UUID)
	config.Name = types.StringValue(c.Name)
	config.Description = types.StringValue(c.Description)
	config.LocationID = types.Int64Value(c.LocationID)
	config.Source = types.StringValue(c.Source)
	config.StudentCount = types.Int64Value(c.StudentCount)
	config.TeacherCount = types.Int64Value(c.TeacherCount)
	config.DeviceGroupID = types.Int64Value(c.DeviceGroupID)
	config.DeviceCount = types.Int64Value(c.DeviceCount)

	devices, err := d.service.GetClassDevices(ctx, config.UUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading class devices", fmt.Sprintf("Could not read devices for class %q: %s", config.UUID.ValueString(), err))
		return
	}
	deviceUDIDs := make([]attr.Value, len(devices))
	for i, dev := range devices {
		deviceUDIDs[i] = types.StringValue(dev.UDID)
	}
	config.DeviceUDIDs, _ = types.ListValue(types.StringType, deviceUDIDs)

	studentIDs := make([]attr.Value, len(c.Students))
	for i, s := range c.Students {
		studentIDs[i] = types.Int64Value(s.ID)
	}
	config.Students, _ = types.SetValue(types.Int64Type, studentIDs)

	teacherIDs := make([]attr.Value, len(c.Teachers))
	for i, t := range c.Teachers {
		teacherIDs[i] = types.Int64Value(t.ID)
	}
	config.Teachers, _ = types.SetValue(types.Int64Type, teacherIDs)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
