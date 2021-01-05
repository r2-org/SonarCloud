# sonarcloud_qualitygate_project_association
Provides a Sonarcloud Quality Gate Project association resource. This can be used to associate a Quality Gate to a Project

## Example: create a quality gate project association
```terraform
resource "sonarcloud_qualitygate" "main" {
    name = "my_qualitygate"
}

resource "sonarcloud_project" "main" {
    name       = "SonarCloud"
    project    = "my_project"
    visibility = "public" 
}

resource "sonarcloud_qualitygate_project_association" "main" {
    gateid     = sonarcloud_qualitygate.main.id
    projectkey = sonarcloud_project.main.project
}
```

## Argument Reference
The following arguments are supported:

- gateid - (Required) The id of the Quality Gate
- projectkey - (Required) Key of the project. Maximum length 400. All letters, digits, dash, underscore, period or colon.
