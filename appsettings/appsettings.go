// Package appsettings provides methods for managing application settings of apps.
package appsettings

import (
	"encoding/json"
	"fmt"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// List lists an app's settings.
func List(c *drycc.Client, app string) (api.AppSettings, error) {
	u := fmt.Sprintf("/v2/apps/%s/settings/", app)

	res, reqErr := c.Request("GET", u, nil)
	if reqErr != nil {
		return api.AppSettings{}, reqErr
	}
	defer res.Body.Close()

	settings := api.AppSettings{}
	if err := json.NewDecoder(res.Body).Decode(&settings); err != nil {
		return api.AppSettings{}, err
	}

	return settings, reqErr
}

// Set sets an app's settings variables.
// This is a patching operation, which means when you call Set() with an api.AppSettings:
//
//   - If the variable does not exist, it will be set.
//   - If the variable exists, it will be overwritten.
//   - If the variable is set to nil, it will be unset.
//   - If the variable was ignored in the api.AppSettings, it will remain unchanged.
//
// Calling Set() with an empty api.AppSettings will return a drycc.ErrConflict.
func Set(c *drycc.Client, app string, appSettings api.AppSettings) (api.AppSettings, error) {
	body, err := json.Marshal(appSettings)

	if err != nil {
		return api.AppSettings{}, err
	}

	u := fmt.Sprintf("/v2/apps/%s/settings/", app)

	res, reqErr := c.Request("POST", u, body)
	if reqErr != nil {
		return api.AppSettings{}, reqErr
	}
	defer res.Body.Close()

	newAppSettings := api.AppSettings{}
	if err = json.NewDecoder(res.Body).Decode(&newAppSettings); err != nil {
		return api.AppSettings{}, err
	}

	return newAppSettings, reqErr
}

// CanaryDelete remove an app's canary settings.
func CanaryRemove(c *drycc.Client, app string, appSettings api.AppSettings) error {
	body, err := json.Marshal(appSettings)

	if err != nil {
		return err
	}

	u := fmt.Sprintf("/v2/apps/%s/settings/", app)
	res, err := c.Request("DELETE", u, body)
	if err == nil {
		res.Body.Close()
	}
	return err
}

// CanaryRelease release an app's canary settings.
func CanaryRelease(c *drycc.Client, app string) error {
	u := fmt.Sprintf("/v2/apps/%s/canary/release/", app)
	res, err := c.Request("POST", u, nil)
	if err == nil {
		res.Body.Close()
	}
	return nil
}

// CanaryRollback rollback an app's canary settings.
func CanaryRollback(c *drycc.Client, app string) error {
	u := fmt.Sprintf("/v2/apps/%s/canary/rollback/", app)
	res, err := c.Request("POST", u, nil)
	if err == nil {
		res.Body.Close()
	}
	return nil
}
