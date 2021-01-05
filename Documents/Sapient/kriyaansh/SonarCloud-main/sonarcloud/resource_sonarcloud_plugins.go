package sonarcloud

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
)

// Returns the resource represented by this file.
func resourceSonarcloudPlugin() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarcloudPluginCreate,
		Read:   resourceSonarcloudPluginRead,
		Delete: resourceSonarcloudPluginDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarcloudPluginImport,
		},

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSonarcloudPluginCreate(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/plugins/install"
	sonarCloudURL.RawQuery = url.Values{
		"key": []string{d.Get("key").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarCloudURL.String(),
		http.StatusNoContent,
		"resourceSonarcloudPluginCreate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	d.SetId(d.Get("key").(string))
	return nil
}

func resourceSonarcloudPluginRead(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/plugins/installed"

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarCloudURL.String(),
		http.StatusOK,
		"resourceSonarcloudPluginRead",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	getInstalledPlugins := GetInstalledPlugins{}
	err = json.NewDecoder(resp.Body).Decode(&getInstalledPlugins)
	if err != nil {
		log.WithError(err).Error("resourceSonarcloudPluginRead: Failed to decode json into struct")
	}

	// Loop over all projects to see if the project we need exists.
	for _, value := range getInstalledPlugins.Plugins {
		if d.Id() == value.Key {
			// If it does, set the values of that project
			d.SetId(value.Key)
			d.Set("key", value.Key)
		}
	}

	return nil
}

func resourceSonarcloudPluginDelete(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/plugins/uninstall"
	sonarCloudURL.RawQuery = url.Values{
		"key": []string{d.Id()},
	}.Encode()

	log.Error(sonarCloudURL.String())
	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarCloudURL.String(),
		http.StatusNoContent,
		"resourceSonarcloudPluginDelete",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func resourceSonarcloudPluginImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarcloudPluginRead(d, m); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
