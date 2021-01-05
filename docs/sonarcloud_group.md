# sonarcloud_group

Provides a Sonarcloud Group resource. This can be used to create and manage Sonarcloud Groups.

## Example: create a group

```terraform
resource "sonarcloud_group" "project_users" {
    name        = "Project-Users"
    description = "This is a group"
}
```

## Argument Reference

The following arguments are supported:

- name - (Required) The name of the Group to create. Changing this forces a new resource to be created.
- description - (Optional) Description of the Group.

## Attributes Reference

The following attributes are exported:

- id - The ID of the Group.

## Import

Groups can be imported using their ID

```terraform
terraform import sonarcloud_group.group 101
```
