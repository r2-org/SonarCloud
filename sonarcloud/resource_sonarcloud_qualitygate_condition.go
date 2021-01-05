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
func resourceSonarcloudQualityGateCondition() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarcloudQualityGateConditionCreate,
		Read:   resourceSonarcloudQualityGateConditionRead,
		Update: resourceSonarcloudQualityGateConditionUpdate,
		Delete: resourceSonarcloudQualityGateConditionDelete,

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"gateid": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"error": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"metric": {
				Type:     schema.TypeString,
				Required: true,
			},
			"op": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceSonarcloudQualityGateConditionCreate(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/qualitygates/create_condition"
	sonarCloudURL.RawQuery = url.Values{
		"gateId": []string{strconv.Itoa(d.Get("gateid").(int))},
		"error":  []string{strconv.Itoa(d.Get("error").(int))},
		"metric": []string{d.Get("metric").(string)},
		"op":     []string{d.Get("op").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarCloudURL.String(),
		http.StatusOK,
		"resourcequalityGateConditionCreate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	qualityGateConditionResponse := CreateQualityGateConditionResponse{}
	err = json.NewDecoder(resp.Body).Decode(&qualityGateConditionResponse)
	if err != nil {
		log.WithError(err).Error("getQualityGateConditionResponse: Failed to decode json into struct")
	}

	d.SetId(strconv.FormatInt(qualityGateConditionResponse.ID, 10))
	return nil
}

func resourceSonarcloudQualityGateConditionRead(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/qualitygates/show"
	sonarCloudURL.RawQuery = url.Values{
		"id": []string{strconv.Itoa(d.Get("gateid").(int))},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarCloudURL.String(),
		http.StatusOK,
		"resourcequalityGateConditionRead",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	getQualityGateConditionResponse := GetQualityGate{}
	err = json.NewDecoder(resp.Body).Decode(&getQualityGateConditionResponse)
	if err != nil {
		log.WithError(err).Error("getQualityGateConditionResponse: Failed to decode json into struct")
	}

	for _, value := range getQualityGateConditionResponse.Conditions {
		if d.Id() == strconv.FormatInt(value.ID, 10) {
			d.SetId(strconv.FormatInt(value.ID, 10))
			d.Set("gateid", getQualityGateConditionResponse.ID)
			d.Set("error", value.Error)
			d.Set("metric", value.Metric)
			d.Set("op", value.OP)
		}
	}

	return nil
}

func resourceSonarcloudQualityGateConditionUpdate(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/qualitygates/update_condition"
	sonarCloudURL.RawQuery = url.Values{
		"gateid": []string{strconv.Itoa(d.Get("gateid").(int))},
		"id":     []string{d.Id()},
		"error":  []string{strconv.Itoa(d.Get("error").(int))},
		"metric": []string{d.Get("metric").(string)},
		"op":     []string{d.Get("op").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarCloudURL.String(),
		http.StatusOK,
		"resourcequalityGateConditionUpdate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return resourceSonarcloudQualityGateConditionRead(d, m)
}

func resourceSonarcloudQualityGateConditionDelete(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/qualitygates/delete_condition"
	sonarCloudURL.RawQuery = url.Values{
		"id": []string{d.Id()},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarCloudURL.String(),
		http.StatusNoContent,
		"resourcequalityGateConditionDelete",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
