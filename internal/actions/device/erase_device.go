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

var _ action.Action = (*EraseDeviceAction)(nil)
var _ action.ActionWithConfigure = (*EraseDeviceAction)(nil)

// EraseDeviceAction schedules a wipe on a device.
type EraseDeviceAction struct {
	deviceAction
}

type EraseDeviceActionModel struct {
	UDID                types.String `tfsdk:"udid"`
	SerialNumber        types.String `tfsdk:"serial_number"`
	ClearActivationLock types.Bool   `tfsdk:"clear_activation_lock"`
}

func NewEraseDeviceAction() action.Action {
	return &EraseDeviceAction{}
}

func (a *EraseDeviceAction) Metadata(_ context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_erase_device"
}

func (a *EraseDeviceAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = actionschema.Schema{
		MarkdownDescription: "Schedules a wipe on a device, erasing all content and settings.",
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
			"clear_activation_lock": actionschema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to clear the activation lock before wiping.",
			},
		},
	}
}

func (a *EraseDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.configure(ctx, req, resp)
}

func (a *EraseDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	if !a.ensureService(resp) {
		return
	}

	var data EraseDeviceActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	udid, ok := a.resolveDeviceUDID(ctx, resp, data.UDID, data.SerialNumber)
	if !ok {
		return
	}

	clearLock := !data.ClearActivationLock.IsNull() && data.ClearActivationLock.ValueBool()

	resp.SendProgress(action.InvokeProgressEvent{Message: fmt.Sprintf("Scheduling erase for device %s", udid)})

	if err := a.service.EraseDevice(ctx, udid, clearLock); err != nil {
		resp.Diagnostics.AddError("Erase Device Failed", fmt.Sprintf("Unable to erase device %s: %s", udid, err))
		return
	}

	resp.SendProgress(action.InvokeProgressEvent{Message: fmt.Sprintf("Erase scheduled for device %s", udid)})
}
