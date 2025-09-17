package api

import (
	"fmt"
)

// Event represents a TLS event as a map of string key-value pairs.
type Event = map[string]string

// TLS is the structure of an app's TLS settings.
type TLS struct {
	// Owner is the app owner. It cannot be updated with TLS.Set(). See app.Transfer().
	Owner string `json:"owner,omitempty"`
	// App is the app the tls settings apply to and cannot be updated.
	App string `json:"app,omitempty"`
	// Created is the time that the TLS settings was created and cannot be updated.
	Created string `json:"created,omitempty"`
	// Updated is the last time the TLS settings was changed and cannot be updated.
	Updated string `json:"updated,omitempty"`
	// UUID is a unique string reflecting the TLS settings in its current state.
	// It changes every time the TLS settings is changed and cannot be updated.
	UUID string `json:"uuid,omitempty"`
	// HTTPSEnforced determines if the router should enable or disable https-only requests.
	HTTPSEnforced *bool `json:"https_enforced,omitempty"`
	// Use ACME to automatically generate certificates if CertsAuto enable
	CertsAutoEnabled *bool   `json:"certs_auto_enabled,omitempty"`
	Issuer           *Issuer `json:"issuer,omitempty"`
	Events           []Event `json:"events,omitempty"`
}

// Issuer is the structure of POST /v2/app/<app id>/tls/.
type Issuer struct {
	Email     string `json:"email"`
	Server    string `json:"server"`
	KeyID     string `json:"key_id"`
	KeySecret string `json:"key_secret"`
}

// NewTLS creates a new TLS object with fields properly zeroed
func NewTLS() *TLS {
	return &TLS{
		HTTPSEnforced:    new(bool),
		CertsAutoEnabled: new(bool),
	}
}

func (t TLS) String() string {
	tpl := `--- HTTPS Enforced: %s
--- Certs Auto: %s
--- Issuer: %s`
	issuerTpl := `
email: %s
server: %s
key-id: %s
key-secret: %s
`
	httpsEnforced := "not set"
	if t.HTTPSEnforced != nil {
		httpsEnforced = fmt.Sprintf("%t", *(t.HTTPSEnforced))
	}
	certsAutoEnabled := "not set"
	if t.CertsAutoEnabled != nil {
		certsAutoEnabled = fmt.Sprintf("%t", *(t.CertsAutoEnabled))
	}
	issuer := "not set"
	if t.Issuer != nil {
		issuer = fmt.Sprintf(issuerTpl, t.Issuer.Email, t.Issuer.Server, t.Issuer.KeyID, t.Issuer.KeySecret)
	}
	return fmt.Sprintf(tpl, httpsEnforced, certsAutoEnabled, issuer)
}
