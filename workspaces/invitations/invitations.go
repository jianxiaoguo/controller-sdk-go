// Package invitations provides methods for managing workspace invitations.
package invitations

import (
	"encoding/json"
	"fmt"

	drycc "github.com/drycc/controller-sdk-go"
	"github.com/drycc/controller-sdk-go/api"
)

// List lists pending invitations in a workspace.
func List(c *drycc.Client, workspace string, results int) (api.WorkspaceInvitations, int, error) {
	u := fmt.Sprintf("/v2/workspaces/%s/invitations", workspace)
	body, count, reqErr := c.LimitedRequest(u, results)

	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return []api.WorkspaceInvitation{}, -1, reqErr
	}

	var invitations []api.WorkspaceInvitation
	if err := json.Unmarshal([]byte(body), &invitations); err != nil {
		return []api.WorkspaceInvitation{}, -1, err
	}

	return invitations, count, reqErr
}

// Create creates a workspace invitation.
func Create(c *drycc.Client, workspace, email string) (api.WorkspaceInvitation, error) {
	u := fmt.Sprintf("/v2/workspaces/%s/invitations", workspace)
	req := api.WorkspaceInvitationCreateRequest{Email: email}
	body, err := json.Marshal(req)
	if err != nil {
		return api.WorkspaceInvitation{}, err
	}

	res, reqErr := c.Request("POST", u, body)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.WorkspaceInvitation{}, reqErr
	}
	defer res.Body.Close()

	invitation := api.WorkspaceInvitation{}
	if err = json.NewDecoder(res.Body).Decode(&invitation); err != nil {
		return api.WorkspaceInvitation{}, err
	}

	return invitation, reqErr
}

// Get fetches an invitation by uid token.
func Get(c *drycc.Client, workspace, uid string) (api.WorkspaceInvitation, error) {
	u := fmt.Sprintf("/v2/workspaces/%s/invitations/%s", workspace, uid)
	res, reqErr := c.Request("GET", u, nil)
	if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
		return api.WorkspaceInvitation{}, reqErr
	}
	defer res.Body.Close()

	invitation := api.WorkspaceInvitation{}
	if err := json.NewDecoder(res.Body).Decode(&invitation); err != nil {
		return api.WorkspaceInvitation{}, err
	}

	return invitation, reqErr
}

// Delete revokes an invitation.
func Delete(c *drycc.Client, workspace, uid string) error {
	u := fmt.Sprintf("/v2/workspaces/%s/invitations/%s", workspace, uid)
	res, err := c.Request("DELETE", u, nil)
	if err == nil {
		res.Body.Close()
	}
	return err
}
