// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package deviceactions

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
)

// deviceAction shares Configure logic across device action implementations.
type deviceAction struct {
	service *jamfschool.Client
}

// configure binds the provider-supplied service to the action.
func (a *deviceAction) configure(_ context.Context, req action.ConfigureRequest, resp *action.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	svc, ok := req.ProviderData.(*jamfschool.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			fmt.Sprintf("Expected *jamfschool.Client, got %T.", req.ProviderData),
		)
		return
	}

	a.service = svc
}

// ensureService guarantees Configure completed successfully before Invoke.
func (a *deviceAction) ensureService(resp *action.InvokeResponse) bool {
	if a.service != nil {
		return true
	}

	resp.Diagnostics.AddError(
		"Provider Not Configured",
		"The Jamf School service was not configured. Re-run terraform init/apply so the provider can configure successfully.",
	)
	return false
}

// resolveDeviceUDID ensures exactly one device identifier is provided and returns the UDID.
func (a *deviceAction) resolveDeviceUDID(ctx context.Context, resp *action.InvokeResponse, udidAttr, serialNumberAttr types.String) (string, bool) {
	hasUDID := !udidAttr.IsNull() && !udidAttr.IsUnknown() && udidAttr.ValueString() != ""
	hasSerial := !serialNumberAttr.IsNull() && !serialNumberAttr.IsUnknown() && serialNumberAttr.ValueString() != ""

	switch {
	case hasUDID && hasSerial:
		resp.Diagnostics.AddError(
			"Multiple Device Identifiers Provided",
			"Specify only one of udid or serial_number when invoking this action.",
		)
		return "", false
	case hasUDID:
		return udidAttr.ValueString(), true
	case hasSerial:
		serial := serialNumberAttr.ValueString()
		resp.SendProgress(action.InvokeProgressEvent{Message: fmt.Sprintf("Resolving serial number %s to UDID", serial)})

		devices, err := a.service.GetDevices(ctx)
		if err != nil {
			resp.Diagnostics.AddError(
				"Device Lookup Failed",
				fmt.Sprintf("Unable to list devices to resolve serial number %s: %s", serial, err),
			)
			return "", false
		}

		for _, d := range devices {
			if d.SerialNumber == serial {
				return d.UDID, true
			}
		}

		resp.Diagnostics.AddError(
			"Device Not Found",
			fmt.Sprintf("No device found with serial number %s.", serial),
		)
		return "", false
	default:
		resp.Diagnostics.AddError(
			"Missing Device Identifier",
			"Specify either udid or serial_number to select the device.",
		)
		return "", false
	}
}
