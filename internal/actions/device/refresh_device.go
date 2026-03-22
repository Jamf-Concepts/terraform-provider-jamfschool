// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package deviceactions

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	actionschema "github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ action.Action = (*RefreshDeviceAction)(nil)
var _ action.ActionWithConfigure = (*RefreshDeviceAction)(nil)

// RefreshDeviceAction schedules a full inventory refresh for a device.
type RefreshDeviceAction struct {
	deviceAction
}

type RefreshDeviceActionModel struct {
	UDID         types.String `tfsdk:"udid"`
	SerialNumber types.String `tfsdk:"serial_number"`
	ClearErrors  types.Bool   `tfsdk:"clear_errors"`
}

func NewRefreshDeviceAction() action.Action {
	return &RefreshDeviceAction{}
}

func (a *RefreshDeviceAction) Metadata(_ context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_refresh_device"
}

func (a *RefreshDeviceAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = actionschema.Schema{
		MarkdownDescription: "Schedules a full inventory refresh for a device.",
		Attributes: map[string]actionschema.Attribute{
			"udid": actionschema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Device UDID. Provide this or `serial_number`.",
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("serial_number")),
				},
			},
			"serial_number": actionschema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "Device serial number. Provide this or `udid`.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("udid")),
				},
			},
			"clear_errors": actionschema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Clear all installation failures so apps get reinstalled when previously failed.",
			},
		},
	}
}

func (a *RefreshDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.configure(ctx, req, resp)
}

func (a *RefreshDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	if !a.ensureService(resp) {
		return
	}

	var data RefreshDeviceActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	udid, ok := a.resolveDeviceUDID(ctx, resp, data.UDID, data.SerialNumber)
	if !ok {
		return
	}

	clearErrs := !data.ClearErrors.IsNull() && data.ClearErrors.ValueBool()

	resp.SendProgress(action.InvokeProgressEvent{Message: fmt.Sprintf("Scheduling inventory refresh for device %s", udid)})

	if err := a.service.RefreshDevice(ctx, udid, clearErrs); err != nil {
		resp.Diagnostics.AddError("Refresh Device Failed", fmt.Sprintf("Unable to refresh device %s: %s", udid, err))
		return
	}

	resp.SendProgress(action.InvokeProgressEvent{Message: fmt.Sprintf("Inventory refresh scheduled for device %s", udid)})
}
