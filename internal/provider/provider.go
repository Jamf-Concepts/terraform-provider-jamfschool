// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
	deviceactions "github.com/Jamf-Concepts/terraform-provider-jamfschool/internal/actions/device"
	appds "github.com/Jamf-Concepts/terraform-provider-jamfschool/internal/resources/app"
	classres "github.com/Jamf-Concepts/terraform-provider-jamfschool/internal/resources/class"
	depds "github.com/Jamf-Concepts/terraform-provider-jamfschool/internal/resources/dep_device"
	deviceds "github.com/Jamf-Concepts/terraform-provider-jamfschool/internal/resources/device"
	devicegroupres "github.com/Jamf-Concepts/terraform-provider-jamfschool/internal/resources/device_group"
	ibeaconres "github.com/Jamf-Concepts/terraform-provider-jamfschool/internal/resources/ibeacon"
	locationds "github.com/Jamf-Concepts/terraform-provider-jamfschool/internal/resources/location"
	profileds "github.com/Jamf-Concepts/terraform-provider-jamfschool/internal/resources/profile"
	userres "github.com/Jamf-Concepts/terraform-provider-jamfschool/internal/resources/user"
	usergroupres "github.com/Jamf-Concepts/terraform-provider-jamfschool/internal/resources/user_group"
)

var (
	_ provider.Provider                  = &JamfSchoolProvider{}
	_ provider.ProviderWithActions       = &JamfSchoolProvider{}
	_ provider.ProviderWithListResources = &JamfSchoolProvider{}
)

// JamfSchoolProvider defines the provider implementation.
type JamfSchoolProvider struct {
	version string
}

// JamfSchoolProviderModel describes the provider data model.
type JamfSchoolProviderModel struct {
	URL       types.String `tfsdk:"url"`
	NetworkID types.String `tfsdk:"network_id"`
	APIKey    types.String `tfsdk:"api_key"`
}

func (p *JamfSchoolProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "jamfschool"
	resp.Version = p.version
}

func (p *JamfSchoolProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "The Jamf School provider allows you to manage resources in Jamf School.",
		Attributes: map[string]schema.Attribute{
			"url": schema.StringAttribute{
				MarkdownDescription: "The base URL of your Jamf School instance (e.g. `https://myschool.jamfcloud.com`). Can also be set with the `JAMFSCHOOL_URL` environment variable.",
				Optional:            true,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "The Network ID used for API authentication, found at Devices > Enroll Device(s) in Jamf School. Can also be set with the `JAMFSCHOOL_NETWORK_ID` environment variable.",
				Optional:            true,
			},
			"api_key": schema.StringAttribute{
				MarkdownDescription: "The API key used for authentication, generated at Organization > Settings > API in Jamf School. Can also be set with the `JAMFSCHOOL_API_KEY` environment variable.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *JamfSchoolProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data JamfSchoolProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	url := os.Getenv("JAMFSCHOOL_URL")
	if !data.URL.IsNull() {
		url = data.URL.ValueString()
	}

	networkID := os.Getenv("JAMFSCHOOL_NETWORK_ID")
	if !data.NetworkID.IsNull() {
		networkID = data.NetworkID.ValueString()
	}

	apiKey := os.Getenv("JAMFSCHOOL_API_KEY")
	if !data.APIKey.IsNull() {
		apiKey = data.APIKey.ValueString()
	}

	if url == "" {
		resp.Diagnostics.AddError("Missing URL", "The Jamf School URL must be set in the provider configuration or JAMFSCHOOL_URL environment variable.")
		return
	}
	if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
		resp.Diagnostics.AddError("Invalid URL", "The Jamf School URL must start with https:// or http://.")
		return
	}
	if networkID == "" {
		resp.Diagnostics.AddError("Missing Network ID", "The Jamf School Network ID must be set in the provider configuration or JAMFSCHOOL_NETWORK_ID environment variable.")
		return
	}
	if apiKey == "" {
		resp.Diagnostics.AddError("Missing API Key", "The Jamf School API key must be set in the provider configuration or JAMFSCHOOL_API_KEY environment variable.")
		return
	}

	var opts []jamfschool.Option
	if shouldEnableHTTPLogging() {
		opts = append(opts, jamfschool.WithLogger(NewTerraformLogger()))
	}
	c := jamfschool.NewClient(url, networkID, apiKey, opts...)

	resp.DataSourceData = c
	resp.ResourceData = c
	resp.ActionData = c
	resp.ListResourceData = c
}

func (p *JamfSchoolProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		userres.NewUserResource,
		usergroupres.NewUserGroupResource,
		devicegroupres.NewDeviceGroupResource,
		classres.NewClassResource,
		ibeaconres.NewIBeaconResource,
		appds.NewAppResource,
	}
}

func (p *JamfSchoolProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		userres.NewUserDataSource,
		usergroupres.NewUserGroupDataSource,
		devicegroupres.NewDeviceGroupDataSource,
		classres.NewClassDataSource,
		ibeaconres.NewIBeaconDataSource,
		deviceds.NewDeviceDataSource,
		appds.NewAppDataSource,
		profileds.NewProfileDataSource,
		locationds.NewLocationDataSource,
		depds.NewDEPDeviceDataSource,
	}
}

func (p *JamfSchoolProvider) Actions(ctx context.Context) []func() action.Action {
	return []func() action.Action{
		deviceactions.NewEraseDeviceAction,
		deviceactions.NewRestartDeviceAction,
		deviceactions.NewRefreshDeviceAction,
		deviceactions.NewUnenrollDeviceAction,
		deviceactions.NewClearDeviceActivationLockAction,
		deviceactions.NewMoveDeviceToTrashAction,
		deviceactions.NewPutDeviceBackAction,
		deviceactions.NewUpdateDeviceESIMAction,
	}
}

func (p *JamfSchoolProvider) ListResources(ctx context.Context) []func() list.ListResource {
	return []func() list.ListResource{
		userres.NewUserListResource,
		usergroupres.NewUserGroupListResource,
		devicegroupres.NewDeviceGroupListResource,
		classres.NewClassListResource,
		ibeaconres.NewIBeaconListResource,
		appds.NewAppListResource,
	}
}

// New returns a new provider factory function.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &JamfSchoolProvider{
			version: version,
		}
	}
}
