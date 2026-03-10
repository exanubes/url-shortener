package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	resolveurl "github.com/exanubes/url-shortener/internal/app/usecases/resolve_url"
	"github.com/exanubes/url-shortener/internal/infrastructure/api/lambda/resolve"
	"github.com/exanubes/url-shortener/internal/infrastructure/clock"
	"github.com/exanubes/url-shortener/internal/infrastructure/persistence/dynamodb"
)

func main() {
	ctx := context.Background()
	client, err := dynamodb.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	clockInstance := clock.NewClock()

	table := dynamodb.NewRepository(client, clockInstance)
	visit_url_use_case := resolveurl.New(table, table, clockInstance)
	handler := resolve.NewHandler(visit_url_use_case, clockInstance)

	lambda.StartWithOptions(func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		return handler.Handle(ctx, req), nil
	}, lambda.WithContext(ctx))
}
