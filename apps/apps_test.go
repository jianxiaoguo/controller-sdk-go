package apps

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
	"github.com/gorilla/websocket"
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
const appRunExpected string = `{"command":"echo hi"}`
const appTransferExpected string = `{"owner":"test"}`

var upgrader = websocket.Upgrader{} // use default options

type fakeHTTPServer struct {
	createID        bool
	createWithoutID bool
}

type fakeLogReqMessage struct {
	Lines   int  `json:"lines"`
	Timeout int  `json:"timeout"`
	Follow  bool `json:"follow"`
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

	// The entire log message is prefixed and suffixed with a few characters (not entirely sure why)
	// We mimic those here
	if req.URL.Path == "/v2/apps/example-go/logs" {
		c, err := upgrader.Upgrade(res, req, nil)
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		defer c.Close()
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Print("upgrade:", err)
			return
		}
		reqJSON := fakeLogReqMessage{}
		json.Unmarshal([]byte(message), &reqJSON)
		for i := 0; i < reqJSON.Lines; i++ {
			c.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("test %d", i)))
		}
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

	expected := api.AppRunResponse{
		Output:     "hi\n",
		ReturnCode: 0,
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := Run(drycc, "example-go", "echo hi", nil)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
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

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	conn, err := Logs(drycc, "example-go", 1, false, 300)
	if err != nil {
		t.Fatal(err)
	}
	_, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	if string(message) != "test 0" {
		t.Errorf("Expected %s, Got %s", "test 0", message)
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
