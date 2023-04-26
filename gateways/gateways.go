// Package gateways provides methods for managing an app's gateways.
package gateways

import (
	"encoding/json"
	"fmt"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// List gateways registered with an app.
func List(c *drycc.Client, appID string, results int) (api.Gateways, int, error) {
	u := fmt.Sprintf("/v2/apps/%s/gateways/", appID)
	body, count, reqErr := c.LimitedRequest(u, results)

	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.Gateway{}, -1, reqErr
	}

	var gateways []api.Gateway
	if err := json.Unmarshal([]byte(body), &gateways); err != nil {
		return []api.Gateway{}, -1, err
	}

	return gateways, count, reqErr
}

// New adds a gateway to an app.
func New(c *drycc.Client, appID string, name string, port int, protocol string) error {
	u := fmt.Sprintf("/v2/apps/%s/gateways/", appID)

	req := api.GatewayCreateRequest{Name: name, Port: port, Protocol: protocol}

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

// Delete removes a gateway or listener of gateway from an app.
func Delete(c *drycc.Client, appID string, name string, port int, protocol string) error {
	u := fmt.Sprintf("/v2/apps/%s/gateways/", appID)

	req := api.GatewayRemoveRequest{Name: name, Port: port, Protocol: protocol}

	body, err := json.Marshal(req)

	if err != nil {
		return err
	}

	res, err := c.Request("DELETE", u, body)
	if err == nil {
		res.Body.Close()
	}
	return err
}
