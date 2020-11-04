// Package allowlist provides methods for managing an app's allowlisted IP's.
package allowlist

import (
	"encoding/json"
	"fmt"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// List IP's allowlisted for an app.
func List(c *drycc.Client, appID string) (api.Allowlist, error) {
	u := fmt.Sprintf("/v2/apps/%s/allowlist/", appID)
	res, reqErr := c.Request("GET", u, nil)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.Allowlist{}, reqErr
	}
	defer res.Body.Close()

	allowlist := api.Allowlist{}
	if err := json.NewDecoder(res.Body).Decode(&allowlist); err != nil {
		return api.Allowlist{}, err
	}

	return allowlist, reqErr
}

// Add adds addresses to an app's allowlist.
func Add(c *drycc.Client, appID string, addresses []string) (api.Allowlist, error) {
	u := fmt.Sprintf("/v2/apps/%s/allowlist/", appID)

	req := api.Allowlist{Addresses: addresses}
	body, err := json.Marshal(req)
	if err != nil {
		return api.Allowlist{}, err
	}
	res, reqErr := c.Request("POST", u, body)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.Allowlist{}, reqErr
	}
	defer res.Body.Close()

	d := api.Allowlist{}
	if err = json.NewDecoder(res.Body).Decode(&d); err != nil {
		return api.Allowlist{}, err
	}

	return d, reqErr
}

// Delete removes addresses from an app's allowlist.
func Delete(c *drycc.Client, appID string, addresses []string) error {
	u := fmt.Sprintf("/v2/apps/%s/allowlist/", appID)

	req := api.Allowlist{Addresses: addresses}
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	_, reqErr := c.Request("DELETE", u, body)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return reqErr
	}
	return nil
}
