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

func TestAccUserGroupResource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-ug")
	resourceName := "jamfschool_user_group.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccUserGroupResourceConfig(rName, "Initial description", "deny", "inherit"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Initial description"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "user_count"),
					resource.TestCheckResourceAttr(resourceName, "jamf_school_teacher", "deny"),
					resource.TestCheckResourceAttr(resourceName, "jamf_parent", "inherit"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccUserGroupResourceConfig(rName, "Updated description", "allow", "deny"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", rName),
					resource.TestCheckResourceAttr(resourceName, "description", "Updated description"),
					resource.TestCheckResourceAttr(resourceName, "jamf_school_teacher", "allow"),
					resource.TestCheckResourceAttr(resourceName, "jamf_parent", "deny"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
		},
	})
}

func testAccUserGroupResourceConfig(name, description, teacher, parent string) string {
	return fmt.Sprintf(`
resource "jamfschool_user_group" "test" {
  name                = %[1]q
  description         = %[2]q
  jamf_school_teacher = %[3]q
  jamf_parent         = %[4]q
}
`, name, description, teacher, parent)
}
