package internal

import (
	"encoding/json"
	"time"
)

type Envelope struct {
	ID        string          `json:"id"`
	Type      EventType       `json:"type"`
	Version   int             `json:"version"`
	Timestamp time.Time       `json:"timestamp"`
	Payload   json.RawMessage `json:"payload"`
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
