// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package device

import (
	"context"
	"fmt"

	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &DeviceDataSource{}

func NewDeviceDataSource() datasource.DataSource {
	return &DeviceDataSource{}
}

type DeviceDataSource struct {
	service *jamfschool.Client
}

func (d *DeviceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

func (d *DeviceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to look up a Jamf School device by UDID.",
		Attributes: map[string]schema.Attribute{
			"udid": schema.StringAttribute{
				Description: "UDID of the device.",
				Required:    true,
			},
			"serial_number": schema.StringAttribute{
				Description: "Serial number.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Device name.",
				Computed:    true,
			},
			"is_managed": schema.BoolAttribute{
				Description: "Whether the device is managed.",
				Computed:    true,
			},
			"is_supervised": schema.BoolAttribute{
				Description: "Whether the device is supervised.",
				Computed:    true,
			},
			"battery_level": schema.Float64Attribute{
				Description: "Battery level.",
				Computed:    true,
			},
			"total_capacity": schema.Float64Attribute{
				Description: "Total storage capacity.",
				Computed:    true,
			},
			"notes": schema.StringAttribute{
				Description: "Notes.",
				Computed:    true,
			},
			"last_checkin": schema.StringAttribute{
				Description: "Last check-in time.",
				Computed:    true,
			},
			"location_id": schema.Int64Attribute{
				Description: "Location ID.",
				Computed:    true,
			},
			"enroll_type": schema.StringAttribute{
				Description: "Enrollment type.",
				Computed:    true,
			},
			"model_name": schema.StringAttribute{
				Description: "Model name.",
				Computed:    true,
			},
			"model_identifier": schema.StringAttribute{
				Description: "Model identifier.",
				Computed:    true,
			},
			"os_version": schema.StringAttribute{
				Description: "OS version.",
				Computed:    true,
			},
			"os_prefix": schema.StringAttribute{
				Description: "OS prefix (e.g. iOS, macOS).",
				Computed:    true,
			},
			"in_trash": schema.BoolAttribute{
				Description: "Whether the device is in trash.",
				Computed:    true,
			},
		},
	}
}

func (d *DeviceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DeviceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config DeviceDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dev, err := d.service.GetDevice(ctx, config.UDID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading device", fmt.Sprintf("Could not read device with UDID %q: %s", config.UDID.ValueString(), err))
		return
	}

	config.UDID = types.StringValue(dev.UDID)
	config.SerialNumber = types.StringValue(dev.SerialNumber)
	config.Name = types.StringValue(dev.Name)
	config.IsManaged = types.BoolValue(dev.IsManaged)
	config.IsSupervised = types.BoolValue(dev.IsSupervised)
	config.BatteryLevel = types.Float64Value(dev.BatteryLevel)
	config.TotalCapacity = types.Float64Value(dev.TotalCapacity)
	config.Notes = types.StringValue(dev.Notes)
	config.LastCheckin = types.StringValue(dev.LastCheckin)
	config.LocationID = types.Int64Value(dev.LocationID)
	config.EnrollType = types.StringValue(dev.DeviceEnrollType)
	config.InTrash = types.BoolValue(dev.InTrash)

	config.ModelName = types.StringValue(dev.Model.Name)
	config.ModelIdentifier = types.StringValue(dev.Model.Identifier)
	config.OSVersion = types.StringValue(dev.OS.Version)
	config.OSPrefix = types.StringValue(dev.OS.Prefix)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
