
data "jamfschool_dep_device" "example" {
  serial_number = "C02ABC123DEF"
}

output "dep_device_model" {
  value = data.jamfschool_dep_device.example.model
}
