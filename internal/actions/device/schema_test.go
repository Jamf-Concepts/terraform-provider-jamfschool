// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package deviceactions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/action"
)

func TestEraseDeviceAction_Metadata(t *testing.T) {
	t.Parallel()
	a := NewEraseDeviceAction()
	req := action.MetadataRequest{ProviderTypeName: "jamfschool"}
	var resp action.MetadataResponse
	a.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "jamfschool_erase_device" {
		t.Errorf("expected type name jamfschool_erase_device, got %s", resp.TypeName)
	}
}

func TestEraseDeviceAction_Schema(t *testing.T) {
	t.Parallel()
	a := NewEraseDeviceAction()
	req := action.SchemaRequest{}
	var resp action.SchemaResponse
	a.Schema(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected schema errors: %v", resp.Diagnostics)
	}
	for _, attr := range []string{"udid", "serial_number", "clear_activation_lock"} {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Errorf("missing expected attribute %q", attr)
		}
	}
}

func TestRestartDeviceAction_Metadata(t *testing.T) {
	t.Parallel()
	a := NewRestartDeviceAction()
	req := action.MetadataRequest{ProviderTypeName: "jamfschool"}
	var resp action.MetadataResponse
	a.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "jamfschool_restart_device" {
		t.Errorf("expected type name jamfschool_restart_device, got %s", resp.TypeName)
	}
}

func TestRestartDeviceAction_Schema(t *testing.T) {
	t.Parallel()
	a := NewRestartDeviceAction()
	req := action.SchemaRequest{}
	var resp action.SchemaResponse
	a.Schema(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected schema errors: %v", resp.Diagnostics)
	}
	for _, attr := range []string{"udid", "serial_number", "clear_passcode"} {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Errorf("missing expected attribute %q", attr)
		}
	}
}

func TestRefreshDeviceAction_Metadata(t *testing.T) {
	t.Parallel()
	a := NewRefreshDeviceAction()
	req := action.MetadataRequest{ProviderTypeName: "jamfschool"}
	var resp action.MetadataResponse
	a.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "jamfschool_refresh_device" {
		t.Errorf("expected type name jamfschool_refresh_device, got %s", resp.TypeName)
	}
}

func TestUnenrollDeviceAction_Metadata(t *testing.T) {
	t.Parallel()
	a := NewUnenrollDeviceAction()
	req := action.MetadataRequest{ProviderTypeName: "jamfschool"}
	var resp action.MetadataResponse
	a.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "jamfschool_unenroll_device" {
		t.Errorf("expected type name jamfschool_unenroll_device, got %s", resp.TypeName)
	}
}

func TestClearDeviceActivationLockAction_Metadata(t *testing.T) {
	t.Parallel()
	a := NewClearDeviceActivationLockAction()
	req := action.MetadataRequest{ProviderTypeName: "jamfschool"}
	var resp action.MetadataResponse
	a.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "jamfschool_clear_device_activation_lock" {
		t.Errorf("expected type name jamfschool_clear_device_activation_lock, got %s", resp.TypeName)
	}
}

func TestMoveDeviceToTrashAction_Metadata(t *testing.T) {
	t.Parallel()
	a := NewMoveDeviceToTrashAction()
	req := action.MetadataRequest{ProviderTypeName: "jamfschool"}
	var resp action.MetadataResponse
	a.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "jamfschool_move_device_to_trash" {
		t.Errorf("expected type name jamfschool_move_device_to_trash, got %s", resp.TypeName)
	}
}

func TestPutDeviceBackAction_Metadata(t *testing.T) {
	t.Parallel()
	a := NewPutDeviceBackAction()
	req := action.MetadataRequest{ProviderTypeName: "jamfschool"}
	var resp action.MetadataResponse
	a.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "jamfschool_put_device_back" {
		t.Errorf("expected type name jamfschool_put_device_back, got %s", resp.TypeName)
	}
}

func TestUpdateDeviceESIMAction_Metadata(t *testing.T) {
	t.Parallel()
	a := NewUpdateDeviceESIMAction()
	req := action.MetadataRequest{ProviderTypeName: "jamfschool"}
	var resp action.MetadataResponse
	a.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "jamfschool_update_device_esim" {
		t.Errorf("expected type name jamfschool_update_device_esim, got %s", resp.TypeName)
	}
}

func TestUpdateDeviceESIMAction_Schema(t *testing.T) {
	t.Parallel()
	a := NewUpdateDeviceESIMAction()
	req := action.SchemaRequest{}
	var resp action.SchemaResponse
	a.Schema(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected schema errors: %v", resp.Diagnostics)
	}
	for _, attr := range []string{"udid", "serial_number", "server_url", "requires_network_tether"} {
		if _, ok := resp.Schema.Attributes[attr]; !ok {
			t.Errorf("missing expected attribute %q", attr)
		}
	}
}
