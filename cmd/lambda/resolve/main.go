package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/exanubes/url-shortener/internal/app/services/analytics"
	visitshorturl "github.com/exanubes/url-shortener/internal/app/usecases/visit_short_url"
	"github.com/exanubes/url-shortener/internal/domain"
	"github.com/exanubes/url-shortener/internal/infrastructure/api/lambda/resolve"
	"github.com/exanubes/url-shortener/internal/infrastructure/event"
	"github.com/exanubes/url-shortener/internal/infrastructure/persistence/dynamodb"
)

func main() {
	ctx := context.Background()
	client, err := dynamodb.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	table := dynamodb.NewRepository(client)
	processor := analytics.NewLinkVisitedProcessor(table)
	//TODO: implement EventPublisher with sqs
	event_bus := event.NewBus(func(event domain.LinkVisited) error { return processor.Handler(event) })
	visit_url_use_case := visitshorturl.New(table, table, event_bus)
	handler := resolve.NewHandler(visit_url_use_case)

	lambda.StartWithOptions(func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		return handler.Handle(ctx, req), nil
	}, lambda.WithContext(ctx))
}
