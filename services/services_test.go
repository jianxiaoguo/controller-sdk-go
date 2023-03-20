package services

import (
	"fmt"
	"io/ioutil"
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
        "procfile_type": "web",
        "port": 5000,
        "protocol": "UDP",
        "target_port": 5000
    },
    {
        "procfile_type": "worker",
        "port": 5000,
        "protocol": "TCP",
        "target_port": 5000
    }
  ]
}`

const serviceFixture string = `
{
    "procfile_type": "web",
    "port": 5000,
    "protocol": "UDP",
    "target_port": 5000
}`

const serviceCreateExpected string = `{"procfile_type":"web","port":5000,"protocol":"UDP","target_port":5000}`

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	if req.URL.Path == "/v2/apps/example-go/services/" && req.Method == "GET" {
		res.Write([]byte(servicesFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/services/" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

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
		res.Write([]byte(serviceFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/services/" && req.Method == "DELETE" {
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
			ProcfileType: "web",
			Port:         5000,
			Protocol:     "UDP",
			TargetPort:   5000,
		},
		{
			ProcfileType: "worker",
			Port:         5000,
			Protocol:     "TCP",
			TargetPort:   5000,
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

	expected := api.Service{
		ProcfileType: "web",
		Port:         5000,
		Protocol:     "UDP",
		TargetPort:   5000,
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := New(drycc, "example-go", "web", 5000, "UDP", 5000)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
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

	if err = Delete(drycc, "example-go", "web"); err != nil {
		t.Fatal(err)
	}
}
