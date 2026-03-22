// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package ibeacon

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

var _ list.ListResource = &IBeaconListResource{}
var _ list.ListResourceWithConfigure = &IBeaconListResource{}
var _ list.ListResourceWithValidateConfig = &IBeaconListResource{}

// NewIBeaconListResource returns a new iBeacon list resource.
func NewIBeaconListResource() list.ListResource {
	return &IBeaconListResource{}
}

// IBeaconListResource lists iBeacons in Jamf School.
type IBeaconListResource struct {
	service *jamfschool.Client
}

func (r *IBeaconListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ibeacon"
}

func (r *IBeaconListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{
		MarkdownDescription: "Lists iBeacons in Jamf School.",
		Attributes: map[string]listschema.Attribute{
			"name_prefix": helpers.NamePrefixSchemaAttribute(),
		},
	}
}

func (r *IBeaconListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IBeaconListResource) ValidateListResourceConfig(_ context.Context, req list.ValidateConfigRequest, resp *list.ValidateConfigResponse) {
	var config helpers.ListConfigModel
	resp.Diagnostics.Append(req.Config.Get(context.Background(), &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	helpers.ValidateNamePrefix(config, &resp.Diagnostics)
}

func (r *IBeaconListResource) List(ctx context.Context, req list.ListRequest, resp *list.ListResultsStream) {
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

	items, err := r.service.GetIBeacons(ctx)
	if err != nil {
		resp.Results = list.ListResultsStreamDiagnostics(diag.Diagnostics{
			diag.NewErrorDiagnostic("Error listing iBeacons", err.Error()),
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
			beacon, err := r.service.GetIBeacon(ctx, item.ID)
			if err != nil {
				result.Diagnostics.AddError("Error reading iBeacon", err.Error())
				results = append(results, result)
				continue
			}

			var data IBeaconResourceModel
			modelFromIBeacon(&data, beacon)
			result.Diagnostics.Append(result.Resource.Set(ctx, &data)...)
			results = append(results, result)
			continue
		}

		result.Resource = nil
		results = append(results, result)
	}

	resp.Results = slices.Values(results)
}
