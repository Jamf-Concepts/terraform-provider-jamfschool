
# Configure the Jamf School provider.
# Prefer environment variables for credentials:
#   JAMFSCHOOL_URL, JAMFSCHOOL_NETWORK_ID, JAMFSCHOOL_API_KEY
provider "jamfschool" {
  url        = "https://myschool.jamfcloud.com"
  network_id = "10482058"
  api_key    = "your-api-key-here"
}
