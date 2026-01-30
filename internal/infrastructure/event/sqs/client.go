package sqs

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Client struct {
	queue *sqs.Client
}

func NewClient(ctx context.Context) *Client {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		queue: sqs.NewFromConfig(cfg),
	}
}
