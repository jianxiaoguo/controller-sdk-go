package auth

import (
	"fmt"
	"github.com/drycc/controller-sdk-go/api"
	"net/http"
	"net/http/httptest"
	"testing"

	drycc "github.com/drycc/controller-sdk-go"
)

type fakeHTTPServer struct {
	regenBodyEmpty    bool
	regenBodyAll      bool
	regenBodyUsername bool
	cancelEmpty       bool
	cancelUsername    bool
}

func (f *fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	if req.URL.Path == "/v2/auth/login/" && req.Method == "POST" {
		res.Header().Add("Location", "/v2/login/drycc/?key=fdbf3b34742e4ed2be4dfa848af13007/")
		res.WriteHeader(http.StatusFound)
		res.Write(nil)
		return
	}

	if req.URL.Path == "/v2/auth/token/fdbf3b34742e4ed2be4dfa848af13007/" && req.Method == "GET" {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(`{"username":"test","token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"}`))
		return
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestLogin(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := Login(drycc)

	if err != nil {
		t.Error(err)
	}

	expected := "/v2/login/drycc/?key=fdbf3b34742e4ed2be4dfa848af13007/"
	if actual != expected {
		t.Errorf("Expected %s, Got %s", expected, actual)
	}
}

func TestToken(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}
	expected := api.AuthLoginResponse{
		Username: "test",
		Token:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
	}

	token, err := Token(drycc, "fdbf3b34742e4ed2be4dfa848af13007")

	if err != nil {
		t.Error(err)
	}

	if token != expected {
		t.Errorf("Expected %s, Got %s", expected, token)
	}
}
