// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package class

import (
	"context"
	"errors"

	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *ClassResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ClassResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := buildCreateInput(&plan)

	uuid, err := r.service.CreateClass(ctx, input)
	if err != nil {
		resp.Diagnostics.AddError("Error creating class", err.Error())
		return
	}

	c, err := r.service.GetClass(ctx, uuid)
	if err != nil {
		resp.Diagnostics.AddError("Error reading created class", err.Error())
		return
	}

	modelFromClass(&plan, c)

	// Assign class members if specified.
	studentIDs := extractInt64Set(plan.Students)
	teacherIDs := extractInt64Set(plan.Teachers)
	if len(studentIDs) > 0 || len(teacherIDs) > 0 {
		if err := r.service.AssignClassUsers(ctx, uuid, studentIDs, teacherIDs); err != nil {
			resp.Diagnostics.AddError("Error assigning class members", err.Error())
			return
		}
		// Re-read to get updated counts and member lists
		c, err = r.service.GetClass(ctx, uuid)
		if err != nil {
			resp.Diagnostics.AddError("Error reading class after member assignment", err.Error())
			return
		}
		modelFromClass(&plan, c)
	}

	resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("uuid"), types.StringValue(plan.UUID.ValueString()))...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ClassResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ClassResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	c, err := r.service.GetClass(ctx, state.UUID.ValueString())
	if err != nil {
		if errors.Is(err, jamfschool.ErrNotFound) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading class", err.Error())
		return
	}

	modelFromClass(&state, c)
	resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("uuid"), types.StringValue(state.UUID.ValueString()))...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *ClassResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ClassResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state ClassResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	input := buildUpdateInput(&plan)

	if err := r.service.UpdateClass(ctx, state.UUID.ValueString(), input); err != nil {
		resp.Diagnostics.AddError("Error updating class", err.Error())
		return
	}

	// Update class members if changed.
	planStudents := extractInt64Set(plan.Students)
	planTeachers := extractInt64Set(plan.Teachers)
	stateStudents := extractInt64Set(state.Students)
	stateTeachers := extractInt64Set(state.Teachers)

	if !equalInt64Sets(planStudents, stateStudents) || !equalInt64Sets(planTeachers, stateTeachers) {
		if err := r.service.AssignClassUsers(ctx, state.UUID.ValueString(), planStudents, planTeachers); err != nil {
			resp.Diagnostics.AddError("Error updating class members", err.Error())
			return
		}
	}

	c, err := r.service.GetClass(ctx, state.UUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error reading updated class", err.Error())
		return
	}

	modelFromClass(&plan, c)
	resp.Diagnostics.Append(resp.Identity.SetAttribute(ctx, path.Root("uuid"), types.StringValue(plan.UUID.ValueString()))...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *ClassResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ClassResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if err := r.service.DeleteClass(ctx, state.UUID.ValueString()); err != nil {
		resp.Diagnostics.AddError("Error deleting class", err.Error())
	}
}

func (r *ClassResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("uuid"), req.ID)...)
}
