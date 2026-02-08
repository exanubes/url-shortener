package main

import (
	"context"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/exanubes/url-shortener/internal/infrastructure/processor/cloudfront"
)

func main() {
	ctx := context.Background()
	handler := cloudfront.NewHandler()
	lambda.StartWithOptions(func(ctx context.Context, req events.KinesisEvent) (events.KinesisEventResponse, error) {
		return handler.Handle(ctx, req)
	}, lambda.WithContext(ctx))
}
