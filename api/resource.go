package api

// Binding is the definition of PATCH /v2/apps/<app_id>/resources/<name>/binding/.
type Binding struct {
	BindAction string `json:"bind_action,omitempty"`
}

// Resource is the structure of an app's resource.
type Resource struct {
	// Owner is the app owner.
	Owner string `json:"owner,omitempty"`
	// App is the app the tls settings apply to and cannot be updated.
	App string `json:"app,omitempty"`
	// Created is the time that the resource was created and cannot be updated.
	Created string `json:"created,omitempty"`
	// Updated is the last time the TLS settings was changed and cannot be updated.
	Updated string `json:"updated,omitempty"`
	// UUID is a unique string reflecting the resource in its current state.
	// It changes every time the resource is changed and cannot be updated.
	UUID string `json:"uuid,omitempty"`
	// Resource's name
	Name string `json:"name,omitempty"`
	// Resource's Plan
	Plan string `json:"plan,omitempty"`
	// Resource connet info
	Data map[string]interface{} `json:"data,omitempty"`
	// Resource's status
	Status string `json:"status,omitempty"`
	// Resource's binding status
	Binding string `json:"binding,omitempty"`
	// Resource Options
	Options map[string]interface{} `json:"options,omitempty"`
}

type Resources []Resource
