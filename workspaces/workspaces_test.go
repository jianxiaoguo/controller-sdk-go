package workspaces

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

const workspaceFixture = `
{
  "id": 1,
  "name": "wsalpha",
  "email": "ws@example.com",
  "created": "2026-03-24T00:00:00Z",
  "updated": "2026-03-24T00:00:00Z"
}`

const workspacesFixture = `
{
  "count": 1,
  "next": null,
  "previous": null,
  "results": [
    {
      "id": 1,
      "name": "wsalpha",
      "email": "ws@example.com",
      "created": "2026-03-24T00:00:00Z",
      "updated": "2026-03-24T00:00:00Z"
    }
  ]
}`

const (
	workspaceCreateExpected = `{"name":"wsalpha","email":"ws@example.com"}`
	workspaceUpdateExpected = `{"email":"ws-new@example.com"}`
)

type fakeWorkspaceServer struct{}

func (f *fakeWorkspaceServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	if req.URL.Path == "/v2/workspaces" && req.Method == "GET" {
		res.Write([]byte(workspacesFixture))
		return
	}
	if req.URL.Path == "/v2/workspaces" && req.Method == "POST" {
		body, _ := io.ReadAll(req.Body)
		if string(body) != workspaceCreateExpected {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(workspaceFixture))
		return
	}
	if req.URL.Path == "/v2/workspaces/wsalpha" && req.Method == "GET" {
		res.Write([]byte(workspaceFixture))
		return
	}
	if req.URL.Path == "/v2/workspaces/wsalpha" && req.Method == "PATCH" {
		body, _ := io.ReadAll(req.Body)
		if string(body) != workspaceUpdateExpected {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.Write([]byte(`{"id":1,"name":"wsalpha","email":"ws-new@example.com","created":"2026-03-24T00:00:00Z","updated":"2026-03-24T00:00:00Z"}`))
		return
	}
	if req.URL.Path == "/v2/workspaces/wsalpha" && req.Method == "DELETE" {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
}

func TestWorkspaces(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(&fakeWorkspaceServer{})
	defer server.Close()

	c, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	expected := api.Workspace{ID: 1, Name: "wsalpha", Email: "ws@example.com", Created: "2026-03-24T00:00:00Z", Updated: "2026-03-24T00:00:00Z"}

	list, _, err := List(c, 100)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(list, api.Workspaces{expected}) {
		t.Fatalf("unexpected list: %#v", list)
	}

	created, err := Create(c, "wsalpha", "ws@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(created, expected) {
		t.Fatalf("unexpected create: %#v", created)
	}

	got, err := Get(c, "wsalpha")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("unexpected get: %#v", got)
	}

	updated, err := Update(c, "wsalpha", "ws-new@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if updated.Email != "ws-new@example.com" {
		t.Fatalf("unexpected update: %#v", updated)
	}

	if err := Delete(c, "wsalpha"); err != nil {
		t.Fatal(err)
	}
}
