// Package tls provides methods for managing tls configuration for apps.
package tls

import (
	"encoding/json"
	"fmt"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// Info displays an app's tls config.
func Info(c *drycc.Client, app string) (api.TLS, error) {
	u := fmt.Sprintf("/v2/apps/%s/tls/", app)

	res, reqErr := c.Request("GET", u, nil)
	if reqErr != nil {
		return api.TLS{}, reqErr
	}
	defer res.Body.Close()

	tls := api.TLS{}
	if err := json.NewDecoder(res.Body).Decode(&tls); err != nil {
		return api.TLS{}, err
	}

	return tls, reqErr
}

// changeTLS enables the router to enforce https-only requests to the application.
func changeTLS(c *drycc.Client, app string, httpsEnforced, certsAutoEnabled *bool, issuer *api.Issuer) (api.TLS, error) {
	t := api.NewTLS()
	t.HTTPSEnforced = httpsEnforced
	t.CertsAutoEnabled = certsAutoEnabled
	t.Issuer = issuer
	body, err := json.Marshal(t)

	if err != nil {
		return api.TLS{}, err
	}

	u := fmt.Sprintf("/v2/apps/%s/tls/", app)

	res, reqErr := c.Request("POST", u, body)
	if reqErr != nil {
		return api.TLS{}, reqErr
	}
	defer res.Body.Close()

	newTLS := api.TLS{}
	if err = json.NewDecoder(res.Body).Decode(&newTLS); err != nil {
		return api.TLS{}, err
	}

	return newTLS, reqErr
}

// EnableHTTPSEnforced enables the router to enforce https-only requests to the application.
func EnableHTTPSEnforced(c *drycc.Client, app string) (api.TLS, error) {
	b := true
	return changeTLS(c, app, &b, nil, nil)
}

// DisableHTTPSEnforced disables the router from enforcing https-only requests to the application.
func DisableHTTPSEnforced(c *drycc.Client, app string) (api.TLS, error) {
	b := false
	return changeTLS(c, app, &b, nil, nil)
}

// EnableCertsAutoEnabled enables ACME to automatically generate certificates.
func EnableCertsAutoEnabled(c *drycc.Client, app string) (api.TLS, error) {
	b := true
	return changeTLS(c, app, nil, &b, nil)
}

// DisableCertsAutoEnabled disables ACME to automatically generate certificates.
func DisableCertsAutoEnabled(c *drycc.Client, app string) (api.TLS, error) {
	b := false
	return changeTLS(c, app, nil, &b, nil)
}

// AddCertsIssuer disables ACME to automatically generate certificates.
func AddCertsIssuer(c *drycc.Client, app string, email string, server string, keyID string, keySecret string) (api.TLS, error) {
	issuer := api.Issuer{Email: email, Server: server, KeyID: keyID, KeySecret: keySecret}
	return changeTLS(c, app, nil, nil, &issuer)
}
