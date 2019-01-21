package whitelist

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

const whitelistFixture string = `
{
    "addresses": ["1.2.3.4", "0.0.0.0/0"]
}`

const whitelistCreateExpected string = `{"addresses":["1.2.3.4","0.0.0.0/0"]}`

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	if req.URL.Path == "/v2/apps/example-go/whitelist/" && req.Method == "GET" {
		res.Write([]byte(whitelistFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/whitelist/" && req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != whitelistCreateExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", whitelistCreateExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(whitelistFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/whitelist/" && req.Method == "DELETE" {
		res.WriteHeader(http.StatusNoContent)
		res.Write([]byte(whitelistFixture))
		return
	}

	if req.URL.Path == "/v2/apps/invalidjson-test/whitelist/" && req.Method == "POST" {
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(`"addresses": "test"`))
		return
	}

	if req.URL.Path == "/v2/apps/invalidjson-test/whitelist/" && req.Method == "GET" {
		res.Write([]byte(`"addresses": "test"`))
		return
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestWhitelistList(t *testing.T) {
	t.Parallel()

	expected := api.Whitelist{
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

func TestWhitelistAdd(t *testing.T) {
	t.Parallel()

	expected := api.Whitelist{
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

func TestWhitelistRemove(t *testing.T) {
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
	expected := "json: cannot unmarshal string into Go value of type api.Whitelist"
	if err == nil || !reflect.DeepEqual(expected, err.Error()) {
		t.Errorf("Expected %v, Got %v", expected, err)
	}

	_, err = Add(drycc, "invalidjson-test", []string{"1.2.3.4"})
	if err == nil || !reflect.DeepEqual(expected, err.Error()) {
		t.Errorf("Expected %v, Got %v", expected, err)
	}
}
