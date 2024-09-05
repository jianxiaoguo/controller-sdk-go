package perms

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

const listUserPermFixture string = `
{
	"results": [
		{"app": "example-go", "username": "foo", "permissions": ["view", "change", "delete"]},
		{"app": "example-go", "username": "bar", "permissions": ["view", "change", "delete"]}
	],
	"count": 2
}`

const createUserPermExpected = `{"username":"foo","permissions":"view,change,delete"}`
const updateUserPermExpected = `{"username":"foo","permissions":"view"}`

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	if req.URL.Path == "/v2/apps/example-go/perms/" && req.Method == "GET" {
		res.Write([]byte(listUserPermFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/perms/" && req.Method == "POST" {
		body, err := io.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != createUserPermExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", createUserPermExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write(nil)
		return
	}

	if req.URL.Path == "/v2/apps/example-go/perms/foo/" && req.Method == "PUT" {
		body, err := io.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != updateUserPermExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", updateUserPermExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write(nil)
		return
	}

	if req.URL.Path == "/v2/apps/example-go/perms/foo/" && req.Method == "DELETE" {
		res.WriteHeader(http.StatusNoContent)
		res.Write(nil)
		return
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestList(t *testing.T) {
	t.Parallel()

	expected := []api.UserPermResponse{
		{App: "example-go", Username: "foo", Permissions: []string{"view", "change", "delete"}},
		{App: "example-go", Username: "bar", Permissions: []string{"view", "change", "delete"}},
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, _, err := List(drycc, "example-go", 300)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}

func TestCreate(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	if err = Create(drycc, "example-go", "foo", "view,change,delete"); err != nil {
		t.Fatal(err)
	}
}

func TestUpdate(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	if err = Update(drycc, "example-go", "foo", "view"); err != nil {
		t.Fatal(err)
	}
}

func TestDelete(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	if err = Delete(drycc, "example-go", "foo"); err != nil {
		t.Fatal(err)
	}
}
