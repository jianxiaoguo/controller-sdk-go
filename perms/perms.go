// Package perms provides methods for managing user app and administrative permissions.
package perms

import (
	"encoding/json"
	"fmt"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// List UserPerm
func List(c *drycc.Client, appID string, results int) ([]api.UserPermResponse, int, error) {
	u := fmt.Sprintf("/v2/apps/%s/perms/", appID)
	body, count, reqErr := c.LimitedRequest(u, results)

	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.UserPermResponse{}, count, reqErr
	}

	var userPerm []api.UserPermResponse
	if err := json.Unmarshal([]byte(body), &userPerm); err != nil {
		return []api.UserPermResponse{}, count, err
	}

	return userPerm, count, reqErr
}

// Create a App user'Perms
func Create(c *drycc.Client, appID, username, permissions string) error {
	req := api.UserPermRequest{Username: username, Permissions: permissions}
	reqBody, err := json.Marshal(req)

	if err != nil {
		return err
	}
	u := fmt.Sprintf("/v2/apps/%s/perms/", appID)
	res, err := c.Request("POST", u, reqBody)
	if err == nil {
		res.Body.Close()
	}

	return err
}

// Update a App user'Perms
func Update(c *drycc.Client, appID, username, permissions string) error {
	req := api.UserPermRequest{Username: username, Permissions: permissions}
	reqBody, err := json.Marshal(req)

	if err != nil {
		return err
	}
	u := fmt.Sprintf("/v2/apps/%s/perms/%s/", appID, username)
	res, err := c.Request("PUT", u, reqBody)
	if err == nil {
		res.Body.Close()
	}

	return err
}

// Delete a App user'Perms.
func Delete(c *drycc.Client, appID, username string) error {
	u := fmt.Sprintf("/v2/apps/%s/perms/%s/", appID, username)
	res, err := c.Request("DELETE", u, nil)
	if err == nil {
		res.Body.Close()
	}
	return err
}
