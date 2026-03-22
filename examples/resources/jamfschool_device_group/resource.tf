
# Create a shared device group shown as an article in the iOS app.
resource "jamfschool_device_group" "example" {
  name            = "Shared iPads"
  description     = "Shared iPad cart devices"
  information     = "These devices are available for classroom use"
  shared          = true
  show_in_ios_app = "article"
}
