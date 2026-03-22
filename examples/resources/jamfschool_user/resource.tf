
# Create a Jamf School user with all configurable attributes.
resource "jamfschool_user" "example" {
  username                       = "jappleseed"
  password                       = "SecurePassword123!"
  email                          = "jappleseed@example.com"
  first_name                     = "Johnny"
  last_name                      = "Appleseed"
  notes                          = "Created via Terraform"
  exclude                        = false
  store_mail_contacts_calendars  = true
  mail_contacts_calendars_domain = "example.com"
  member_of                      = [jamfschool_user_group.teachers.id]
}
