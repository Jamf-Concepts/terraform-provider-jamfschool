// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package device_group_test

import (
	"fmt"
	"testing"

	"github.com/Jamf-Concepts/terraform-provider-jamfschool/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceGroupDataSource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-dg-ds")
	resourceName := "data.jamfschool_device_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceGroupDataSourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Test device group"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "members"),
					resource.TestCheckResourceAttrSet(resourceName, "is_smart_group"),
					resource.TestCheckResourceAttrSet(resourceName, "location_id"),
					resource.TestCheckResourceAttrSet(resourceName, "shared"),
					resource.TestCheckResourceAttr(resourceName, "information", ""),
				),
			},
		},
	})
}

func testAccDeviceGroupDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "jamfschool_device_group" "test" {
  name        = %[1]q
  description = "Test device group"
}

data "jamfschool_device_group" "test" {
  id = jamfschool_device_group.test.id
}
`, name)
}
