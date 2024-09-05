// Package routes provides methods for managing an app's routes.
package routes

import (
	"encoding/json"
	"fmt"
	"io"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// List routes registered with an app.
func List(c *drycc.Client, appID string, results int) (api.Routes, int, error) {
	u := fmt.Sprintf("/v2/apps/%s/routes/", appID)
	body, count, reqErr := c.LimitedRequest(u, results)

	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.Route{}, -1, reqErr
	}

	var routes []api.Route
	if err := json.Unmarshal([]byte(body), &routes); err != nil {
		return []api.Route{}, -1, err
	}

	return routes, count, reqErr
}

// New adds a route to an app.
func New(c *drycc.Client, appID string, name string, Ptype string, kind string, port int) error {
	u := fmt.Sprintf("/v2/apps/%s/routes/", appID)

	req := api.RouteCreateRequest{Name: name, Ptype: Ptype, Kind: kind, Port: port}

	body, err := json.Marshal(req)

	if err != nil {
		return err
	}

	res, reqErr := c.Request("POST", u, body)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return reqErr
	}
	defer res.Body.Close()
	return reqErr
}

// AttachGateway route attach a gateway.
func AttachGateway(c *drycc.Client, appID string, name string, port int, gateway string) error {
	u := fmt.Sprintf("/v2/apps/%s/routes/%s/attach/", appID, name)

	req := api.RouteAttackRequest{Port: port, Gateway: gateway}

	body, err := json.Marshal(req)

	if err != nil {
		return err
	}

	res, reqErr := c.Request("PATCH", u, body)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return reqErr
	}
	defer res.Body.Close()

	return reqErr
}

// DetachGateway route attach a gateway.
func DetachGateway(c *drycc.Client, appID string, name string, port int, gateway string) error {
	u := fmt.Sprintf("/v2/apps/%s/routes/%s/detach/", appID, name)

	req := api.RouteDetackRequest{Port: port, Gateway: gateway}

	body, err := json.Marshal(req)

	if err != nil {
		return err
	}

	res, reqErr := c.Request("PATCH", u, body)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return reqErr
	}
	defer res.Body.Close()

	return reqErr
}

// GetRoute get info rule of a route from an app.
func GetRule(c *drycc.Client, appID string, name string) (string, error) {
	u := fmt.Sprintf("/v2/apps/%s/routes/%s/rules/", appID, name)
	res, err := c.Request("GET", u, nil)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	respBytes, err := io.ReadAll(res.Body)
	return string(respBytes), err
}

// SetRule set rule of a route.
func SetRule(c *drycc.Client, appID string, name string, rules string) error {
	u := fmt.Sprintf("/v2/apps/%s/routes/%s/rules/", appID, name)
	body, err := json.Marshal(rules)
	if err != nil {
		return err
	}
	res, err := c.Request("PUT", u, body)
	if err == nil {
		res.Body.Close()
	}
	return err
}

// Delete Delete a route from an app.
func Delete(c *drycc.Client, appID string, name string) error {
	u := fmt.Sprintf("/v2/apps/%s/routes/%s/", appID, name)
	res, err := c.Request("DELETE", u, nil)
	if err == nil {
		res.Body.Close()
	}
	return err
}
