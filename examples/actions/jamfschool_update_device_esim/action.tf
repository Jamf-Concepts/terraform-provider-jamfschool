
# Refresh the eSIM cellular plan on a device.
action "jamfschool_update_device_esim" "refresh_esim" {
  config {
    serial_number           = "C02ABC123DEF"
    server_url              = "https://carrier.example.com"
    requires_network_tether = false
  }
}
