// Package builds provides methods for managing app builds.
package builds

import (
	"encoding/json"
	"fmt"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// List lists an app's builds.
func List(c *drycc.Client, appID string, results int) ([]api.Build, int, error) {
	u := fmt.Sprintf("/v2/apps/%s/builds/", appID)
	body, count, reqErr := c.LimitedRequest(u, results)

	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.Build{}, -1, reqErr
	}

	var builds []api.Build
	if err := json.Unmarshal([]byte(body), &builds); err != nil {
		return []api.Build{}, -1, err
	}

	return builds, count, reqErr
}

// New creates a build for an app from an container image.
// By default this will create a cmd process that runs the CMD command from the Dockerfile.
// If you want to define more process types, you can pass a Procfile map,
// where the key is the process name and the value is the command for that process.
// To pull from a private container registry, a custom username and password must be set in the app's
// configuration object. This can be done with `drycc registry:set` or by using this SDK.
//
// This example adds custom registry credentials to an app:
//
//	import (
//		"github.com/drycc/controller-sdk-go/api"
//		"github.com/drycc/controller-sdk-go/config"
//	)
//
//	// Create username/password map
//	registryMap := map[string]string{
//		"username": "password"
//	}
//
//	// Create a new configuration, assign the credentials, and set it.
//	// Note that config setting is a patching operation, it doesn't overwrite or unset
//	// unrelated configuration.
//	newConfig := api.Config{}
//	newConfig.Registry = registryMap
//	_, err := config.Set(<client>, "appname", newConfig)
//	if err != nil {
//	    log.Fatal(err)
//	}
func New(c *drycc.Client, appID string, image string, stack string,
	procfile map[string]string, dryccfile map[string]interface{}) (api.Build, error) {

	u := fmt.Sprintf("/v2/apps/%s/builds/", appID)

	req := api.CreateBuildRequest{
		Image:     image,
		Stack:     stack,
		Procfile:  procfile,
		Dryccfile: dryccfile,
	}

	body, err := json.Marshal(req)

	if err != nil {
		return api.Build{}, err
	}

	res, reqErr := c.Request("POST", u, body)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.Build{}, reqErr
	}
	defer res.Body.Close()

	build := api.Build{}
	if err = json.NewDecoder(res.Body).Decode(&build); err != nil {
		return api.Build{}, err
	}

	return build, reqErr
}
