// Package services provides methods for managing an app's services.
package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// List services registered with an app.
func List(c *drycc.Client, appID string) (api.Services, error) {
	u := fmt.Sprintf("/v2/apps/%s/services/", appID)
	res, reqErr := c.Request("GET", u, nil)

	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.Service{}, reqErr
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []api.Service{}, err
	}

	r := make(map[string]interface{})
	if err = json.Unmarshal([]byte(body), &r); err != nil {
		return []api.Service{}, err
	}

	out, err := json.Marshal(r["services"].([]interface{}))
	if err != nil {
		return []api.Service{}, err
	}

	var services []api.Service
	if err := json.Unmarshal([]byte(out), &services); err != nil {
		return []api.Service{}, err
	}

	return services, reqErr
}

// New adds a new service to an app. App should already exists.
// Service is the way to route some traffic matching given URL pattern to worker different than `web`
// procfileType - name of the process in Procfile (i.e. <process type> from the `<process type>: <command>`), e.g. `webhooks`
// for more about Procfile see this https://devcenter.heroku.com/articles/procfile
// procfileType and pathPattern are mandatory and should have valid values.
func New(c *drycc.Client, appID string, procfileType string, port int, protocol string, targetPort int) (api.Service, error) {
	u := fmt.Sprintf("/v2/apps/%s/services/", appID)

	req := api.ServiceCreateUpdateRequest{ProcfileType: procfileType, Port: port, Protocol: protocol, TargetPort: targetPort}

	body, err := json.Marshal(req)

	if err != nil {
		return api.Service{}, err
	}

	res, reqErr := c.Request("POST", u, body)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.Service{}, reqErr
	}
	defer res.Body.Close()

	d := api.Service{ProcfileType: procfileType, Port: port, Protocol: protocol, TargetPort: targetPort}
	return d, reqErr
}

// Delete service from app
// If given service for the app doesn't exists then error returned
func Delete(c *drycc.Client, appID string, procfileType string) error {
	u := fmt.Sprintf("/v2/apps/%s/services/", appID)

	req := api.ServiceDeleteRequest{ProcfileType: procfileType}

	body, err := json.Marshal(req)

	if err != nil {
		return err
	}

	_, err = c.Request("DELETE", u, body)

	return err
}
