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
	Ptype      string      `json:"ptype,omitempty"`
	Kind       string      `json:"kind,omitempty"`
	Port       int         `json:"port,omitempty"`
	ParentRefs []ParentRef `json:"parent_refs,omitempty"`
}

type ParentRef struct {
	Name string `json:"name,omitempty"`
	Port int    `json:"port,omitempty"`
}

// // Routes defines a collection of Route objects.
type Routes []Route

// RouteCreateRequest is the structure of POST /v2/app/<app id>/routes/.
type RouteCreateRequest struct {
	Name  string `json:"name,omitempty"`
	Ptype string `json:"ptype,omitempty"`
	Port  int    `json:"port,omitempty"`
	Kind  string `json:"kind,omitempty"`
}

// RouteAttackRequest is the structure of PATCH /v2/apps/(?P<id>{})/routes/(?P<name>{})/attach/?$.
type RouteAttackRequest struct {
	Port    int    `json:"port,omitempty"`
	Gateway string `json:"gateway,omitempty"`
}

// RouteDetackRequest is the structure of PATCH /v2/apps/(?P<id>{})/routes/(?P<name>{})/detach/?$.
type RouteDetackRequest struct {
	Port    int    `json:"port,omitempty"`
	Gateway string `json:"gateway,omitempty"`
}

// RouteRule is the structure of GET RESPONSE /v2/apps/(?P<id>{})/routes/(?P<name>{})/rules/?$.
type RouteRule struct {
	Name  string      `json:"name,omitempty"`
	Ptype string      `json:"ptype,omitempty"`
	Kind  string      `json:"kind,omitempty"`
	Rules interface{} `json:"rules,omitempty"`
}
