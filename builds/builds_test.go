package builds

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

const buildsFixture string = `
{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
            "app": "example-go",
            "created": "2014-01-01T00:00:00UTC",
            "dockerfile": "FROM drycc/slugrunner RUN mkdir -p /app WORKDIR /app ENTRYPOINT [\"/runner/init\"] ADD slug.tgz /app ENV GIT_SHA 060da68f654e75fac06dbedd1995d5f8ad9084db",
            "image": "example-go",
            "stack": "container",
            "owner": "test",
            "procfile": {
                "web": "example-go"
            },
			"dryccfile": {},
            "sha": "060da68f",
            "updated": "2014-01-01T00:00:00UTC",
            "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
        }
    ]
}`

const buildFixture string = `
{
    "app": "example-go",
    "created": "2014-01-01T00:00:00UTC",
    "dockerfile": "",
    "image": "drycc/example-go:latest",
    "stack": "heroku-18",
    "owner": "test",
    "procfile": {
        "web": "example-go"
    },
	"dryccfile": {},
    "sha": "",
    "updated": "2014-01-01T00:00:00UTC",
    "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
}`

const buildExpected string = `{"image":"drycc/example-go","stack":"heroku-18","procfile":{"web":"example-go"}}`

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	if req.URL.Path == "/v2/apps/example-go/builds/" && req.Method == "GET" {
		res.Write([]byte(buildsFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/builds/" && req.Method == "POST" {
		body, err := io.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != buildExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", buildExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(buildFixture))
		return
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestBuildsList(t *testing.T) {
	t.Parallel()

	expected := []api.Build{
		{
			App:        "example-go",
			Created:    "2014-01-01T00:00:00UTC",
			Dockerfile: "FROM drycc/slugrunner RUN mkdir -p /app WORKDIR /app ENTRYPOINT [\"/runner/init\"] ADD slug.tgz /app ENV GIT_SHA 060da68f654e75fac06dbedd1995d5f8ad9084db",
			Image:      "example-go",
			Stack:      "container",
			Owner:      "test",
			Procfile:   map[string]string{"web": "example-go"},
			Dryccfile:  map[string]interface{}{},
			Sha:        "060da68f",
			Updated:    "2014-01-01T00:00:00UTC",
			UUID:       "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
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

func TestBuildCreate(t *testing.T) {
	t.Parallel()

	expected := api.Build{
		App:       "example-go",
		Created:   "2014-01-01T00:00:00UTC",
		Image:     "drycc/example-go:latest",
		Stack:     "heroku-18",
		Owner:     "test",
		Procfile:  map[string]string{"web": "example-go"},
		Dryccfile: map[string]interface{}{},
		Updated:   "2014-01-01T00:00:00UTC",
		UUID:      "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	procfile := map[string]string{
		"web": "example-go",
	}

	actual, err := New(drycc, "example-go", "drycc/example-go", "heroku-18", procfile, map[string]interface{}{})

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}
