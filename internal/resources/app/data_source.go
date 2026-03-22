// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package app

import (
	"context"
	"fmt"

	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &AppDataSource{}

// NewAppDataSource creates a new app data source.
func NewAppDataSource() datasource.DataSource {
	return &AppDataSource{}
}

// AppDataSource defines the data source implementation.
type AppDataSource struct {
	service *jamfschool.Client
}

func (d *AppDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app"
}

func (d *AppDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to look up a Jamf School app by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Unique identifier.",
				Required:    true,
			},
			"bundle_id": schema.StringAttribute{
				Description: "Bundle identifier.",
				Computed:    true,
			},
			"adam_id": schema.Int64Attribute{
				Description: "Adam ID (App Store ID).",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "App name.",
				Computed:    true,
			},
			"vendor": schema.StringAttribute{
				Description: "Vendor.",
				Computed:    true,
			},
			"version": schema.StringAttribute{
				Description: "Version.",
				Computed:    true,
			},
			"platform": schema.StringAttribute{
				Description: "Platform.",
				Computed:    true,
			},
		},
	}
}

func (d *AppDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AppDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config AppDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	app, err := d.service.GetApp(ctx, config.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error reading app", fmt.Sprintf("Could not read app with ID %d: %s", config.ID.ValueInt64(), err))
		return
	}

	config.ID = types.Int64Value(app.ID)
	config.BundleID = types.StringValue(app.BundleID)
	config.AdamID = types.Int64Value(app.AdamID)
	config.Name = types.StringValue(app.Name)
	config.Vendor = types.StringValue(app.Vendor)
	config.Version = types.StringValue(app.Version)
	config.Platform = types.StringValue(app.Platform)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
