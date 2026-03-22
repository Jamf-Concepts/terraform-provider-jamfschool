
data "jamfschool_class" "example" {
  uuid = "550e8400-e29b-41d4-a716-446655440000"
}

output "class_name" {
  value = data.jamfschool_class.example.name
}
