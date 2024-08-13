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

const codenamesFixture string = `
{
	"results": [
		{"codename": "use_app", "description": "Can use app"},
		{"codename": "use_cert", "description": "Can use cert"}
	],
	"count": 2
}`

const listUserPermFixture string = `
{
	"results": [
		{"id": 1, "codename": "use_app", "uniqueid": "autotest-app", "username": "foo"},
		{"id": 2, "codename": "use_cert", "uniqueid": "autotest-cert-1", "username": "foo"}
	],
	"count": 2
}`

const createUserPermExpected = `{"codename":"use_app","uniqueid":"autotest","username":"foo"}`

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	if req.URL.Path == "/v2/perms/codes/" && req.Method == "GET" {
		res.Write([]byte(codenamesFixture))
		return
	}

	if req.URL.Path == "/v2/perms/rules/" && req.Method == "GET" {
		res.Write([]byte(listUserPermFixture))
		return
	}

	if req.URL.Path == "/v2/perms/rules/" && req.Method == "POST" {
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

	if req.URL.Path == "/v2/perms/rules/1/" && req.Method == "DELETE" {
		res.WriteHeader(http.StatusNoContent)
		res.Write(nil)
		return
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestNames(t *testing.T) {
	t.Parallel()

	expected := []api.PermCodeResponse{
		{Codename: "use_app", Description: "Can use app"},
		{Codename: "use_cert", Description: "Can use cert"},
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, _, err := Codes(drycc, 300)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}

func TestList(t *testing.T) {
	t.Parallel()

	expected := []api.UserPermResponse{
		{ID: 1, Codename: "use_app", Uniqueid: "autotest-app", Username: "foo"},
		{ID: 2, Codename: "use_cert", Uniqueid: "autotest-cert-1", Username: "foo"},
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, _, err := List(drycc, "", 300)

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

	if err = Create(drycc, "use_app", "autotest", "foo"); err != nil {
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

	if err = Delete(drycc, "1"); err != nil {
		t.Fatal(err)
	}
}
