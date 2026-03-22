// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package location_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Jamf-Concepts/terraform-provider-jamfschool/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLocationDataSource(t *testing.T) {
	locationID := os.Getenv("JAMFSCHOOL_TEST_LOCATION_ID")
	if locationID == "" {
		t.Skip("JAMFSCHOOL_TEST_LOCATION_ID must be set for location data source acceptance tests")
	}
	resourceName := "data.jamfschool_location.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccLocationDataSourceConfig(locationID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "is_district"),
				),
			},
		},
	})
}

func testAccLocationDataSourceConfig(id string) string {
	return fmt.Sprintf(`
data "jamfschool_location" "test" {
  id = %[1]s
}
`, id)
}
