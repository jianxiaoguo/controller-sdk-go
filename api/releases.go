package api

// Release is the definition of the release object.
type Release struct {
	App        string      `json:"app"`
	State      string      `json:"state"`
	Build      string      `json:"build,omitempty"`
	Config     string      `json:"config"`
	Created    string      `json:"created"`
	Owner      string      `json:"owner"`
	Summary    string      `json:"summary"`
	Exception  string      `json:"exception"`
	Conditions []Condition `json:"conditions"`
	Updated    string      `json:"updated"`
	UUID       string      `json:"uuid"`
	Version    int         `json:"version"`
}

// ReleaseRollback is the defenition of POST /v2/apps/<app id>/releases/.
type ReleaseRollback struct {
	Version int    `json:"version"`
	Ptypes  string `json:"ptypes"`
}

type Condition struct {
	State     string   `json:"state"`
	Action    string   `json:"action"`
	Ptypes    []string `json:"ptypes"`
	Exception string   `json:"exception"`
	Created   string   `json:"created"`
}
