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

var _ action.Action = (*RestartDeviceAction)(nil)
var _ action.ActionWithConfigure = (*RestartDeviceAction)(nil)

// RestartDeviceAction schedules a restart on a device.
type RestartDeviceAction struct {
	deviceAction
}

type RestartDeviceActionModel struct {
	UDID          types.String `tfsdk:"udid"`
	SerialNumber  types.String `tfsdk:"serial_number"`
	ClearPasscode types.Bool   `tfsdk:"clear_passcode"`
}

func NewRestartDeviceAction() action.Action {
	return &RestartDeviceAction{}
}

func (a *RestartDeviceAction) Metadata(_ context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_restart_device"
}

func (a *RestartDeviceAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = actionschema.Schema{
		MarkdownDescription: "Schedules a restart on a device.",
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
			"clear_passcode": actionschema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether to clear the passcode before restarting.",
			},
		},
	}
}

func (a *RestartDeviceAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.configure(ctx, req, resp)
}

func (a *RestartDeviceAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	if !a.ensureService(resp) {
		return
	}

	var data RestartDeviceActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	udid, ok := a.resolveDeviceUDID(ctx, resp, data.UDID, data.SerialNumber)
	if !ok {
		return
	}

	clearPass := !data.ClearPasscode.IsNull() && data.ClearPasscode.ValueBool()

	resp.SendProgress(action.InvokeProgressEvent{Message: fmt.Sprintf("Scheduling restart for device %s", udid)})

	if err := a.service.RestartDevice(ctx, udid, clearPass); err != nil {
		resp.Diagnostics.AddError("Restart Device Failed", fmt.Sprintf("Unable to restart device %s: %s", udid, err))
		return
	}

	resp.SendProgress(action.InvokeProgressEvent{Message: fmt.Sprintf("Restart scheduled for device %s", udid)})
}
