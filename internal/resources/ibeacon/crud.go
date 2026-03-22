// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package ibeacon

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

func (r *IBeaconResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan IBeaconResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := buildCreateInput(&plan)

	beacon, err := r.service.CreateIBeacon(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error creating iBeacon", err.Error())
		return
	}

	modelFromIBeacon(&plan, beacon)
	resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("id"), types.Int64Value(plan.ID.ValueInt64()))...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *IBeaconResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state IBeaconResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	beacon, err := r.service.GetIBeacon(ctx, state.ID.ValueInt64())
	if err != nil {
		if errors.Is(err, jamfschool.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading iBeacon", err.Error())
		return
	}

	modelFromIBeacon(&state, beacon)
	resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("id"), types.Int64Value(state.ID.ValueInt64()))...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *IBeaconResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan IBeaconResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state IBeaconResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := buildUpdateInput(&plan)

	beacon, err := r.service.UpdateIBeacon(ctx, state.ID.ValueInt64(), input)
	if err != nil {
		resp.Diagnostics.AddError("Error updating iBeacon", err.Error())
		return
	}

	modelFromIBeacon(&plan, beacon)
	resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("id"), types.Int64Value(plan.ID.ValueInt64()))...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *IBeaconResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state IBeaconResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.service.DeleteIBeacon(ctx, state.ID.ValueInt64()); err != nil {
		resp.Diagnostics.AddError("Error deleting iBeacon", err.Error())
	}
}

func (r *IBeaconResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", fmt.Sprintf("Expected integer, got: %s", req.ID))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
