// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package ibeacon_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/Jamf-Concepts/terraform-provider-jamfschool/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccIBeaconResource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-ibeacon")
	resourceName := "jamfschool_ibeacon.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories(),
		CheckDestroy:             testAccCheckIBeaconDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIBeaconResourceConfig(rName, "Initial beacon", 100, 200),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Initial beacon"),
					resource.TestCheckResourceAttr(resourceName, "major", "100"),
					resource.TestCheckResourceAttr(resourceName, "minor", "200"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccIBeaconResourceConfig(rName, "Updated beacon", 300, 400),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated beacon"),
					resource.TestCheckResourceAttr(resourceName, "major", "300"),
					resource.TestCheckResourceAttr(resourceName, "minor", "400"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
				),
			},
		},
	})
}

func testAccCheckIBeaconDestroy(s *terraform.State) error {
	svc := provider.TestAccService()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "jamfschool_ibeacon" {
			continue
		}

		id, err := strconv.ParseInt(rs.Primary.ID, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid ID: %s", rs.Primary.ID)
		}

		_, err = svc.GetIBeacon(context.Background(), id)
		if err == nil {
			return fmt.Errorf("ibeacon %d still exists after destroy", id)
		}
	}
	return nil
}

func testAccIBeaconResourceConfig(name, description string, major, minor int) string {
	return fmt.Sprintf(`
resource "jamfschool_ibeacon" "test" {
  name        = %[1]q
  description = %[2]q
  major       = %[3]d
  minor       = %[4]d
}
`, name, description, major, minor)
}
