variable user {
  description = "contains the sonarcloud username"
}

variable pass {
  description = "contains the sonarcloud password"
  default     = ""
}

variable host {
  description = "contains the sonarcloud hostname"
  default     = "sonarcloud.io"
}

variable scheme {
  description = "contains the sonarcloud scheme name"
  default     = "https"
}

variable organization {
  description = "contains the sonarcloud organization name"
  default     = "meetdpv"
}

variable visibility {
  description = "contains the sonarcloud project visibility"
  default     = "public"
}

variable project_name {
  description = "contains the sonarcloud project name"
  default     = "sonarcloud_test_project"
}

variable project_key {
  description = "contains the sonarcloud project key"
  default     = "sonarcloud_test_project_key"
}

variable group_name {
  description = "contains the sonarcloud group name"
  default     = "sonarcloud_test_group"
}

variable group_description {
  description = "contains the sonarcloud group description"
  default     = "this is a test group created for IaC validation"
}

variable language {
  description = "contains the sonarcloud analysis language name"
  default     = "java"
}

variable profile_name {
  description = "contains the sonarcloud profile name"
  default     = "sonarcloud_test_profile"
}

variable copy_profile_name {
  description = "contains the sonarcloud profile name"
  default     = "sonarcloud_copy_profile"
}
