package api

import "github.com/drycc/controller-sdk-go/pkg/time"

// Event defines the structure of event.
type AppEvent struct {
	Reason  string    `json:"reason"`
	Message string    `json:"message"`
	Created time.Time `json:"created"`
}

type AppEvents []AppEvent
