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
    "values": [{
        "name": "NEW_URL2", 
        "value": "http://localhost:8080/", 
        "group": "global"
      },
      {
        "name": "NEW_URL", 
        "value": "http://localhost:8080", 
        "ptype": "web"
      }
    ],
    "limits": {
      "web": "std1.xlarge.c1m1"
    },
    "tags": {
      "web": {
        "test": "tests"
      }
    },
    "registry": {
      "web": {
        "username": "bob"
      }
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
    "values": [{
		"name":  "NEW_URL2",
		"value": "http://localhost:8080/",
		"group": "global"
	},
	{
		"name":  "NEW_URL",
		"value": "http://localhost:8080",
		"ptype": "web"
	}],
	"limits": {
	  "web": "std1.xlarge.c1m1"
	},
    "tags": {
	  "web": {
        "test": "tests"
	  }
    },
    "registry": {
      "web": {
        "username": "bob"
      }
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
    "values": [],
    "limits": {},
    "tags": {},
	"registry": {},
    "created": "2014-01-01T00:00:00UTC",
    "updated": "2014-01-01T00:00:00UTC",
    "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
}
`

const configSetRefsFixture string = `
{
    "owner": "test",
    "app": "setrefs-test",
    "values": [],
	"values_refs":{"web":["myconfig1"]},
    "limits": {},
    "tags": {},
	"registry": {},
    "created": "2014-01-01T00:00:00UTC",
    "updated": "2014-01-01T00:00:00UTC",
    "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf76"
}
`

const (
	configSetExpected     string = `{"values":[{"group":"global","name":"NEW_URL2","value":"http://localhost:8080/"},{"ptype":"web","name":"NEW_URL","value":"http://localhost:8080"}],"limits":{"web":"std1.xlarge.c1m1"},"tags":{"web":{"test":"tests"}},"registry":{"web":{"username":"bob"}}}`
	configUnsetExpected   string = `{"values":[{"group":"global","name":"TEST","value":""}],"limits":{"web":null},"tags":{"web":{"test":null}},"registry":{"web":{"username":null}}}`
	configSetRefsExpected string = `{"values_refs":{"web":["myconfig1"]}}`
)

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

	if req.URL.Path == "/v2/apps/setrefs-test/config/" && req.Method == "POST" {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != configSetRefsExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", configSetRefsExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(configSetRefsFixture))
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
		Values: []api.ConfigValue{
			{
				Group: "global",
				ConfigVar: api.ConfigVar{
					Name:  "NEW_URL2",
					Value: "http://localhost:8080/",
				},
			},
			{
				Ptype: "web",
				ConfigVar: api.ConfigVar{
					Name:  "NEW_URL",
					Value: "http://localhost:8080",
				},
			},
		},
		Limits: map[string]interface{}{
			"web": "std1.xlarge.c1m1",
		},
		Tags: map[string]api.ConfigTags{
			"web": {
				"test": "tests",
			},
		},
		Registry: map[string]map[string]interface{}{
			"web": {
				"username": "bob",
			},
		},
		Created: "2014-01-01T00:00:00UTC",
		Updated: "2014-01-01T00:00:00UTC",
		UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	}

	configVars := api.Config{
		Values: []api.ConfigValue{
			{
				Group: "global",
				ConfigVar: api.ConfigVar{
					Name:  "NEW_URL2",
					Value: "http://localhost:8080/",
				},
			},
			{
				Ptype: "web",
				ConfigVar: api.ConfigVar{
					Name:  "NEW_URL",
					Value: "http://localhost:8080",
				},
			},
		},
		Limits: map[string]interface{}{
			"web": "std1.xlarge.c1m1",
		},
		Tags: map[string]api.ConfigTags{
			"web": {
				"test": "tests",
			},
		},
		Registry: map[string]map[string]interface{}{
			"web": {
				"username": "bob",
			},
		},
	}

	actual, err := Set(drycc, "example-go", configVars, true)
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
		Values:   []api.ConfigValue{},
		Limits:   map[string]interface{}{},
		Tags:     map[string]api.ConfigTags{},
		Registry: map[string]map[string]interface{}{},
		Created:  "2014-01-01T00:00:00UTC",
		Updated:  "2014-01-01T00:00:00UTC",
		UUID:     "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	}

	configVars := api.Config{
		Values: []api.ConfigValue{
			{
				Group: "global",
				ConfigVar: api.ConfigVar{
					Name:  "TEST",
					Value: "",
				},
			},
		},
		Limits: map[string]interface{}{
			"web": nil,
		},
		Tags: map[string]api.ConfigTags{
			"web": {"test": nil},
		},
		Registry: map[string]map[string]interface{}{
			"web": {
				"username": nil,
			},
		},
	}

	actual, err := Set(drycc, "unset-test", configVars, true)
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
		Values: []api.ConfigValue{
			{
				Group: "global",
				ConfigVar: api.ConfigVar{
					Name:  "NEW_URL2",
					Value: "http://localhost:8080/",
				},
			},
			{
				Ptype: "web",
				ConfigVar: api.ConfigVar{
					Name:  "NEW_URL",
					Value: "http://localhost:8080",
				},
			},
		},
		Limits: map[string]interface{}{
			"web": "std1.xlarge.c1m1",
		},
		Tags: map[string]api.ConfigTags{
			"web": {"test": "tests"},
		},
		Registry: map[string]map[string]interface{}{
			"web": {
				"username": "bob",
			},
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
	url := actual.Values[0].Value
	if url != "" {
		if !reflect.DeepEqual(url, "http://localhost:8080/") {
			t.Errorf("Expected %v, Got %v", "http://localhost:8080/", url)
		}
	} else {
		t.Errorf("version not found %v", actual)
	}
}

func TestConfigRefs(t *testing.T) {
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
		App:    "setrefs-test",
		Values: []api.ConfigValue{},
		ValuesRefs: map[string][]string{
			"web": {
				"myconfig1",
			},
		},
		Limits:   map[string]interface{}{},
		Tags:     map[string]api.ConfigTags{},
		Registry: map[string]map[string]interface{}{},
		Created:  "2014-01-01T00:00:00UTC",
		Updated:  "2014-01-01T00:00:00UTC",
		UUID:     "de1bf5b5-4a72-4f94-a10c-d2a3741cdf76",
	}

	configVars := api.Config{
		ValuesRefs: map[string][]string{
			"web": {
				"myconfig1",
			},
		},
	}

	actual, err := Set(drycc, "setrefs-test", configVars, true)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}
