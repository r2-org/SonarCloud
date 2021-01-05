# terraform-sonarcloud
Terraform provider for managing SonarCloud configuration

## Installation
To install this provider based on terraform version follow under mentioned steps:
1. Terraform Version <= 0.12
2. Terraform Version >= 0.13

## Usage
[example](example) contains a sample code of how to use this provider.

Consult the docs below for more details.

## Docs
[Provider configuration](docs/provider.md)

Resources:
- [sonarcloud_group](docs/sonarcloud_group.md)
- [sonarcloud_permissions](docs/sonarcloud_permissions.md)
- [sonarcloud_permission_template](docs/sonarcloud_permission_template.md)
- [sonarcloud_project](docs/sonarcloud_project.md)
- [sonarcloud_qualityprofile](docs/sonarcloud_qualityprofile.md)
- [sonarcloud_qualitygate](docs/sonarcloud_qualitygate.md)
- [sonarcloud_qualitygate_condition](docs/sonarcloud_qualitygate_condition.md)
- [sonarcloud_qualitygate_project_association](docs/sonarcloud_qualitygate_project_association.md)
- [sonarcloud_user](docs/sonarcloud_user.md)
- [sonarcloud_user_token](docs/sonarcloud_user_token.md)

TODO:
- rules
- settings
- webhooks
