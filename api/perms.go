package api

// PermCodeResponse is the definition of GET /v2/perms/codes/.
type PermCodeResponse struct {
	Codename    string `json:"codename"`
	Description string `json:"description"`
}

// UserPermResponse is the definition of GET /v2/perms/rules/.
type UserPermResponse struct {
	ID       uint64 `json:"id"`
	Codename string `json:"codename"`
	Uniqueid string `json:"uniqueid"`
	Username string `json:"username"`
}

// UserPermRequest is the definition of a requst on /v2/perms/rules/.
type UserPermRequest struct {
	Codename string `json:"codename"`
	Uniqueid string `json:"uniqueid"`
	Username string `json:"username"`
}
