// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package user_group

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

func (r *UserGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UserGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := buildCreateInput(&plan)

	id, err := r.service.CreateGroup(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error creating user group", err.Error())
		return
	}

	group, err := r.service.GetGroup(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error reading created user group", err.Error())
		return
	}

	modelFromGroup(&plan, group)
	resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("id"), types.Int64Value(plan.ID.ValueInt64()))...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state UserGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	group, err := r.service.GetGroup(ctx, state.ID.ValueInt64())
	if err != nil {
		if errors.Is(err, jamfschool.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading user group", err.Error())
		return
	}

	modelFromGroup(&state, group)
	resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("id"), types.Int64Value(state.ID.ValueInt64()))...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *UserGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan UserGroupResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state UserGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := buildUpdateInput(&plan)

	if err := r.service.UpdateGroup(ctx, state.ID.ValueInt64(), input); err != nil {
		resp.Diagnostics.AddError("Error updating user group", err.Error())
		return
	}

	group, err := r.service.GetGroup(ctx, state.ID.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError("Error reading updated user group", err.Error())
		return
	}

	modelFromGroup(&plan, group)
	resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("id"), types.Int64Value(plan.ID.ValueInt64()))...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *UserGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UserGroupResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.service.DeleteGroup(ctx, state.ID.ValueInt64()); err != nil {
		resp.Diagnostics.AddError("Error deleting user group", err.Error())
	}
}

func (r *UserGroupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError("Invalid import ID", fmt.Sprintf("Expected integer, got: %s", req.ID))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
}
