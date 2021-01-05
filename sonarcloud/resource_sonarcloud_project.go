package sonarcloud

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
)

// Returns the resource represented by this file.
func resourceSonarcloudProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarcloudProjectCreate,
		Read:   resourceSonarcloudProjectRead,
		Delete: resourceSonarcloudProjectDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarcloudProjectImport,
		},

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"visibility": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "public",
				ForceNew: true,
			},
                        "organization": {
                                Type:     schema.TypeString,
                                Required: true,
                                ForceNew: true,
                        },
		},
	}
}

func resourceSonarcloudProjectCreate(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/projects/create"
	sonarCloudURL.RawQuery = url.Values{
		"name":       []string{d.Get("name").(string)},
		"project":    []string{d.Get("project").(string)},
		"visibility": []string{d.Get("visibility").(string)},
                "organization": []string{d.Get("organization").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarCloudURL.String(),
		http.StatusOK,
		"resourceSonarcloudProjectCreate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	projectResponse := CreateProjectResponse{}
	err = json.NewDecoder(resp.Body).Decode(&projectResponse)
	if err != nil {
		log.WithError(err).Error("resourceSonarcloudProjectCreate: Failed to decode json into struct")
	}

	d.SetId(projectResponse.Project.Key)
	return nil
}

func resourceSonarcloudProjectRead(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/projects/search"
	sonarCloudURL.RawQuery = url.Values{
		"project": []string{d.Id()},
                "organization": []string{d.Get("organization").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarCloudURL.String(),
		http.StatusOK,
		"resourceSonarcloudProjectRead",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	projectReadResponse := GetProject{}
	err = json.NewDecoder(resp.Body).Decode(&projectReadResponse)
	if err != nil {
		log.WithError(err).Error("resourceSonarcloudProjectRead: Failed to decode json into struct")
	}

	// Loop over all projects to see if the project we need exists.
	readSuccess := false
	for _, value := range projectReadResponse.Components {
		if d.Id() == value.Key {
			// If it does, set the values of that project
			d.SetId(value.Key)
			d.Set("name", value.Name)
			d.Set("key", value.Key)
			d.Set("visibility", value.Visibility)
			readSuccess = true
		}
	}

	if !readSuccess {
		// Project not found
		d.SetId("")
	}

	return nil
}

func resourceSonarcloudProjectDelete(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/projects/delete"
	sonarCloudURL.RawQuery = url.Values{
		"project": []string{d.Id()},
                "organization": []string{d.Get("organization").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarCloudURL.String(),
		http.StatusNoContent,
		"resourceSonarcloudProjectDelete",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func resourceSonarcloudProjectImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarcloudProjectRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
