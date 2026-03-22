// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package user_group_test

import (
	"fmt"
	"testing"

	"github.com/Jamf-Concepts/terraform-provider-jamfschool/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserGroupDataSource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-ug-ds")
	resourceName := "data.jamfschool_user_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroupDataSourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Test group"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "user_count"),
					resource.TestCheckResourceAttrSet(resourceName, "jamf_school_teacher"),
					resource.TestCheckResourceAttrSet(resourceName, "jamf_parent"),
				),
			},
		},
	})
}

func testAccUserGroupDataSourceConfig(name string) string {
	return fmt.Sprintf(`
resource "jamfschool_user_group" "test" {
  name        = %[1]q
  description = "Test group"
}

data "jamfschool_user_group" "test" {
  id = jamfschool_user_group.test.id
}
`, name)
}
