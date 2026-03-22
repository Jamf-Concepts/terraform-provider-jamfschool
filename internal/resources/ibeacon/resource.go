// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package ibeacon

import (
	"context"

	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var (
	_ resource.Resource                = &IBeaconResource{}
	_ resource.ResourceWithImportState = &IBeaconResource{}
	_ resource.ResourceWithIdentity    = &IBeaconResource{}
)

// NewIBeaconResource creates a new iBeacon resource.
func NewIBeaconResource() resource.Resource {
	return &IBeaconResource{}
}

// IBeaconResource defines the resource implementation.
type IBeaconResource struct {
	service *jamfschool.Client
}

func (r *IBeaconResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ibeacon"
}

func (r *IBeaconResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Jamf School iBeacon.",
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "Unique identifier of the iBeacon.",
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Name of the iBeacon.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Description: "Description.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"uuid": schema.StringAttribute{
				Description: "iBeacon proximity UUID. A random UUID is generated if not specified.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						uuidRegexp,
						"must be a valid UUID (e.g. A1B2C3D4-E5F6-7890-ABCD-EF1234567890)",
					),
				},
			},
			"major": schema.Int64Attribute{
				Description: "Major value (0-65535).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Int64{
					int64validator.Between(0, 65535),
				},
			},
			"minor": schema.Int64Attribute{
				Description: "Minor value (0-65535).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Validators: []validator.Int64{
					int64validator.Between(0, 65535),
				},
			},
		},
	}
}

func (r *IBeaconResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *IBeaconResource) IdentitySchema(_ context.Context, _ resource.IdentitySchemaRequest, resp *resource.IdentitySchemaResponse) {
	resp.IdentitySchema = identityschema.Schema{
		Attributes: map[string]identityschema.Attribute{
			"id": identityschema.Int64Attribute{
				RequiredForImport: true,
				Description:       "Unique identifier of the iBeacon.",
			},
		},
	}
}
