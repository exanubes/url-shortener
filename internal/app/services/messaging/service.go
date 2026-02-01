package messaging

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/exanubes/url-shortener/internal/domain"
)

type MessagingService struct {
	publisher Publisher
}

func NewService(
	publisher Publisher,
) *MessagingService {
	return &MessagingService{publisher}
}

func (service MessagingService) Publish(ctx context.Context, event domain.Event) error {
	var event_type EventType
	var version int
	var payload json.RawMessage
	var err error

	switch event := event.(type) {
	case domain.LinkVisited:
		event_type = LINK_VISITED_EVENT
		version = 1
		payload, err = json.Marshal(LinkVisitedV1{
			ShortCode: event.ShortCode,
			VisitedAt: event.VisitedAt,
			IpAddress: event.IpAddress,
		})
	default:
		return errors.New("Unknown event type")
	}

	if err != nil {
		return err
	}
	envelope := NewEnvelope(event_type, version, payload)

	return service.publisher.Publish(ctx, envelope)
}
