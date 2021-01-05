# Provider configuration

The sonarcloud provider is used to configure sonarcloud. The provider needs to be configured with a url, user and password.

## Example Usage
```terraform
provider "sonarcloud" {
    api_token   = "xxxxxxxxxxxxxxxx"
    host        = "sonarcloud.io"
    scheme      = "https"
}
```

## Argument Reference
The following arguments are supported:

- api_token - (Required) Sonarcloud access token. This can also be set via the SONARCLOUD_API_TOKEN environment variable.
- host - (Required) Sonarcloud url. This can be also be set via the SONARCLOUD_HOST environment variable.
- scheme - (Required) Http scheme to use. Either http or https. This can be also be set via the SONARCLOUD_SCHEME environment variable.
