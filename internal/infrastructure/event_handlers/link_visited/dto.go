package link_visited

import (
	"encoding/json"
	"errors"

	"github.com/exanubes/url-shortener/internal/app/services/messaging"
	"github.com/exanubes/url-shortener/internal/domain"
)

func parse_message(msg string) (messaging.LinkVisitedV1, error) {
	var envelope messaging.Envelope
	var message messaging.LinkVisitedV1

	err := json.Unmarshal([]byte(msg), &envelope)

	if err != nil {
		return message, err
	}

	if envelope.Type != messaging.LINK_VISITED_EVENT {
		return message, ErrWrongMessageType
	}

	err = json.Unmarshal(envelope.Payload, &message)

	return message, err
}

var ErrWrongMessageType = errors.New("Invalid message type")

func map_to_domain_event(msg messaging.LinkVisitedV1) domain.LinkVisited {
	return domain.LinkVisited{
		ShortCode: msg.ShortCode,
		VisitedAt: msg.VisitedAt,
		IpAddress: msg.IpAddress,
	}
}
