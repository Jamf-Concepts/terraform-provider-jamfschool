// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package location

import (
	"context"
	"fmt"

	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
	"github.com/Jamf-Concepts/terraform-provider-jamfschool/internal/common/helpers"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &LocationDataSource{}

func NewLocationDataSource() datasource.DataSource {
	return &LocationDataSource{}
}

type LocationDataSource struct {
	service *jamfschool.Client
}

func (d *LocationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_location"
}

func (d *LocationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to look up a Jamf School location by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Unique identifier.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Location name.",
				Computed:    true,
			},
			"is_district": schema.BoolAttribute{
				Description: "Whether this is a district.",
				Computed:    true,
			},
			"street": schema.StringAttribute{
				Description: "Street.",
				Computed:    true,
			},
			"street_number": schema.StringAttribute{
				Description: "Street number.",
				Computed:    true,
			},
			"postal_code": schema.StringAttribute{
				Description: "Postal code.",
				Computed:    true,
			},
			"city": schema.StringAttribute{
				Description: "City.",
				Computed:    true,
			},
			"source": schema.StringAttribute{
				Description: "Source.",
				Computed:    true,
			},
			"asm_identifier": schema.StringAttribute{
				Description: "ASM identifier.",
				Computed:    true,
			},
			"school_number": schema.StringAttribute{
				Description: "School number.",
				Computed:    true,
			},
		},
	}
}

func (d *LocationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LocationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config LocationDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	loc, err := d.service.GetLocation(ctx, config.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error reading location", fmt.Sprintf("Could not read location with ID %d: %s", config.ID.ValueInt64(), err))
		return
	}

	config.ID = types.Int64Value(loc.ID)
	config.Name = types.StringValue(loc.Name)
	config.IsDistrict = types.BoolValue(loc.IsDistrict)
	config.Street = helpers.StringPtrValueOrNull(loc.Street)
	config.StreetNumber = helpers.StringPtrValueOrNull(loc.StreetNumber)
	config.PostalCode = helpers.StringPtrValueOrNull(loc.PostalCode)
	config.City = helpers.StringPtrValueOrNull(loc.City)
	config.Source = types.StringValue(loc.Source)
	config.ASMIdentifier = helpers.StringPtrValueOrNull(loc.ASMIdentifier)
	config.SchoolNumber = types.StringValue(loc.SchoolNumber)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
