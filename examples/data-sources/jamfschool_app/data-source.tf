
data "jamfschool_app" "example" {
  id = 1
}

output "app_name" {
  value = data.jamfschool_app.example.name
}
