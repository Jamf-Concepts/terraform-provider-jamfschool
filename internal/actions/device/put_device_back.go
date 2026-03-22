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

var _ action.Action = (*PutDeviceBackAction)(nil)
var _ action.ActionWithConfigure = (*PutDeviceBackAction)(nil)

// PutDeviceBackAction restores a device from trash.
type PutDeviceBackAction struct {
	deviceAction
}

type PutDeviceBackActionModel struct {
	UDID         types.String `tfsdk:"udid"`
	SerialNumber types.String `tfsdk:"serial_number"`
}

func NewPutDeviceBackAction() action.Action {
	return &PutDeviceBackAction{}
}

func (a *PutDeviceBackAction) Metadata(_ context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_put_device_back"
}

func (a *PutDeviceBackAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = actionschema.Schema{
		MarkdownDescription: "Restores a device from trash.",
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

func (a *PutDeviceBackAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.configure(ctx, req, resp)
}

func (a *PutDeviceBackAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	if !a.ensureService(resp) {
		return
	}

	var data PutDeviceBackActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	udid, ok := a.resolveDeviceUDID(ctx, resp, data.UDID, data.SerialNumber)
	if !ok {
		return
	}

	resp.SendProgress(action.InvokeProgressEvent{Message: fmt.Sprintf("Restoring device %s from trash", udid)})

	if err := a.service.RestoreDevice(ctx, udid); err != nil {
		resp.Diagnostics.AddError("Put Device Back Failed", fmt.Sprintf("Unable to restore device %s: %s", udid, err))
		return
	}

	resp.SendProgress(action.InvokeProgressEvent{Message: fmt.Sprintf("Device %s restored from trash", udid)})
}
