
data "jamfschool_user" "example" {
  id = 123
}

output "user_email" {
  value = data.jamfschool_user.example.email
}
