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

func TestAccDeviceGroupResource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-dg")
	resourceName := "jamfschool_device_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceGroupResourceConfig(rName, "Initial description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Initial description"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "members"),
					resource.TestCheckResourceAttrSet(resourceName, "is_smart_group"),
					resource.TestCheckResourceAttrSet(resourceName, "location_id"),
					resource.TestCheckResourceAttr(resourceName, "shared", "false"),
					resource.TestCheckResourceAttr(resourceName, "show_in_ios_app", "none"),
					resource.TestCheckResourceAttr(resourceName, "member_udids.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"show_in_ios_app", "member_udids"},
			},
			{
				Config: testAccDeviceGroupResourceConfig(rName, "Updated description"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated description"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "shared", "false"),
				),
			},
		},
	})
}

func testAccDeviceGroupResourceConfig(name, description string) string {
	return fmt.Sprintf(`
resource "jamfschool_device_group" "test" {
  name        = %[1]q
  description = %[2]q
}
`, name, description)
}
