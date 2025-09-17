package tokens

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

const tokensFixture string = `
{
    "count": 2,
    "next": null,
    "previous": null,
    "results": [
		{
			"uuid": "f71e3b18-e702-409e-bd7f-8fb0a66d7b12",
			"owner": "test",
			"alias": "",
			"fuzzy_key": "c8e74fa4cbf...e4954d602ec5ed19ba",
			"created": "2023-04-18T00:00:00UTC",
			"updated": "2023-04-19T00:00:00UTC"
		},
		{
			"uuid": "f71e3b18-e702-499e-bd7f-8fb0a66d7b12",
			"owner": "test",
			"alias": "test",
			"fuzzy_key": "c8e74fa4cbf...e4954d60cec5ed19ba",
			"created": "2023-04-18T10:00:00UTC",
			"updated": "2023-04-19T12:00:00UTC"
		}
    ]
}`

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	if req.URL.Path == "/v2/tokens/" && req.Method == "GET" {
		res.Write([]byte(tokensFixture))
		return
	}

	if req.URL.Path == "/v2/tokens/f71e3b18-e702-499e-bd7f-8fb0a66d7b12/" && req.Method == "DELETE" {
		res.WriteHeader(204)
		return
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestTokensList(t *testing.T) {
	t.Parallel()

	expected := []api.Token{
		{
			UUID:    "f71e3b18-e702-409e-bd7f-8fb0a66d7b12",
			Owner:   "test",
			Alias:   "",
			Key:     "c8e74fa4cbf...e4954d602ec5ed19ba",
			Created: "2023-04-18T00:00:00UTC",
			Updated: "2023-04-19T00:00:00UTC",
		},
		{
			UUID:    "f71e3b18-e702-499e-bd7f-8fb0a66d7b12",
			Owner:   "test",
			Alias:   "test",
			Key:     "c8e74fa4cbf...e4954d60cec5ed19ba",
			Created: "2023-04-18T10:00:00UTC",
			Updated: "2023-04-19T12:00:00UTC",
		},
	}
	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, _, err := List(drycc, 100)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
	}
}

func TestTokensRemove(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	err = Delete(drycc, "f71e3b18-e702-499e-bd7f-8fb0a66d7b12")
	if err != nil {
		t.Fatal(err)
	}
}
