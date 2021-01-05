package sonarcloud

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// Returns the resource represented by this file.
func resourceSonarcloudUserToken() *schema.Resource {
	return &schema.Resource{
		Create: resourceSonarcloudUserTokenCreate,
		Read:   resourceSonarcloudUserTokenRead,
		Delete: resourceSonarcloudUserTokenDelete,

		// Define the fields of this schema.
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"login_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceSonarcloudUserTokenCreate(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/user_tokens/generate"

	rawQuery := url.Values{
		"login": []string{d.Get("login_name").(string)},
		"name":  []string{d.Get("name").(string)},
	}

	sonarCloudURL.RawQuery = rawQuery.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarCloudURL.String(),
		http.StatusOK,
		"resourceSonarcloudUserTokenCreate",
	)
	if err != nil {
		return fmt.Errorf("Error creating Sonarcloud user token: %+v", err)
	}
	defer resp.Body.Close()

	// Decode response into struct
	tokenResponse := Token{}
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarcloudUserTokenCreate: Failed to decode json into struct: %+v", err)
	}

	if tokenResponse.Login != "" {
		// the ID consists of the login_name and the token name (foo/bar)
		d.SetId(fmt.Sprintf("%s/%s", d.Get("login_name").(string), d.Get("name").(string)))
		// we set the token value here as the API wont return it later
		if tokenResponse.Token != "" {
			d.Set("token", tokenResponse.Token)
		} else {
			return fmt.Errorf("resourceSonarcloudUserTokenCreate: Create response didn't contain the token")
		}
	} else {
		return fmt.Errorf("resourceSonarcloudUserTokenCreate: Create response didn't contain the user login")
	}

	return resourceSonarcloudUserTokenRead(d, m)
}

func resourceSonarcloudUserTokenRead(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/user_tokens/search"
	sonarCloudURL.RawQuery = url.Values{
		"login": []string{d.Get("login_name").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"GET",
		sonarCloudURL.String(),
		http.StatusOK,
		"resourceSonarcloudUserTokenRead",
	)
	if err != nil {
		return fmt.Errorf("Error reading Sonarcloud user tokens: %+v", err)
	}
	defer resp.Body.Close()

	// Decode response into struct
	getTokensResponse := GetTokens{}
	err = json.NewDecoder(resp.Body).Decode(&getTokensResponse)
	if err != nil {
		return fmt.Errorf("resourceSonarcloudUserTokenCreate: Failed to decode json into struct: %+v", err)
	}

	// Loop over all user token to see if the current token exists.
	readSuccess := false
	if getTokensResponse.Tokens != nil {
		for _, value := range getTokensResponse.Tokens {
			if d.Get("name").(string) == value.Name {
				d.SetId(fmt.Sprintf("%s/%s", d.Get("login_name").(string), d.Get("name").(string)))
				d.Set("login_name", getTokensResponse.Login)
				d.Set("name", value.Name)
				readSuccess = true
			}
		}
	} else {
		// the user has no tokens
		d.SetId("")
	}

	if !readSuccess {
		// Token not found
		d.SetId("")
	}

	return nil
}

func resourceSonarcloudUserTokenDelete(d *schema.ResourceData, m interface{}) error {
	sonarCloudURL := m.(*ProviderConfiguration).sonarCloudURL
	sonarCloudURL.Path = "api/user_tokens/revoke"
	sonarCloudURL.RawQuery = url.Values{
		"login": []string{d.Get("login_name").(string)},
		"name":  []string{d.Get("name").(string)},
	}.Encode()

	resp, err := httpRequestHelper(
		m.(*ProviderConfiguration).httpClient,
		"POST",
		sonarCloudURL.String(),
		http.StatusNoContent,
		"resourceSonarcloudUserTokenDelete",
	)
	if err != nil {
		return fmt.Errorf("Error deleting Sonarcloud user token: %+v", err)
	}
	defer resp.Body.Close()

	return nil
}
