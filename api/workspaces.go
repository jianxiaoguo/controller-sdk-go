package api

// Workspace is the definition of the workspace object.
type Workspace struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Created string `json:"created"`
	Updated string `json:"updated"`
}

// Workspaces defines a collection of workspaces.
type Workspaces []Workspace

func (w Workspaces) Len() int           { return len(w) }
func (w Workspaces) Swap(i, j int)      { w[i], w[j] = w[j], w[i] }
func (w Workspaces) Less(i, j int) bool { return w[i].Name < w[j].Name }

// WorkspaceCreateRequest is the definition of POST /v2/workspaces.
type WorkspaceCreateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// WorkspaceUpdateRequest is the definition of PATCH /v2/workspaces/<name>.
type WorkspaceUpdateRequest struct {
	Email string `json:"email,omitempty"`
}

// WorkspaceMember is the definition of the workspace member object.
type WorkspaceMember struct {
	ID        int    `json:"id"`
	User      string `json:"user"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Alerts    bool   `json:"alerts"`
	Workspace string `json:"workspace"`
	Created   string `json:"created"`
	Updated   string `json:"updated"`
}

// WorkspaceMembers defines a collection of workspace member objects.
type WorkspaceMembers []WorkspaceMember

func (w WorkspaceMembers) Len() int           { return len(w) }
func (w WorkspaceMembers) Swap(i, j int)      { w[i], w[j] = w[j], w[i] }
func (w WorkspaceMembers) Less(i, j int) bool { return w[i].User < w[j].User }

// WorkspaceMemberUpdateRequest is the definition of PATCH /v2/workspaces/<name>/members/<user>.
type WorkspaceMemberUpdateRequest struct {
	Role   string `json:"role,omitempty"`
	Alerts *bool  `json:"alerts,omitempty"`
}

// WorkspaceInvitation is the definition of the workspace invitation object.
type WorkspaceInvitation struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Token     string `json:"token"`
	Inviter   string `json:"inviter"`
	Created   string `json:"created"`
	Accepted  bool   `json:"accepted"`
	Workspace string `json:"workspace"`
}

// WorkspaceInvitations defines a collection of workspace invitation objects.
type WorkspaceInvitations []WorkspaceInvitation

func (w WorkspaceInvitations) Len() int           { return len(w) }
func (w WorkspaceInvitations) Swap(i, j int)      { w[i], w[j] = w[j], w[i] }
func (w WorkspaceInvitations) Less(i, j int) bool { return w[i].Email < w[j].Email }

// WorkspaceInvitationCreateRequest is the definition of POST /v2/workspaces/<name>/invitations.
type WorkspaceInvitationCreateRequest struct {
	Email string `json:"email"`
}
