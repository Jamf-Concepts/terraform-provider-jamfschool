// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
)

func TestProviderMetadata(t *testing.T) {
	t.Parallel()

	p := New("1.0.0")()
	resp := &provider.MetadataResponse{}
	p.Metadata(context.Background(), provider.MetadataRequest{}, resp)

	if resp.TypeName != "jamfschool" {
		t.Errorf("expected type name %q, got %q", "jamfschool", resp.TypeName)
	}
	if resp.Version != "1.0.0" {
		t.Errorf("expected version %q, got %q", "1.0.0", resp.Version)
	}
}

func TestProviderSchema(t *testing.T) {
	t.Parallel()

	p := New("test")()
	resp := &provider.SchemaResponse{}
	p.Schema(context.Background(), provider.SchemaRequest{}, resp)

	s := resp.Schema

	requiredAttrs := []string{"url", "network_id", "api_key"}
	for _, attr := range requiredAttrs {
		a, ok := s.Attributes[attr]
		if !ok {
			t.Errorf("expected attribute %q in schema", attr)
			continue
		}
		if !a.(schema.StringAttribute).Optional {
			t.Errorf("expected attribute %q to be optional", attr)
		}
	}

	// api_key should be sensitive
	apiKey := s.Attributes["api_key"].(schema.StringAttribute)
	if !apiKey.Sensitive {
		t.Error("expected api_key to be sensitive")
	}
}

func TestProviderResources(t *testing.T) {
	t.Parallel()

	p := New("test")()
	jp, ok := p.(*JamfSchoolProvider)
	if !ok {
		t.Fatal("expected provider to be *JamfSchoolProvider")
	}
	resList := jp.Resources(context.Background())

	expectedResources := []string{
		"jamfschool_user",
		"jamfschool_user_group",
		"jamfschool_device_group",
		"jamfschool_class",
		"jamfschool_ibeacon",
		"jamfschool_app",
	}
	if len(resList) != len(expectedResources) {
		t.Errorf("expected %d resources, got %d", len(expectedResources), len(resList))
	}
}

func TestProviderDataSources(t *testing.T) {
	t.Parallel()

	p := New("test")()
	jp, ok := p.(*JamfSchoolProvider)
	if !ok {
		t.Fatal("expected provider to be *JamfSchoolProvider")
	}
	dsList := jp.DataSources(context.Background())

	expectedDataSources := []string{
		"jamfschool_user",
		"jamfschool_user_group",
		"jamfschool_device_group",
		"jamfschool_class",
		"jamfschool_ibeacon",
		"jamfschool_device",
		"jamfschool_app",
		"jamfschool_profile",
		"jamfschool_location",
		"jamfschool_dep_device",
	}
	if len(dsList) != len(expectedDataSources) {
		t.Errorf("expected %d data sources, got %d", len(expectedDataSources), len(dsList))
	}
}

func TestProviderActions(t *testing.T) {
	t.Parallel()

	p := New("test")()
	jp, ok := p.(*JamfSchoolProvider)
	if !ok {
		t.Fatal("expected provider to be *JamfSchoolProvider")
	}
	actionsList := jp.Actions(context.Background())

	expectedCount := 8
	if len(actionsList) != expectedCount {
		t.Errorf("expected %d actions, got %d", expectedCount, len(actionsList))
	}
}

func TestProviderListResources(t *testing.T) {
	t.Parallel()

	p := New("test")()
	jp, ok := p.(*JamfSchoolProvider)
	if !ok {
		t.Fatal("expected provider to be *JamfSchoolProvider")
	}
	listResList := jp.ListResources(context.Background())

	expectedCount := 6
	if len(listResList) != expectedCount {
		t.Errorf("expected %d list resources, got %d", expectedCount, len(listResList))
	}
}
