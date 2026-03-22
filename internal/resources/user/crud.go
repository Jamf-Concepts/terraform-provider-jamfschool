// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package user

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

func (r *UserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := buildCreateInput(&plan)

	id, err := r.service.CreateUser(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error creating user", err.Error())
		return
	}

	user, err := r.service.GetUser(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error reading created user", err.Error())
		return
	}

	modelFromUser(&plan, user)
	resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("id"), types.Int64Value(plan.ID.ValueInt64()))...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	user, err := r.service.GetUser(ctx, state.ID.ValueInt64())
	if err != nil {
		if errors.Is(err, jamfschool.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading user", err.Error())
		return
	}

	modelFromUser(&state, user)
	resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("id"), types.Int64Value(state.ID.ValueInt64()))...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan UserResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If location changed, migrate the user to the new location.
	locationChanged := !plan.LocationID.IsUnknown() && !state.LocationID.IsUnknown() &&
		plan.LocationID.ValueInt64() != state.LocationID.ValueInt64()

	if locationChanged {
		onlyUser := !plan.MoveDevicesOnLocationChange.ValueBool()
		if err := r.service.MigrateUser(ctx, state.ID.ValueInt64(), plan.LocationID.ValueInt64(), onlyUser); err != nil {
			resp.Diagnostics.AddError("Error migrating user",
				fmt.Sprintf("Failed to move user %d to location %d: %s", state.ID.ValueInt64(), plan.LocationID.ValueInt64(), err))
			return
		}
	}

	// Apply regular field updates.
	input := buildUpdateInput(&plan, &state)
	if err := r.service.UpdateUser(ctx, state.ID.ValueInt64(), input); err != nil {
		resp.Diagnostics.AddError("Error updating user", err.Error())
		return
	}

	user, err := r.service.GetUser(ctx, state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error reading updated user", err.Error())
		return
	}

	modelFromUser(&plan, user)
	resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("id"), types.Int64Value(plan.ID.ValueInt64()))...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.service.DeleteUser(ctx, state.ID.ValueInt64()); err != nil {
		resp.Diagnostics.AddError("Error deleting user", err.Error())
	}
}

func (r *UserResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", fmt.Sprintf("Expected integer, got: %s", req.ID))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
