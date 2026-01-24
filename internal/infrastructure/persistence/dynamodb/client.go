package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
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

func NewLocalStackClient(ctx context.Context) (*client, error) {
	// Hardcoded LocalStack configuration
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               "http://localhost:4566",
			HostnameImmutable: true,
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			"test", // Access Key ID
			"test", // Secret Access Key
			"",     // Session Token (empty)
		)),
	)

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
