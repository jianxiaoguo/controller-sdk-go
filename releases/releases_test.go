package releases

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

const releasesFixture string = `
{
    "count": 3,
    "next": null,
    "previous": null,
    "results": [
        {
            "app": "example-go",
            "state": "succeed",
            "build": null,
            "config": "95bd6dea-1685-4f78-a03d-fd7270b058d1",
            "created": "2014-01-01T00:00:00UTC",
            "owner": "test",
            "summary": "test created initial release",
            "exception": null,
            "conditions": [{
                "state":   "succeed",
                "action":  "pipeline",
                "ptypes":  ["web"],
                "created": "2024-08-27T08:31:36Z"
            }],
            "updated": "2014-01-01T00:00:00UTC",
            "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
            "version": 1
        }
    ]
}`

const releaseFixture string = `
{
    "app": "example-go",
    "build": null,
    "state": "succeed",
    "config": "95bd6dea-1685-4f78-a03d-fd7270b058d1",
    "created": "2014-01-01T00:00:00UTC",
    "owner": "test",
    "summary": "test created initial release",
    "exception": null,
    "conditions": [{
        "state":   "succeed",
        "action":  "pipeline",
        "ptypes":  ["web"],
        "created": "2024-08-27T08:31:36Z"
    }],
    "updated": "2014-01-01T00:00:00UTC",
    "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
    "version": 1
}
`

const deployExpected string = `{"types":"web,task"}`

const rollbackFixture string = `
{"ptypes":"web,task", "version": 5}
`
const rollbackerFixture string = `
{"ptypes":"web,task", "version": 7}
`

const rollbackExpected string = `{"version":2,"ptypes":"web,task"}`
const rollbackerExpected string = ``

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	if req.URL.Path == "/v2/apps/example-go/releases/" && req.Method == "GET" {
		res.Write([]byte(releasesFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/releases/v1/" && req.Method == "GET" {
		res.Write([]byte(releaseFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/releases/deploy/" && req.Method == "POST" {
		body, err := io.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != deployExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", deployExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusCreated)
		return
	}

	if req.URL.Path == "/v2/apps/example-go/releases/rollback/" && req.Method == "POST" {
		body, err := io.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != rollbackExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", rollbackExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(rollbackFixture))
		return
	}

	if req.URL.Path == "/v2/apps/rollbacker/releases/rollback/" && req.Method == "POST" {
		body, err := io.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != rollbackerExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", rollbackerExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(rollbackerFixture))
		return
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestReleasesList(t *testing.T) {
	t.Parallel()

	expected := []api.Release{
		{
			App:       "example-go",
			Build:     "",
			State:     "succeed",
			Config:    "95bd6dea-1685-4f78-a03d-fd7270b058d1",
			Created:   "2014-01-01T00:00:00UTC",
			Owner:     "test",
			Summary:   "test created initial release",
			Exception: "",
			Conditions: []api.Condition{
				{
					State:   "succeed",
					Action:  "pipeline",
					Ptypes:  []string{"web"},
					Created: "2024-08-27T08:31:36Z",
				},
			},
			Updated: "2014-01-01T00:00:00UTC",
			UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
			Version: 1,
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

func TestReleasesGet(t *testing.T) {
	t.Parallel()

	expected := api.Release{
		App:       "example-go",
		State:     "succeed",
		Build:     "",
		Config:    "95bd6dea-1685-4f78-a03d-fd7270b058d1",
		Created:   "2014-01-01T00:00:00UTC",
		Owner:     "test",
		Summary:   "test created initial release",
		Exception: "",
		Conditions: []api.Condition{
			{
				State:   "succeed",
				Action:  "pipeline",
				Ptypes:  []string{"web"},
				Created: "2024-08-27T08:31:36Z",
			},
		},
		Updated: "2014-01-01T00:00:00UTC",
		UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
		Version: 1,
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := Get(drycc, "example-go", 1)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}

func TestDeploy(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}
	targets := map[string]interface{}{"types": "web,task"}

	err = Deploy(drycc, "example-go", targets)

	if err != nil {
		t.Fatal(err)
	}
}

func TestRollback(t *testing.T) {
	t.Parallel()

	expected := 5

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := Rollback(drycc, "example-go", "web,task", 2)

	if err != nil {
		t.Fatal(err)
	}

	if expected != actual {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}

func TestRollbacker(t *testing.T) {
	t.Parallel()

	expected := 7

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := Rollback(drycc, "rollbacker", "web,task", -1)

	if err != nil {
		t.Fatal(err)
	}

	if expected != actual {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}
