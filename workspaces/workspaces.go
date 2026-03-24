// Package workspaces provides methods for managing controller workspaces.
package workspaces

import (
	"encoding/json"
	"fmt"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// List lists workspaces visible to the current user.
func List(c *drycc.Client, results int) (api.Workspaces, int, error) {
	body, count, reqErr := c.LimitedRequest("/v2/workspaces", results)

	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.Workspace{}, -1, reqErr
	}

	var workspaces []api.Workspace
	if err := json.Unmarshal([]byte(body), &workspaces); err != nil {
		return []api.Workspace{}, -1, err
	}

	return workspaces, count, reqErr
}

// Create creates a workspace.
func Create(c *drycc.Client, name, email string) (api.Workspace, error) {
	req := api.WorkspaceCreateRequest{Name: name, Email: email}
	body, err := json.Marshal(req)
	if err != nil {
		return api.Workspace{}, err
	}

	res, reqErr := c.Request("POST", "/v2/workspaces", body)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.Workspace{}, reqErr
	}
	defer res.Body.Close()

	workspace := api.Workspace{}
	if err = json.NewDecoder(res.Body).Decode(&workspace); err != nil {
		return api.Workspace{}, err
	}

	return workspace, reqErr
}

// Get fetches a workspace by name.
func Get(c *drycc.Client, name string) (api.Workspace, error) {
	u := fmt.Sprintf("/v2/workspaces/%s", name)
	res, reqErr := c.Request("GET", u, nil)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.Workspace{}, reqErr
	}
	defer res.Body.Close()

	workspace := api.Workspace{}
	if err := json.NewDecoder(res.Body).Decode(&workspace); err != nil {
		return api.Workspace{}, err
	}

	return workspace, reqErr
}

// Update updates workspace attributes.
func Update(c *drycc.Client, name, email string) (api.Workspace, error) {
	u := fmt.Sprintf("/v2/workspaces/%s", name)
	req := api.WorkspaceUpdateRequest{Email: email}
	body, err := json.Marshal(req)
	if err != nil {
		return api.Workspace{}, err
	}

	res, reqErr := c.Request("PATCH", u, body)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.Workspace{}, reqErr
	}
	defer res.Body.Close()

	workspace := api.Workspace{}
	if err = json.NewDecoder(res.Body).Decode(&workspace); err != nil {
		return api.Workspace{}, err
	}

	return workspace, reqErr
}

// Delete removes a workspace.
func Delete(c *drycc.Client, name string) error {
	u := fmt.Sprintf("/v2/workspaces/%s", name)
	res, err := c.Request("DELETE", u, nil)
	if err == nil {
		res.Body.Close()
	}
	return err
}
