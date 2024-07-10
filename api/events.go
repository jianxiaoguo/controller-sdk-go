package api

// Event defines the structure of event.
type AppEvent struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
	Created string `json:"created"`
}

type AppEvents []AppEvent
