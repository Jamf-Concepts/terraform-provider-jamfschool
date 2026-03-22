
data "jamfschool_ibeacon" "example" {
  id = 42
}

output "ibeacon_uuid" {
  value = data.jamfschool_ibeacon.example.uuid
}
