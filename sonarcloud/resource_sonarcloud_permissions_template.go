package sonarcloud

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Returns the resource represented by this file.
func resourceSonarcloudPermissionTemplate() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarcloudPermissionTemplateCreate,
		Read:   resourceSonarcloudPermissionTemplateRead,
		Update: resourceSonarcloudPermissionTemplateUpdate,
		Delete: resourceSonarcloudPermissionTemplateDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarcloudPermissionTemplateImport,
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
			"project_key_pattern": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceSonarcloudPermissionTemplateCreate(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/permissions/create_template"
	sonarCloudURL.RawQuery = url.Values{
		"name":              []string{d.Get("name").(string)},
		"description":       []string{d.Get("description").(string)},
		"projectKeyPattern": []string{d.Get("project_key_pattern").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarCloudURL.String(),
		http.StatusOK,
		"resourceSonarcloudPermissionTemplateCreate",
	)
	if err != nil {
		return fmt.Errorf("Error creating Sonarcloud permission template: %+v", err)
	}
	defer resp.Body.Close()

	// Decode response into struct
	permissionTemplateResponse := CreatePermissionTemplateResponse{}
	err = json.NewDecoder(resp.Body).Decode(&permissionTemplateResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarcloudPermissionTemplateCreate: Failed to decode json into struct: %+v", err)
	}

	if permissionTemplateResponse.PermissionTemplate.ID != "" {
		d.SetId(permissionTemplateResponse.PermissionTemplate.ID)
	} else {
		return fmt.Errorf("resourceSonarcloudPermissionTemplateCreate: Create response didn't contain an ID")
	}

	return resourceSonarcloudPermissionTemplateRead(d, m)
}

func resourceSonarcloudPermissionTemplateRead(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/permissions/search_templates"
	sonarCloudURL.RawQuery = url.Values{
		"q": []string{d.Get("name").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarCloudURL.String(),
		http.StatusOK,
		"resourceSonarcloudPermissionTemplateRead",
	)
	if err != nil {
		return fmt.Errorf("Error reading Sonarcloud permission templates: %+v", err)
	}
	defer resp.Body.Close()

	// Decode response into struct
	permissionTemplateReadResponse := GetPermissionTemplates{}
	err = json.NewDecoder(resp.Body).Decode(&permissionTemplateReadResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarcloudPermissionTemplateRead: Failed to decode json into struct: %+v", err)
	}

	// Loop over all permission templates to see if the template we look for exists.
	readSuccess := false
	for _, value := range permissionTemplateReadResponse.PermissionTemplates {
		log.Printf("[DEBUG][resourceSonarcloudPermissionTemplateRead] Comparing '%s' with '%s'", d.Id(), value.ID)
		if d.Id() == value.ID {
			log.Printf("[DEBUG][resourceSonarcloudPermissionTemplateRead] Found PermissionTemplate with ID '%s'", value.ID)
			// If it does, set the values of that template
			d.SetId(value.ID)
			d.Set("name", value.Name)
			d.Set("description", value.Description)
			d.Set("project_key_pattern", value.ProjectKeyPattern)
			readSuccess = true
		}
	}

	if !readSuccess {
		// Resource not found
		log.Printf("[DEBUG][resourceSonarcloudPermissionTemplateRead] No permission template with ID '%s' found, removing from state", d.Id())
		d.SetId("")
	}

	return nil
}

func resourceSonarcloudPermissionTemplateUpdate(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/permissions/update_template"

	rawQuery := url.Values{
		"templateId": []string{d.Id()},
	}

	if _, ok := d.GetOk("description"); ok {
		rawQuery.Add("description", d.Get("description").(string))
	} else {
		rawQuery.Add("description", "")
	}

	if _, ok := d.GetOk("project_key_pattern"); ok {
		rawQuery.Add("projectKeyPattern", d.Get("project_key_pattern").(string))
	} else {
		rawQuery.Add("projectKeyPattern", "")
	}

	sonarCloudURL.RawQuery = rawQuery.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarCloudURL.String(),
		http.StatusOK,
		"resourceSonarcloudPermissionTemplateUpdate",
	)
	if err != nil {
		return fmt.Errorf("Error updating Sonarcloud permission template: %+v", err)
	}
	defer resp.Body.Close()

	return resourceSonarcloudPermissionTemplateRead(d, m)
}

func resourceSonarcloudPermissionTemplateDelete(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/permissions/delete_template"
	sonarCloudURL.RawQuery = url.Values{
		"templateId": []string{d.Id()},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarCloudURL.String(),
		http.StatusNoContent,
		"resourceSonarcloudPermissionTemplateDelete",
	)
	if err != nil {
		return fmt.Errorf("Error deleting Sonarcloud permission template: %+v", err)
	}
	defer resp.Body.Close()

	return nil
}

func resourceSonarcloudPermissionTemplateImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarcloudPermissionTemplateRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
