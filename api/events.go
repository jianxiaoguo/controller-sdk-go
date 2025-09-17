package api

// AppEvent represents an event in the system.
type AppEvent struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
	Created string `json:"created"`
}

// AppEvents is a collection of AppEvent.
type AppEvents []AppEvent
