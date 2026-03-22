// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package class_test

import (
	"fmt"
	"testing"

	"github.com/Jamf-Concepts/terraform-provider-jamfschool/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccClassDataSource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-class-ds")
	resourceName := "data.jamfschool_class.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccClassDataSourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Test class"),
					resource.TestCheckResourceAttrSet(resourceName, "uuid"),
					resource.TestCheckResourceAttrSet(resourceName, "source"),
					resource.TestCheckResourceAttrSet(resourceName, "student_count"),
					resource.TestCheckResourceAttrSet(resourceName, "teacher_count"),
					resource.TestCheckResourceAttrSet(resourceName, "location_id"),
					resource.TestCheckResourceAttr(resourceName, "device_count", "0"),
					resource.TestCheckResourceAttr(resourceName, "device_udids.#", "0"),
				),
			},
		},
	})
}

func testAccClassDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "jamfschool_class" "test" {
  name        = %[1]q
  description = "Test class"
}

data "jamfschool_class" "test" {
  uuid = jamfschool_class.test.uuid
}
`, name)
}
