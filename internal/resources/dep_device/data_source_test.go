// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package dep_device_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Jamf-Concepts/terraform-provider-jamfschool/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDEPDeviceDataSource(t *testing.T) {
	serialNumber := os.Getenv("JAMFSCHOOL_TEST_DEP_SERIAL")
	if serialNumber == "" {
		t.Skip("JAMFSCHOOL_TEST_DEP_SERIAL must be set for DEP device data source acceptance tests")
	}
	resourceName := "data.jamfschool_dep_device.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDEPDeviceDataSourceConfig(serialNumber),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "serial_number", serialNumber),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
		},
	})
}

func testAccDEPDeviceDataSourceConfig(serial string) string {
	return fmt.Sprintf(`
data "jamfschool_dep_device" "test" {
  serial_number = %[1]q
}
`, serial)
}
