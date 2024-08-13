// Package perms provides methods for managing user app and administrative permissions.
package perms

import (
	"encoding/json"
	"fmt"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// List Codename
func Codes(c *drycc.Client, results int) ([]api.PermCodeResponse, int, error) {
	body, count, reqErr := c.LimitedRequest("/v2/perms/codes/", results)

	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.PermCodeResponse{}, count, reqErr
	}

	var codenames []api.PermCodeResponse
	if err := json.Unmarshal([]byte(body), &codenames); err != nil {
		return []api.PermCodeResponse{}, count, err
	}

	return codenames, count, reqErr
}

// List UserPerm
func List(c *drycc.Client, codename string, results int) ([]api.UserPermResponse, int, error) {
	u := fmt.Sprintf("/v2/perms/rules/?codename=%s", codename)
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

// Create a UserPerm
func Create(c *drycc.Client, codename, uniqueid, username string) error {
	req := api.UserPermRequest{Codename: codename, Uniqueid: uniqueid, Username: username}
	reqBody, err := json.Marshal(req)

	if err != nil {
		return err
	}

	res, err := c.Request("POST", "/v2/perms/rules/", reqBody)
	if err == nil {
		res.Body.Close()
	}

	return err
}

// Delete a UserPerm.
func Delete(c *drycc.Client, userPermID string) error {
	res, err := c.Request("DELETE", fmt.Sprintf("/v2/perms/rules/%s/", userPermID), nil)
	if err == nil {
		res.Body.Close()
	}
	return err
}
