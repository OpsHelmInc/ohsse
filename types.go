package ohsse

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// Struct that is passed to a handler
type SSE_Entry struct {
	Comment string
	Data    string
	Event   string
	ID      string
	Retry   string
	Unknown map[string]string
}

// Define handlers type so that a function can be passed to receieve events
type StreamHandler func(SSE_Entry)
type CloudEventHandler func(cloudevents.Event)
