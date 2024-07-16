package pts

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

const ptypesFixture string = `
{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
			"name": "example-go-web",
			"release": "v1",
			"ready": "1/1",
			"up_to_date": 1,
            "available_replicas": 1,
            "started": "2024-07-03T16:28:00"
        }
    ]
}`

const ptypeStateFixture string = `
{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [{
        "container": "example-go-web",
        "image": "registry.drycc.cc/base/base",
        "command": ["bash", "-c"],
        "args": ["sleep", "3600s"],
        "readiness_probe": {
            "exec": {
                "command": ["ls", "-la"]
            },
            "failureThreshold": 3,
            "initialDelaySeconds": 50,
            "periodSeconds": 10,
            "successThreshold": 1,
            "timeoutSeconds": 50
        },
        "limits": {
            "cpu": "1",
            "memory": "2Gi"
        },
        "volume_mounts": [
            {
                "mountPath": "/data",
                "name": "myvolume"
            }
        ],
		"node_selector": ["kubernetes.io/os=linux"]
    }]
}`

const scaleExpected string = `{"web":2}`

const restartExpected string = `{"types":"web,worker"}`

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	if req.URL.Path == "/v2/apps/example-go/ptypes/" && req.Method == "GET" {
		res.Write([]byte(ptypesFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/ptypes/example-go-web/describe/" && req.Method == "GET" {
		res.Write([]byte(ptypeStateFixture))
		return
	}
	if req.URL.Path == "/v2/apps/example-go/ptypes/restart/" && req.Method == "POST" {
		body, err := io.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}
		if string(body) != restartExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", restartExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}
		res.WriteHeader(http.StatusNoContent)
		return
	}

	if req.URL.Path == "/v2/apps/example-go/ptypes/scale/" && req.Method == "POST" {
		body, err := io.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != scaleExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", scaleExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusNoContent)
		res.Write(nil)
		return
	}
	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestPtype(t *testing.T) {
	t.Parallel()

	started := "2024-07-03T16:28:00"
	expected := api.Ptypes{
		{
			Name:              "example-go-web",
			Release:           "v1",
			Ready:             "1/1",
			UpToDate:          1,
			AvailableReplicas: 1,
			Started:           started,
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

func TestDescribe(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}
	actual, _, err := Describe(drycc, "example-go", "web", 100)
	if err != nil {
		t.Error(err)
	}
	expected := api.PtypeStates{
		{
			Container: "example-go-web",
			Image:     "registry.drycc.cc/base/base",
			Command:   []string{"bash", "-c"},
			Args:      []string{"sleep", "3600s"},
			ReadinessProbe: api.Healthcheck{
				Exec: &api.ExecProbe{
					Command: []string{"ls", "-la"},
				},
				FailureThreshold:    3,
				InitialDelaySeconds: 50,
				PeriodSeconds:       10,
				SuccessThreshold:    1,
				TimeoutSeconds:      50,
			},
			Limits: map[string]string{
				"cpu":    "1",
				"memory": "2Gi",
			},
			VolumeMounts: []api.VolumeMount{
				{
					Name:      "myvolume",
					MountPath: "/data",
				},
			},
			NodeSelector: []string{"kubernetes.io/os=linux"},
		},
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Error(fmt.Errorf("Expected %v, Got %v", actual, expected))
	}
}

func TestAppsRestart(t *testing.T) {
	t.Parallel()

	types := map[string]string{
		"types": "web,worker",
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	err = Restart(drycc, "example-go", types)

	if err != nil {
		t.Error(err)
	}
}

func TestScale(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	if err = Scale(drycc, "example-go", map[string]int{"web": 2}); err != nil {
		t.Fatal(err)
	}
}
