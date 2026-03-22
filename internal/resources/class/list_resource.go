// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package class

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

var _ list.ListResource = &ClassListResource{}
var _ list.ListResourceWithConfigure = &ClassListResource{}
var _ list.ListResourceWithValidateConfig = &ClassListResource{}

// NewClassListResource returns a new class list resource.
func NewClassListResource() list.ListResource {
	return &ClassListResource{}
}

// ClassListResource lists classes in Jamf School.
type ClassListResource struct {
	service *jamfschool.Client
}

func (r *ClassListResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_class"
}

func (r *ClassListResource) ListResourceConfigSchema(_ context.Context, _ list.ListResourceSchemaRequest, resp *list.ListResourceSchemaResponse) {
	resp.Schema = listschema.Schema{
		MarkdownDescription: "Lists classes in Jamf School.",
		Attributes: map[string]listschema.Attribute{
			"name_prefix": helpers.NamePrefixSchemaAttribute(),
		},
	}
}

func (r *ClassListResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *ClassListResource) ValidateListResourceConfig(_ context.Context, req list.ValidateConfigRequest, resp *list.ValidateConfigResponse) {
	var config helpers.ListConfigModel
	resp.Diagnostics.Append(req.Config.Get(context.Background(), &config)...)
	if resp.Diagnostics.HasError() {
		return
	}
	helpers.ValidateNamePrefix(config, &resp.Diagnostics)
}

func (r *ClassListResource) List(ctx context.Context, req list.ListRequest, resp *list.ListResultsStream) {
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

	items, err := r.service.GetClasses(ctx)
	if err != nil {
		resp.Results = list.ListResultsStreamDiagnostics(diag.Diagnostics{
			diag.NewErrorDiagnostic("Error listing classes", err.Error()),
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
		result.Diagnostics.Append(result.Identity.SetAttribute(ctx, path.Root("uuid"), types.StringValue(item.UUID))...)
		if result.Diagnostics.HasError() {
			results = append(results, result)
			continue
		}

		if req.IncludeResource {
			cls, err := r.service.GetClass(ctx, item.UUID)
			if err != nil {
				result.Diagnostics.AddError("Error reading class", err.Error())
				results = append(results, result)
				continue
			}

			var data ClassResourceModel
			modelFromClass(&data, cls)
			result.Diagnostics.Append(result.Resource.Set(ctx, &data)...)
			results = append(results, result)
			continue
		}

		result.Resource = nil
		results = append(results, result)
	}

	resp.Results = slices.Values(results)
}
