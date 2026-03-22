
# Clear the activation lock from a device.
action "jamfschool_clear_device_activation_lock" "clear_lock" {
  config {
    serial_number = "C02ABC123DEF"
  }
}
