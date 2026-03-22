// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package app

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *AppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan AppResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := jamfschool.AppCreateInput{
		AdamID:      plan.AdamID.ValueInt64(),
		CountryCode: plan.CountryCode.ValueString(),
	}

	id, err := r.service.CreateApp(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error creating app", err.Error())
		return
	}

	app, err := r.service.GetApp(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error reading created app", err.Error())
		return
	}

	modelFromApp(&plan, app)
	resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("id"), types.Int64Value(plan.ID.ValueInt64()))...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *AppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state AppResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	app, err := r.service.GetApp(ctx, state.ID.ValueInt64())
	if err != nil {
		if errors.Is(err, jamfschool.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading app", err.Error())
		return
	}

	modelFromApp(&state, app)
	resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("id"), types.Int64Value(state.ID.ValueInt64()))...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *AppResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	// No update endpoint — all mutable fields use RequiresReplace.
	resp.Diagnostics.AddError("Update not supported", "App resources cannot be updated. Changes require replacement.")
}

func (r *AppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state AppResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.service.TrashApp(ctx, state.ID.ValueInt64()); err != nil {
		resp.Diagnostics.AddError("Error trashing app", err.Error())
	}
}

func (r *AppResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", fmt.Sprintf("Expected integer, got: %s", req.ID))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
