package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/exanubes/url-shortener/internal/infrastructure/event_handlers/link_visited"
	"github.com/exanubes/url-shortener/internal/infrastructure/persistence/dynamodb"
)

func main() {
	ctx := context.Background()
	client, err := dynamodb.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	table := dynamodb.NewRepository(client)
	handler := link_visited.NewHandler(table)

	lambda.StartWithOptions(func(ctx context.Context, req events.SQSEvent) (events.SQSEventResponse, error) {
		return handler.Handle(ctx, req)
	}, lambda.WithContext(ctx))
}
