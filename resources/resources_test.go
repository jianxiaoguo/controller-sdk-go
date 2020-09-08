package resources

import (
	"fmt"
	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

const resourceCreateExpected string = `{"name":"mysql","plan":"mysql:5.6"}`

const resourcePutExpected string = `{"plan":"mysql:5.7"}`

const resourceCreateFixture string = `
{
	"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	"owner": "test",
	"app": "example-go",
	"name": "mysql",
	"plan": "mysql:5.6",
	"data": {},
	"options": {},
	"status": null,
	"binding": null,
	"created": "2020-09-08T00:00:00UTC",
	"updated": "2020-09-08T00:00:00UTC"
}
`

const resourcePutFixture string = `
{
	"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	"owner": "test",
	"app": "example-go",
	"name": "mysql",
	"plan": "mysql:5.7",
	"data": {},
	"options": {},
	"status": null,
	"binding": null,
	"created": "2020-09-08T00:00:00UTC",
	"updated": "2020-09-08T00:00:00UTC"
}
`

const resourcesFixture string = `
{
   "count": 1,
   "next": null,
   "previous": null,
   "results": [
		{
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
			"owner": "test",
			"app": "example-go",
			"name": "mysql",
			"plan": "mysql:5.6",
			"data": {},
			"options": {},
			"status": null,
			"binding": null,
			"created": "2020-09-08T00:00:00UTC",
			"updated": "2020-09-08T00:00:00UTC"
		}
   ]
}
`

const resourceFixture string = `
{
	"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	"owner": "test",
	"app": "example-go",
	"name": "mysql",
	"plan": "mysql:5.6",
	"data": {},
	"options": {},
	"status": null,
	"binding": null,
	"created": "2020-09-08T00:00:00UTC",
	"updated": "2020-09-08T00:00:00UTC"
}
`

//
const resourceBindingFixture string = `
{
	"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	"owner": "test",
	"app": "example-bind",
	"name": "mysql",
	"plan": "mysql:5.6",
	"data": {},
	"options": {},
	"status": null,
	"binding": null,
	"created": "2020-09-08T00:00:00UTC",
	"updated": "2020-09-08T00:00:00UTC"
}
`

const resourceUnbindFixture string = `
{
	"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	"owner": "test",
	"app": "example-unbind",
	"name": "mysql",
	"plan": "mysql:5.6",
	"data": {},
	"options": {},
	"status": null,
	"binding": null,
	"created": "2020-09-08T00:00:00UTC",
	"updated": "2020-09-08T00:00:00UTC"
}
`

const resourceBindExpected string = `{"bind_action":"bind"}`
const resourceUnbindExpected string = `{"bind_action":"unbind"}`

type fakeHTTPServer struct{}

func (f *fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	// Create
	if req.URL.Path == "/v2/apps/example-go/resources/" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != resourceCreateExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", resourceCreateExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(resourceCreateFixture))
		return
	}
	// List
	if req.URL.Path == "/v2/apps/example-go/resources/" && req.Method == "GET" {
		res.Write([]byte(resourcesFixture))
		return
	}

	// Delete
	if req.URL.Path == "/v2/apps/example-go/resources/mysql/" && req.Method == "DELETE" {
		res.WriteHeader(http.StatusNoContent)
		return
	}
	// Get
	if req.URL.Path == "/v2/apps/example-go/resources/mysql/" && req.Method == "GET" {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(resourceFixture))
		return
	}
	// Put
	if req.URL.Path == "/v2/apps/example-go/resources/mysql/" && req.Method == "PUT" {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != resourcePutExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", resourcePutExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(resourcePutFixture))
		return
	}
	// Patch bind
	if req.URL.Path == "/v2/apps/example-bind/resources/mysql/binding/" && req.Method == "PATCH" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != resourceBindExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", resourceBindExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(resourceBindingFixture))
		return
	}
	// Patch unbind
	if req.URL.Path == "/v2/apps/example-unbind/resources/mysql/binding/" && req.Method == "PATCH" {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != resourceUnbindExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", resourceUnbindExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(resourceUnbindFixture))
		return
	}
	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestResourcesCreate(t *testing.T) {
	t.Parallel()

	expected := api.Resource{
		UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
		Owner:   "test",
		App:     "example-go",
		Name:    "mysql",
		Plan:    "mysql:5.6",
		Status:  "",
		Binding: "",
		Data:    map[string]interface{}{},
		Options: map[string]interface{}{},
		Created: "2020-09-08T00:00:00UTC",
		Updated: "2020-09-08T00:00:00UTC",
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}
	resource := api.Resource{
		Name: "mysql",
		Plan: "mysql:5.6",
	}
	actual, err := Create(drycc, "example-go", resource)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}

//
func TestResourcesList(t *testing.T) {
	t.Parallel()

	expected := api.Resources{
		{
			UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
			Owner:   "test",
			App:     "example-go",
			Name:    "mysql",
			Plan:    "mysql:5.6",
			Status:  "",
			Binding: "",
			Data:    map[string]interface{}{},
			Options: map[string]interface{}{},
			Created: "2020-09-08T00:00:00UTC",
			Updated: "2020-09-08T00:00:00UTC",
		},
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
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

func TestResourceDelete(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	if err = Delete(drycc, "example-go", "mysql"); err != nil {
		t.Fatal(err)
	}
}

func TestResourceGet(t *testing.T) {
	t.Parallel()

	expected := api.Resource{
		UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
		Owner:   "test",
		App:     "example-go",
		Name:    "mysql",
		Plan:    "mysql:5.6",
		Status:  "",
		Binding: "",
		Data:    map[string]interface{}{},
		Options: map[string]interface{}{},
		Created: "2020-09-08T00:00:00UTC",
		Updated: "2020-09-08T00:00:00UTC",
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := Get(drycc, "example-go", "mysql")

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}

func TestResourcePut(t *testing.T) {
	t.Parallel()

	expected := api.Resource{
		UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
		Owner:   "test",
		App:     "example-go",
		Name:    "mysql",
		Plan:    "mysql:5.7",
		Status:  "",
		Binding: "",
		Data:    map[string]interface{}{},
		Options: map[string]interface{}{},
		Created: "2020-09-08T00:00:00UTC",
		Updated: "2020-09-08T00:00:00UTC",
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}
	resource := api.Resource{
		Plan: "mysql:5.7",
	}
	actual, err := Put(drycc, "example-go", "mysql", resource)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}

func TestResourceBind(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	expected := api.Resource{
		UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
		Owner:   "test",
		App:     "example-bind",
		Name:    "mysql",
		Plan:    "mysql:5.6",
		Status:  "",
		Binding: "",
		Data:    map[string]interface{}{},
		Options: map[string]interface{}{},
		Created: "2020-09-08T00:00:00UTC",
		Updated: "2020-09-08T00:00:00UTC",
	}

	resourceVars := api.Binding{
		BindAction: "bind",
	}
	actual, err := Binding(drycc, "example-bind", "mysql", resourceVars)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}

func TestResourceUnbind(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	expected := api.Resource{
		UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
		Owner:   "test",
		App:     "example-unbind",
		Name:    "mysql",
		Plan:    "mysql:5.6",
		Status:  "",
		Binding: "",
		Data:    map[string]interface{}{},
		Options: map[string]interface{}{},
		Created: "2020-09-08T00:00:00UTC",
		Updated: "2020-09-08T00:00:00UTC",
	}

	resourceVars := api.Binding{
		BindAction: "unbind",
	}
	actual, err := Binding(drycc, "example-unbind", "mysql", resourceVars)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}
