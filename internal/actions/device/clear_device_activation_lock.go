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

var _ action.Action = (*ClearDeviceActivationLockAction)(nil)
var _ action.ActionWithConfigure = (*ClearDeviceActivationLockAction)(nil)

// ClearDeviceActivationLockAction clears the activation lock from a device.
type ClearDeviceActivationLockAction struct {
	deviceAction
}

type ClearDeviceActivationLockActionModel struct {
	UDID         types.String `tfsdk:"udid"`
	SerialNumber types.String `tfsdk:"serial_number"`
}

func NewClearDeviceActivationLockAction() action.Action {
	return &ClearDeviceActivationLockAction{}
}

func (a *ClearDeviceActivationLockAction) Metadata(_ context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_clear_device_activation_lock"
}

func (a *ClearDeviceActivationLockAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = actionschema.Schema{
		MarkdownDescription: "Clears the activation lock from a device.",
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
		},
	}
}

func (a *ClearDeviceActivationLockAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.configure(ctx, req, resp)
}

func (a *ClearDeviceActivationLockAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	if !a.ensureService(resp) {
		return
	}

	var data ClearDeviceActivationLockActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	udid, ok := a.resolveDeviceUDID(ctx, resp, data.UDID, data.SerialNumber)
	if !ok {
		return
	}

	resp.SendProgress(action.InvokeProgressEvent{Message: fmt.Sprintf("Clearing activation lock for device %s", udid)})

	if err := a.service.ClearDeviceActivationLock(ctx, udid); err != nil {
		resp.Diagnostics.AddError("Clear Activation Lock Failed", fmt.Sprintf("Unable to clear activation lock for device %s: %s", udid, err))
		return
	}

	resp.SendProgress(action.InvokeProgressEvent{Message: fmt.Sprintf("Activation lock cleared for device %s", udid)})
}
