// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package user

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func TestUserListResource_Metadata(t *testing.T) {
	t.Parallel()
	r := NewUserListResource()
	req := resource.MetadataRequest{ProviderTypeName: "jamfschool"}
	var resp resource.MetadataResponse
	r.Metadata(context.Background(), req, &resp)
	if resp.TypeName != "jamfschool_user" {
		t.Errorf("expected type name jamfschool_user, got %s", resp.TypeName)
	}
}

func TestUserListResource_Schema(t *testing.T) {
	t.Parallel()
	r := NewUserListResource()
	req := list.ListResourceSchemaRequest{}
	var resp list.ListResourceSchemaResponse
	r.ListResourceConfigSchema(context.Background(), req, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected schema errors: %v", resp.Diagnostics)
	}
	if _, ok := resp.Schema.Attributes["name_prefix"]; !ok {
		t.Error("missing expected attribute name_prefix")
	}
}
