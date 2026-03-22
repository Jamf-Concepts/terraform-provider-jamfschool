// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package profile

import (
	"context"
	"fmt"

	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &ProfileDataSource{}

func NewProfileDataSource() datasource.DataSource {
	return &ProfileDataSource{}
}

type ProfileDataSource struct {
	service *jamfschool.Client
}

func (d *ProfileDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_profile"
}

func (d *ProfileDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to look up a Jamf School configuration profile by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Unique identifier.",
				Required:    true,
			},
			"location_id": schema.Int64Attribute{
				Description: "Location ID.",
				Computed:    true,
			},
			"identifier": schema.StringAttribute{
				Description: "Profile identifier.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "Profile name.",
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description.",
				Computed:    true,
			},
			"platform": schema.StringAttribute{
				Description: "Platform.",
				Computed:    true,
			},
		},
	}
}

func (d *ProfileDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ProfileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config ProfileDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	p, err := d.service.GetProfile(ctx, config.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error reading profile", fmt.Sprintf("Could not read profile with ID %d: %s", config.ID.ValueInt64(), err))
		return
	}

	config.ID = types.Int64Value(p.ID)
	config.LocationID = types.Int64Value(p.LocationID)
	config.Identifier = types.StringValue(p.Identifier)
	config.Name = types.StringValue(p.Name)
	config.Description = types.StringValue(p.Description)
	config.Platform = types.StringValue(p.Platform)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
