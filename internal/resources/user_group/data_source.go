// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package user_group

import (
	"context"
	"fmt"

	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &UserGroupDataSource{}

// NewUserGroupDataSource creates a new user group data source.
func NewUserGroupDataSource() datasource.DataSource {
	return &UserGroupDataSource{}
}

// UserGroupDataSource defines the data source implementation.
type UserGroupDataSource struct {
	service *jamfschool.Client
}

func (d *UserGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user_group"
}

func (d *UserGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to look up a Jamf School user group by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Unique identifier of the group.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the group.",
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
			"user_count": schema.Int64Attribute{
				Description: "Number of users.",
				Computed:    true,
			},
			"jamf_school_teacher": schema.StringAttribute{
				Description: "Jamf School Teacher feature status.",
				Computed:    true,
			},
			"jamf_parent": schema.StringAttribute{
				Description: "Jamf Parent feature status.",
				Computed:    true,
			},
		},
	}
}

func (d *UserGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UserGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config UserGroupDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := d.service.GetGroup(ctx, config.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error reading user group", fmt.Sprintf("Could not read user group with ID %d: %s", config.ID.ValueInt64(), err))
		return
	}

	config.ID = types.Int64Value(group.ID)
	config.Name = types.StringValue(group.Name)
	config.Description = types.StringValue(group.Description)
	config.LocationID = types.Int64Value(group.LocationID)
	config.UserCount = types.Int64Value(group.UserCount)
	config.JamfSchoolTeacher = types.StringValue(group.ACL.Teacher)
	config.JamfParent = types.StringValue(group.ACL.Parent)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
