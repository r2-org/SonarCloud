package sonarcloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/satori/uuid"
)

// Returns the resource represented by this file.
func resourceSonarcloudPermissions() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarcloudPermissionsCreate,
		Read:   resourceSonarcloudPermissionsRead,
		Delete: resourceSonarcloudPermissionsDelete,

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"login_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"login_name", "group_name"},
			},
			"group_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"login_name", "group_name"},
			},
			"project_key": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				ConflictsWith: []string{"template_id"},
			},
			"template_id": {
				Type:          schema.TypeString,
				ForceNew:      true,
				Optional:      true,
				ConflictsWith: []string{"project_key"},
			},
			"permissions": {
				Type:     schema.TypeList,
				MinItems: 1,
				Required: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceSonarcloudPermissionsCreate(d *schema.ResourceData, m interface{}) error {

	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	permissions := expandPermissions(d)

	// build the base query
	RawQuery := url.Values{}

	// if the permissions should be applied to a project
	// we append the project_key to the request
	if projectKey, ok := d.GetOk("project_key"); ok {
		RawQuery.Add("projectKey", projectKey.(string))
	}

	// we use different API endpoints and request params
	// based on the target principal type (group or user)
	// and if its a direct or template permission
	if _, ok := d.GetOk("login_name"); ok {
		// user permission
		RawQuery.Add("login", d.Get("login_name").(string))
		if templateID, ok := d.GetOk("template_id"); ok {
			// template user permission
			sonarCloudURL.Path = "api/permissions/add_user_to_template"
			RawQuery.Add("templateId", templateID.(string))
		} else {
			// direct user permission
			sonarCloudURL.Path = "api/permissions/add_user"
		}
	} else {
		// group permission
		RawQuery.Add("groupName", d.Get("group_name").(string))
		if templateID, ok := d.GetOk("template_id"); ok {
			// template user permission
			sonarCloudURL.Path = "api/permissions/add_group_to_template"
			RawQuery.Add("templateId", templateID.(string))
		} else {
			// direct user permission
			sonarCloudURL.Path = "api/permissions/add_group"
		}
	}

	// loop through all permissions that should be applied
	for _, permission := range *&permissions {
		CurrentRawQuery := RawQuery
		CurrentRawQuery.Del("permission")
		CurrentRawQuery.Add("permission", permission)
		sonarCloudURL.RawQuery = CurrentRawQuery.Encode()

		resp, err := httpRequestHelper(
			m.(*ProviderConfiguration).httpClient,
			"POST",
			sonarCloudURL.String(),
			http.StatusNoContent,
			"resourceSonarcloudPermissionsCreate",
		)
		if err != nil {
			return fmt.Errorf("Error creating Sonarcloud permission: %+v", err)
		}
		defer resp.Body.Close()
	}

	// generate a unique ID
	d.SetId(uuid.NewV4().String())
	return resourceSonarcloudPermissionsRead(d, m)
}

func resourceSonarcloudPermissionsRead(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL

	// build the base query
	RawQuery := url.Values{
		// set the page size to 100
		"ps": []string{"100"},
	}

	// if the permissions should be applied to a project
	// we append the project_key to the request
	if projectKey, ok := d.GetOk("project_key"); ok {
		RawQuery.Add("projectKey", projectKey.(string))
	}

	// we use different API endpoints and request params
	// based on the target principal type (group or user)
	// and if its a direct or template permission
	if _, ok := d.GetOk("login_name"); ok {
		// permission target is USER
		if templateID, ok := d.GetOk("template_id"); ok {
			// template user permission
			sonarCloudURL.Path = "api/permissions/template_users"
			RawQuery.Add("templateId", templateID.(string))
		} else {
			// direct user permission
			sonarCloudURL.Path = "api/permissions/users"
		}
		sonarCloudURL.RawQuery = RawQuery.Encode()

		resp, err := httpRequestHelper(
			m.(*ProviderConfiguration).httpClient,
			"GET",
			sonarCloudURL.String(),
			http.StatusOK,
			"resourceSonarcloudPermissionsRead",
		)
		if err != nil {
			return fmt.Errorf("Error reading Sonarcloud permissions: %+v", err)
		}
		defer resp.Body.Close()

		// Decode response into struct
		users := GetUser{}
		err = json.NewDecoder(resp.Body).Decode(&users)
		if err != nil {
			return fmt.Errorf("resourceSonarcloudPermissionsRead: Failed to decode json into struct: %+v", err)
		}

		// Loop over all groups to see if the group we need exists.
		readSuccess := false
		loginName := d.Get("login_name").(string)
		for _, value := range users.Users {
			if strings.EqualFold(value.Login, loginName) {
				d.Set("login_name", value.Login)
				d.Set("permissions", flattenPermissions(&value.Permissions))
				readSuccess = true
			}
		}

		if !readSuccess {
			// Permission not found
			d.SetId("")
		}

	} else {
		// permission target is GROUP
		if templateID, ok := d.GetOk("template_id"); ok {
			// template group permission
			sonarCloudURL.Path = "api/permissions/template_groups"
			RawQuery.Add("templateId", templateID.(string))
		} else {
			// direct group permission
			sonarCloudURL.Path = "api/permissions/groups"
		}
		sonarCloudURL.RawQuery = RawQuery.Encode()

		resp, err := httpRequestHelper(
			m.(*ProviderConfiguration).httpClient,
			"GET",
			sonarCloudURL.String(),
			http.StatusOK,
			"resourceSonarcloudPermissionsRead",
		)
		if err != nil {
			return fmt.Errorf("Error reading Sonarcloud permissions: %+v", err)
		}
		defer resp.Body.Close()

		// Decode response into struct
		groups := GetGroupPermissions{}
		err = json.NewDecoder(resp.Body).Decode(&groups)
		if err != nil {
			return fmt.Errorf("resourceSonarcloudPermissionsRead: Failed to decode json into struct: %+v", err)
		}

		// Loop over all groups to see if the group we need exists.
		readSuccess := false
		groupName := d.Get("group_name").(string)
		for _, value := range groups.Groups {
			if strings.EqualFold(value.Name, groupName) {
				d.Set("group_name", value.Name)
				d.Set("permissions", flattenPermissions(&value.Permissions))

				readSuccess = true
			}
		}

		if !readSuccess {
			d.SetId("")
			return fmt.Errorf("resourceSonarcloudPermissionsRead: Unable to find group permissions for group: %s", groupName)
		}
	}

	return nil
}

func resourceSonarcloudPermissionsDelete(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	permissions := expandPermissions(d)

	// build the base query
	RawQuery := url.Values{}

	// if the permissions should be applied to a project
	// we append the project_key to the request
	if projectKey, ok := d.GetOk("project_key"); ok {
		RawQuery.Add("projectKey", projectKey.(string))
	}

	// we use different API endpoints and request params
	// based on the target principal type (group or user)
	if _, ok := d.GetOk("login_name"); ok {
		// permission target is USER
		if templateID, ok := d.GetOk("template_id"); ok {
			// template user permission
			sonarCloudURL.Path = "api/permissions/remove_user_from_template"
			RawQuery.Add("templateId", templateID.(string))
		} else {
			// direct user permission
			sonarCloudURL.Path = "api/permissions/remove_user"
		}
		RawQuery.Add("login", d.Get("login_name").(string))
		sonarCloudURL.RawQuery = RawQuery.Encode()

	} else {
		// permission target is GROUP
		if templateID, ok := d.GetOk("template_id"); ok {
			// template group permission
			sonarCloudURL.Path = "api/permissions/remove_group_from_template"
			RawQuery.Add("templateId", templateID.(string))
		} else {
			// direct group permission
			sonarCloudURL.Path = "api/permissions/remove_group"
		}
		RawQuery.Add("groupName", d.Get("group_name").(string))
		sonarCloudURL.RawQuery = RawQuery.Encode()
	}

	// loop through all permissions that should be applied
	for _, permission := range *&permissions {
		CurrentRawQuery := RawQuery
		CurrentRawQuery.Del("permission")
		CurrentRawQuery.Add("permission", permission)
		sonarCloudURL.RawQuery = CurrentRawQuery.Encode()

		resp, err := httpRequestHelper(
			m.(*ProviderConfiguration).httpClient,
			"POST",
			sonarCloudURL.String(),
			http.StatusNoContent,
			"resourceSonarcloudPermissionsDelete",
		)
		if err != nil {
			return fmt.Errorf("Error creating Sonarcloud permission: %+v", err)
		}
		defer resp.Body.Close()
	}

	return nil
}

func expandPermissions(d *schema.ResourceData) []string {
	expandedPermissions := make([]string, 0)
	flatPermissions := d.Get("permissions").([]interface{})
	for _, permission := range flatPermissions {
		expandedPermissions = append(expandedPermissions, permission.(string))
	}

	return expandedPermissions
}

func flattenPermissions(input *[]string) []interface{} {
	flatPermissions := make([]interface{}, 0)
	if input == nil {
		return flatPermissions
	}

	for _, permission := range *input {
		flatPermissions = append(flatPermissions, permission)
	}

	return flatPermissions
}
