package sqs

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/exanubes/url-shortener/internal/app/services/messaging"
)

type Transport struct {
	client    *Client
	queue_url string
}

func NewSqsTransport(client *Client, queue_url string) *Transport {
	return &Transport{client, queue_url}
}

func (transport *Transport) Publish(ctx context.Context, message messaging.Envelope) error {
	raw_message, err := json.Marshal(message)

	if err != nil {
		return err
	}
	body := string(raw_message)

	_, err = transport.client.queue.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    &transport.queue_url,
		MessageBody: &body,
	})

	return err
}
