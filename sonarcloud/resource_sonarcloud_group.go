package sonarcloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Returns the resource represented by this file.
func resourceSonarcloudGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarcloudGroupCreate,
		Read:   resourceSonarcloudGroupRead,
		Update: resourceSonarcloudGroupUpdate,
		Delete: resourceSonarcloudGroupDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarcloudGroupImport,
		},

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
                        "organization": {
                                Type:     schema.TypeString,
                                Required: true,
                        },
		},
	}
}

func resourceSonarcloudGroupCreate(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/user_groups/create"
	sonarCloudURL.RawQuery = url.Values{
		"name":        []string{d.Get("name").(string)},
		"description": []string{d.Get("description").(string)},
                "organization": []string{d.Get("organization").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarCloudURL.String(),
		http.StatusOK,
		"resourceSonarcloudGroupCreate",
	)
	if err != nil {
		return fmt.Errorf("Error creating Sonarcloud group: %+v", err)
	}
	defer resp.Body.Close()

	// Decode response into struct
	groupResponse := CreateGroupResponse{}
	err = json.NewDecoder(resp.Body).Decode(&groupResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarcloudGroupRead: Failed to decode json into struct: %+v", err)
	}

	d.SetId(strconv.Itoa(groupResponse.Group.ID))
	return resourceSonarcloudGroupRead(d, m)
}

func resourceSonarcloudGroupRead(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/user_groups/search"
	sonarCloudURL.RawQuery = url.Values{
		"name": []string{d.Get("name").(string)},
                "organization": []string{d.Get("organization").(string)},

	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarCloudURL.String(),
		http.StatusOK,
		"resourceSonarcloudGroupRead",
	)
	if err != nil {
		return fmt.Errorf("Error reading Sonarcloud group: %+v", err)
	}
	defer resp.Body.Close()

	// Decode response into struct
	groupReadResponse := GetGroup{}
	err = json.NewDecoder(resp.Body).Decode(&groupReadResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarcloudGroupRead: Failed to decode json into struct: %+v", err)
	}

	// Loop over all groups to see if the group we need exists.
	readSuccess := false
	for _, value := range groupReadResponse.Groups {
		if d.Id() == strconv.Itoa(value.ID) {
			// If it does, set the values of that group
			d.SetId(strconv.Itoa(value.ID))
			d.Set("name", value.Name)
			d.Set("description", value.Description)
			readSuccess = true
		}
	}

	if !readSuccess {
		// Group not found
		d.SetId("")
	}

	return nil
}

func resourceSonarcloudGroupUpdate(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/user_groups/update"

	rawQuery := url.Values{
		"id": []string{d.Id()},
	}

	if _, ok := d.GetOk("description"); ok {
		rawQuery.Add("description", d.Get("description").(string))
	} else {
		rawQuery.Add("description", "")
	}

	sonarCloudURL.RawQuery = rawQuery.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarCloudURL.String(),
		http.StatusOK,
		"resourceSonarcloudGroupUpdate",
	)
	if err != nil {
		return fmt.Errorf("Error updating Sonarcloud group: %+v", err)
	}
	defer resp.Body.Close()

	return resourceSonarcloudGroupRead(d, m)
}

func resourceSonarcloudGroupDelete(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/user_groups/delete"
	sonarCloudURL.RawQuery = url.Values{
		"id": []string{d.Id()},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarCloudURL.String(),
		http.StatusNoContent,
		"resourceSonarcloudGroupDelete",
	)
	if err != nil {
		return fmt.Errorf("Error deleting Sonarcloud group: %+v", err)
	}
	defer resp.Body.Close()

	return nil
}

func resourceSonarcloudGroupImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarcloudGroupRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
