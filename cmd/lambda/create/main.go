package main

import (
	"context"
	"hash/fnv"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/exanubes/url-shortener/internal/app/services/expiration"
	"github.com/exanubes/url-shortener/internal/app/services/shortcode"
	createshorturl "github.com/exanubes/url-shortener/internal/app/usecases/create_short_url"
	"github.com/exanubes/url-shortener/internal/infrastructure/api/lambda/create"
	"github.com/exanubes/url-shortener/internal/infrastructure/clock"
	encoding "github.com/exanubes/url-shortener/internal/infrastructure/encoding/base_62"
	"github.com/exanubes/url-shortener/internal/infrastructure/persistence/dynamodb"
)

const feistel_key uint64 = 0x8c3f19d2e4a761b5
const epoch = "2026-02-21T11:04:57.497"

func main() {
	ctx := context.Background()
	client, err := dynamodb.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	hash := fnv.New64a()
	log_stream_name := os.Getenv("AWS_LAMBDA_LOG_STREAM_NAME")
	hash.Write([]byte(log_stream_name))

	epoch_date, err := time.Parse("2006-01-02T15:04:05.000", epoch)
	if err != nil {
		log.Fatal(err)
	}

	clock := clock.NewClock()
	token_generator := shortcode.NewSnowflakeGenerator(hash, epoch_date, clock)
	table := dynamodb.NewRepository(client, clock)
	policy_factory := createshorturl.NewRetryPolicyFactory(3)
	expiration_factory := expiration.NewFactory()
	encoder := encoding.New()
	scrambler := shortcode.NewFeistel(feistel_key)
	shortcodes_service := shortcode.NewService(token_generator, scrambler, encoder)
	create_short_url_use_case := createshorturl.New(table, shortcodes_service, policy_factory, expiration_factory, clock)
	handler := create.NewHandler(create_short_url_use_case)
	lambda.StartWithOptions(func(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
		return handler.Handle(ctx, req), nil
	}, lambda.WithContext(ctx))
}
