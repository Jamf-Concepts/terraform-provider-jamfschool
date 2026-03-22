// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package user_test

import (
	"fmt"
	"testing"

	"github.com/Jamf-Concepts/terraform-provider-jamfschool/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccUserDataSource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-user-ds")
	resourceName := "data.jamfschool_user.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccUserDataSourceConfig(rName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "username", rName),
					resource.TestCheckResourceAttr(resourceName, "email", rName+"@example.com"),
					resource.TestCheckResourceAttr(resourceName, "first_name", "Test"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "User"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "device_count"),
					resource.TestCheckResourceAttrSet(resourceName, "exclude"),
					resource.TestCheckResourceAttr(resourceName, "member_of.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "group_names.#", "1"),
				),
			},
		},
	})
}

func testAccUserDataSourceConfig(username string) string {
	return fmt.Sprintf(`
resource "jamfschool_user_group" "test" {
  name = "%[1]s-group"
}

resource "jamfschool_user" "test" {
  username   = %[1]q
  password   = "TestPass123!"
  email      = "%[1]s@example.com"
  first_name = "Test"
  last_name  = "User"
  member_of  = [jamfschool_user_group.test.id]
}

data "jamfschool_user" "test" {
  id = jamfschool_user.test.id
}
`, username)
}
