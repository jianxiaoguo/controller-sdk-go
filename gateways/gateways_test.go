package gateways

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

const gatewaysFixture string = `
{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
            "app": "example-go",
            "name": "example-go",
            "created": "2023-04-19T00:00:00UTC",
            "owner": "test",
            "updated": "2023-04-19T00:00:00UTC",
            "listeners": [
                {
                    "name": "example-go-80-http",
                    "port": 80,
                    "protocol": "HTTP",
                    "allowedRoutes": {"namespaces": {"from": "All"}}
                },
                {
                    "name": "example-go-443-https",
                    "port": 443,
                    "protocol": "HTTPS",
                    "allowedRoutes": {"namespaces": {"from": "All"}}
                }
            ],
            "addresses": [
                {
                    "type": "IPAddress",
                    "value": "172.22.108.207"
                }
            ]
        }
    ]
}`

const (
	gatewayCreateExpected string = `{"name":"example-go","port":443,"protocol":"HTTPS"}`
	gatewayRemoveExpected string = `{"name":"example-go","port":443,"protocol":"HTTPS"}`
)

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	if req.URL.Path == "/v2/apps/example-go/gateways/" && req.Method == "GET" {
		res.Write([]byte(gatewaysFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/gateways/" && req.Method == "POST" {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}
		if string(body) != gatewayCreateExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", gatewayCreateExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusCreated)
		return
	}

	if req.URL.Path == "/v2/apps/example-go/gateways/" && req.Method == "DELETE" {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}
		if string(body) != gatewayRemoveExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", gatewayRemoveExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}
		res.WriteHeader(http.StatusNoContent)
		return
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestGatewaysList(t *testing.T) {
	t.Parallel()

	expected := api.Gateways{
		{
			App:     "example-go",
			Created: "2023-04-19T00:00:00UTC",
			Name:    "example-go",
			Owner:   "test",
			Updated: "2023-04-19T00:00:00UTC",
			Listeners: []api.Listener{
				{
					Name:          "example-go-80-http",
					Port:          80,
					Protocol:      "HTTP",
					AllowedRoutes: map[string]interface{}{"namespaces": map[string]interface{}{"from": "All"}},
				},
				{
					Name:          "example-go-443-https",
					Port:          443,
					Protocol:      "HTTPS",
					AllowedRoutes: map[string]interface{}{"namespaces": map[string]interface{}{"from": "All"}},
				},
			},
			Addresses: []api.Address{
				{
					Type:  "IPAddress",
					Value: "172.22.108.207",
				},
			},
		},
	}
	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, _, err := List(drycc, "example-go", 100)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}

func TestGatewaysAdd(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	err = New(drycc, "example-go", "example-go", 443, "HTTPS")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGatewaysRemove(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	if err = Delete(drycc, "example-go", "example-go", 443, "HTTPS"); err != nil {
		t.Fatal(err)
	}
}
