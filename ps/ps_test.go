package ps

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync"
	"testing"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
	"github.com/drycc/controller-sdk-go/pkg/time"
	"golang.org/x/net/websocket"
)

const processesFixture string = `
{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
            "release": "v2",
            "type": "web",
            "name": "example-go-v2-web-45678",
            "state": "up",
            "started": "2016-02-13T00:47:52"
        }
    ]
}`

const scaleExpected string = `{"web":2}`

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	if req.URL.Path == "/v2/apps/example-go/pods/" && req.Method == "GET" {
		res.Write([]byte(processesFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/pods/restart/" && req.Method == "POST" {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	if req.URL.Path == "/v2/apps/example-go/pods/web/restart/" && req.Method == "POST" {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	if req.URL.Path == "/v2/apps/example-go/scale/" && req.Method == "POST" {
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

func TestProcessesList(t *testing.T) {
	t.Parallel()

	started := time.Time{}
	started.UnmarshalText([]byte("2016-02-13T00:47:52"))
	expected := api.PodsList{
		{
			Release: "v2",
			Type:    "web",
			Name:    "example-go-v2-web-45678",
			State:   "up",
			Started: started,
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

func TestExec(t *testing.T) {
	var once sync.Once
	var addr string

	once.Do(func() {
		http.Handle(
			"/v2/apps/example-go/pods/example-go-web-111/exec/",
			websocket.Handler(func(conn *websocket.Conn) {
				io.Copy(conn, conn)
			}),
		)
		server := httptest.NewServer(nil)
		addr = server.Listener.Addr().String()
		log.Print("Test WebSocket server listening on ", addr)
	})

	drycc, err := drycc.New(false, addr, "abc")
	if err != nil {
		t.Fatal(err)
	}

	expected := api.Command{
		Tty:     true,
		Stdin:   true,
		Command: []string{"/bin/sh"},
	}

	conn, err := Exec(drycc, "example-go", "example-go-web-111", expected)
	if err != nil {
		t.Fatal(err)
	}
	var actual api.Command
	err = websocket.JSON.Receive(conn, &actual)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected: %v, Got %v", expected, actual)
	}
}

func TestPodLogs(t *testing.T) {
	var once sync.Once
	var addr string

	once.Do(func() {
		http.Handle(
			"/v2/apps/example-go/pods/example-go-web-111/logs/",
			websocket.Handler(func(conn *websocket.Conn) {
				io.Copy(conn, conn)
			}),
		)
		server := httptest.NewServer(nil)
		addr = server.Listener.Addr().String()
		log.Print("Test WebSocket server listening on ", addr)
	})

	drycc, err := drycc.New(false, addr, "abc")
	if err != nil {
		t.Fatal(err)
	}

	expected := api.PodLogsRequest{
		Lines:     100,
		Follow:    true,
		Container: "runner",
	}

	conn, err := Logs(drycc, "example-go", "example-go-web-111", expected)
	if err != nil {
		t.Fatal(err)
	}
	var actual api.PodLogsRequest
	err = websocket.JSON.Receive(conn, &actual)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected: %v, Got %v", expected, actual)
	}
}

type testExpected struct {
	Name     string
	Type     string
	Expected api.PodsList
}

func TestAppsRestart(t *testing.T) {
	t.Parallel()

	started := time.Time{}
	started.UnmarshalText([]byte("2016-02-13T00:47:52"))
	tests := []testExpected{
		{
			Type: "",
		},
		{
			Type: "web",
		},
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	for _, test := range tests {
		err := Restart(drycc, "example-go", test.Type)

		if err != nil {
			t.Error(err)
		}
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

func TestByType(t *testing.T) {
	t.Parallel()

	started := time.Time{}
	started.UnmarshalText([]byte("2016-02-13T00:47:52"))

	expected := api.PodTypes{
		{
			Type: "abc",
			PodsList: api.PodsList{
				{Type: "abc", Name: "1", Started: started},
				{Type: "abc", Name: "2", Started: started},
				{Type: "abc", Name: "3", Started: started},
			},
		},
		{
			Type: "web",
			PodsList: api.PodsList{
				{Type: "web", Name: "test1", Started: started},
				{Type: "web", Name: "test2", Started: started},
				{Type: "web", Name: "test3", Started: started},
			},
		},
		{
			Type: "worker",
			PodsList: api.PodsList{
				{Type: "worker", Name: "a", Started: started},
				{Type: "worker", Name: "b", Started: started},
				{Type: "worker", Name: "c", Started: started},
			},
		},
	}

	input := api.PodsList{
		{Type: "worker", Name: "c", Started: started},
		{Type: "abc", Name: "2", Started: started},
		{Type: "worker", Name: "b", Started: started},
		{Type: "web", Name: "test1", Started: started},
		{Type: "web", Name: "test3", Started: started},
		{Type: "abc", Name: "1", Started: started},
		{Type: "worker", Name: "a", Started: started},
		{Type: "abc", Name: "3", Started: started},
		{Type: "web", Name: "test2", Started: started},
	}

	actual := ByType(input)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected: %v, Got %v", expected, actual)
	}
}
