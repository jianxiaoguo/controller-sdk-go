package events

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

const podEventsFixture string = `
{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
            "reason": "Scheduled",
            "message": "Successfully assigned example-go/example-go-web-6b44dbd6c8-h89cg to node1",
            "created": "2024-07-03T16:28:00"
        }
    ]
}`

const ptypeEventsFixture string = `
{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
            "reason": "ScalingReplicaSet",
            "message": "Scaled up replica set example-go-web-6b44dbd6c8 to 2 from 1",
            "created": "2024-07-03T16:28:00"
        }
    ]
}`

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)
	if req.URL.Path == "/v2/apps/example-go/events/" && req.Method == "GET" && req.URL.RawQuery == "ptype=example-go-web" {
		res.Write([]byte(ptypeEventsFixture))
		return
	}
	if req.URL.Path == "/v2/apps/example-go/events/" && req.Method == "GET" && req.URL.RawQuery == "pod_name=example-go-web-6b44dbd6c8-h89cg" {
		res.Write([]byte(podEventsFixture))
		return
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestEvents(t *testing.T) {
	t.Parallel()

	created := "2024-07-03T16:28:00"
	podExpected := api.AppEvents{
		{
			Reason:  "Scheduled",
			Message: "Successfully assigned example-go/example-go-web-6b44dbd6c8-h89cg to node1",
			Created: created,
		},
	}
	ptypeExpected := api.AppEvents{
		{
			Reason:  "ScalingReplicaSet",
			Message: "Scaled up replica set example-go-web-6b44dbd6c8 to 2 from 1",
			Created: created,
		},
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, _, err := ListPodEvents(drycc, "example-go", "example-go-web-6b44dbd6c8-h89cg", 100)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(podExpected, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", podExpected, actual))
	}

	actual, _, err = ListPtypeEvents(drycc, "example-go", "web", 100)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(ptypeExpected, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", ptypeExpected, actual))
	}
}
