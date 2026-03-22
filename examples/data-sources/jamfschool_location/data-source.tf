
data "jamfschool_location" "example" {
  id = 1
}

output "location_name" {
  value = data.jamfschool_location.example.name
}
