# sonarcloud_user

Provides a Sonarcloud User resource. This can be used to manage Sonarcloud Users.

## Example: create a local user

```terraform
resource "sonarcloud_user" "user" {
  login_name = "terraform-test"
  name       = "terraform-test"
  password   = "secret-sauce37!"
}
```

## Example: create a remote user

```terraform
resource "sonarcloud_user" "remote_user" {
  login_name = "terraform-test"
  name       = "terraform-test"
  email      = "terraform-test@sonarcloud.com"
  is_local   = false
}
```

## Argument Reference

The following arguments are supported:

- login_name - (Required) The login name of the User to create. Changing this forces a new resource to be created.
- name - (Required) The name of the User to create. Changing this forces a new resource to be created.
- email - (Optional) The email of the User to create.
- password - (Optional) The password of User to create. This is only used if the user is of type `local`.
- is_local - (Optional) `True` if the User should be of type `local`. Defaults to `true`.

## Attributes Reference

The following attributes are exported:

- id - The ID of the User.

## Import

Users can be imported using their `login_name`:

```terraform
terraform import sonarcloud_user.user terraform-test
```
