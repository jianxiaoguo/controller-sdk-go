package members

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

const memberFixture = `
{
  "id": 1,
  "user": "autotest",
  "email": "auto@test.com",
  "role": "member",
  "alerts": true,
  "workspace": "wsalpha",
  "created": "2026-03-24T00:00:00Z",
  "updated": "2026-03-24T00:00:00Z"
}`

const membersFixture = `
{
  "count": 1,
  "next": null,
  "previous": null,
  "results": [
    {
      "id": 1,
      "user": "autotest",
      "email": "auto@test.com",
      "role": "member",
      "alerts": true,
      "workspace": "wsalpha",
      "created": "2026-03-24T00:00:00Z",
      "updated": "2026-03-24T00:00:00Z"
    }
  ]
}`

const memberUpdateExpected = `{"role":"viewer"}`

type fakeWorkspaceMemberServer struct{}

func (f *fakeWorkspaceMemberServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	if req.URL.Path == "/v2/workspaces/wsalpha/members" && req.Method == "GET" {
		res.Write([]byte(membersFixture))
		return
	}
	if req.URL.Path == "/v2/workspaces/wsalpha/members/autotest" && req.Method == "GET" {
		res.Write([]byte(memberFixture))
		return
	}
	if req.URL.Path == "/v2/workspaces/wsalpha/members/autotest" && req.Method == "PATCH" {
		body, _ := io.ReadAll(req.Body)
		if string(body) != memberUpdateExpected {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.Write([]byte(`{"id":1,"user":"autotest","email":"auto@test.com","role":"viewer","alerts":true,"workspace":"wsalpha","created":"2026-03-24T00:00:00Z","updated":"2026-03-24T00:00:00Z"}`))
		return
	}
	if req.URL.Path == "/v2/workspaces/wsalpha/members/autotest" && req.Method == "DELETE" {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
}

func TestWorkspaceMembers(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(&fakeWorkspaceMemberServer{})
	defer server.Close()

	c, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	expected := api.WorkspaceMember{ID: 1, User: "autotest", Email: "auto@test.com", Role: "member", Alerts: true, Workspace: "wsalpha", Created: "2026-03-24T00:00:00Z", Updated: "2026-03-24T00:00:00Z"}

	list, _, err := List(c, "wsalpha", 100)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(list, api.WorkspaceMembers{expected}) {
		t.Fatalf("unexpected list: %#v", list)
	}

	got, err := Get(c, "wsalpha", "autotest")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("unexpected get: %#v", got)
	}

	updated, err := Update(c, "wsalpha", "autotest", "viewer", nil)
	if err != nil {
		t.Fatal(err)
	}
	if updated.Role != "viewer" {
		t.Fatalf("unexpected update: %#v", updated)
	}

	if err := Delete(c, "wsalpha", "autotest"); err != nil {
		t.Fatal(err)
	}
}
