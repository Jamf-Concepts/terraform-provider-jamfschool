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

func TestAccUserResource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-user")
	resourceName := "jamfschool_user.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccUserResourceConfig(rName, "initial@example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "username", rName),
					resource.TestCheckResourceAttr(resourceName, "email", "initial@example.com"),
					resource.TestCheckResourceAttr(resourceName, "first_name", "Test"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "User"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
					resource.TestCheckResourceAttrSet(resourceName, "device_count"),
					resource.TestCheckResourceAttr(resourceName, "exclude", "false"),
					resource.TestCheckResourceAttr(resourceName, "store_mail_contacts_calendars", "false"),
					resource.TestCheckResourceAttr(resourceName, "move_devices_on_location_change", "true"),
					resource.TestCheckResourceAttr(resourceName, "member_of.#", "0"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "store_mail_contacts_calendars", "mail_contacts_calendars_domain", "move_devices_on_location_change"},
			},
			{
				Config: testAccUserResourceConfig(rName, "updated@example.com"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "username", rName),
					resource.TestCheckResourceAttr(resourceName, "email", "updated@example.com"),
					resource.TestCheckResourceAttr(resourceName, "first_name", "Test"),
					resource.TestCheckResourceAttr(resourceName, "last_name", "User"),
					resource.TestCheckResourceAttr(resourceName, "exclude", "false"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "status"),
				),
			},
		},
	})
}

func TestAccUserResource_MemberOf(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-user-mo")
	groupName := acctest.RandomWithPrefix("tf-acc-grp-mo")
	groupName2 := acctest.RandomWithPrefix("tf-acc-grp-mo2")
	resourceName := "jamfschool_user.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			// Step 1: Create user with one group membership
			{
				Config: testAccUserResourceMemberOfConfig(rName, groupName, groupName2, true, false),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "username", rName),
					resource.TestCheckResourceAttr(resourceName, "member_of.#", "1"),
				),
			},
			// Step 2: Import
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"password", "store_mail_contacts_calendars", "mail_contacts_calendars_domain", "move_devices_on_location_change"},
			},
			// Step 3: Update to add second group
			{
				Config: testAccUserResourceMemberOfConfig(rName, groupName, groupName2, true, true),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "username", rName),
					resource.TestCheckResourceAttr(resourceName, "member_of.#", "2"),
				),
			},
			// Step 4: Update to remove all groups
			{
				Config: testAccUserResourceMemberOfEmptyConfig(rName, groupName, groupName2),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "username", rName),
					resource.TestCheckResourceAttr(resourceName, "member_of.#", "0"),
				),
			},
		},
	})
}

func testAccUserResourceConfig(username, email string) string {
	return fmt.Sprintf(`
resource "jamfschool_user" "test" {
  username   = %[1]q
  password   = "TestPass123!"
  email      = %[2]q
  first_name = "Test"
  last_name  = "User"
}
`, username, email)
}

func testAccUserResourceMemberOfConfig(username, groupName, groupName2 string, group1, group2 bool) string {
	memberOf := "member_of = [jamfschool_user_group.group1.id]"
	if group1 && group2 {
		memberOf = "member_of = [jamfschool_user_group.group1.id, jamfschool_user_group.group2.id]"
	}
	return fmt.Sprintf(`
resource "jamfschool_user_group" "group1" {
  name = %[2]q
}

resource "jamfschool_user_group" "group2" {
  name = %[3]q
}

resource "jamfschool_user" "test" {
  username   = %[1]q
  password   = "TestPass123!"
  email      = "%[1]s@example.com"
  first_name = "Test"
  last_name  = "User"
  %[4]s
}
`, username, groupName, groupName2, memberOf)
}

func testAccUserResourceMemberOfEmptyConfig(username, groupName, groupName2 string) string {
	return fmt.Sprintf(`
resource "jamfschool_user_group" "group1" {
  name = %[2]q
}

resource "jamfschool_user_group" "group2" {
  name = %[3]q
}

resource "jamfschool_user" "test" {
  username   = %[1]q
  password   = "TestPass123!"
  email      = "%[1]s@example.com"
  first_name = "Test"
  last_name  = "User"
  member_of  = []
}
`, username, groupName, groupName2)
}
