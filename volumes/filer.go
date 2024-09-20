// Package config provides methods for managing configuration of apps.
package volumes

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

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
	u := fmt.Sprintf("/v2/apps/%s/volumes/%s/client/", appID, volumeID)
	r, err := c.NewRequest("POST", u, reader)
	if err != nil {
		return nil, err
	}
	r.ContentLength = size
	r.Header.Add("Content-Type", "filer/octet-stream")
	disposition := fmt.Sprintf(
		`attachment;filename*=utf-8''%s;filepath*=utf-8''%s`,
		url.QueryEscape(name),
		url.QueryEscape(volumePath),
	)
	r.Header.Add("Content-Disposition", disposition)
	return c.Do(r)
}

// Get file to an app's volume.
func DeleteFile(c *drycc.Client, appID, volumeID, path string) (*http.Response, error) {
	u := fmt.Sprintf("/v2/apps/%s/volumes/%s/client/%s", appID, volumeID, path)
	return c.Request("DELETE", u, nil)
}
