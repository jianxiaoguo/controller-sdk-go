// Package config provides methods for managing configuration of apps.
package resources

import (
	"encoding/json"
	"fmt"
	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// List list an app's resources.
func List(c *drycc.Client, appID string, results int) (api.Resources, int, error) {
	u := fmt.Sprintf("/v2/apps/%s/resources/", appID)
	body, count, reqErr := c.LimitedRequest(u, results)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.Resource{}, -1, reqErr
	}
	var resources []api.Resource
	if err := json.Unmarshal([]byte(body), &resources); err != nil {
		return []api.Resource{}, -1, err
	}
	return resources, count, reqErr
}

// Create create an app's resource.
func Create(c *drycc.Client, appID string, resource api.Resource) (api.Resource, error) {
	body, err := json.Marshal(resource)
	if err != nil {
		return api.Resource{}, err
	}
	u := fmt.Sprintf("/v2/apps/%s/resources/", appID)
	res, reqErr := c.Request("POST", u, body)
	if reqErr != nil {
		return api.Resource{}, reqErr
	}
	defer res.Body.Close()
	newResource := api.Resource{}
	if err = json.NewDecoder(res.Body).Decode(&newResource); err != nil {
		return api.Resource{}, err
	}
	return newResource, reqErr
}

// Get retrieves information about a resource
func Get(c *drycc.Client, appID string, name string) (api.Resource, error) {
	u := fmt.Sprintf("/v2/apps/%s/resources/%s/", appID, name)
	res, reqErr := c.Request("GET", u, nil)
	if reqErr != nil {
		return api.Resource{}, reqErr
	}
	defer res.Body.Close()
	resResource := api.Resource{}
	if err := json.NewDecoder(res.Body).Decode(&resResource); err != nil {
		return api.Resource{}, err
	}
	return resResource, reqErr
}

// Delete delete an app's resource.
func Delete(c *drycc.Client, appID string, name string) error {
	u := fmt.Sprintf("/v2/apps/%s/resources/%s/", appID, name)
	res, err := c.Request("DELETE", u, nil)
	if err == nil {
		res.Body.Close()
	}
	return err
}

// Put update resource
func Put(c *drycc.Client, appID string, name string, resource api.Resource) (api.Resource, error) {
	body, err := json.Marshal(resource)
	if err != nil {
		return api.Resource{}, err
	}
	u := fmt.Sprintf("/v2/apps/%s/resources/%s/", appID, name)
	res, reqErr := c.Request("PUT", u, body)
	if reqErr != nil {
		return api.Resource{}, reqErr
	}
	defer res.Body.Close()
	newResource := api.Resource{}
	if err = json.NewDecoder(res.Body).Decode(&newResource); err != nil {
		return api.Resource{}, err
	}
	return newResource, reqErr
}

// Binding servicebinding binding with a serviceinstance
func Binding(c *drycc.Client, appID string, name string, resource api.Binding) (api.Resource, error) {
	body, err := json.Marshal(resource)
	if err != nil {
		return api.Resource{}, err
	}
	u := fmt.Sprintf("/v2/apps/%s/resources/%s/binding/", appID, name)
	res, reqErr := c.Request("PATCH", u, body)
	if reqErr != nil {
		return api.Resource{}, reqErr
	}
	defer res.Body.Close()
	newResource := api.Resource{}
	if err = json.NewDecoder(res.Body).Decode(&newResource); err != nil {
		return api.Resource{}, err
	}
	return newResource, reqErr
}
