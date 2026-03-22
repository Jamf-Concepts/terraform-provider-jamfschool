
data "jamfschool_device_group" "example" {
  id = 789
}

output "device_group_name" {
  value = data.jamfschool_device_group.example.name
}
