// Package certs manages SSL keys and certificates on the drycc platform
package certs

import (
	"encoding/json"
	"fmt"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// List lists certificates added to drycc.
func List(c *drycc.Client, appID string, results int) ([]api.Cert, int, error) {
	u := fmt.Sprintf("/v2/apps/%s/certs/", appID)
	body, count, reqErr := c.LimitedRequest(u, results)

	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.Cert{}, -1, reqErr
	}

	var res []api.Cert
	if err := json.Unmarshal([]byte(body), &res); err != nil {
		return []api.Cert{}, -1, err
	}

	return res, count, reqErr
}

// New creates a new certificate.
// Certificates are created independently from apps and are applied on a per domain basis.
// So to enable SSL for an app with the domain test.com, you would first create the certificate,
// then use the attach method to attach test.com to the certificate.
func New(c *drycc.Client, appID string, cert string, key string, name string) (api.Cert, error) {
	req := api.CertCreateRequest{Certificate: cert, Key: key, Name: name}
	reqBody, err := json.Marshal(req)
	if err != nil {
		return api.Cert{}, err
	}
	u := fmt.Sprintf("/v2/apps/%s/certs/", appID)
	res, reqErr := c.Request("POST", u, reqBody)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.Cert{}, reqErr
	}
	defer res.Body.Close()

	resCert := api.Cert{}
	if err = json.NewDecoder(res.Body).Decode(&resCert); err != nil {
		return api.Cert{}, err
	}

	return resCert, reqErr
}

// Get retrieves information about a certificate
func Get(c *drycc.Client, appID string, name string) (api.Cert, error) {
	url := fmt.Sprintf("/v2/apps/%s/certs/%s", appID, name)
	res, reqErr := c.Request("GET", url, nil)
	if reqErr != nil {
		return api.Cert{}, reqErr
	}
	defer res.Body.Close()

	resCert := api.Cert{}
	if err := json.NewDecoder(res.Body).Decode(&resCert); err != nil {
		return api.Cert{}, err
	}

	return resCert, reqErr
}

// Delete removes a certificate.
func Delete(c *drycc.Client, appID string, name string) error {
	url := fmt.Sprintf("/v2/apps/%s/certs/%s", appID, name)
	res, err := c.Request("DELETE", url, nil)
	if err == nil {
		res.Body.Close()
	}
	return err
}

// Attach adds a domain to a certificate.
func Attach(c *drycc.Client, appID string, name string, domain string) error {
	req := api.CertAttachRequest{Domain: domain}
	reqBody, err := json.Marshal(req)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("/v2/apps/%s/certs/%s/domain/", appID, name)
	res, err := c.Request("POST", url, reqBody)
	if err == nil {
		res.Body.Close()
	}
	return err
}

// Detach removes a domain from a certificate.
func Detach(c *drycc.Client, appID string, name string, domain string) error {
	url := fmt.Sprintf("/v2/apps/%s/certs/%s/domain/%s", appID, name, domain)
	res, err := c.Request("DELETE", url, nil)
	if err == nil {
		res.Body.Close()
	}
	return err
}
