// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/Jamf-Concepts/jamfschool-go-sdk/jamfschool"
)

// TestAccProtoV6ProviderFactories returns provider factories for use in
// acceptance tests across all resource packages.
func TestAccProtoV6ProviderFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"jamfschool": providerserver.NewProtocol6WithError(New("test")()),
	}
}

// TestAccPreCheck validates that the required environment variables are set
// for acceptance testing.
func TestAccPreCheck(t *testing.T) {
	t.Helper()

	if os.Getenv("JAMFSCHOOL_URL") == "" {
		t.Fatal("JAMFSCHOOL_URL must be set for acceptance tests")
	}
	if os.Getenv("JAMFSCHOOL_NETWORK_ID") == "" {
		t.Fatal("JAMFSCHOOL_NETWORK_ID must be set for acceptance tests")
	}
	if os.Getenv("JAMFSCHOOL_API_KEY") == "" {
		t.Fatal("JAMFSCHOOL_API_KEY must be set for acceptance tests")
	}
}

// TestAccService returns a Client configured from environment
// variables, for use in CheckDestroy functions.
func TestAccService() *jamfschool.Client {
	return jamfschool.NewClient(
		os.Getenv("JAMFSCHOOL_URL"),
		os.Getenv("JAMFSCHOOL_NETWORK_ID"),
		os.Getenv("JAMFSCHOOL_API_KEY"),
	)
}
