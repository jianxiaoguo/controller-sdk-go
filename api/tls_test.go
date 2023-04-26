package api

import (
	"strings"
	"testing"
)

func TestTLSString(t *testing.T) {
	tls := &TLS{}

	expected := "--- HTTPS Enforced: not set\n--- Certs Auto: not set\n--- Issuer: not set"

	if strings.TrimSpace(tls.String()) != expected {
		t.Errorf("Expected:\n\n%s\n\nGot:\n\n%s", expected, tls.String())
	}

	tls = NewTLS()

	expected = "--- HTTPS Enforced: false\n--- Certs Auto: false\n--- Issuer: not set"

	if strings.TrimSpace(tls.String()) != expected {
		t.Errorf("Expected:\n\n%s\n\nGot:\n\n%s", expected, tls.String())
	}

	b := true
	tls.HTTPSEnforced = &b

	expected = "--- HTTPS Enforced: true\n--- Certs Auto: false\n--- Issuer: not set"

	if strings.TrimSpace(tls.String()) != expected {
		t.Errorf("Expected:\n\n%s\n\nGot:\n\n%s", expected, tls.String())
	}

	issuer := Issuer{
		Email:     "anonymous@cert-manager.io",
		Server:    "https://acme-v02.api.letsencrypt.org/directory",
		KeyID:     "",
		KeySecret: "",
	}
	tls.Issuer = &issuer

	expected = `--- HTTPS Enforced: true
--- Certs Auto: false
--- Issuer: 
email: anonymous@cert-manager.io
server: https://acme-v02.api.letsencrypt.org/directory
key-id: 
key-secret:`

	if strings.TrimSpace(tls.String()) != expected {
		t.Errorf("Expected:\n\n%s\n\nGot:\n\n%s", expected, tls.String())
	}
}
