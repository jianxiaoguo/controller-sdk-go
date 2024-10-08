// Package config provides methods for managing configuration of apps.
package volumes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// ListDir to an app's volume.
func ListDir(c *drycc.Client, appID, volumeID, path string, results int) (api.FilerDirEntries, int, error) {
	u := fmt.Sprintf("/v2/apps/%s/volumes/%s/client/?path=%s", appID, volumeID, url.QueryEscape(path))

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
func PostFile(c *drycc.Client, appID, volumeID, volumePath, name string, size int64, reader io.Reader) (*http.Response, error) {

	buffer := new(bytes.Buffer)
	writer := multipart.NewWriter(buffer)
	if err := writer.WriteField("path", volumePath); err != nil {
		return nil, err
	}
	if _, err := writer.CreateFormFile("file", name); err != nil {
		return nil, err
	}
	size += int64(buffer.Len())
	head := strings.NewReader(buffer.String())
	buffer.Reset()
	writer.Close()
	bottom := strings.NewReader(buffer.String())
	size += int64(buffer.Len())
	u := fmt.Sprintf("/v2/apps/%s/volumes/%s/client/", appID, volumeID)

	r, err := c.NewRequest("POST", u, io.MultiReader(head, reader, bottom))
	if err != nil {
		return nil, err
	}
	r.ContentLength = size
	r.Header.Add("Content-Type", writer.FormDataContentType())
	return c.Do(r)
}

// Get file to an app's volume.
func DeleteFile(c *drycc.Client, appID, volumeID, path string) (*http.Response, error) {
	u := fmt.Sprintf("/v2/apps/%s/volumes/%s/client/%s", appID, volumeID, path)
	return c.Request("DELETE", u, nil)
}
