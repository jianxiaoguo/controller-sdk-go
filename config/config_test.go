package config

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

const configFixture string = `
{
    "owner": "test",
    "app": "example-go",
    "values": {
      "TEST": "testing",
      "FOO": "bar"
    },
	"limits": {
	  "web": "std1.xlarge.c1m1"
	},
    "tags": {
      "test": "tests"
    },
    "registry": {
      "username": "bob"
    },
    "created": "2014-01-01T00:00:00UTC",
    "updated": "2014-01-01T00:00:00UTC",
    "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
}
`

const configUnsetFixture string = `
{
    "owner": "test",
    "app": "unset-test",
    "values": {},
    "limits": {},
    "tags": {},
	"registry": {},
    "created": "2014-01-01T00:00:00UTC",
    "updated": "2014-01-01T00:00:00UTC",
    "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
}
`

const configSetExpected string = `{"values":{"FOO":"bar","TEST":"testing"},"limits":{"web":"std1.xlarge.c1m1"},"tags":{"test":"tests"},"registry":{"username":"bob"}}`
const configUnsetExpected string = `{"values":{"FOO":null,"TEST":null},"limits":{"web":null},"tags":{"test":null},"registry":{"username":null}}`

type fakeHTTPServer struct{}

func (f *fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	if req.URL.Path == "/v2/apps/example-go/config/" && req.Method == "POST" {
		body, err := io.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != configSetExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", configSetExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(configFixture))
		return
	}

	if req.URL.Path == "/v2/apps/unset-test/config/" && req.Method == "POST" {
		body, err := io.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != configUnsetExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", configUnsetExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(configUnsetFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/config/" && req.Method == "GET" {
		res.Write([]byte(configFixture))
		return
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestConfigSet(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	expected := api.Config{
		Owner: "test",
		App:   "example-go",
		Values: map[string]interface{}{
			"TEST": "testing",
			"FOO":  "bar",
		},
		Limits: map[string]interface{}{
			"web": "std1.xlarge.c1m1",
		},
		Tags: map[string]interface{}{
			"test": "tests",
		},
		Registry: map[string]interface{}{
			"username": "bob",
		},
		Created: "2014-01-01T00:00:00UTC",
		Updated: "2014-01-01T00:00:00UTC",
		UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	}

	configVars := api.Config{
		Values: map[string]interface{}{
			"TEST": "testing",
			"FOO":  "bar",
		},
		Limits: map[string]interface{}{
			"web": "std1.xlarge.c1m1",
		},
		Tags: map[string]interface{}{
			"test": "tests",
		},
		Registry: map[string]interface{}{
			"username": "bob",
		},
	}

	actual, err := Set(drycc, "example-go", configVars)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}

func TestConfigUnset(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	expected := api.Config{
		Owner:    "test",
		App:      "unset-test",
		Values:   map[string]interface{}{},
		Limits:   map[string]interface{}{},
		Tags:     map[string]interface{}{},
		Registry: map[string]interface{}{},
		Created:  "2014-01-01T00:00:00UTC",
		Updated:  "2014-01-01T00:00:00UTC",
		UUID:     "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	}

	configVars := api.Config{
		Values: map[string]interface{}{
			"TEST": nil,
			"FOO":  nil,
		},
		Limits: map[string]interface{}{
			"web": nil,
		},
		Tags: map[string]interface{}{
			"test": nil,
		},
		Registry: map[string]interface{}{
			"username": nil,
		},
	}

	actual, err := Set(drycc, "unset-test", configVars)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}

func TestConfigList(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	expected := api.Config{
		Owner: "test",
		App:   "example-go",
		Values: map[string]interface{}{
			"TEST": "testing",
			"FOO":  "bar",
		},
		Limits: map[string]interface{}{
			"web": "std1.xlarge.c1m1",
		},
		Tags: map[string]interface{}{
			"test": "tests",
		},
		Registry: map[string]interface{}{
			"username": "bob",
		},
		Created: "2014-01-01T00:00:00UTC",
		Updated: "2014-01-01T00:00:00UTC",
		UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	}

	actual, err := List(drycc, "example-go")

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}
