package tls

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
	"github.com/stretchr/testify/require"
)

const (
	tlsDisabledFixture string = `{
	"uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
	"app": "foo",
	"owner": "test",
	"created": "2016-08-22T17:40:16Z",
	"updated": "2016-08-22T17:40:16Z",
	"https_enforced": false,
	"certs_auto_enabled": false,
	"issuer": {}
}`
	tlsEnabledFixture string = `{
	"uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
	"app": "foo",
	"owner": "test",
	"created": "2016-08-22T17:40:16Z",
	"updated": "2016-08-22T17:40:16Z",
	"https_enforced": true,
	"certs_auto_enabled": null
}`
	issuerFixture string = `{
    "uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
    "app": "foo",
    "owner": "test",
    "created": "2016-08-22T17:40:16Z",
    "updated": "2016-08-22T17:40:16Z",
    "https_enforced": null,
    "certs_auto_enabled": null,
    "issuer": {
        "email":"anonymous@cert-manager.io",
        "server":"https://acme-v02.api.letsencrypt.org/directory",
        "key_id":"keyID",
        "key_secret":"keySecret"
    }
}`

	tlsEnableExpected   string = `{"https_enforced":true}`
	tlsDisableExpected  string = `{"https_enforced":false}`
	tlsCertsAutoEnabled string = `{"certs_auto_enabled":true}`
	tlsCertsAutoFixture string = `{
	"uuid": "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
	"app": "foo",
	"owner": "test",
	"created": "2016-08-22T17:40:16Z",
	"updated": "2016-08-22T17:40:16Z",
	"https_enforced": null,
	"certs_auto_enabled": true,
	"issuer": {
		"email":"anonymous@cert-manager.io",
		"server":"https://acme-v02.api.letsencrypt.org/directory",
		"key_id":"keyID",
		"key_secret":"keySecret"
	},
	"events": [{"name": "foo", "kind": "Issuer", "time": "2024-04-08T01:14:49Z", "type": "Ready", "status": "True", "message": "ready message"}]
}`
	issuerExpected string = `{"issuer":{"email":"anonymous@cert-manager.io","server":"https://acme-v02.api.letsencrypt.org/directory","key_id":"keyID","key_secret":"keySecret"}}`
)

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	if req.URL.Path == "/v2/apps/foo/tls/" && req.Method == "GET" {
		res.Write([]byte(tlsDisabledFixture))
		return
	}

	if req.URL.Path == "/v2/apps/foo/tls/" && req.Method == "POST" {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) == tlsEnableExpected {
			res.WriteHeader(http.StatusCreated)
			res.Write([]byte(tlsEnabledFixture))
			return
		} else if string(body) == tlsDisableExpected {
			res.WriteHeader(http.StatusCreated)
			res.Write([]byte(tlsDisableExpected))
			return
		} else if string(body) == issuerExpected {
			res.WriteHeader(http.StatusCreated)
			res.Write([]byte(issuerFixture))
			return
		} else if string(body) == tlsCertsAutoEnabled {
			res.WriteHeader(http.StatusCreated)
			res.Write([]byte(tlsCertsAutoFixture))
			return
		}
		fmt.Printf("Expected '%s', %s or '%s', Got '%s'\n",
			tlsEnableExpected,
			tlsDisableExpected,
			issuerExpected,
			body)
		res.WriteHeader(http.StatusInternalServerError)
		res.Write(nil)
		return
	}

	fmt.Printf("Unrecongized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

type badJSONFakeHTTPServer struct{}

func (badJSONFakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	if req.URL.Path == "/v2/apps/foo/tls/" && req.Method == "GET" {
		res.Write([]byte(tlsDisabledFixture))
		return
	}

	if req.URL.Path == "/v2/apps/foo/tls/" && req.Method == "POST" {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) == tlsEnableExpected {
			res.WriteHeader(http.StatusCreated)
			res.Write([]byte(tlsEnabledFixture + "blarg"))
			return
		} else if string(body) == tlsDisableExpected {
			res.WriteHeader(http.StatusCreated)
			res.Write([]byte(tlsDisableExpected + "blarg"))
			return
		} else if string(body) == issuerExpected {
			res.WriteHeader(http.StatusCreated)
			res.Write([]byte(issuerFixture + "blarg"))
			return
		}
		fmt.Printf("Expected '%s' or '%s', Got '%s'\n",
			tlsEnableExpected,
			tlsDisableExpected,
			body)
		res.WriteHeader(http.StatusInternalServerError)
		res.Write(nil)
		return
	}

	fmt.Printf("Unrecongized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestTLSInfo(t *testing.T) {
	t.Parallel()

	expected := api.TLS{
		Created:          "2016-08-22T17:40:16Z",
		Updated:          "2016-08-22T17:40:16Z",
		App:              "foo",
		Owner:            "test",
		UUID:             "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
		HTTPSEnforced:    new(bool),
		CertsAutoEnabled: new(bool),
		Issuer:           new(api.Issuer),
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	dClient, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := Info(dClient, "foo")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}

	// now test with bad JSON in response, expecting command to return an error
	badHandler := badJSONFakeHTTPServer{}
	badServer := httptest.NewServer(badHandler)
	defer badServer.Close()

	dClient, err = drycc.New(false, badServer.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	if _, err = Info(dClient, "foo"); err != nil {
		t.Errorf("Expected Info() with poorly JSON response to fail")
	}
}

func TestTLSEnable(t *testing.T) {
	t.Parallel()

	b := true
	expected := api.TLS{
		Created:       "2016-08-22T17:40:16Z",
		Updated:       "2016-08-22T17:40:16Z",
		App:           "foo",
		Owner:         "test",
		UUID:          "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
		HTTPSEnforced: &b,
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	dClient, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := EnableHTTPSEnforced(dClient, "foo")
	if err != nil {
		t.Fatal(err)
	}

	if expected.String() != actual.String() {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}

	// now test with bad JSON in response, expecting command to return an error
	badHandler := badJSONFakeHTTPServer{}
	badServer := httptest.NewServer(badHandler)
	defer badServer.Close()

	dClient, err = drycc.New(false, badServer.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	if _, err = EnableHTTPSEnforced(dClient, "foo"); err != nil {
		t.Errorf("Expected Enable() with poorly JSON response to fail")
	}
}

func TestTLSDisable(t *testing.T) {
	t.Parallel()

	b := false
	expected := api.TLS{
		Created:       "2016-08-22T17:40:16Z",
		Updated:       "2016-08-22T17:40:16Z",
		App:           "foo",
		Owner:         "test",
		UUID:          "c4aed81c-d1ca-4ff1-ab89-d2151264e1a3",
		HTTPSEnforced: &b,
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	dClient, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := DisableHTTPSEnforced(dClient, "foo")
	if err != nil {
		t.Fatal(err)
	}

	if expected.String() != actual.String() {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}

	// now test with bad JSON in response, expecting command to return an error
	badHandler := badJSONFakeHTTPServer{}
	badServer := httptest.NewServer(badHandler)
	defer badServer.Close()

	dClient, err = drycc.New(false, badServer.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	if _, err = DisableHTTPSEnforced(dClient, "foo"); err != nil {
		t.Errorf("Expected Disable() with poorly JSON response to fail")
	}
}

func TestEnableCertsAutoEnabled(t *testing.T) {
	t.Parallel()
	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	dClient, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	tls, err := EnableCertsAutoEnabled(dClient, "foo")
	if err != nil {
		t.Fatal(err)
	}

	expected := []api.Event{
		{
			"name":    "foo",
			"kind":    "Issuer",
			"time":    "2024-04-08T01:14:49Z",
			"type":    "Ready",
			"status":  "True",
			"message": "ready message",
		},
	}
	actual := tls.Events
	a, _ := json.Marshal(expected)
	b, _ := json.Marshal(actual)
	require.JSONEq(t, string(a), string(b), "Expected %v, Got %v", expected, actual)
}

func TestAddIssuer(t *testing.T) {
	t.Parallel()
	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	dClient, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	_, err = AddCertsIssuer(dClient, "foo", "anonymous@cert-manager.io", "https://acme-v02.api.letsencrypt.org/directory", "keyID", "keySecret")
	if err != nil {
		t.Fatal(err)
	}
}
