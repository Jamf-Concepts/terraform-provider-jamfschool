// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package ibeacon

import (
	"context"
	"fmt"

	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &IBeaconDataSource{}

// NewIBeaconDataSource creates a new iBeacon data source.
func NewIBeaconDataSource() datasource.DataSource {
	return &IBeaconDataSource{}
}

// IBeaconDataSource defines the data source implementation.
type IBeaconDataSource struct {
	service *jamfschool.Client
}

func (d *IBeaconDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ibeacon"
}

func (d *IBeaconDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to look up a Jamf School iBeacon by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Unique identifier.",
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
			"uuid": schema.StringAttribute{
				Description: "iBeacon UUID.",
				Computed:    true,
			},
			"major": schema.Int64Attribute{
				Description: "Major value.",
				Computed:    true,
			},
			"minor": schema.Int64Attribute{
				Description: "Minor value.",
				Computed:    true,
			},
		},
	}
}

func (d *IBeaconDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *IBeaconDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config IBeaconDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	beacon, err := d.service.GetIBeacon(ctx, config.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error reading iBeacon", fmt.Sprintf("Could not read iBeacon with ID %d: %s", config.ID.ValueInt64(), err))
		return
	}

	config.ID = types.Int64Value(beacon.ID)
	config.Name = types.StringValue(beacon.Name)
	config.Description = types.StringValue(beacon.Description)
	config.UUID = types.StringValue(beacon.UUID)
	config.Major = types.Int64Value(beacon.Major)
	config.Minor = types.Int64Value(beacon.Minor)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
