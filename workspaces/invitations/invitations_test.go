package invitations

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

const invitationFixture = `
{
  "id": 1,
  "email": "invitee@example.com",
  "token": "abc",
  "inviter": "autotest",
  "created": "2026-03-24T00:00:00Z",
  "accepted": false,
  "workspace": "wsalpha"
}`

const invitationsFixture = `
{
  "count": 1,
  "next": null,
  "previous": null,
  "results": [
    {
      "id": 1,
      "email": "invitee@example.com",
      "token": "abc",
      "inviter": "autotest",
      "created": "2026-03-24T00:00:00Z",
      "accepted": false,
      "workspace": "wsalpha"
    }
  ]
}`

const invitationCreateExpected = `{"email":"invitee@example.com"}`

type fakeWorkspaceInvitationServer struct{}

func (f *fakeWorkspaceInvitationServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	if req.URL.Path == "/v2/workspaces/wsalpha/invitations" && req.Method == "GET" {
		res.Write([]byte(invitationsFixture))
		return
	}
	if req.URL.Path == "/v2/workspaces/wsalpha/invitations" && req.Method == "POST" {
		body, _ := io.ReadAll(req.Body)
		if string(body) != invitationCreateExpected {
			res.WriteHeader(http.StatusInternalServerError)
			return
		}
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(invitationFixture))
		return
	}
	if req.URL.Path == "/v2/workspaces/wsalpha/invitations/abc" && req.Method == "GET" {
		res.Write([]byte(invitationFixture))
		return
	}
	if req.URL.Path == "/v2/workspaces/wsalpha/invitations/abc" && req.Method == "DELETE" {
		res.WriteHeader(http.StatusNoContent)
		return
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
}

func TestWorkspaceInvitations(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(&fakeWorkspaceInvitationServer{})
	defer server.Close()

	c, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	expected := api.WorkspaceInvitation{ID: 1, Email: "invitee@example.com", Token: "abc", Inviter: "autotest", Created: "2026-03-24T00:00:00Z", Accepted: false, Workspace: "wsalpha"}

	list, _, err := List(c, "wsalpha", 100)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(list, api.WorkspaceInvitations{expected}) {
		t.Fatalf("unexpected list: %#v", list)
	}

	created, err := Create(c, "wsalpha", "invitee@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(created, expected) {
		t.Fatalf("unexpected create: %#v", created)
	}

	got, err := Get(c, "wsalpha", "abc")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("unexpected get: %#v", got)
	}

	if err := Delete(c, "wsalpha", "abc"); err != nil {
		t.Fatal(err)
	}
}
