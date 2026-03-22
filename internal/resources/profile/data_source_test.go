// Copyright Jamf Software LLC 2026
// SPDX-License-Identifier: MPL-2.0

package profile_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Jamf-Concepts/terraform-provider-jamfschool/internal/provider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccProfileDataSource(t *testing.T) {
	profileID := os.Getenv("JAMFSCHOOL_TEST_PROFILE_ID")
	if profileID == "" {
		t.Skip("JAMFSCHOOL_TEST_PROFILE_ID must be set for profile data source acceptance tests")
	}
	resourceName := "data.jamfschool_profile.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: provider.TestAccProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccProfileDataSourceConfig(profileID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "name"),
					resource.TestCheckResourceAttrSet(resourceName, "identifier"),
					resource.TestCheckResourceAttrSet(resourceName, "platform"),
				),
			},
		},
	})
}

func testAccProfileDataSourceConfig(id string) string {
	return fmt.Sprintf(`
data "jamfschool_profile" "test" {
  id = %[1]s
}
`, id)
}
