package ohsse

import (
	"encoding/json"

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

type OHEvent struct {
	Current  OHResourceVersion `json:"current,omitempty"`
	Previous OHResourceVersion `json:"previous,omitempty"`
}

type OHIaC struct {
	Status    string `json:"status,omitempty"`
	Framework string `json:"framework,omitempty"`
}

type OHResourceVersion struct {
	Resource    OHVagueResource `json:"resource"`
	Attribution OHAttribution   `json:"attribution"`
	IaC         OHIaC           `json:"iac,omitempty"`
	Meta        *Metadata       `json:"meta,omitempty"` // Replaces OH__Meta, but temporarily both here for legacy records

}

type OHVagueResource struct {
	MetaLegacy *Metadata `json:"OH__Meta,omitempty"` // Replaces OH__Meta, but temporarily both here for legacy records
	Meta       *Metadata `json:"Meta,omitempty"`     // This will disappear, and is replaced by the Meta field, but temporarily both here for legacy records
	Raw        json.RawMessage
}

// Metadata, previously OH__Meta, metadata we have added
type Metadata struct {
	ARN          string `json:"arn,omitempty"`
	Region       string `json:"region,omitempty"`
	ResourceType string `json:"resource_type,omitempty"`
}

type OHAttribution struct {
	IP        string `json:"IP,omitempty"`
	Agent     string `json:"Agent,omitempty"`
	Platform  string `json:"Platform,omitempty"`
	UserAgent string `json:"UserAgent,omitempty"`
	Principal struct {
		Type          string   `json:"Type,omitempty"`
		ID            string   `json:"ID,omitempty"`
		ARN           string   `json:"ARN,omitempty"`
		RoleHistory   []string `json:"RoleHistory,omitempty"`
		CloudProvider string   `json:"CloudProvider,omitempty"`
	} `json:"Principal,omitempty"`
	Version int `json:"version"`
}

// Define handlers type so that a function can be passed to receieve events
type StreamHandler func(SSE_Entry)
type CloudEventHandler func(cloudevents.Event)
