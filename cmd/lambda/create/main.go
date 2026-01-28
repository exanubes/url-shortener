package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/exanubes/url-shortener/internal/app/services/expiration"
	"github.com/exanubes/url-shortener/internal/app/services/shortcode"
	createshorturl "github.com/exanubes/url-shortener/internal/app/usecases/create_short_url"
	"github.com/exanubes/url-shortener/internal/infrastructure/api/lambda/create"
	encoding "github.com/exanubes/url-shortener/internal/infrastructure/encoding/base_62"
	"github.com/exanubes/url-shortener/internal/infrastructure/persistence/dynamodb"
)

func main() {
	ctx := context.Background()
	client, err := dynamodb.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	table := dynamodb.NewRepository(client)
	token_generator := shortcode.NewGenerator(int64(7))
	policy_factory := createshorturl.NewRetryPolicyFactory(3)
	expiration_factory := expiration.NewFactory()
	encoder := encoding.New()
	shortcodes_service := shortcode.NewService(token_generator, encoder)
	create_short_url_use_case := createshorturl.New(table, shortcodes_service, policy_factory, expiration_factory)
	handler := create.NewHandler(create_short_url_use_case)
	lambda.StartWithOptions(func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		return handler.Handle(ctx, req), nil
	}, lambda.WithContext(ctx))
}
