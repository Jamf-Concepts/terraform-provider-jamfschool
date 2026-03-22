// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package device_group

import (
	"context"
	"fmt"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/list"
	listschema "github.com/hashicorp/terraform-plugin-framework/list/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
	"github.com/Jamf-Concepts/terraform-provider-jamfschool/internal/common/helpers"
)

var _ list.ListResource = &DeviceGroupListResource{}
var _ list.ListResourceWithConfigure = &DeviceGroupListResource{}
var _ list.ListResourceWithValidateConfig = &DeviceGroupListResource{}

// NewDeviceGroupListResource returns a new device group list resource.
func NewDeviceGroupListResource() list.ListResource {
	return &DeviceGroupListResource{}
}

// DeviceGroupListResource lists device groups in Jamf School.
type DeviceGroupListResource struct {
	service *jamfschool.Client
}

func (r *DeviceGroupListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_group"
}

func (r *DeviceGroupListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{
		MarkdownDescription: "Lists device groups in Jamf School.",
		Attributes: map[string]listschema.Attribute{
			"name_prefix": helpers.NamePrefixSchemaAttribute(),
		},
	}
}

func (r *DeviceGroupListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	svc, ok := req.ProviderData.(*jamfschool.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Configure Type", fmt.Sprintf("Expected *jamfschool.Client, got %T", req.ProviderData))
		return
	}
	r.service = svc
}

func (r *DeviceGroupListResource) ValidateListResourceConfig(_ context.Context, req list.ValidateConfigRequest, resp *list.ValidateConfigResponse) {
	var config helpers.ListConfigModel
	resp.Diagnostics.Append(req.Config.Get(context.Background(), &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	helpers.ValidateNamePrefix(config, &resp.Diagnostics)
}

func (r *DeviceGroupListResource) List(ctx context.Context, req list.ListRequest, resp *list.ListResultsStream) {
	if r.service == nil {
		resp.Results = list.ListResultsStreamDiagnostics(diag.Diagnostics{
			diag.NewErrorDiagnostic("Missing Jamf School service", "The provider was not configured for list resources."),
		})
		return
	}

	var config helpers.ListConfigModel
	configDiags := req.Config.Get(ctx, &config)
	if configDiags.HasError() {
		resp.Results = list.ListResultsStreamDiagnostics(configDiags)
		return
	}

	items, err := r.service.GetDeviceGroups(ctx)
	if err != nil {
		resp.Results = list.ListResultsStreamDiagnostics(diag.Diagnostics{
			diag.NewErrorDiagnostic("Error listing device groups", err.Error()),
		})
		return
	}

	results := make([]list.ListResult, 0, len(items))
	for _, item := range items {
		if !helpers.MatchesNamePrefix(config, item.Name) {
			continue
		}
		if req.Limit > 0 && int64(len(results)) >= req.Limit {
			break
		}

		result := req.NewListResult(ctx)
		result.DisplayName = item.Name
		result.Diagnostics.Append(result.Identity.SetAttribute(ctx, path.Root("id"), types.Int64Value(item.ID))...)
		if result.Diagnostics.HasError() {
			results = append(results, result)
			continue
		}

		if req.IncludeResource {
			dg, err := r.service.GetDeviceGroup(ctx, item.ID)
			if err != nil {
				result.Diagnostics.AddError("Error reading device group", err.Error())
				results = append(results, result)
				continue
			}

			var data DeviceGroupResourceModel
			modelFromDeviceGroup(&data, dg)
			result.Diagnostics.Append(result.Resource.Set(ctx, &data)...)
			results = append(results, result)
			continue
		}

		result.Resource = nil
		results = append(results, result)
	}

	resp.Results = slices.Values(results)
}
