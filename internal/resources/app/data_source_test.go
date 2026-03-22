// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package app_test

import (
	"testing"

	"github.com/Jamf-Concepts/terraform-provider-jamfschool/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccAppDataSource(t *testing.T) {
	resourceName := "data.jamfschool_app.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
resource "jamfschool_app" "test" {
  adam_id      = 409183694
  country_code = "us"
}

data "jamfschool_app" "test" {
  id = jamfschool_app.test.id
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "bundle_id"),
					resource.TestCheckResourceAttrSet(resourceName, "platform"),
					resource.TestCheckResourceAttrSet(resourceName, "vendor"),
					resource.TestCheckResourceAttrSet(resourceName, "adam_id"),
				),
			},
		},
	})
}
