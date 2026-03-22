
data "jamfschool_device" "example" {
  udid = "00008020-001234567890002E"
}

output "device_name" {
  value = data.jamfschool_device.example.name
}
