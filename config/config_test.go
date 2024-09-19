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

const configFixtureV1 string = `
{
    "owner": "test",
    "app": "example-go",
    "values": {
      "TEST": "testing",
      "FOO": "bar"
    },
    "typed_values": {
		"web": {
		  "PORT": "9000"
		}
	},
	"limits": {
	  "web": "std1.xlarge.c1m1"
	},
    "tags": {
	  "web": {
        "test": "tests"
	  }
    },
    "registry": {
      "username": "bob"
    },
    "created": "2014-01-01T00:00:00UTC",
    "updated": "2014-01-01T00:00:00UTC",
    "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
}
`

const configFixtureV2 string = `
{
    "owner": "test",
    "app": "example-go",
    "values": {
      "TEST": "testing",
      "FOO": "bar",
	  "VERSION": "2"
    },
    "typed_values": {
		"web": {
		  "PORT": "9000"
		}
	},
	"limits": {
	  "web": "std1.xlarge.c1m1"
	},
    "tags": {
	  "web": {
        "test": "tests"
	  }
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
	"typed_values": {"web":{"PORT":null}},
    "limits": {},
    "tags": {},
	"registry": {},
    "created": "2014-01-01T00:00:00UTC",
    "updated": "2014-01-01T00:00:00UTC",
    "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
}
`

const configSetExpected string = `{"values":{"FOO":"bar","TEST":"testing"},"typed_values":{"web":{"PORT":"9000"}},"limits":{"web":"std1.xlarge.c1m1"},"tags":{"web":{"test":"tests"}},"registry":{"username":"bob"}}`
const configUnsetExpected string = `{"values":{"FOO":null,"TEST":null},"typed_values":{"web":{"PORT":null}},"limits":{"web":null},"tags":{"web":{"test":null}},"registry":{"username":null}}`

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
		res.Write([]byte(configFixtureV1))
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

		if req.URL.RawQuery == "version=v2" {
			res.Write([]byte(configFixtureV2))
		} else {
			res.Write([]byte(configFixtureV1))
		}
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
		Values: api.ConfigValues{
			"TEST": "testing",
			"FOO":  "bar",
		},
		TypedValues: map[string]api.ConfigValues{
			"web": {"PORT": "9000"},
		},
		Limits: map[string]interface{}{
			"web": "std1.xlarge.c1m1",
		},
		Tags: map[string]api.ConfigTags{
			"web": {
				"test": "tests",
			},
		},
		Registry: map[string]interface{}{
			"username": "bob",
		},
		Created: "2014-01-01T00:00:00UTC",
		Updated: "2014-01-01T00:00:00UTC",
		UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	}

	configVars := api.Config{
		Values: api.ConfigValues{
			"TEST": "testing",
			"FOO":  "bar",
		},
		TypedValues: map[string]api.ConfigValues{
			"web": {"PORT": "9000"},
		},
		Limits: map[string]interface{}{
			"web": "std1.xlarge.c1m1",
		},
		Tags: map[string]api.ConfigTags{
			"web": {
				"test": "tests",
			},
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
		Owner:  "test",
		App:    "unset-test",
		Values: map[string]interface{}{},
		TypedValues: map[string]api.ConfigValues{
			"web": {"PORT": nil},
		},
		Limits:   map[string]interface{}{},
		Tags:     map[string]api.ConfigTags{},
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
		TypedValues: map[string]api.ConfigValues{
			"web": {"PORT": nil},
		},
		Limits: map[string]interface{}{
			"web": nil,
		},
		Tags: map[string]api.ConfigTags{
			"web": {"test": nil},
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
		TypedValues: map[string]api.ConfigValues{
			"web": {"PORT": "9000"},
		},
		Limits: map[string]interface{}{
			"web": "std1.xlarge.c1m1",
		},
		Tags: map[string]api.ConfigTags{
			"web": {"test": "tests"},
		},
		Registry: map[string]interface{}{
			"username": "bob",
		},
		Created: "2014-01-01T00:00:00UTC",
		Updated: "2014-01-01T00:00:00UTC",
		UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	}

	actual, err := List(drycc, "example-go", -1)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}

	actual, err = List(drycc, "example-go", 2)
	if err != nil {
		t.Error(err)
	}
	if version, ok := actual.Values["VERSION"]; ok {
		if !reflect.DeepEqual(version, "2") {
			t.Errorf("Expected %v, Got %v", "2", version)
		}
	} else {
		t.Errorf("version not found %v", actual)
	}
}
