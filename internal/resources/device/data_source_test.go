// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package device_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Jamf-Concepts/terraform-provider-jamfschool/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDeviceDataSource(t *testing.T) {
	udid := os.Getenv("JAMFSCHOOL_TEST_DEVICE_UDID")
	if udid == "" {
		t.Skip("JAMFSCHOOL_TEST_DEVICE_UDID must be set for device data source acceptance tests")
	}
	resourceName := "data.jamfschool_device.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccDeviceDataSourceConfig(udid),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "udid", udid),
					resource.TestCheckResourceAttrSet(resourceName, "serial_number"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "is_managed"),
					resource.TestCheckResourceAttrSet(resourceName, "is_supervised"),
					resource.TestCheckResourceAttrSet(resourceName, "model_name"),
					resource.TestCheckResourceAttrSet(resourceName, "os_version"),
					resource.TestCheckResourceAttrSet(resourceName, "os_prefix"),
				),
			},
		},
	})
}

func testAccDeviceDataSourceConfig(udid string) string {
	return fmt.Sprintf(`
data "jamfschool_device" "test" {
  udid = %[1]q
}
`, udid)
}
