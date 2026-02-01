package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/exanubes/url-shortener/internal/app/services/messaging"
	resolveurl "github.com/exanubes/url-shortener/internal/app/usecases/resolve_url"
	"github.com/exanubes/url-shortener/internal/infrastructure/api/lambda/resolve"
	"github.com/exanubes/url-shortener/internal/infrastructure/event/sqs"
	"github.com/exanubes/url-shortener/internal/infrastructure/persistence/dynamodb"
)

func main() {
	ctx := context.Background()
	client, err := dynamodb.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	table := dynamodb.NewRepository(client)
	sqs_client := sqs.NewClient(ctx)
	sqs_transport := sqs.NewSqsTransport(sqs_client, get_queue_url())
	messaging_service := messaging.NewService(sqs_transport)
	visit_url_use_case := resolveurl.New(table, table, messaging_service)
	handler := resolve.NewHandler(visit_url_use_case)

	lambda.StartWithOptions(func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		return handler.Handle(ctx, req), nil
	}, lambda.WithContext(ctx))
}

func get_queue_url() string {
	queue_url, ok := os.LookupEnv("LINK_VISITED_QUEUE_URL")
	if !ok {
		log.Fatal("LINK_VISITED_QUEUE_URL is required")
	}

	return queue_url
}
