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
func resourceSonarcloudQualityProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarcloudQualityProfileCreate,
		Read:   resourceSonarcloudQualityProfileRead,
		Delete: resourceSonarcloudQualityProfileDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarcloudQualityProfileImport,
		},

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
                        "language": {
                                Type:     schema.TypeString,
                                Required: true,
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

func resourceSonarcloudQualityProfileCreate(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/qualityprofiles/create"
	sonarCloudURL.RawQuery = url.Values{
		"name": []string{d.Get("name").(string)},
                "language": []string{d.Get("language").(string)},
                "organization": []string{d.Get("organization").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarCloudURL.String(),
		http.StatusOK,
		"resourceQualityProfileCreate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	qualityProfileResponse := CreateQualityProfileResponse{}
	err = json.NewDecoder(resp.Body).Decode(&qualityProfileResponse)
	if err != nil {
		log.WithError(err).Error("resourceQualityProfileCreate: Failed to decode json into struct")
	}

	d.SetId(strconv.FormatInt(qualityProfileResponse.ID, 10))
	return nil
}

func resourceSonarcloudQualityProfileRead(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/qualityprofiles/search"
	sonarCloudURL.RawQuery = url.Values{
		"qualityProfile": []string{d.Get("name").(string)},
                "organization": []string{d.Get("organization").(string)},
                "language": []string{d.Get("language").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarCloudURL.String(),
		http.StatusOK,
		"resourceQualityProfileRead",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	qualityProfileReadResponse := GetQualityProfile{}
	err = json.NewDecoder(resp.Body).Decode(&qualityProfileReadResponse)
	if err != nil {
		log.WithError(err).Error("resourceQualityProfileRead: Failed to decode json into struct")
	}

	return nil
}

func resourceSonarcloudQualityProfileDelete(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/qualityprofiles/delete"
	sonarCloudURL.RawQuery = url.Values{
		"qualityProfile": []string{d.Get("name").(string)},
                "organization": []string{d.Get("organization").(string)},
                "language": []string{d.Get("language").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarCloudURL.String(),
		http.StatusNoContent,
		"resourceQualityProfileDelete",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	qualityProfileReadResponse := GetQualityProfile{}
	err = json.NewDecoder(resp.Body).Decode(&qualityProfileReadResponse)
	if err != nil {
		log.WithError(err).Error("resourceQualityProfileDelete: Failed to decode json into struct")
	}

	return nil
}

func resourceSonarcloudQualityProfileImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarcloudQualityProfileRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
