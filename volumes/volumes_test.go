package volumes

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

const volumeCreateExpected string = `{"name":"myvolume","size":"500G"}`

const volumeCreateFixture string = `
{
	"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	"owner": "test",
	"app": "example-go",
	"name": "myvolume",
	"size": "500G",
	"path": {},
    "type": "csi",
    "parameters": {},
	"created": "2020-08-26T00:00:00UTC",
	"updated": "2020-08-26T00:00:00UTC"
}
`

const volumeExpandExpected string = `{"name":"myvolume","size":"500G"}`

const volumeExpandFixture string = `
{
	"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	"owner": "test",
	"app": "example-go",
	"name": "myvolume",
	"size": "500G",
	"path": {},
    "type": "csi",
    "parameters": {},
	"created": "2020-08-26T00:00:00UTC",
	"updated": "2020-08-26T00:00:00UTC"
}
`

const volumesFixture string = `
{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
		{
			"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
			"owner": "test",
			"app": "example-go",
			"name": "myvolume",
			"size": "500G",
			"path": {},
			"type": "csi",
			"parameters": {},
			"created": "2020-08-26T00:00:00UTC",
			"updated": "2020-08-26T00:00:00UTC"
		}
    ]
}
`

const volumeGetFixture string = `
{
	"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	"owner": "test",
	"app": "example-go",
	"name": "myvolume",
	"size": "500G",
	"path": {},
	"type": "csi",
	"parameters": {},
	"created": "2020-08-26T00:00:00UTC",
	"updated": "2020-08-26T00:00:00UTC"
}
`

const volumeMountFixture string = `
{
	"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	"owner": "test",
	"app": "example-go",
	"name": "myvolume",
	"size": "500G",
	"path": {
		"cmd":  "/data/cmd1",
		"web": "/data/web1"
	},
	"type": "csi",
    "parameters": {},
	"created": "2020-08-26T00:00:00UTC",
	"updated": "2020-08-26T00:00:00UTC"
}
`

const volumeUnmountFixture string = `
{
	"uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	"owner": "test",
	"app": "unmount-test",
	"name": "myvolume",
	"size": "500G",
	"path": {},
	"type": "csi",
    "parameters": {},
	"created": "2020-08-26T00:00:00UTC",
	"updated": "2020-08-26T00:00:00UTC"
}
`

const volumeMountExpected string = `{"path":{"cmd":"/data/cmd1","web":"/data/web1"}}`
const volumeUnmountExpected string = `{"path":{"cmd":null,"web":null}}`

type fakeHTTPServer struct{}

func (f *fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	// Create
	if req.URL.Path == "/v2/apps/example-go/volumes/" && req.Method == "POST" {
		body, err := io.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != volumeCreateExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", volumeCreateExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(volumeCreateFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/volumes/myvolume/" && req.Method == "GET" {
		res.Write([]byte(volumeGetFixture))
		return
	}
	// Client
	if strings.Contains(req.URL.Path, "/v2/apps/example-go/volumes/myvolume/files/") {
		if req.Method == "GET" {
			if req.URL.Query().Get("action") == "get" {
				res.Header().Add("Content-Type", "application/octet-stream")
				res.Write([]byte(volumeFileContentExpected))
				return
			} else if req.URL.Query().Get("action") == "list" {
				res.Header().Add("Content-Type", "application/json")
				res.Write([]byte("[]"))
				return
			}
		} else if req.Method == "POST" {
			if err := req.ParseMultipartForm(1024 * 1024); err != nil {
				return
			}
			for _, tmpFiles := range req.MultipartForm.File {
				for _, tmpFile := range tmpFiles {
					srcFile, err := tmpFile.Open()
					if err != nil {
						return
					}
					body, err := io.ReadAll(srcFile)
					if err != nil {
						return
					}
					if string(body) != volumeFileContentExpected {
						fmt.Printf("Expected '%s', Got '%s'\n", volumeFileContentExpected, body)
						res.WriteHeader(http.StatusInternalServerError)
						res.Write(nil)
						return
					}
				}
			}
			return
		} else if req.Method == "DELETE" {
			res.WriteHeader(http.StatusNoContent)
			res.Write(nil)
			return
		}
	}

	// Expand
	if req.URL.Path == "/v2/apps/example-go/volumes/myvolume/" && req.Method == "PATCH" {
		body, err := io.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != volumeExpandExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", volumeExpandExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusOK)
		res.Write([]byte(volumeExpandFixture))
		return
	}

	// Delete
	if req.URL.Path == "/v2/apps/example-go/volumes/myvolume/" && req.Method == "DELETE" {
		res.WriteHeader(http.StatusNoContent)
		//res.Write([]byte(volumeMountFixture))
		return
	}

	// List
	if req.URL.Path == "/v2/apps/example-go/volumes/" && req.Method == "GET" {
		res.Write([]byte(volumesFixture))
		return
	}

	//　Mount
	if req.URL.Path == "/v2/apps/example-go/volumes/myvolume/path/" && req.Method == "PATCH" {
		body, err := io.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != volumeMountExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", volumeMountExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(volumeMountFixture))
		return
	}

	//　Unmount
	if req.URL.Path == "/v2/apps/unmount-test/volumes/myvolume/path/" && req.Method == "PATCH" {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != volumeUnmountExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", volumeUnmountExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(volumeUnmountFixture))
		return
	}
	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestVolumesCreate(t *testing.T) {
	t.Parallel()

	expected := api.Volume{
		UUID:       "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
		Owner:      "test",
		App:        "example-go",
		Name:       "myvolume",
		Size:       "500G",
		Path:       map[string]interface{}{},
		Type:       "csi",
		Parameters: map[string]interface{}{},
		Created:    "2020-08-26T00:00:00UTC",
		Updated:    "2020-08-26T00:00:00UTC",
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}
	volume := api.Volume{
		Name: "myvolume",
		Size: "500G",
	}
	actual, err := Create(drycc, "example-go", volume)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}

func TestVolumesExpand(t *testing.T) {
	t.Parallel()

	expected := api.Volume{
		UUID:       "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
		Owner:      "test",
		App:        "example-go",
		Name:       "myvolume",
		Size:       "500G",
		Path:       map[string]interface{}{},
		Type:       "csi",
		Parameters: map[string]interface{}{},
		Created:    "2020-08-26T00:00:00UTC",
		Updated:    "2020-08-26T00:00:00UTC",
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}
	volume := api.Volume{
		Name: "myvolume",
		Size: "500G",
	}
	actual, err := Expand(drycc, "example-go", volume)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}

func TestVolumesDelete(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	if err = Delete(drycc, "example-go", "myvolume"); err != nil {
		t.Fatal(err)
	}
}

func TestVolumesList(t *testing.T) {
	t.Parallel()

	expected := api.Volumes{
		{
			UUID:       "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
			App:        "example-go",
			Owner:      "test",
			Name:       "myvolume",
			Path:       map[string]interface{}{},
			Size:       "500G",
			Type:       "csi",
			Parameters: map[string]interface{}{},
			Created:    "2020-08-26T00:00:00UTC",
			Updated:    "2020-08-26T00:00:00UTC",
		},
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
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

func TestVolumeGet(t *testing.T) {
	expected := api.Volume{
		UUID:       "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
		App:        "example-go",
		Owner:      "test",
		Name:       "myvolume",
		Path:       map[string]interface{}{},
		Size:       "500G",
		Type:       "csi",
		Parameters: map[string]interface{}{},
		Created:    "2020-08-26T00:00:00UTC",
		Updated:    "2020-08-26T00:00:00UTC",
	}
	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := Get(drycc, "example-go", "myvolume")

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}

func TestVolumeMount(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	expected := api.Volume{
		Name:  "myvolume",
		Owner: "test",
		App:   "example-go",
		Path: map[string]interface{}{
			"cmd": "/data/cmd1",
			"web": "/data/web1",
		},
		Size:       "500G",
		Type:       "csi",
		Parameters: map[string]interface{}{},
		Created:    "2020-08-26T00:00:00UTC",
		Updated:    "2020-08-26T00:00:00UTC",
		UUID:       "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	}

	volumeVars := api.Volume{
		Path: map[string]interface{}{
			"cmd": "/data/cmd1",
			"web": "/data/web1",
		},
	}
	actual, err := Mount(drycc, "example-go", "myvolume", volumeVars)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}

func TestVolumeUnmount(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	expected := api.Volume{
		Name:       "myvolume",
		Owner:      "test",
		App:        "unmount-test",
		Path:       map[string]interface{}{},
		Size:       "500G",
		Type:       "csi",
		Parameters: map[string]interface{}{},
		Created:    "2020-08-26T00:00:00UTC",
		Updated:    "2020-08-26T00:00:00UTC",
		UUID:       "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	}

	volumeVars := api.Volume{
		Path: map[string]interface{}{
			"cmd": nil,
			"web": nil,
		},
	}
	actual, err := Mount(drycc, "unmount-test", "myvolume", volumeVars)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}
