// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package class

import (
	"context"

	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &ClassResource{}
	_ resource.ResourceWithImportState = &ClassResource{}
	_ resource.ResourceWithIdentity    = &ClassResource{}
)

// NewClassResource creates a new class resource.
func NewClassResource() resource.Resource {
	return &ClassResource{}
}

// ClassResource defines the resource implementation.
type ClassResource struct {
	service *jamfschool.Client
}

func (r *ClassResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_class"
}

func (r *ClassResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Jamf School class. Classes are identified by UUID.",
		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				Description: "UUID of the class.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the class.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description of the class.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"location_id": schema.Int64Attribute{
				Description: "Location ID.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"source": schema.StringAttribute{
				Description: "Source of the class.",
				Computed:    true,
			},
			"student_count": schema.Int64Attribute{
				Description: "Number of students.",
				Computed:    true,
			},
			"teacher_count": schema.Int64Attribute{
				Description: "Number of teachers.",
				Computed:    true,
			},
			"device_group_id": schema.Int64Attribute{
				Description: "Associated device group ID.",
				Computed:    true,
			},
			"students": schema.SetAttribute{
				Description: "Set of user IDs assigned as students.",
				Optional:    true,
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"teachers": schema.SetAttribute{
				Description: "Set of user IDs assigned as teachers.",
				Optional:    true,
				Computed:    true,
				ElementType: types.Int64Type,
			},
		},
	}
}

func (r *ClassResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	svc, ok := req.ProviderData.(*jamfschool.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", "Expected *jamfschool.Client")
		return
	}
	r.service = svc
}

func (r *ClassResource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"uuid": identityschema.StringAttribute{
				RequiredForImport: true,
				Description:       "UUID of the class.",
			},
		},
	}
}
