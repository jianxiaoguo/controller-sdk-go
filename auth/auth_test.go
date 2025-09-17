package auth

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

const keyFixture = "61a7907bf5b34659a14f96371fed2ebc"

type fakeHTTPServer struct{}

func (f *fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	if req.URL.Path == "/v2/auth/login/" && req.Method == "POST" {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}
		if len(body) == 0 {
			res.Header().Add("Location", fmt.Sprintf("/v2/login/drycc/?key=%s/?alias=test", keyFixture))
			res.WriteHeader(http.StatusFound)
			res.Write(nil)
		} else {
			res.WriteHeader(http.StatusFound)
			fmt.Fprintf(res, `{"key": "%s"}`, keyFixture)
		}
	}

	if req.URL.Path == "/v2/auth/token/61a7907bf5b34659a14f96371fed2ebc/" && req.Method == "GET" {
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

	actual, err := Login(drycc, "", "")
	if err != nil {
		t.Error(err)
	}

	expected := fmt.Sprintf("/v2/login/drycc/?key=%s/?alias=test", keyFixture)
	if actual != expected {
		t.Errorf("Expected %s, Got %s", expected, actual)
	}

	actual, err = Login(drycc, "admin", "admin")
	if err != nil {
		t.Error(err)
	}
	if actual != keyFixture {
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
	expected := api.AuthTokenResponse{
		Username: "test",
		Token:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
	}

	token, err := Token(drycc, keyFixture, "test")
	if err != nil {
		t.Error(err)
	}

	if token != expected {
		t.Errorf("Expected %s, Got %s", expected, token)
	}
}
