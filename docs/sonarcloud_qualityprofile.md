# sonarcloud_qualityprofile
Provides a Sonarcloud Quality Profile resource. This can be used to create and manage Sonarcloud Quality Profiles.

## Example: create a quality profile
```terraform
resource "sonarcloud_qualityprofile" "main" {
    name = "example"
}
```

## Argument Reference
The following arguments are supported:

- name - (Required) The name of the Quality Profile to create.
- organization - (Required) The name of the organization
- language - (Required) The name of the language

## Attributes Reference
The following attributes are exported:

- name - Name of the Sonarcloud Quality Profile
- id - ID of the Sonarcloud Quality Gate
