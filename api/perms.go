package api

// UserPermResponse is the definition of GET /v2/perms/rules/.
type UserPermResponse struct {
	App         string   `json:"app"`
	Username    string   `json:"username"`
	Permissions []string `json:"permissions"`
}

// UserPermRequest is the definition of a requst on /v2/perms/rules/.
type UserPermRequest struct {
	Username    string `json:"username"`
	Permissions string `json:"permissions"`
}
