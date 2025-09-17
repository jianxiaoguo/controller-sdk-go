// Package builds provides methods for managing app builds.
package builds

import (
	"encoding/json"
	"fmt"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// Get a build of an app.
func Get(c *drycc.Client, appID string, version int) (api.Build, error) {
	u := fmt.Sprintf("/v2/apps/%s/build/", appID)
	if version > 0 {
		u = fmt.Sprintf("%s?version=v%d", u, version)
	}

	res, reqErr := c.Request("GET", u, nil)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.Build{}, reqErr
	}
	defer res.Body.Close()

	var build api.Build
	if err := json.NewDecoder(res.Body).Decode(&build); err != nil {
		return api.Build{}, err
	}
	return build, reqErr
}

// New a build of an app.
func New(c *drycc.Client, appID string, image string, stack string,
	procfile map[string]string, dryccfile map[string]interface{},
) (api.Build, error) {
	u := fmt.Sprintf("/v2/apps/%s/build/", appID)

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
