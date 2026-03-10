package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/exanubes/url-shortener/internal/infrastructure/clock"
	"github.com/exanubes/url-shortener/internal/infrastructure/persistence/dynamodb"
	"github.com/exanubes/url-shortener/internal/infrastructure/processor/cloudfront"
)

func main() {
	ctx := context.Background()
	client, err := dynamodb.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	clockInstance := clock.NewClock()
	table := dynamodb.NewRepository(client, clockInstance)
	handler := cloudfront.NewHandler(table)

	lambda.StartWithOptions(func(ctx context.Context, req events.KinesisEvent) (events.KinesisEventResponse, error) {
		return handler.Handle(ctx, req)
	}, lambda.WithContext(ctx))
}
