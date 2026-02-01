package messaging

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Envelope struct {
	ID        string          `json:"id"`
	Type      EventType       `json:"type"`
	Version   int             `json:"version"`
	Timestamp time.Time       `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
}

func NewEnvelope(event_type EventType, version int, payload json.RawMessage) Envelope {
	return Envelope{
		ID:        uuid.NewString(),
		Type:      event_type,
		Version:   version,
		Timestamp: time.Now(),
		Payload:   payload,
	}
}

type LinkVisitedV1 struct {
	ShortCode string    `json:"id"`
	VisitedAt time.Time `json:"visited_at"`
	IpAddress string    `json:"ip_address,omitempty"`
}

type EventType string

const (
	LINK_VISITED_EVENT EventType = "LINK_VISITED"
)

type Publisher interface {
	Publish(context.Context, Envelope) error
}
