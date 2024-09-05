// Package domains provides methods for managing an app's domains.
package domains

import (
	"encoding/json"
	"fmt"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// List domains registered with an app.
func List(c *drycc.Client, appID string, results int) (api.Domains, int, error) {
	u := fmt.Sprintf("/v2/apps/%s/domains/", appID)
	body, count, reqErr := c.LimitedRequest(u, results)

	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.Domain{}, -1, reqErr
	}

	var domains []api.Domain
	if err := json.Unmarshal([]byte(body), &domains); err != nil {
		return []api.Domain{}, -1, err
	}

	return domains, count, reqErr
}

// New adds a domain to an app.
func New(c *drycc.Client, appID, domain, Ptype string) (api.Domain, error) {
	u := fmt.Sprintf("/v2/apps/%s/domains/", appID)

	req := api.DomainCreateRequest{Domain: domain, Ptype: Ptype}

	body, err := json.Marshal(req)

	if err != nil {
		return api.Domain{}, err
	}

	res, reqErr := c.Request("POST", u, body)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.Domain{}, reqErr
	}
	defer res.Body.Close()

	d := api.Domain{}
	if err = json.NewDecoder(res.Body).Decode(&d); err != nil {
		return api.Domain{}, err
	}

	return d, reqErr
}

// Delete removes a domain from an app.
func Delete(c *drycc.Client, appID string, domain string) error {
	u := fmt.Sprintf("/v2/apps/%s/domains/%s", appID, domain)
	res, err := c.Request("DELETE", u, nil)
	if err == nil {
		res.Body.Close()
	}
	return err
}
