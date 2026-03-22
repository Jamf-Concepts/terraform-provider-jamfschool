
# Restart a device by serial number.
action "jamfschool_restart_device" "nightly_restart" {
  config {
    serial_number = "C02ABC123DEF"
  }
}
