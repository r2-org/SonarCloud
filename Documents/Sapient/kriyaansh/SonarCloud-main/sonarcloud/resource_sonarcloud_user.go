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
func resourceSonarcloudUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarcloudUserCreate,
		Read:   resourceSonarcloudUserRead,
		Update: resourceSonarcloudUserUpdate,
		Delete: resourceSonarcloudUserDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSonarcloudUserImport,
		},

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"login_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"is_local": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
		},
	}
}

func resourceSonarcloudUserCreate(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/users/create"
	isLocal := d.Get("is_local").(bool)

	rawQuery := url.Values{
		"login": []string{d.Get("login_name").(string)},
		"name":  []string{d.Get("name").(string)},
		"local": []string{strconv.FormatBool(isLocal)},
	}

	if password, ok := d.GetOk("password"); ok {
		rawQuery.Add("password", password.(string))
	}

	if email, ok := d.GetOk("email"); ok {
		rawQuery.Add("email", email.(string))
	}

	sonarCloudURL.RawQuery = rawQuery.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarCloudURL.String(),
		http.StatusOK,
		"resourceSonarcloudUserCreate",
	)
	if err != nil {
		return fmt.Errorf("Error creating Sonarcloud user: %+v", err)
	}
	defer resp.Body.Close()

	// Decode response into struct
	userResponse := CreateUserResponse{}
	err = json.NewDecoder(resp.Body).Decode(&userResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarcloudUserCreate: Failed to decode json into struct: %+v", err)
	}

	if userResponse.User.Login != "" {
		d.SetId(userResponse.User.Login)
	} else {
		return fmt.Errorf("resourceSonarcloudUserCreate: Create response didn't contain the user login")
	}

	return resourceSonarcloudUserRead(d, m)
}

func resourceSonarcloudUserRead(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/users/search"
	sonarCloudURL.RawQuery = url.Values{
		"q": []string{d.Get("login_name").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarCloudURL.String(),
		http.StatusOK,
		"resourceSonarcloudUserRead",
	)
	if err != nil {
		return fmt.Errorf("Error reading Sonarcloud user: %+v", err)
	}
	defer resp.Body.Close()

	// Decode response into struct
	userResponse := GetUser{}
	err = json.NewDecoder(resp.Body).Decode(&userResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarcloudUserCreate: Failed to decode json into struct: %+v", err)
	}

	// Loop over all users to see if the current user exists.
	readSuccess := false
	for _, value := range userResponse.Users {
		if d.Id() == value.Login {
			d.SetId(value.Login)
			d.Set("login_name", value.Login)
			d.Set("name", value.Name)
			d.Set("email", value.Email)
			d.Set("is_local", value.IsLocal)
			readSuccess = true
		}
	}

	if !readSuccess {
		// user not found
		d.SetId("")
	}

	return nil
}

func resourceSonarcloudUserUpdate(d *schema.ResourceData, m interface{}) error {

	// handle default updates (api/users/update)
	if d.HasChange("email") {
		sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
		sonarCloudURL.Path = "api/users/update"
		sonarCloudURL.RawQuery = url.Values{
			"login": []string{d.Id()},
			"email": []string{d.Get("email").(string)},
		}.Encode()

		resp, err := httpRequestHelper(
			m.(*ProviderConfiguration).httpClient,
			"POST",
			sonarCloudURL.String(),
			http.StatusOK,
			"resourceSonarcloudUserUpdate",
		)
		if err != nil {
			return fmt.Errorf("Error updating Sonarcloud user: %+v", err)
		}
		defer resp.Body.Close()
	}

	// handle password updates (api/users/change_password)
	if d.HasChange("password") {

		sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
		sonarCloudURL.Path = "api/users/change_password"
		sonarCloudURL.RawQuery = url.Values{
			"login":    []string{d.Id()},
			"password": []string{d.Get("password").(string)},
		}.Encode()

		resp, err := httpRequestHelper(
			m.(*ProviderConfiguration).httpClient,
			"POST",
			sonarCloudURL.String(),
			http.StatusNoContent,
			"resourceSonarcloudUserUpdate",
		)
		if err != nil {
			return fmt.Errorf("Error updating Sonarcloud user: %+v", err)
		}
		defer resp.Body.Close()
	}

	return resourceSonarcloudUserRead(d, m)
}

func resourceSonarcloudUserDelete(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/users/deactivate"
	sonarCloudURL.RawQuery = url.Values{
		"login": []string{d.Id()},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarCloudURL.String(),
		http.StatusOK,
		"resourceSonarcloudUserDelete",
	)
	if err != nil {
		return fmt.Errorf("Error deleting (deactivating) Sonarcloud user: %+v", err)
	}
	defer resp.Body.Close()

	return nil
}

func resourceSonarcloudUserImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceSonarcloudUserRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
