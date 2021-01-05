terraform {
  required_providers {
    sonarcloud = {
     # source  = "terraform.kriyaansh.com/kriyaansh/sonarcloud"
     # source  = "github.com/meetdpv/SonarCloud"
     source  = "terraform.meetdpv.com/meetdpv/sonarcloud"
      version = ">= 1.0"
    }
  }
}

provider "sonarcloud" {
    user         = var.user
    pass         = var.pass
    host         = var.host
    scheme       = var.scheme
}

resource "sonarcloud_project" "main" {
    name         = var.project_name
    project      = var.project_key
    visibility   = var.visibility
    organization = var.organization
}

# resource "sonarcloud_group" "main" {
#     name         = var.group_name
#     description  = var.group_description
#     organization = var.organization
# }

resource "sonarcloud_qualityprofile" "main" {
    name         = var.profile_name
    organization = var.organization
    language     = var.language
}

# resource "sonarcloud_qualityprofile_copy" "main" {
#     name         = var.copy_profile_name
#     organization = var.organization
#     language     = var.language
# }
