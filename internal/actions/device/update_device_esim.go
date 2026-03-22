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

var _ action.Action = (*UpdateDeviceESIMAction)(nil)
var _ action.ActionWithConfigure = (*UpdateDeviceESIMAction)(nil)

// UpdateDeviceESIMAction queries a carrier URL for active eSIM cellular-plan profiles on a device.
type UpdateDeviceESIMAction struct {
	deviceAction
}

type UpdateDeviceESIMActionModel struct {
	UDID                  types.String `tfsdk:"udid"`
	SerialNumber          types.String `tfsdk:"serial_number"`
	ServerURL             types.String `tfsdk:"server_url"`
	RequiresNetworkTether types.Bool   `tfsdk:"requires_network_tether"`
}

func NewUpdateDeviceESIMAction() action.Action {
	return &UpdateDeviceESIMAction{}
}

func (a *UpdateDeviceESIMAction) Metadata(_ context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_update_device_esim"
}

func (a *UpdateDeviceESIMAction) Schema(_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = actionschema.Schema{
		MarkdownDescription: "Queries a carrier URL to refresh eSIM cellular-plan profiles on a device.",
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
			"server_url": actionschema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The carrier's eSIM server URL to query.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"requires_network_tether": actionschema.BoolAttribute{
				Optional:            true,
				MarkdownDescription: "Whether network tethering is required for executing this command.",
			},
		},
	}
}

func (a *UpdateDeviceESIMAction) Configure(ctx context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	a.configure(ctx, req, resp)
}

func (a *UpdateDeviceESIMAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	if !a.ensureService(resp) {
		return
	}

	var data UpdateDeviceESIMActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	udid, ok := a.resolveDeviceUDID(ctx, resp, data.UDID, data.SerialNumber)
	if !ok {
		return
	}

	tether := !data.RequiresNetworkTether.IsNull() && data.RequiresNetworkTether.ValueBool()

	resp.SendProgress(action.InvokeProgressEvent{Message: fmt.Sprintf("Updating eSIM cellular plan for device %s", udid)})

	if err := a.service.UpdateDeviceESIM(ctx, udid, data.ServerURL.ValueString(), tether); err != nil {
		resp.Diagnostics.AddError("Update Device eSIM Failed", fmt.Sprintf("Unable to update eSIM for device %s: %s", udid, err))
		return
	}

	resp.SendProgress(action.InvokeProgressEvent{Message: fmt.Sprintf("eSIM cellular plan refresh scheduled for device %s", udid)})
}
