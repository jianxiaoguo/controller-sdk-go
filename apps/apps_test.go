package apps

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
	"golang.org/x/net/websocket"
)

const appFixture string = `
{
    "created": "2014-01-01T00:00:00UTC",
    "id": "example-go",
    "owner": "test",
    "structure": {},
    "updated": "2014-01-01T00:00:00UTC",
    "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
}`

const appsFixture string = `
{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
            "created": "2014-01-01T00:00:00UTC",
            "id": "example-go",
            "owner": "test",
            "structure": {},
            "updated": "2014-01-01T00:00:00UTC",
            "uuid": "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75"
        }
    ]
}`

const appCreateExpected string = `{"id":"example-go"}`
const appRunExpected string = `{"command":"echo hi","timeout":3600,"expires":3600}`
const appTransferExpected string = `{"owner":"test"}`

type fakeHTTPServer struct {
	createID        bool
	createWithoutID bool
}

func (f *fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	if req.URL.Path == "/v2/apps/" && req.Method == "POST" {
		body, err := io.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) == appCreateExpected && !f.createID {
			f.createID = true
			res.WriteHeader(http.StatusCreated)
			res.Write([]byte(appFixture))
			return
		} else if string(body) == "" && !f.createWithoutID {
			f.createWithoutID = true
			res.WriteHeader(http.StatusCreated)
			res.Write([]byte(appFixture))
			return
		}

		fmt.Printf("Unexpected Body: %s'\n", body)
		res.WriteHeader(http.StatusInternalServerError)
		res.Write(nil)
		return
	}

	if req.URL.Path == "/v2/apps/" && req.Method == "GET" {
		res.Write([]byte(appsFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/" && req.Method == "GET" {
		res.Write([]byte(appFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/" && req.Method == "DELETE" {
		res.WriteHeader(http.StatusNoContent)
		res.Write(nil)
		return
	}

	if req.URL.Path == "/v2/apps/example-go/run" && req.Method == "POST" {
		body, err := io.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != appRunExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", appRunExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.Write([]byte(`{"exit_code":0,"output":"hi\n"}`))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/" && req.Method == "POST" {
		body, err := io.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != appTransferExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", appTransferExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.WriteHeader(http.StatusNoContent)
		res.Write(nil)
		return
	}

	fmt.Printf("Unrecongized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestAppsCreate(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{createID: false, createWithoutID: false}
	server := httptest.NewServer(&handler)
	defer server.Close()

	expected := api.App{
		ID:      "example-go",
		Created: "2014-01-01T00:00:00UTC",
		Owner:   "test",
		Updated: "2014-01-01T00:00:00UTC",
		UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	}

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	for _, id := range []string{"example-go", ""} {
		actual, err := New(drycc, id)

		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("Expected %v, Got %v", expected, actual)
		}
	}
}

func TestAppsGet(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	expected := api.App{
		ID:      "example-go",
		Created: "2014-01-01T00:00:00UTC",
		Owner:   "test",
		Updated: "2014-01-01T00:00:00UTC",
		UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
	}

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := Get(drycc, "example-go")

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}

func TestAppsDestroy(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	if err = Delete(drycc, "example-go"); err != nil {
		t.Fatal(err)
	}
}

func TestAppsRun(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	if err := Run(drycc, "example-go", "echo hi", nil, 3600, 3600); err != nil {
		t.Fatal(err)
	}
}

func TestAppsList(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	expected := api.Apps{
		{
			ID:      "example-go",
			Created: "2014-01-01T00:00:00UTC",
			Owner:   "test",
			Updated: "2014-01-01T00:00:00UTC",
			UUID:    "de1bf5b5-4a72-4f94-a10c-d2a3741cdf75",
		},
	}

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, _, err := List(drycc, 100)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}

func TestAppsLogs(t *testing.T) {
	t.Parallel()
	var once sync.Once
	var addr string
	once.Do(func() {
		http.Handle(
			"/v2/apps/example-go/logs",
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
	expected := api.AppLogsRequest{
		Lines:   1,
		Follow:  false,
		Timeout: 300,
	}
	conn, err := Logs(drycc, "example-go", expected)
	if err != nil {
		t.Fatal(err)
	}
	var actual api.AppLogsRequest
	err = websocket.JSON.Receive(conn, &actual)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected: %v, Got %v", expected, actual)
	}
}

func TestAppsTransfer(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	if err = Transfer(drycc, "example-go", "test"); err != nil {
		t.Fatal(err)
	}
}
