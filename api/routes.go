package api

// Route is the structure of an app's route.
type Route struct {
	// Owner is the app owner. It cannot be updated with AppSettings.Set(). See app.Transfer().
	Owner string `json:"owner,omitempty"`
	// App is the app name. It cannot be updated at all right now.
	App string `json:"app,omitempty"`
	// Created is the time that the application settings was created and cannot be updated.
	Created string `json:"created,omitempty"`
	// Updated is the last time the application settings was changed and cannot be updated.
	Updated string `json:"updated,omitempty"`
	// UUID is a unique string reflecting the application settings in its current state.
	// It changes every time the application settings is changed and cannot be updated.
	UUID       string      `json:"uuid,omitempty"`
	Name       string      `json:"name,omitempty"`
	Kind       string      `json:"kind,omitempty"`
	ParentRefs []ParentRef `json:"parent_refs,omitempty"`
	Rules      []RouteRule `json:"rules,omitempty"`
}

// ParentRef represents a reference to a parent gateway.
type ParentRef struct {
	Name string `json:"name,omitempty"`
	Port int    `json:"port,omitempty"`
}

// RouteRule represents a rule in a route configuration.
type RouteRule map[string]interface{}

// Routes defines a collection of Route objects.
type Routes []Route

// RouteCreateRequest is the structure of POST /v2/app/<app_id>/routes/.

// RouteCreateRequest is the structure of POST /v2/app/<app_id>/routes/.
type RouteCreateRequest struct {
	Name  string             `json:"name,omitempty"`
	Kind  string             `json:"kind,omitempty"`
	Rules []RequestRouteRule `json:"rules,omitempty"`
}

// BackendRefRequest represents a backend reference in a route request.
type BackendRefRequest struct {
	Kind   string `json:"kind,omitempty"`
	Name   string `json:"name,omitempty"`
	Port   int32  `json:"port,omitempty"`
	Weight int32  `json:"weight,omitempty"`
}

// RequestRouteRule represents a route rule in a request.
type RequestRouteRule struct {
	BackendRefs []BackendRefRequest `json:"backendRefs,omitempty"`
}

// RouteAttachRequest is the structure of PATCH /v2/apps/(?P<id>{})/routes/(?P<name>{})/attach/?$.
type RouteAttachRequest struct {
	Port    int    `json:"port,omitempty"`
	Gateway string `json:"gateway,omitempty"`
}

// RouteDetachRequest is the structure of PATCH /v2/apps/(?P<id>{})/routes/(?P<name>{})/detach/?$.
type RouteDetachRequest struct {
	Port    int    `json:"port,omitempty"`
	Gateway string `json:"gateway,omitempty"`
}
