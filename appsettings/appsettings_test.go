package appsettings

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

const appSettingsFixture string = `
{
    "owner": "test",
    "app": "example-go",
    "routable": true,
    "autodeploy": true,
    "autorollback": true,
    "allowlist": ["1.2.3.4", "0.0.0.0/0"],
    "autoscale": {"cmd": {"min": 3, "max": 8, "cpu_percent": 40}},
    "label": {"git_repo": "https://github.com/drycc/controller-sdk-go", "team" : "drycc"},
    "created": "2014-01-01T00:00:00UTC",
    "updated": "2014-01-01T00:00:00UTC",
    "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
}
`

const appSettingsUnsetFixture string = `
{
    "owner": "test",
    "app": "unset-test",
    "routable": true,
    "autodeploy": true,
    "autorollback": true,
    "allowlist": ["1.2.3.4", "0.0.0.0/0"],
    "autoscale": {"cmd": {"min": 3, "max": 8, "cpu_percent": 40}},
    "label": {"git_repo": "https://github.com/drycc/controller-sdk-go", "team" : "drycc"},
    "created": "2014-01-01T00:00:00UTC",
    "updated": "2014-01-01T00:00:00UTC",
    "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
}
`
const appSettingsSetExpected string = `{"routable":true,"allowlist":["1.2.3.4","0.0.0.0/0"],"autodeploy":true,"autorollback":true,"autoscale":{"cmd":{"min":3,"max":8,"cpu_percent":40}},"label":{"git_repo":"https://github.com/drycc/controller-sdk-go","team":"drycc"}}`
const appSettingsUnsetExpected string = `{"routable":true,"allowlist":["1.2.3.4","0.0.0.0/0"],"autodeploy":true,"autorollback":true,"autoscale":{"cmd":{"min":3,"max":8,"cpu_percent":40}},"label":{"git_repo":"https://github.com/drycc/controller-sdk-go","team":"drycc"}}`

type fakeHTTPServer struct{}

func (f *fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	if req.URL.Path == "/v2/apps/example-go/settings/" && req.Method == "POST" {
		body, err := io.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != appSettingsSetExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", appSettingsSetExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(appSettingsFixture))
		return
	}

	if req.URL.Path == "/v2/apps/unset-test/settings/" && req.Method == "POST" {
		body, err := io.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != appSettingsUnsetExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", appSettingsUnsetExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(appSettingsUnsetFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/settings/" && req.Method == "GET" {
		res.Write([]byte(appSettingsFixture))
		return
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestAppSettingsSet(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	expected := api.AppSettings{
		Owner:        "test",
		App:          "example-go",
		Routable:     api.NewRoutable(),
		Autodeploy:   api.NewAutodeploy(),
		Autorollback: api.NewAutorollback(),
		Allowlist:    []string{"1.2.3.4", "0.0.0.0/0"},
		Autoscale: map[string]*api.Autoscale{
			"cmd": {
				Min:        3,
				Max:        8,
				CPUPercent: 40,
			},
		},
		Label: map[string]interface{}{
			"git_repo": "https://github.com/drycc/controller-sdk-go",
			"team":     "drycc",
		},
		Created: "2014-01-01T00:00:00UTC",
		Updated: "2014-01-01T00:00:00UTC",
		UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	}

	appSettingsVars := api.AppSettings{
		Routable:     api.NewRoutable(),
		Autodeploy:   api.NewAutodeploy(),
		Autorollback: api.NewAutorollback(),
		Allowlist:    []string{"1.2.3.4", "0.0.0.0/0"},
		Autoscale: map[string]*api.Autoscale{
			"cmd": {
				Min:        3,
				Max:        8,
				CPUPercent: 40,
			},
		},
		Label: map[string]interface{}{
			"git_repo": "https://github.com/drycc/controller-sdk-go",
			"team":     "drycc",
		},
	}

	actual, err := Set(drycc, "example-go", appSettingsVars)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}

func TestAppSettingsUnset(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	expected := api.AppSettings{
		Owner:        "test",
		App:          "unset-test",
		Routable:     api.NewRoutable(),
		Autodeploy:   api.NewAutodeploy(),
		Autorollback: api.NewAutorollback(),
		Allowlist:    []string{"1.2.3.4", "0.0.0.0/0"},
		Autoscale: map[string]*api.Autoscale{
			"cmd": {
				Min:        3,
				Max:        8,
				CPUPercent: 40,
			},
		},
		Label: map[string]interface{}{
			"git_repo": "https://github.com/drycc/controller-sdk-go",
			"team":     "drycc",
		},
		Created: "2014-01-01T00:00:00UTC",
		Updated: "2014-01-01T00:00:00UTC",
		UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	}

	appSettingsVars := api.AppSettings{
		Routable:     api.NewRoutable(),
		Autodeploy:   api.NewAutodeploy(),
		Autorollback: api.NewAutorollback(),
		Allowlist:    []string{"1.2.3.4", "0.0.0.0/0"},
		Autoscale: map[string]*api.Autoscale{
			"cmd": {
				Min:        3,
				Max:        8,
				CPUPercent: 40,
			},
		},
		Label: map[string]interface{}{
			"git_repo": "https://github.com/drycc/controller-sdk-go",
			"team":     "drycc",
		},
	}

	actual, err := Set(drycc, "unset-test", appSettingsVars)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}

func TestAppSettingsList(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	expected := api.AppSettings{
		Owner:        "test",
		App:          "example-go",
		Routable:     api.NewRoutable(),
		Autodeploy:   api.NewAutodeploy(),
		Autorollback: api.NewAutorollback(),
		Allowlist:    []string{"1.2.3.4", "0.0.0.0/0"},
		Autoscale: map[string]*api.Autoscale{
			"cmd": {
				Min:        3,
				Max:        8,
				CPUPercent: 40,
			},
		},
		Label: map[string]interface{}{
			"git_repo": "https://github.com/drycc/controller-sdk-go",
			"team":     "drycc",
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
