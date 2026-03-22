
data "jamfschool_profile" "example" {
  id = 10
}

output "profile_name" {
  value = data.jamfschool_profile.example.name
}
