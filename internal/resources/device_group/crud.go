// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package device_group

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

func (r *DeviceGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan DeviceGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := buildCreateInput(&plan)

	id, err := r.service.CreateDeviceGroup(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error creating device group", err.Error())
		return
	}

	dg, err := r.service.GetDeviceGroup(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error reading created device group", err.Error())
		return
	}

	modelFromDeviceGroup(&plan, dg)

	// Add member devices if specified.
	plannedUDIDs := extractStringSet(plan.MemberUDIDs)
	if len(plannedUDIDs) > 0 {
		if err := r.service.AddDevicesToGroup(ctx, id, plannedUDIDs); err != nil {
			resp.Diagnostics.AddError("Error adding devices to group", err.Error())
			return
		}
	}
	// Read back current members.
	memberUDIDs, err := r.service.GetDeviceGroupMembers(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error reading device group members", err.Error())
		return
	}
	plan.MemberUDIDs = stringSliceToSet(memberUDIDs)

	resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("id"), types.Int64Value(plan.ID.ValueInt64()))...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DeviceGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state DeviceGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dg, err := r.service.GetDeviceGroup(ctx, state.ID.ValueInt64())
	if err != nil {
		if errors.Is(err, jamfschool.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading device group", err.Error())
		return
	}

	modelFromDeviceGroup(&state, dg)

	memberUDIDs, err := r.service.GetDeviceGroupMembers(ctx, state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error reading device group members", err.Error())
		return
	}
	state.MemberUDIDs = stringSliceToSet(memberUDIDs)

	resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("id"), types.Int64Value(state.ID.ValueInt64()))...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *DeviceGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DeviceGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state DeviceGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := buildUpdateInput(&plan)

	if err := r.service.UpdateDeviceGroup(ctx, state.ID.ValueInt64(), input); err != nil {
		resp.Diagnostics.AddError("Error updating device group", err.Error())
		return
	}

	// Sync device group members.
	planUDIDs := extractStringSet(plan.MemberUDIDs)
	stateUDIDs := extractStringSet(state.MemberUDIDs)

	toAdd := stringSetDiff(planUDIDs, stateUDIDs)
	toRemove := stringSetDiff(stateUDIDs, planUDIDs)

	if len(toAdd) > 0 {
		if err := r.service.AddDevicesToGroup(ctx, state.ID.ValueInt64(), toAdd); err != nil {
			resp.Diagnostics.AddError("Error adding devices to group", err.Error())
			return
		}
	}
	if len(toRemove) > 0 {
		if err := r.service.RemoveDevicesFromGroup(ctx, state.ID.ValueInt64(), toRemove); err != nil {
			resp.Diagnostics.AddError("Error removing devices from group", err.Error())
			return
		}
	}

	dg, err := r.service.GetDeviceGroup(ctx, state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error reading updated device group", err.Error())
		return
	}

	modelFromDeviceGroup(&plan, dg)

	// Read back current members.
	memberUDIDs, err := r.service.GetDeviceGroupMembers(ctx, state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error reading device group members", err.Error())
		return
	}
	plan.MemberUDIDs = stringSliceToSet(memberUDIDs)

	resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("id"), types.Int64Value(plan.ID.ValueInt64()))...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *DeviceGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state DeviceGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.service.DeleteDeviceGroup(ctx, state.ID.ValueInt64()); err != nil {
		resp.Diagnostics.AddError("Error deleting device group", err.Error())
	}
}

func (r *DeviceGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", fmt.Sprintf("Expected integer, got: %s", req.ID))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
