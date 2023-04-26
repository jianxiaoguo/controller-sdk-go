package api

// Gateway is the structure of an app's gateways.
type Gateway struct {
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
	UUID      string     `json:"uuid,omitempty"`
	Name      string     `json:"name,omitempty"`
	Listeners []Listener `json:"listeners,omitempty"`
}

type Listener struct {
	Name          string      `json:"name,omitempty"`
	Port          int         `json:"port,omitempty"`
	Protocol      string      `json:"protocol,omitempty"`
	AllowedRoutes interface{} `json:"allowedRoutes,omitempty"`
}

// Gateways defines a collection of gateway objects.
type Gateways []Gateway

// GatewayCreateRequest is the structure of POST /v2/app/<app id>/gateways/.
type GatewayCreateRequest struct {
	Name     string `json:"name,omitempty"`
	Port     int    `json:"port,omitempty"`
	Protocol string `json:"protocol,omitempty"`
}

// GatewayRemoteRequest is the structure of Delete /v2/app/<app id>/gateways/.
type GatewayRemoveRequest struct {
	Name     string `json:"name,omitempty"`
	Port     int    `json:"port,omitempty"`
	Protocol string `json:"protocol,omitempty"`
}
