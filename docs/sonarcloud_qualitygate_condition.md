# sonarcloud_qualitygate_condition
Provides a Sonarcloud Quality Gate Condition resource. This can be used to create and manage Sonarcloud Quality Gate conditions.

## Example: create a quality gate condition
```terraform
resource "sonarcloud_qualitygate" "main" {
    name = "my_qualitygate"
}

resource "sonarcloud_qualitygate_condition" "main" {
    gateid = sonarcloud_qualitygate.main.id
    metric = "vulnerabilities"
    error  = 10
    op     = "GT"
}
```

## Argument Reference
The following arguments are supported:

- gateid - (Required) The id of the Quality Gate
- metric - (Required) Condition metric. Only metric of the following types are allowed: INT, MILLISEC, RATING, WORK_DUR, FLOAT, PERCENT and LEVEL. Following metrics are forbidden: alert_status, security_hotspots and new_security_hotspots
- error - (Required) Condition error threshold
- op - (Required) Condition operator. Possible values are: LT and GT

## Attributes Reference
The following attributes are exported:

- id - ID of the Sonarcloud Quality Gate
- metric - Condition metric
- error - Condition error threshold
- warning - Condition warning threshold
- op - Condition operator

