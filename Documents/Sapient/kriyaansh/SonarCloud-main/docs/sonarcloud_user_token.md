# sonarcloud_user_token

Provides a Sonarcloud User token resource. This can be used to manage Sonarcloud User tokens.

## Example: create a user, user token and output the token value

```terraform
resource "sonarcloud_user" "user" {
  login_name = "terraform-test"
  name       = "terraform-test"
  password   = "secret-sauce37!"
}

resource "sonarcloud_user_token" "token" {
  login_name = sonarcloud_user.user.login_name
  name       = "my-token"
}

output "user_token" {
  value = sonarcloud_user_token.token.token
}
```

## Argument Reference

The following arguments are supported:

- login_name - (Required) The login name of the User for which the token should be created. Changing this forces a new resource to be created.
- name - (Required) The name of the Token to create. Changing this forces a new resource to be created.

## Attributes Reference

The following attributes are exported:

- id - The ID of the Token.
- token - The Token value.

## Import

Import is not supported for this resource.
