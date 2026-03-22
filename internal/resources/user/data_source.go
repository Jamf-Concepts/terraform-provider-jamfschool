// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package user

import (
	"context"
	"fmt"

	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ datasource.DataSource = &UserDataSource{}

// NewUserDataSource creates a new user data source.
func NewUserDataSource() datasource.DataSource {
	return &UserDataSource{}
}

// UserDataSource defines the data source implementation.
type UserDataSource struct {
	service *jamfschool.Client
}

func (d *UserDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *UserDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Use this data source to look up a Jamf School user by ID.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Unique identifier of the user.",
				Required:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username.",
				Computed:    true,
			},
			"email": schema.StringAttribute{
				Description: "Email address.",
				Computed:    true,
			},
			"first_name": schema.StringAttribute{
				Description: "First name.",
				Computed:    true,
			},
			"last_name": schema.StringAttribute{
				Description: "Last name.",
				Computed:    true,
			},
			"notes": schema.StringAttribute{
				Description: "Notes.",
				Computed:    true,
			},
			"exclude": schema.BoolAttribute{
				Description: "Don't apply Teacher restrictions.",
				Computed:    true,
			},
			"location_id": schema.Int64Attribute{
				Description: "Location ID.",
				Computed:    true,
			},
			"member_of": schema.SetAttribute{
				Description: "Set of group IDs the user is a member of.",
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"group_names": schema.ListAttribute{
				Description: "Names of groups the user is a member of, in the same order as member_of.",
				Computed:    true,
				ElementType: types.StringType,
			},
			"status": schema.StringAttribute{
				Description: "User status.",
				Computed:    true,
			},
			"device_count": schema.Int64Attribute{
				Description: "Number of devices.",
				Computed:    true,
			},
		},
	}
}

func (d *UserDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *UserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config UserDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := d.service.GetUser(ctx, config.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error reading user", fmt.Sprintf("Could not read user with ID %d: %s", config.ID.ValueInt64(), err))
		return
	}

	config.ID = types.Int64Value(user.ID)
	config.Username = types.StringValue(user.Username)
	config.Email = types.StringValue(user.Email)
	config.FirstName = types.StringValue(user.FirstName)
	config.LastName = types.StringValue(user.LastName)
	config.Notes = types.StringValue(user.Notes)
	config.Exclude = types.BoolValue(user.Exclude)
	config.LocationID = types.Int64Value(user.LocationID)
	config.MemberOf = int64SliceToSet(user.GroupIDs)
	config.GroupNames = stringSliceToList(user.Groups)
	config.Status = types.StringValue(user.Status)
	config.DeviceCount = types.Int64Value(user.DeviceCount)

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
