# Terraform Provider for Jamf School

> [!NOTE]
> This provider is in early development. All resources have been tested via acceptance tests against a real Jamf School tenant. However, the API surface is subject to change as we gather feedback from the community.

The Jamf School Terraform provider allows you to manage [Jamf School](https://www.jamf.com/products/jamf-school/) resources via the Jamf School REST API. Built using the [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework) v1.18.0 (Protocol v6).

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.26 (for building from source)

## Installation

The provider is published to the [Terraform Registry](https://registry.terraform.io/providers/Jamf-Concepts/jamfschool). Add the following to your Terraform configuration:

```hcl
terraform {
  required_providers {
    jamfschool = {
      source = "Jamf-Concepts/jamfschool"
    }
  }
}
```

## Authentication

The provider authenticates against the Jamf School REST API using HTTP Basic Auth. Configure credentials via the provider block or environment variables:

| Provider Attribute | Environment Variable      | Description                                                                                                  |
| ------------------ | ------------------------- | ------------------------------------------------------------------------------------------------------------ |
| `url`              | `JAMFSCHOOL_URL`          | Base URL (e.g. `https://myschool.jamfcloud.com`)                                                             |
| `network_id`       | `JAMFSCHOOL_NETWORK_ID`   | Network ID — found at Devices > Enroll Device(s) in Jamf School                                              |
| `api_key`          | `JAMFSCHOOL_API_KEY`      | API key — generated at Organization > Settings > API in Jamf School                                          |

### Example Provider Configuration

```hcl
provider "jamfschool" {
  url        = "https://myschool.jamfcloud.com"
  network_id = var.jamfschool_network_id
  api_key    = var.jamfschool_api_key
}
```

Or use environment variables and leave the provider block empty:

```hcl
provider "jamfschool" {}
```

## Supported Resources

| Resource                    | Description                          |
| --------------------------- | ------------------------------------ |
| `jamfschool_user`           | Manage users                         |
| `jamfschool_user_group`     | Manage user groups                   |
| `jamfschool_device_group`   | Manage device groups                 |
| `jamfschool_class`          | Manage classes                       |
| `jamfschool_ibeacon`        | Manage iBeacon regions               |

All resources support full CRUD operations and `terraform import`.

## Supported Data Sources

| Data Source                   | Description                                      |
| ----------------------------- | ------------------------------------------------ |
| `jamfschool_user`             | Look up a user by ID                             |
| `jamfschool_user_group`       | Look up a user group by ID                       |
| `jamfschool_device_group`     | Look up a device group by ID                     |
| `jamfschool_class`            | Look up a class by UUID                          |
| `jamfschool_ibeacon`          | Look up an iBeacon by ID                         |
| `jamfschool_device`           | Look up an enrolled device by UDID               |
| `jamfschool_app`              | Look up an app by ID                             |
| `jamfschool_profile`          | Look up a configuration profile by ID            |
| `jamfschool_location`         | Look up a location by ID                         |
| `jamfschool_dep_device`       | Look up a DEP device by serial number            |

## Usage Examples

### User

```hcl
resource "jamfschool_user" "student" {
  username   = "jsmith"
  password   = var.user_password
  email      = "jsmith@school.example.com"
  first_name = "John"
  last_name  = "Smith"
}
```

### User Group

```hcl
resource "jamfschool_user_group" "teachers" {
  name        = "Teachers"
  description = "All teaching staff"
}
```

### Device Group

```hcl
resource "jamfschool_device_group" "lab_ipads" {
  name        = "Computer Lab iPads"
  description = "Shared iPads in the computer lab"
}
```

### Class

```hcl
resource "jamfschool_class" "maths" {
  name        = "Year 10 Maths"
  description = "Mathematics class for Year 10 students"
}
```

### iBeacon

```hcl
resource "jamfschool_ibeacon" "library" {
  name        = "Library"
  description = "iBeacon region for the school library"
}
```

### Data Source

```hcl
data "jamfschool_device" "macbook" {
  udid = "AAAAAAAA-BBBB-CCCC-DDDD-EEEEEEEEEEEE"
}

data "jamfschool_location" "main_campus" {
  id = 1
}
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for development setup, testing instructions, and contribution guidelines.

## Feedback & Discussion

Please file issues and feature requests via [GitHub Issues](https://github.com/Jamf-Concepts/terraform-provider-jamfschool/issues).

The Jamf Terraform community has discussions in #terraform on [MacAdmins Slack](https://www.macadmins.org/).

## Included Components

The following third party acknowledgements and licenses are incorporated by reference:

- [Jamf School Go SDK](https://github.com/Jamf-Concepts/jamfschool-go-sdk) ([MIT](https://github.com/Jamf-Concepts/jamfschool-go-sdk?tab=MIT-1-ov-file))
- [Terraform Plugin Framework](https://github.com/hashicorp/terraform-plugin-framework) ([MPL](https://github.com/hashicorp/terraform-plugin-framework?tab=MPL-2.0-1-ov-file))
- [Terraform Plugin Go](https://github.com/hashicorp/terraform-plugin-go) ([MPL](https://github.com/hashicorp/terraform-plugin-go?tab=MPL-2.0-1-ov-file))
- [Terraform Plugin Log](https://github.com/hashicorp/terraform-plugin-log) ([MPL](https://github.com/hashicorp/terraform-plugin-log?tab=MPL-2.0-1-ov-file))
- [Terraform Plugin Testing](https://github.com/hashicorp/terraform-plugin-testing) ([MPL](https://github.com/hashicorp/terraform-plugin-testing?tab=MPL-2.0-1-ov-file))
