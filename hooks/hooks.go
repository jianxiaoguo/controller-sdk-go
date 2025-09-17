// Package hooks implements the controller's builder hooks api.
//
// This is primarily intended to be consumed by the builder to communicate with the controller.
package hooks

import (
	"encoding/json"
	"fmt"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// UserFromKey retrives a user from their SSH key fingerprint.
func UserFromKey(c *drycc.Client, fingerprint string) (api.UserApps, error) {
	res, reqErr := c.Request("GET", fmt.Sprintf("/v2/hooks/key/%s", fingerprint), nil)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.UserApps{}, reqErr
	}

	defer res.Body.Close()

	resUser := api.UserApps{}
	if err := json.NewDecoder(res.Body).Decode(&resUser); err != nil {
		return api.UserApps{}, err
	}

	return resUser, reqErr
}

// GetAppConfig retrives an app's configuration from the controller.
func GetAppConfig(c *drycc.Client, username, app string) (api.Config, error) {
	req := api.ConfigHookRequest{User: username, App: app}
	b, err := json.Marshal(req)
	if err != nil {
		return api.Config{}, err
	}

	res, reqErr := c.Request("POST", "/v2/hooks/config/", b)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.Config{}, reqErr
	}
	defer res.Body.Close()

	config := api.Config{}
	if err := json.NewDecoder(res.Body).Decode(&config); err != nil {
		return api.Config{}, err
	}

	return config, reqErr
}

// CreateBuild creates a new release of an application. It returns the version of the new release.
// gitSha should be the first 8 characters of the git commit sha. Image is either the container image
// location for the dockerfile app the absolute url to the tar file for a buldpack app.
func CreateBuild(c *drycc.Client, username, app, image, stack, gitSha string,
	procfile api.ProcessType, dryccfile map[string]interface{}, dockerfile string,
) (int, error) {
	req := api.BuildHookRequest{
		Sha:        gitSha,
		User:       username,
		App:        app,
		Image:      image,
		Stack:      stack,
		Procfile:   procfile,
		Dockerfile: dockerfile,
		Dryccfile:  dryccfile,
	}

	b, err := json.Marshal(req)
	if err != nil {
		return -1, err
	}

	res, reqErr := c.Request("POST", "/v2/hooks/build/", b)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return -1, reqErr
	}
	defer res.Body.Close()

	resMap := make(map[string]map[string]int)
	if err := json.NewDecoder(res.Body).Decode(&resMap); err != nil {
		return -1, err
	}

	return resMap["release"]["version"], reqErr
}
