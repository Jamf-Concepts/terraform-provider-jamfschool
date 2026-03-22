// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package dep_device

import (
	"context"
	"fmt"

	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &DEPDeviceDataSource{}

func NewDEPDeviceDataSource() datasource.DataSource {
	return &DEPDeviceDataSource{}
}

type DEPDeviceDataSource struct {
	service *jamfschool.Client
}

func (d *DEPDeviceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_dep_device"
}

func (d *DEPDeviceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to look up a Jamf School DEP device by serial number.",
		Attributes: map[string]schema.Attribute{
			"serial_number": schema.StringAttribute{
				Description: "Serial number of the DEP device.",
				Required:    true,
			},
			"id": schema.Int64Attribute{
				Description: "Unique identifier.",
				Computed:    true,
			},
			"user_id": schema.Int64Attribute{
				Description: "Associated user ID.",
				Computed:    true,
			},
			"location_id": schema.Int64Attribute{
				Description: "Location ID.",
				Computed:    true,
			},
			"model": schema.StringAttribute{
				Description: "Device model.",
				Computed:    true,
			},
			"color": schema.StringAttribute{
				Description: "Device color.",
				Computed:    true,
			},
			"status": schema.StringAttribute{
				Description: "DEP status.",
				Computed:    true,
			},
			"device_name": schema.StringAttribute{
				Description: "Device name.",
				Computed:    true,
			},
			"profile_name": schema.StringAttribute{
				Description: "DEP profile name.",
				Computed:    true,
			},
			"placeholder_name": schema.StringAttribute{
				Description: "Placeholder name.",
				Computed:    true,
			},
		},
	}
}

func (d *DEPDeviceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DEPDeviceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config DEPDeviceDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dev, err := d.service.GetDEPDevice(ctx, config.SerialNumber.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading DEP device", fmt.Sprintf("Could not read DEP device with serial number %q: %s", config.SerialNumber.ValueString(), err))
		return
	}

	config.SerialNumber = types.StringValue(dev.SerialNumber)
	config.ID = types.Int64Value(dev.ID)
	config.UserID = types.Int64Value(dev.UserID)
	config.LocationID = types.Int64Value(dev.LocationID)
	config.Model = types.StringValue(dev.Model)
	config.Color = types.StringValue(dev.Color)
	config.Status = types.StringValue(dev.Status)
	config.DeviceName = types.StringValue(dev.DeviceName)
	config.ProfileName = types.StringValue(dev.ProfileName)
	config.PlaceholderName = types.StringValue(dev.PlaceholderName)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
