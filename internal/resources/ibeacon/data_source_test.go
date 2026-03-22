// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package ibeacon_test

import (
	"fmt"
	"testing"

	"github.com/Jamf-Concepts/terraform-provider-jamfschool/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIBeaconDataSource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-ibeacon-ds")
	resourceName := "data.jamfschool_ibeacon.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccIBeaconDataSourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Test iBeacon"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttrSet(resourceName, "major"),
					resource.TestCheckResourceAttrSet(resourceName, "minor"),
				),
			},
		},
	})
}

func testAccIBeaconDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "jamfschool_ibeacon" "test" {
  name        = %[1]q
  description = "Test iBeacon"
}

data "jamfschool_ibeacon" "test" {
  id = jamfschool_ibeacon.test.id
}
`, name)
}
