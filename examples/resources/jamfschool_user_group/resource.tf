
# Create a user group with Jamf School Teacher enabled.
resource "jamfschool_user_group" "example" {
  name                = "Teaching Staff"
  description         = "All teaching staff members"
  jamf_school_teacher = "allow"
  jamf_parent         = "deny"
}
