package routes

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

const routesFixture string = `
{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
            "app": "example-go",
            "created": "2023-04-19T00:00:00UTC",
            "owner": "test",
            "updated": "2023-04-19T00:00:00UTC",
            "name": "example-go",
            "ptype": "web",
            "kind": "HTTPRoute",
            "port": 80,
            "parent_refs": [
                {
                    "name": "example-go",
                    "port": 80
                }
            ]
        }
    ]
}`

const routerulesFixture string = `
[
  {
    "backendRefs": [
      {
        "kind": "Service",
        "name": "py3django3",
        "port": 80
      }
    ]
  }
]`

const routerulesSetFixture string = `"[{\"backendRefs\": [{\"kind\": \"Service\",\"name\": \"py3django3\",\"port\": 80}]}]"`

const routeCreateExpected string = `{"name":"example-go","ptype":"web","port":80,"kind":"HTTPRoute"}`

const routeRulesSetExpected string = `[{"backendRefs": [{"kind": "Service","name": "py3django3","port": 80}]}]`

const routeAttachExpected string = `{"port":80,"gateway":"example-go"}`

const routeDetachExpected string = `{"port":80,"gateway":"example-go"}`

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	if req.URL.Path == "/v2/apps/example-go/routes/" && req.Method == "GET" {
		res.Write([]byte(routesFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/routes/example-go/rules/" && req.Method == "GET" {
		res.Write([]byte(routerulesFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/routes/example-go/rules/" && req.Method == "PUT" {
		body, err := io.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}
		if string(body) != routerulesSetFixture {
			fmt.Printf("Expected '%s', Got '%s'\n", routerulesSetFixture, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusNoContent)
		return
	}

	if req.URL.Path == "/v2/apps/example-go/routes/" && req.Method == "POST" {
		body, err := io.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}
		if string(body) != routeCreateExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", routeCreateExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusCreated)
		return
	}

	if req.URL.Path == "/v2/apps/example-go/routes/example-go/attach/" && req.Method == "PATCH" {
		body, err := io.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}
		if string(body) != routeAttachExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", routeAttachExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusNoContent)
		return
	}

	if req.URL.Path == "/v2/apps/example-go/routes/example-go/detach/" && req.Method == "PATCH" {
		body, err := io.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}
		if string(body) != routeDetachExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", routeDetachExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusNoContent)
		return
	}

	if req.URL.Path == "/v2/apps/example-go/routes/example-go/" && req.Method == "DELETE" {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestRoutesList(t *testing.T) {
	t.Parallel()

	expected := api.Routes{
		{
			App:     "example-go",
			Created: "2023-04-19T00:00:00UTC",
			Name:    "example-go",
			Owner:   "test",
			Updated: "2023-04-19T00:00:00UTC",
			Ptype:   "web",
			Kind:    "HTTPRoute",
			Port:    80,
			ParentRefs: []api.ParentRef{
				{
					Name: "example-go",
					Port: 80,
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

func TestRouteGet(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := GetRule(drycc, "example-go", "example-go")

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(routerulesFixture, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", routerulesFixture, actual))
	}
}

func TestRouteSet(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	err = SetRule(drycc, "example-go", "example-go", routeRulesSetExpected)

	if err != nil {
		t.Fatal(err)
	}
}

func TestRoutesAdd(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	err = New(drycc, "example-go", "example-go", "web", "HTTPRoute", 80)

	if err != nil {
		t.Fatal(err)
	}
}

func TestRoutesAttachGateway(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	err = AttachGateway(drycc, "example-go", "example-go", 80, "example-go")

	if err != nil {
		t.Fatal(err)
	}
}

func TestRoutesDetachGateway(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	err = DetachGateway(drycc, "example-go", "example-go", 80, "example-go")

	if err != nil {
		t.Fatal(err)
	}
}

func TestRoutesRemove(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	if err = Delete(drycc, "example-go", "example-go"); err != nil {
		t.Fatal(err)
	}
}
