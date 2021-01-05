package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/kriyaanshtechnologies/SonarCloud/sonarcloud"
)

func main() {
	plugin.Serve(
		&plugin.ServeOpts{
			ProviderFunc: sonarcloud.Provider,
		},
	)
}
