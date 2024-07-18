package volumes

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	drycc "github.com/drycc/controller-sdk-go"
)

const volumeFileContentExpected string = `hello world`

type fakeFilerServer struct{}

func (f *fakeFilerServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)
	// get/delete file
	if strings.Contains(req.URL.Path, "/v2/apps/example-go/volumes/myvolume/client/tmp/helloword.txt") {
		if req.Method == "GET" {
			res.Header().Add("Content-Type", "application/octet-stream")
			res.Write([]byte(volumeFileContentExpected))
			return
		} else if req.Method == "DELETE" {
			res.WriteHeader(http.StatusNoContent)
			res.Write(nil)
			return
		}
	}
	// post file or list dir
	if strings.Contains(req.URL.Path, "/v2/apps/example-go/volumes/myvolume/client/") {
		if req.Method == "GET" {
			res.Header().Add("Content-Type", "application/json")
			res.Write([]byte(`{"results":[], "count": 0}`))
			return
		} else if req.Method == "POST" {
			if err := req.ParseMultipartForm(1024 * 1024); err != nil {
				http.Error(res, fmt.Sprintf("Parse multipart form error: %v", err), http.StatusBadRequest)
				return
			}
			for _, tmpFiles := range req.MultipartForm.File {
				for _, tmpFile := range tmpFiles {
					srcFile, err := tmpFile.Open()
					if err != nil {
						return
					}
					body, err := io.ReadAll(srcFile)
					if err != nil {
						return
					}
					if string(body) != volumeFileContentExpected {
						fmt.Printf("Expected '%s', Got '%s'\n", volumeFileContentExpected, body)
						res.WriteHeader(http.StatusInternalServerError)
						res.Write(nil)
						return
					}
				}
			}
			return
		}
	}

	fmt.Printf("Unrecognized URL %s\n", req.URL)
	res.WriteHeader(http.StatusNotFound)
	res.Write(nil)
}

func TestVolumesListDir(t *testing.T) {
	t.Parallel()

	handler := fakeFilerServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	_, counts, err := ListDir(drycc, "example-go", "myvolume", "tmp", 3000)
	if err != nil {
		t.Fatal(err)
	}
	if counts != 0 {
		t.Error(fmt.Errorf("Expected %v, Got %v", 0, counts))
	}
}

func TestVolumesGetFile(t *testing.T) {
	t.Parallel()

	handler := fakeFilerServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	res, err := GetFile(drycc, "example-go", "myvolume", "tmp/helloword.txt")
	if err != nil {
		t.Fatal(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(volumeFileContentExpected, string(body)) {
		t.Error(fmt.Errorf("Expected %v, Got %v", volumeFileContentExpected, string(body)))
	}
}

func TestVolumesPostFile(t *testing.T) {
	t.Parallel()

	handler := fakeFilerServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	testFile := "helloword.txt"

	err := os.WriteFile(testFile, []byte(volumeFileContentExpected), 0644)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Remove(testFile); err != nil {
			t.Fatal(err)
		}
	}()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	file, err := os.Open(testFile)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	if _, err := PostFile(drycc, "example-go", "myvolume", "tmp/", file.Name(), file); err != nil {
		t.Fatal(err)
	}

}

func TestVolumesDeleteFile(t *testing.T) {
	t.Parallel()

	handler := fakeFilerServer{}
	server := httptest.NewServer(&handler)
	defer server.Close()

	drycc, err := drycc.New(false, server.URL, "abc")
	if err != nil {
		t.Fatal(err)
	}

	_, err = DeleteFile(drycc, "example-go", "myvolume", "tmp/helloword.txt")
	if err != nil {
		t.Fatal(err)
	}
}
