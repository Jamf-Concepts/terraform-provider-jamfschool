
# Erase a device, clearing the activation lock.
action "jamfschool_erase_device" "return_device" {
  config {
    serial_number         = "C02ABC123DEF"
    clear_activation_lock = true
  }
}
