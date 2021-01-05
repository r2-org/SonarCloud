package sonarcloud

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
)

// Returns the resource represented by this file.
func resourceSonarcloudQualityGate() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarcloudQualityGateCreate,
		Read:   resourceSonarcloudQualityGateRead,
		Delete: resourceSonarcloudQualityGateDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarcloudQualityGateImport,
		},

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSonarcloudQualityGateCreate(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/qualitygates/create"
	sonarCloudURL.RawQuery = url.Values{
		"name": []string{d.Get("name").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarCloudURL.String(),
		http.StatusOK,
		"resourceQualityGateCreate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	qualityGateResponse := CreateQualityGateResponse{}
	err = json.NewDecoder(resp.Body).Decode(&qualityGateResponse)
	if err != nil {
		log.WithError(err).Error("resourceQualityGateCreate: Failed to decode json into struct")
	}

	d.SetId(strconv.FormatInt(qualityGateResponse.ID, 10))
	return nil
}

func resourceSonarcloudQualityGateRead(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/qualitygates/show"
	sonarCloudURL.RawQuery = url.Values{
		"id": []string{d.Id()},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarCloudURL.String(),
		http.StatusOK,
		"resourceQualityGateRead",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	qualityGateReadResponse := GetQualityGate{}
	err = json.NewDecoder(resp.Body).Decode(&qualityGateReadResponse)
	if err != nil {
		log.WithError(err).Error("resourceQualityGateRead: Failed to decode json into struct")
	}

	d.SetId(strconv.FormatInt(qualityGateReadResponse.ID, 10))
	d.Set("name", qualityGateReadResponse.Name)
	return nil
}

func resourceSonarcloudQualityGateDelete(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/qualitygates/destroy"
	sonarCloudURL.RawQuery = url.Values{
		"id": []string{d.Id()},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarCloudURL.String(),
		http.StatusNoContent,
		"resourceQualityGateDelete",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	qualityGateReadResponse := GetQualityGate{}
	err = json.NewDecoder(resp.Body).Decode(&qualityGateReadResponse)
	if err != nil {
		log.WithError(err).Error("resourceQualityGateDelete: Failed to decode json into struct")
	}

	return nil
}

func resourceSonarcloudQualityGateImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarcloudQualityGateRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
