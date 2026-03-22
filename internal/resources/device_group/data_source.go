// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package device_group

import (
	"context"
	"fmt"

	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &DeviceGroupDataSource{}

// NewDeviceGroupDataSource creates a new device group data source.
func NewDeviceGroupDataSource() datasource.DataSource {
	return &DeviceGroupDataSource{}
}

// DeviceGroupDataSource defines the data source implementation.
type DeviceGroupDataSource struct {
	service *jamfschool.Client
}

func (d *DeviceGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_group"
}

func (d *DeviceGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to look up a Jamf School device group by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Unique identifier of the device group.",
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
			"information": schema.StringAttribute{
				Description: "Information text for the group.",
				Computed:    true,
			},
			"location_id": schema.Int64Attribute{
				Description: "Location ID.",
				Computed:    true,
			},
			"shared": schema.BoolAttribute{
				Description: "Whether the group is shared.",
				Computed:    true,
			},
			"members": schema.Int64Attribute{
				Description: "Number of members.",
				Computed:    true,
			},
			"is_smart_group": schema.BoolAttribute{
				Description: "Whether this is a smart group.",
				Computed:    true,
			},
		},
	}
}

func (d *DeviceGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DeviceGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config DeviceGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dg, err := d.service.GetDeviceGroup(ctx, config.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error reading device group", fmt.Sprintf("Could not read device group with ID %d: %s", config.ID.ValueInt64(), err))
		return
	}

	config.ID = types.Int64Value(dg.ID)
	config.Name = types.StringValue(dg.Name)
	config.Description = types.StringValue(dg.Description)
	config.Information = types.StringValue(dg.Information)
	config.LocationID = types.Int64Value(dg.LocationID)
	config.Shared = types.BoolValue(dg.Shared)
	config.Members = types.Int64Value(dg.Members)
	config.IsSmartGroup = types.BoolValue(dg.IsSmartGroup)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
