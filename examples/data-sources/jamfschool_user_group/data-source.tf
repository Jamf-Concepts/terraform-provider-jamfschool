
data "jamfschool_user_group" "example" {
  id = 456
}

output "group_name" {
  value = data.jamfschool_user_group.example.name
}
