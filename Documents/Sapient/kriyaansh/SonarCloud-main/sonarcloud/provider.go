package sonarcloud

import (
	"encoding/binary"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	log "github.com/sirupsen/logrus"
)

var sonarcloudProvider *schema.Provider

// Provider for sonarcloud
func Provider() terraform.ResourceProvider {
	sonarcloudProvider = &schema.Provider{
		// Provider configuration
		Schema: map[string]*schema.Schema{
			"user": {
				Type:        schema.TypeString,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SONAR_USER", "SONARCLOUD_USER"}, nil),
				Required:    true,
			},
			"pass": {
				Type:        schema.TypeString,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SONAR_PASS", "SONARCLOUD_PASS"}, nil),
				Required:    true,
			},
			"host": {
				Type:        schema.TypeString,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SONAR_HOST", "SONARCLOUD_HOST"}, nil),
				Required:    true,
			},
			"scheme": {
				Type:        schema.TypeString,
				Default:     "https",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{"SONAR_SCHEME", "SONARCLOUD_SCHEME"}, nil),
				Optional:    true,
			},
		},
		// Add the resources supported by this provider to this map.
		ResourcesMap: map[string]*schema.Resource{
			"sonarcloud_group":                           resourceSonarcloudGroup(),
			"sonarcloud_permission_template":             resourceSonarcloudPermissionTemplate(),
			"sonarcloud_permissions":                     resourceSonarcloudPermissions(),
			"sonarcloud_plugin":                          resourceSonarcloudPlugin(),
			"sonarcloud_project":                         resourceSonarcloudProject(),
			"sonarcloud_qualitygate":                     resourceSonarcloudQualityGate(),
			"sonarcloud_qualityprofile":                  resourceSonarcloudQualityProfile(),
			"sonarcloud_qualityprofile_copy":             resourceSonarcloudQualityProfileCopy(),
			"sonarcloud_qualitygate_condition":           resourceSonarcloudQualityGateCondition(),
			"sonarcloud_qualitygate_project_association": resourceSonarcloudQualityGateProjectAssociation(),
			"sonarcloud_user":                            resourceSonarcloudUser(),
			"sonarcloud_user_token":                      resourceSonarcloudUserToken(),
		},
		ConfigureFunc: configureProvider,
	}
	return sonarcloudProvider
}

//ProviderConfiguration contains the sonarcloud providers configuration
type ProviderConfiguration struct {
	httpClient    *retryablehttp.Client
	sonarCloudURL url.URL
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	client := retryablehttp.NewClient()

	sonarCloudURL := url.URL{
		Scheme:     d.Get("scheme").(string),
		Host:       d.Get("host").(string),
		User:       url.UserPassword(d.Get("user").(string), d.Get("pass").(string)),
		ForceQuery: true,
	}

	// Check that the sonarcloud api is available and a supported version
	err := sonarcloudHealth(client, sonarCloudURL)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &ProviderConfiguration{
		httpClient:    client,
		sonarCloudURL: sonarCloudURL,
	}, nil
}

func sonarcloudHealth(client *retryablehttp.Client, sonarcloud url.URL) error {
	// Make request to sonarcloud version endpoint
	sonarcloud.Path = "api/server/version"
	req, err := retryablehttp.NewRequest("GET", sonarcloud.String(), http.NoBody)
	if err != nil {
		log.Error(err)
		return errors.New("Unable to construct sonarcloud version request")
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return errors.New("Unable to reach sonarcloud")
	}
	defer resp.Body.Close()

	// Check response code
	if resp.StatusCode != http.StatusOK {
		return errors.New("Sonarcloud version api did not return a 200")
	}

	// Read in the response
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return errors.New("Failed to parse response body on GET sonarcloud version api")
	}

	// Convert response to a int.
	version := binary.BigEndian.Uint64(bodyBytes)
	if version < 8 {
		log.Error(err)
		return errors.New("Unsupported version of sonarcloud. Minimum supported version is 8")
	}

	return nil
}
