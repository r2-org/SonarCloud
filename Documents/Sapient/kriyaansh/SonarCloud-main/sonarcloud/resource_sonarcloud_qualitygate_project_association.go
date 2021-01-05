package sonarcloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	log "github.com/sirupsen/logrus"
)

// Returns the resource represented by this file.
func resourceSonarcloudQualityGateProjectAssociation() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarcloudQualityGateProjectAssociationCreate,
		Read:   resourceSonarcloudQualityGateProjectAssociationRead,
		Delete: resourceSonarcloudQualityGateProjectAssociationDelete,

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"gateid": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"projectkey": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSonarcloudQualityGateProjectAssociationCreate(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/qualitygates/select"
	sonarCloudURL.RawQuery = url.Values{
		"gateId":     []string{d.Get("gateid").(string)},
		"projectKey": []string{d.Get("projectkey").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarCloudURL.String(),
		http.StatusNoContent,
		"resourceSonarcloudQualityGateProjectAssociationCreate",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	id := fmt.Sprintf("%v/%v", d.Get("gateid").(string), d.Get("projectkey").(string))
	d.SetId(id)
	return nil
}

func resourceSonarcloudQualityGateProjectAssociationRead(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/qualitygates/search"
	sonarCloudURL.RawQuery = url.Values{
		"gateId": []string{d.Get("gateid").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarCloudURL.String(),
		http.StatusOK,
		"resourceSonarcloudQualityGateProjectAssociationRead",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Decode response into struct
	qualityGateAssociationReadResponse := GetQualityGateAssociation{}
	err = json.NewDecoder(resp.Body).Decode(&qualityGateAssociationReadResponse)
	if err != nil {
		log.WithError(err).Error("resourceSonarcloudQualityGateProjectAssociationRead: Failed to decode json into struct")
	}

	// ID is in format <gateid>/<projectkey>. This splits the id into gateid and projectkey
	// EG: "1/my_project" >> ["1", "my_project"]
	idSlice := strings.Split(d.Id(), "/")

	for _, value := range qualityGateAssociationReadResponse.Results {
		if idSlice[1] == value.Key {
			d.Set("gateid", idSlice[0])
			d.Set("projectkey", value.Key)
		}
	}

	return nil
}

func resourceSonarcloudQualityGateProjectAssociationDelete(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/qualitygates/deselect"
	sonarCloudURL.RawQuery = url.Values{
		"gateId":     []string{d.Get("gateid").(string)},
		"projectKey": []string{d.Get("projectkey").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarCloudURL.String(),
		http.StatusNoContent,
		"resourceSonarcloudQualityGateProjectAssociationDelete",
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
