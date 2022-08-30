// Package users provides methods for viewing users.
package users

import (
	"encoding/json"
	"fmt"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// List lists users registered with the controller.
func List(c *drycc.Client, results int) (api.Users, int, error) {
	body, count, reqErr := c.LimitedRequest("/v2/users/", results)

	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.User{}, -1, reqErr
	}

	var users []api.User
	if err := json.Unmarshal([]byte(body), &users); err != nil {
		return []api.User{}, -1, err
	}

	return users, count, reqErr
}

// Enable user with the controller.
func Enable(c *drycc.Client, username string) error {
	u := fmt.Sprintf("/v2/users/%s/enable/", username)
	res, err := c.Request("PATCH", u, nil)

	if err == nil {
		return res.Body.Close()
	}
	return err
}

// Disable user with the controller.
func Disable(c *drycc.Client, username string) error {
	u := fmt.Sprintf("/v2/users/%s/disable/", username)
	res, err := c.Request("PATCH", u, nil)

	if err == nil {
		return res.Body.Close()
	}
	return err
}
