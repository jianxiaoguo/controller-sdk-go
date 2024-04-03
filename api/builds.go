package api

// Build is the structure of the build object.
type Build struct {
	App        string                 `json:"app"`
	Created    string                 `json:"created"`
	Dockerfile string                 `json:"dockerfile,omitempty"`
	Image      string                 `json:"image,omitempty"`
	Stack      string                 `json:"stack,omitempty"`
	Owner      string                 `json:"owner"`
	Procfile   map[string]string      `json:"procfile"`
	Dryccfile  map[string]interface{} `json:"dryccfile"`
	Sha        string                 `json:"sha,omitempty"`
	Updated    string                 `json:"updated"`
	UUID       string                 `json:"uuid"`
}

// CreateBuildRequest is the structure of POST /v2/apps/<app id>/builds/.
type CreateBuildRequest struct {
	Image     string                 `json:"image"`
	Stack     string                 `json:"stack,omitempty"`
	Procfile  map[string]string      `json:"procfile,omitempty"`
	Dryccfile map[string]interface{} `json:"dryccfile,omitempty"`
}

// BuildHookRequest is a hook request to create a new build.
type BuildHookRequest struct {
	Sha        string                 `json:"sha"`
	User       string                 `json:"receive_user"`
	App        string                 `json:"receive_repo"`
	Image      string                 `json:"image"`
	Stack      string                 `json:"stack"`
	Procfile   ProcessType            `json:"procfile"`
	Dockerfile string                 `json:"dockerfile"`
	Dryccfile  map[string]interface{} `json:"dryccfile"`
}
