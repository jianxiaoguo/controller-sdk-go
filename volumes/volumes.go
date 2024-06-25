// Package config provides methods for managing configuration of apps.
package volumes

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// List list an app's volumes.
func List(c *drycc.Client, appID string, results int) (api.Volumes, int, error) {
	u := fmt.Sprintf("/v2/apps/%s/volumes/", appID)
	body, count, reqErr := c.LimitedRequest(u, results)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.Volume{}, -1, reqErr
	}
	var volumes []api.Volume
	if err := json.Unmarshal([]byte(body), &volumes); err != nil {
		return []api.Volume{}, -1, err
	}
	return volumes, count, reqErr
}

// Get an app's volume.
func Get(c *drycc.Client, appID string, name string) (api.Volume, error) {
	u := fmt.Sprintf("/v2/apps/%s/volumes/%s/", appID, name)
	res, reqErr := c.Request("GET", u, nil)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.Volume{}, reqErr
	}
	defer res.Body.Close()

	volume := api.Volume{}
	if err := json.NewDecoder(res.Body).Decode(&volume); err != nil {
		return volume, err
	}

	return volume, nil
}

// ListDir to an app's volume.
func ListDir(c *drycc.Client, appID, volumeID, path string) (*http.Response, error) {
	u := fmt.Sprintf("/v2/apps/%s/volumes/%s/files%s", appID, volumeID, path)
	req, err := c.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	return c.Do(req)
}

// Getfile to an app's volume.
func GetFile(c *drycc.Client, appID, volumeID, path string) (*http.Response, error) {
	u := fmt.Sprintf("/v2/apps/%s/volumes/%s/files%s", appID, volumeID, path)
	req, err := c.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/octet-stream")
	return c.Do(req)
}

// Put file to an app's volume.
func PostFile(c *drycc.Client, appID, volumeID, path string, files ...string) (*http.Response, error) {
	buffer := new(bytes.Buffer)
	writer := multipart.NewWriter(buffer)
	for _, file := range files {
		if f, err := os.Open(file); err != nil {
			return nil, err
		} else if part, err := writer.CreateFormFile("file", f.Name()); err != nil {
			return nil, err
		} else if _, err = io.Copy(part, f); err != nil {
			return nil, err
		} else {
			defer f.Close()
		}
	}
	writer.Close()

	if !strings.HasPrefix(path, "/") {
		return nil, errors.New("please use an absolute path starting with /")
	}
	u := fmt.Sprintf("/v2/apps/%s/volumes/%s/files%s", appID, volumeID, path)
	r, err := c.NewRequest("POST", u, buffer)
	if err != nil {
		return nil, err
	}
	return c.Do(r)
}

// Get file to an app's volume.
func DeleteFile(c *drycc.Client, appID, volumeID, path string) (*http.Response, error) {
	if !strings.HasPrefix(path, "/") {
		return nil, errors.New("please use an absolute path starting with /")
	}
	u := fmt.Sprintf("/v2/apps/%s/volumes/%s/files%s", appID, volumeID, path)
	return c.Request("DELETE", u, nil)
}

// Create create an app's Volume.
func Create(c *drycc.Client, appID string, volume api.Volume) (api.Volume, error) {
	body, err := json.Marshal(volume)
	if err != nil {
		return api.Volume{}, err
	}
	u := fmt.Sprintf("/v2/apps/%s/volumes/", appID)
	res, reqErr := c.Request("POST", u, body)
	if reqErr != nil {
		return api.Volume{}, reqErr
	}
	defer res.Body.Close()
	newVolume := api.Volume{}
	if err = json.NewDecoder(res.Body).Decode(&newVolume); err != nil {
		return api.Volume{}, err
	}
	return newVolume, reqErr
}

// Expand create an app's Volume.
func Expand(c *drycc.Client, appID string, volume api.Volume) (api.Volume, error) {
	body, err := json.Marshal(volume)
	if err != nil {
		return api.Volume{}, err
	}
	u := fmt.Sprintf("/v2/apps/%s/volumes/%s/", appID, volume.Name)
	res, reqErr := c.Request("PATCH", u, body)
	if reqErr != nil {
		return api.Volume{}, reqErr
	}
	defer res.Body.Close()
	newVolume := api.Volume{}
	if err = json.NewDecoder(res.Body).Decode(&newVolume); err != nil {
		return api.Volume{}, err
	}
	return newVolume, reqErr
}

// Delete delete an app's Volume.
func Delete(c *drycc.Client, appID string, name string) error {
	u := fmt.Sprintf("/v2/apps/%s/volumes/%s/", appID, name)
	res, err := c.Request("DELETE", u, nil)
	if err == nil {
		res.Body.Close()
	}
	return err
}

// Mount mount an app's volume and creates a new release.
// This is a patching operation, which means when you call Mount() with an api.Volumes:
//
//   - If the variable does not exist, it will be set.
//   - If the variable exists, it will be overwritten.
//   - If the variable is set to nil, it will be unmount.
//   - If the variable was ignored in the api.Volumes, it will remain unchanged.
//
// Calling Mount() with an empty api.Volume will return a drycc.ErrConflict.
// Trying to Unmount a key that does not exist returns a drycc.ErrUnprocessable.
func Mount(c *drycc.Client, appID string, name string, volume api.Volume) (api.Volume, error) {
	body, err := json.Marshal(volume)
	if err != nil {
		return api.Volume{}, err
	}
	u := fmt.Sprintf("/v2/apps/%s/volumes/%s/path/", appID, name)
	res, reqErr := c.Request("PATCH", u, body)
	if reqErr != nil {
		return api.Volume{}, reqErr
	}
	defer res.Body.Close()
	newVolume := api.Volume{}
	if err = json.NewDecoder(res.Body).Decode(&newVolume); err != nil {
		return api.Volume{}, err
	}
	return newVolume, reqErr
}
