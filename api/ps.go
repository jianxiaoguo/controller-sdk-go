package api

// ProcessType represents the key/value mappings of a process type to a process inside
// a Heroku Procfile.
//
// See https://devcenter.heroku.com/articles/procfile
type ProcessType map[string]string

// Pods defines the structure of a process.
type Pods struct {
	Release  string `json:"release"`
	Type     string `json:"type"`
	Name     string `json:"name"`
	State    string `json:"state"`
	Ready    string `json:"ready"`
	Restarts int    `json:"restarts"`
	Started  string `json:"started"`
}

// PodsList defines a collection of app pods.
type PodsList []Pods

func (p PodsList) Len() int           { return len(p) }
func (p PodsList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PodsList) Less(i, j int) bool { return p[i].Name < p[j].Name }

// Command defines a command of app exec.
type Command struct {
	Tty     bool     `json:"tty"`
	Stdin   bool     `json:"stdin"`
	Command []string `json:"command"`
}

// PodType holds pods of the same type.
type PodType struct {
	Ptype    string
	PodsList PodsList
}

// PodTypes holds groups of pods organized by type.
type PodTypes []PodType

// AppLogsRequest is the definition of websocket /v2/apps/<app id>/logs
type PodLogsRequest struct {
	Lines     int    `json:"lines"`
	Follow    bool   `json:"follow"`
	Container string `json:"container"`
}

// Start is the definition of POST /v2/apps/<app_id>/stop or POST /v2/apps/<app_id>/start.
type Types struct {
	Types []string `json:"types,omitempty"`
}

func (p PodTypes) Len() int           { return len(p) }
func (p PodTypes) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PodTypes) Less(i, j int) bool { return p[i].Ptype < p[j].Ptype }

// ContainerState defines a container state.
type ContainerState struct {
	Container    string                            `json:"container"`
	Image        string                            `json:"image"`
	Command      []string                          `json:"command"`
	Args         []string                          `json:"args"`
	State        map[string]map[string]interface{} `json:"state"`
	LastState    map[string]map[string]interface{} `json:"lastState"`
	Ready        bool                              `json:"ready"`
	RestartCount int                               `json:"restartCount"`
	Status       string                            `json:"status"`
	Reason       string                            `json:"reason"`
	Message      string                            `json:"message"`
}

// PodState defines a collection of container state.
type PodState []ContainerState

type PodIDs struct {
	PodIDs string `json:"pod_ids"`
}
