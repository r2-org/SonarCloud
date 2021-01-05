# sonarcloud_qualitygate
Provides a Sonarcloud Quality Gate resource. This can be used to create and manage Sonarcloud Quality Gates.

## Example: create a quality gate
```terraform
resource "sonarcloud_qualitygate" "main" {
    name = "example"
}
```

## Argument Reference
The following arguments are supported:

- name - (Required) The name of the Quality Gate to create. Maximum length 100

## Attributes Reference
The following attributes are exported:

- name - Name of the Sonarcloud Quality Gate
- id - ID of the Sonarcloud Quality Gate

## Import 
Quality Gates can be imported using their numeric value

```terraform
terraform import sonarcloud_qualitygate.main 11
```

