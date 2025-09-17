package allowlist

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

const allowlistFixture string = `
{
    "addresses": ["1.2.3.4", "0.0.0.0/0"]
}`

const allowlistCreateExpected string = `{"addresses":["1.2.3.4","0.0.0.0/0"]}`

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	if req.URL.Path == "/v2/apps/example-go/allowlist/" && req.Method == "GET" {
		res.Write([]byte(allowlistFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/allowlist/" && req.Method == "POST" {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != allowlistCreateExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", allowlistCreateExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(allowlistFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/allowlist/" && req.Method == "DELETE" {
		res.WriteHeader(http.StatusNoContent)
		res.Write([]byte(allowlistFixture))
		return
	}

	if req.URL.Path == "/v2/apps/invalidjson-test/allowlist/" && req.Method == "POST" {
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(`"addresses": "test"`))
		return
	}

	if req.URL.Path == "/v2/apps/invalidjson-test/allowlist/" && req.Method == "GET" {
		res.Write([]byte(`"addresses": "test"`))
		return
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestAllowlistList(t *testing.T) {
	t.Parallel()

	expected := api.Allowlist{
		Addresses: []string{"1.2.3.4", "0.0.0.0/0"},
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

func TestAllowlistAdd(t *testing.T) {
	t.Parallel()

	expected := api.Allowlist{
		Addresses: []string{"1.2.3.4", "0.0.0.0/0"},
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := Add(drycc, "example-go", []string{"1.2.3.4", "0.0.0.0/0"})
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}

func TestAllowlistRemove(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	if err = Delete(drycc, "example-go", []string{"1.2.3.4"}); err != nil {
		t.Fatal(err)
	}
}

func TestAppSettingsInvalidJson(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	_, err = List(drycc, "invalidjson-test")
	expected := "json: cannot unmarshal string into Go value of type api.Allowlist"
	if err == nil || !reflect.DeepEqual(expected, err.Error()) {
		t.Errorf("Expected %v, Got %v", expected, err)
	}

	_, err = Add(drycc, "invalidjson-test", []string{"1.2.3.4"})
	if err == nil || !reflect.DeepEqual(expected, err.Error()) {
		t.Errorf("Expected %v, Got %v", expected, err)
	}
}
