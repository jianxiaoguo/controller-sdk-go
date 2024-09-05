// Package releases provides methods for managing app releases.
package releases

import (
	"encoding/json"
	"fmt"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// List lists an app's releases.
func List(c *drycc.Client, appID string, results int) ([]api.Release, int, error) {
	u := fmt.Sprintf("/v2/apps/%s/releases/", appID)

	body, count, reqErr := c.LimitedRequest(u, results)

	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.Release{}, -1, reqErr
	}

	var releases []api.Release
	if err := json.Unmarshal([]byte(body), &releases); err != nil {
		return []api.Release{}, -1, err
	}

	return releases, count, reqErr
}

// Get retrieves a release of an app.
func Get(c *drycc.Client, appID string, version int) (api.Release, error) {
	u := fmt.Sprintf("/v2/apps/%s/releases/v%d/", appID, version)

	res, reqErr := c.Request("GET", u, nil)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.Release{}, reqErr
	}
	defer res.Body.Close()

	release := api.Release{}
	if err := json.NewDecoder(res.Body).Decode(&release); err != nil {
		return api.Release{}, err
	}

	return release, nil
}

// Deploy deploy an app's processes. To deploy all app processes, pass empty strings for
// procType and name. To deploy an specific process, pass an procType by leave name empty.
// To deploy a specific instance, pass a procType and a name.
func Deploy(c *drycc.Client, appID string, targets map[string]interface{}) error {
	u := fmt.Sprintf("/v2/apps/%s/releases/deploy/", appID)
	body, err := json.Marshal(targets)
	if err != nil {
		return err
	}
	res, err := c.Request("POST", u, body)
	if err != nil && !drycc.IsErrAPIMismatch(err) {
		return err
	}
	defer res.Body.Close()
	return err
}

// Rollback rolls back an app to a previous release. If version is -1, this rolls back to
// the previous release. Otherwise, roll back to the specified version.
func Rollback(c *drycc.Client, appID string, ptypes string, version int) (int, error) {
	u := fmt.Sprintf("/v2/apps/%s/releases/rollback/", appID)

	req := api.ReleaseRollback{Ptypes: ptypes, Version: version}

	var err error
	var reqBody []byte
	if version != -1 {
		reqBody, err = json.Marshal(req)

		if err != nil {
			return -1, err
		}
	}

	res, reqErr := c.Request("POST", u, reqBody)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return -1, reqErr
	}
	defer res.Body.Close()

	response := api.ReleaseRollback{}

	if err = json.NewDecoder(res.Body).Decode(&response); err != nil {
		return -1, err
	}

	return response.Version, reqErr
}
