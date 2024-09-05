package certs

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/drycc/controller-sdk-go/pkg/time"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

const certsFixture string = `
{
    "count": 1,
    "next": null,
    "previous": null,
    "results": [
        {
            "app": "example-go",
            "name": "test-example-com",
            "common_name": "test.example.com",
            "expires": "2014-01-01T00:00:00UTC",
            "fingerprint": "12:34:56:78:90"
        }
    ]
}`

const certFixture string = `
{
    "updated": "2014-01-01T00:00:00UTC",
    "created": "2014-01-01T00:00:00UTC",
    "expires": "2015-01-01T00:00:00UTC",
    "starts": "2014-01-01T00:00:00UTC",
    "fingerprint": "12:34:56:78:90",
    "name": "test-example-com",
    "app": "example-go",
    "owner": "test",
    "id": 1
}`

const certExpected string = `{"certificate":"test","key":"foo","name":"test-example-com"}`

type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

	if req.URL.Path == "/v2/apps/example-go/certs/" && req.Method == "GET" {
		res.Write([]byte(certsFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/certs/test-example-com" && req.Method == "GET" {
		res.Write([]byte(certFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/certs/" && req.Method == "POST" {
		body, err := io.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		if string(body) != certExpected {
			fmt.Printf("Expected '%s', Got '%s'\n", certExpected, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(certFixture))
		return
	}

	if req.URL.Path == "/v2/apps/example-go/certs/test-example-com/domain/" && req.Method == "POST" {
		res.WriteHeader(http.StatusCreated)
		res.Write(nil)
		return
	}

	if req.URL.Path == "/v2/apps/example-go/certs/test-example-com" && req.Method == "DELETE" {
		res.WriteHeader(http.StatusNoContent)
		res.Write(nil)
		return
	}

	if req.URL.Path == "/v2/apps/example-go/certs/test-example-com/domain/foo.com" && req.Method == "DELETE" {
		res.WriteHeader(http.StatusNoContent)
		res.Write(nil)
		return
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestCertsList(t *testing.T) {
	t.Parallel()

	expires := time.Time{}

	if err := expires.UnmarshalText([]byte("2014-01-01T00:00:00UTC")); err != nil {
		t.Fatalf("could not unmarshal time (%s)", err)
	}

	expected := []api.Cert{
		{
			Name:        "test-example-com",
			App:         "example-go",
			Expires:     expires,
			CommonName:  "test.example.com",
			Fingerprint: "12:34:56:78:90",
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
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}

func TestCert(t *testing.T) {
	t.Parallel()

	created := "2014-01-01T00:00:00UTC"
	starts := time.Time{}
	starts.UnmarshalText([]byte("2014-01-01T00:00:00UTC"))
	expires := time.Time{}
	expires.UnmarshalText([]byte("2015-01-01T00:00:00UTC"))

	expected := api.Cert{
		Updated:     created,
		Created:     created,
		Starts:      starts,
		Expires:     expires,
		Fingerprint: "12:34:56:78:90",
		App:         "example-go",
		Name:        "test-example-com",
		Owner:       "test",
		ID:          1,
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := New(drycc, "example-go", "test", "foo", "test-example-com")

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}

func TestCertInfo(t *testing.T) {
	t.Parallel()

	created := "2014-01-01T00:00:00UTC"
	starts := time.Time{}
	starts.UnmarshalText([]byte("2014-01-01T00:00:00UTC"))
	expires := time.Time{}
	expires.UnmarshalText([]byte("2015-01-01T00:00:00UTC"))

	expected := api.Cert{
		Updated:     created,
		Created:     created,
		Starts:      starts,
		Expires:     expires,
		Fingerprint: "12:34:56:78:90",
		App:         "example-go",
		Name:        "test-example-com",
		Owner:       "test",
		ID:          1,
	}

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	actual, err := Get(drycc, "example-go", "test-example-com")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected %v, Got %v", expected, actual)
	}
}

func TestCertDeletion(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	if err = Delete(drycc, "example-go", "test-example-com"); err != nil {
		t.Fatal(err)
	}

	if err := Delete(drycc, "example-go", "non-existent-cert"); err == nil {
		t.Fatal("An Error should have resulted from the attempt to delete a non-existent-cert")
	}
}

func TestCertAttach(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	if err = Attach(drycc, "example-go", "test-example-com", "foo.com"); err != nil {
		t.Fatal(err)
	}

	if err := Attach(drycc, "example-go", "non-existent-cert", "foo.com"); err == nil {
		t.Fatal("An Error should have resulted from the attempt to attach a non-existent cert to a valid domain")
	}

	// TODO: #475
	// if err := Attach(&drycc, "example-go", "test-example-com", "non-existent.domain.com"); err == nil {
	// 	t.Fatal("An Error should have resulted from the attempt to attach a valid cert to a non-existent domain")
	// }
}

func TestCertDetach(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	if err = Detach(drycc, "example-go", "test-example-com", "foo.com"); err != nil {
		t.Fatal(err)
	}

	if err := Detach(drycc, "example-go", "non-existent-cert", "foo.com"); err == nil {
		t.Fatal("An Error should have resulted from the attempt to detach a non-existent cert from a valid domain")
	}

	if err := Detach(drycc, "example-go", "test-example-com", "non-existent.domain.com"); err == nil {
		t.Fatal("An Error should have resulted from the attempt to detach a valid cert from a non-existent domain")
	}
}
