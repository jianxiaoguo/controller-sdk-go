// Package members provides methods for managing workspace members.
package members

import (
	"encoding/json"
	"fmt"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// List lists members in a workspace.
func List(c *drycc.Client, workspace string, results int) (api.WorkspaceMembers, int, error) {
	u := fmt.Sprintf("/v2/workspaces/%s/members", workspace)
	body, count, reqErr := c.LimitedRequest(u, results)

	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.WorkspaceMember{}, -1, reqErr
	}

	var members []api.WorkspaceMember
	if err := json.Unmarshal([]byte(body), &members); err != nil {
		return []api.WorkspaceMember{}, -1, err
	}

	return members, count, reqErr
}

// Get fetches a workspace member by username.
func Get(c *drycc.Client, workspace, user string) (api.WorkspaceMember, error) {
	u := fmt.Sprintf("/v2/workspaces/%s/members/%s", workspace, user)
	res, reqErr := c.Request("GET", u, nil)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.WorkspaceMember{}, reqErr
	}
	defer res.Body.Close()

	member := api.WorkspaceMember{}
	if err := json.NewDecoder(res.Body).Decode(&member); err != nil {
		return api.WorkspaceMember{}, err
	}

	return member, reqErr
}

// Update updates a workspace member role/alerts.
func Update(c *drycc.Client, workspace, user, role string, alerts *bool) (api.WorkspaceMember, error) {
	u := fmt.Sprintf("/v2/workspaces/%s/members/%s", workspace, user)
	req := api.WorkspaceMemberUpdateRequest{Role: role, Alerts: alerts}
	body, err := json.Marshal(req)
	if err != nil {
		return api.WorkspaceMember{}, err
	}

	res, reqErr := c.Request("PATCH", u, body)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.WorkspaceMember{}, reqErr
	}
	defer res.Body.Close()

	member := api.WorkspaceMember{}
	if err = json.NewDecoder(res.Body).Decode(&member); err != nil {
		return api.WorkspaceMember{}, err
	}

	return member, reqErr
}

// Delete removes a workspace member.
func Delete(c *drycc.Client, workspace, user string) error {
	u := fmt.Sprintf("/v2/workspaces/%s/members/%s", workspace, user)
	res, err := c.Request("DELETE", u, nil)
	if err == nil {
		res.Body.Close()
	}
	return err
}
