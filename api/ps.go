package api

import "github.com/drycc/controller-sdk-go/pkg/time"

// ProcessType represents the key/value mappings of a process type to a process inside
// a Heroku Procfile.
//
// See https://devcenter.heroku.com/articles/procfile
type ProcessType map[string]string

// Pods defines the structure of a process.
type Pods struct {
	Release  string    `json:"release,omitempty"`
	Type     string    `json:"type,omitempty"`
	Name     string    `json:"name,omitempty"`
	State    string    `json:"state,omitempty"`
	Started  time.Time `json:"started,omitempty"`
	Replicas string    `json:"replicas,omitempty"`
}

// PodsList defines a collection of app pods.
type PodsList []Pods

func (p PodsList) Len() int           { return len(p) }
func (p PodsList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PodsList) Less(i, j int) bool { return p[i].Name < p[j].Name }

// PodType holds pods of the same type.
type PodType struct {
	Type     string
	PodsList PodsList
	Replicas string
	Status   string
}

// PodTypes holds groups of pods organized by type.
type PodTypes []PodType

// Start is the definition of POST /v2/apps/<app_id>/stop or POST /v2/apps/<app_id>/start.
type Types struct {
	Types []string `json:"types,omitempty"`
}

func (p PodTypes) Len() int           { return len(p) }
func (p PodTypes) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p PodTypes) Less(i, j int) bool { return p[i].Type < p[j].Type }
