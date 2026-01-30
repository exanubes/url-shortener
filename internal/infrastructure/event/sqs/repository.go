package sqs

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/exanubes/url-shortener/internal/domain"
	"github.com/exanubes/url-shortener/internal/infrastructure/event/sqs/internal"
	"github.com/google/uuid"
)

type Repository struct {
	client    *Client
	queue_url string
}

func NewRepository(client *Client, queue_url string) *Repository {
	return &Repository{client, queue_url}
}

func (repository *Repository) Publish(ctx context.Context, event domain.LinkVisited) error {

	envelope, err := marshal(event)

	if err != nil {
		return err
	}

	raw_message, err := json.Marshal(envelope)

	if err != nil {
		return err
	}

	message := string(raw_message)

	_, err = repository.client.queue.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &repository.queue_url,
		MessageBody: &message,
	})

	return err
}

func marshal(event domain.LinkVisited) (internal.Envelope, error) {
	payload := internal.LinkVisitedV1{
		ShortCode: event.ShortCode,
		VisitedAt: event.VisitedAt,
		IpAddress: event.IpAddress,
	}

	raw_message, err := json.Marshal(payload)

	if err != nil {
		return internal.Envelope{}, err
	}

	return internal.Envelope{
		ID:        uuid.NewString(),
		Type:      internal.LINK_VISITED_EVENT,
		Version:   1,
		Timestamp: time.Now(),
		Payload:   raw_message,
	}, err
}
