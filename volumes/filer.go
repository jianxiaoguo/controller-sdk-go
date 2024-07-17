// Package config provides methods for managing configuration of apps.
package volumes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// ListDir to an app's volume.
func ListDir(c *drycc.Client, appID, volumeID, path string, results int) (api.FilerDirEntries, int, error) {
	u := fmt.Sprintf("/v2/apps/%s/volumes/%s/client/?path=%s", appID, volumeID, path)

	body, count, reqErr := c.LimitedRequest(u, results)

	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.FilerDirEntry{}, -1, reqErr
	}

	var filerDirEntries []api.FilerDirEntry
	if err := json.Unmarshal([]byte(body), &filerDirEntries); err != nil {
		return []api.FilerDirEntry{}, -1, err
	}

	return filerDirEntries, count, reqErr
}

// Getfile to an app's volume.
func GetFile(c *drycc.Client, appID, volumeID, path string) (*http.Response, error) {
	u := fmt.Sprintf("/v2/apps/%s/volumes/%s/client/%s", appID, volumeID, path)
	req, err := c.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Put file to an app's volume.
func PostFile(c *drycc.Client, appID, volumeID, path string, files ...string) (*http.Response, error) {
	buffer := new(bytes.Buffer)
	writer := multipart.NewWriter(buffer)
	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		if part, err := writer.CreateFormFile("file", f.Name()); err != nil {
			return nil, err
		} else if _, err = io.Copy(part, f); err != nil {
			return nil, err
		}
	}

	if err := writer.WriteField("path", path); err != nil {
		return nil, err
	}
	writer.Close()

	u := fmt.Sprintf("/v2/apps/%s/volumes/%s/client/", appID, volumeID)
	r, err := c.NewRequest("POST", u, buffer)
	if err != nil {
		return nil, err
	}
	r.Header.Add("Content-Type", writer.FormDataContentType())
	return c.Do(r)
}

// Get file to an app's volume.
func DeleteFile(c *drycc.Client, appID, volumeID, path string) (*http.Response, error) {
	u := fmt.Sprintf("/v2/apps/%s/volumes/%s/client/%s", appID, volumeID, path)
	return c.Request("DELETE", u, nil)
}
