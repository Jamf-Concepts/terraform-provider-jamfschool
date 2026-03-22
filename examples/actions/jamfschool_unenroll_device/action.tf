
# Remove the management profile from a device.
action "jamfschool_unenroll_device" "unenroll" {
  config {
    serial_number = "C02ABC123DEF"
  }
}
