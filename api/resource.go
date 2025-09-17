package api

// ResourcePlan is the structure of an app's resource plan.
type ResourcePlan struct {
	// ID is a unique string for resource plan.
	ID string `json:"id,omitempty"`
	// Name is a unique string for resource plan.
	Name string `json:"name,omitempty"`
	// Description is a detailed description of the resource plan
	Description string `json:"description,omitempty"`
}

// ResourcePlans is a collection of ResourcePlan.
type ResourcePlans []ResourcePlan

// ResourceService is the structure of an app's resource service.
type ResourceService struct {
	// ID is a unique string for resource service.
	ID string `json:"id,omitempty"`
	// Name is a unique string for resource service.
	Name string `json:"name,omitempty"`
	// Updatable is the plans of the current resource can be upgraded
	Updateable bool `json:"updateable,omitempty"`
}

// ResourceServices is a collection of ResourceService.
type ResourceServices []ResourceService

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
	// Resource instance message
	Message string `json:"message,omitempty"`
}

// Resources is a collection of Resource.
type Resources []Resource

// ResourceBinding is the definition of PATCH /v2/apps/<app_id>/resources/<name>/binding/.
type ResourceBinding struct {
	BindAction string `json:"bind_action,omitempty"`
}
