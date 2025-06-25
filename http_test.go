package drycc

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeHTTPServer struct {
	Version         string
	PlatformVersion string
}

const limitedFixture string = `
{
    "count": 4,
    "next": "http://replaced.com/limited2/",
    "previous": null,
    "results": [
        {
            "test": "foo"
        },
        {
            "test": "bar"
        }
    ]
}
`

func (f fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", f.Version)
	res.Header().Add("DRYCC_PLATFORM_VERSION", f.PlatformVersion)

	eA := "test"

	if req.Header.Get("User-Agent") != eA {
		fmt.Printf("User Agent Wrong: Expected %s, Got %s\n", eA, req.Header.Get("User-Agent"))
		res.WriteHeader(http.StatusInternalServerError)
		res.Write(nil)
		return
	}

	if req.URL.Path == "/v2/" {
		res.WriteHeader(http.StatusUnauthorized)
		res.Write(nil)
		return
	}

	if req.URL.Path == "/healthz" {
		res.WriteHeader(http.StatusOK)
		res.Write(nil)
		return
	}

	if req.URL.Path == "/limited/" && req.Method == "GET" && req.URL.RawQuery == "limit=2" {
		res.Write([]byte(limitedFixture))
		return
	}

	if req.URL.Path == "/request/" && req.Method == "POST" {
		eT := "token abc"
		if req.Header.Get("Authorization") != eT {
			fmt.Printf("Token Wrong: Expected %s, Got %s\n", eT, req.Header.Get("Authorization"))
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		bT := "testing"
		if req.Header.Get("X-Drycc-Service-Key") != bT {
			fmt.Printf("Hook Token Wrong: Expected %s, Got %s\n", bT, req.Header.Get("X-Drycc-Service-Key"))
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		eC := "application/json"
		if req.Header.Get("Content-Type") != eC {
			fmt.Printf("Content Type Wrong: Expected %s, Got %s\n", eC, req.Header.Get("Content-Type"))
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		body, err := io.ReadAll(req.Body)

		if err != nil {
			fmt.Println(err)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
		}

		eB := "test"
		if string(body) != eB {
			fmt.Printf("Body Wrong: Expected %s, Got %s\n", eB, body)
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(nil)
			return
		}

		res.Write([]byte("request"))
		return
	}

	fmt.Printf("Unrecongized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestCheckConnection(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{Version: APIVersion}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := New(false, server.URL, "")
	if err != nil {
		t.Fatal(err)
	}
	drycc.UserAgent = "test"

	if err = drycc.CheckConnection(); err != nil {
		t.Error(err)
	}
}

func TestAPIMistmatch(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{Version: "3.0"}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := New(false, server.URL, "")
	if err != nil {
		t.Fatal(err)
	}
	drycc.UserAgent = "test"

	if err = drycc.CheckConnection(); err != ErrAPIMismatch {
		t.Error("Expected ErrAPIMismatch error")
	}

	if drycc.ControllerAPIVersion != handler.Version {
		t.Errorf("Expected %s, Got %s", handler.Version, drycc.ControllerAPIVersion)
	}
}

func TestBasicRequest(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{Version: APIVersion, PlatformVersion: "v9000"}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}
	drycc.UserAgent = "test"
	drycc.ServiceKey = "testing"

	res, err := drycc.Request("POST", "/request/", []byte("test"))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	expected := "request"
	if string(body) != expected {
		t.Errorf("Expected %s, Got %s", expected, string(body))
	}

	if drycc.ControllerAPIVersion != handler.Version {
		t.Errorf("Expected %s, Got %s", handler.Version, drycc.ControllerAPIVersion)
	}

	if drycc.ControllerVersion != handler.PlatformVersion {
		t.Errorf("Expected %s, Got %s", handler.PlatformVersion, drycc.ControllerVersion)
	}

	// Make sure the request doesn't modify the URL
	if drycc.ControllerURL.String() != server.URL {
		t.Errorf("Expected %s, Got %s", server.URL, drycc.ControllerURL.String())
	}
}

func TestLimitedRequest(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{Version: APIVersion}
	server := httptest.NewServer(handler)
	defer server.Close()

	drycc, err := New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}
	drycc.UserAgent = "test"

	expected := `[{"test":"foo"},{"test":"bar"}]`
	expectedC := 4

	actual, count, err := drycc.LimitedRequest("/limited/", 2)

	if err != nil {
		t.Fatal(err)
	}

	if count != expectedC {
		t.Errorf("Expected %d, Got %d", expectedC, count)
	}

	if actual != expected {
		t.Errorf("Expected %s, Got %s", expected, actual)
	}

	if drycc.ControllerAPIVersion != handler.Version {
		t.Errorf("Expected %s, Got %s", handler.Version, drycc.ControllerAPIVersion)
	}

	// Make sure the request doesn't modify the URL
	if drycc.ControllerURL.String() != server.URL {
		t.Errorf("Expected %s, Got %s", server.URL, drycc.ControllerURL.String())
	}
}

func TestHealthcheck(t *testing.T) {
	t.Parallel()

	handler := fakeHTTPServer{Version: APIVersion}
	server := httptest.NewServer(handler)
	defer server.Close()

	// Test with a trailing slash
	drycc, err := New(false, server.URL+"/", "")
	if err != nil {
		t.Fatal(err)
	}
	drycc.UserAgent = "test"

	if err = drycc.Healthcheck(); err != nil {
		t.Error(err)
	}

	// Test without a trailing slash
	drycc, err = New(false, server.URL, "")
	if err != nil {
		t.Fatal(err)
	}
	drycc.UserAgent = "test"

	if err = drycc.Healthcheck(); err != nil {
		t.Error(err)
	}
}
