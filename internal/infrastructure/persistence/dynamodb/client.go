package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type client struct {
	db *dynamodb.Client
}

func NewClient(ctx context.Context) (*client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		return nil, err
	}

	return &client{
		db: dynamodb.NewFromConfig(cfg),
	}, nil
}

func (c *client) Queries() *queries {
	return new_queries(c.db)
}
