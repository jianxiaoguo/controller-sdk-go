package services

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

const servicesFixture string = `
{
    "services": [
        {
            "domain": "example-go.example-go.svc.cluster.local",
            "ptype": "web",
            "ports": [
                {
                    "name": "example-go-web-udp-5000",
                    "port": 5000,
                    "protocol": "UDP",
                    "targetPort": 5000
                },
                {
                    "name": "example-go-web-tcp-2379",
                    "port": 2379,
                    "protocol": "TCP",
                    "targetPort": 2379
                }
            ]
        },
        {
            "domain": "example-go-worker.example-go.svc.cluster.local",
            "ptype": "worker",
            "ports": [
                {
                    "name": "example-go-worker-tcp-5000",
                    "port": 5000,
                    "protocol": "TCP",
                    "targetPort": 5000
                }
            ]
        }
    ]
}`

const serviceCreateExpected string = `{"ptype":"web","port":5000,"protocol":"UDP","target_port":5000}`
const serviceDeleteExpected string = `{"ptype":"web","port":5000,"protocol":"UDP"}`

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	if req.URL.Path == "/v2/apps/example-go/services/" && req.Method == "GET" {
		res.Write([]byte(servicesFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/services/" && req.Method == "POST" {
		body, err := io.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != serviceCreateExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", serviceCreateExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusCreated)
		return
	}

	if req.URL.Path == "/v2/apps/example-go/services/" && req.Method == "DELETE" {

		body, err := io.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != serviceDeleteExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", serviceDeleteExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusNoContent)
		res.Write([]byte(servicesFixture))
		return
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestServicesList(t *testing.T) {
	t.Parallel()

	expected := api.Services{
		{
			Domain: "example-go.example-go.svc.cluster.local",
			Ptype:  "web",
			Ports: []api.Port{
				{
					Name:       "example-go-web-udp-5000",
					Port:       5000,
					Protocol:   "UDP",
					TargetPort: 5000,
				},
				{
					Name:       "example-go-web-tcp-2379",
					Port:       2379,
					Protocol:   "TCP",
					TargetPort: 2379,
				},
			},
		},
		{
			Domain: "example-go-worker.example-go.svc.cluster.local",
			Ptype:  "worker",
			Ports: []api.Port{
				{
					Name:       "example-go-worker-tcp-5000",
					Port:       5000,
					Protocol:   "TCP",
					TargetPort: 5000,
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

	actual, err := List(drycc, "example-go")

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}

func TestServicesAdd(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	err = New(drycc, "example-go", "web", 5000, "UDP", 5000)

	if err != nil {
		t.Fatal(err)
	}
}

func TestServicesRemove(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	if err = Delete(drycc, "example-go", "web", "UDP", 5000); err != nil {
		t.Fatal(err)
	}
}
