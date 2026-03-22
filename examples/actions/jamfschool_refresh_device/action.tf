
# Refresh a device's inventory, clearing any installation errors.
action "jamfschool_refresh_device" "refresh_inventory" {
  config {
    serial_number = "C02ABC123DEF"
    clear_errors  = true
  }
}
