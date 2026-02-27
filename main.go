package main

import (
	"context"
	"hash/fnv"
	"log"
	"time"

	"github.com/exanubes/url-shortener/internal/app/services/expiration"
	"github.com/exanubes/url-shortener/internal/app/services/shortcode"
	createshorturl "github.com/exanubes/url-shortener/internal/app/usecases/create_short_url"
	resolveurl "github.com/exanubes/url-shortener/internal/app/usecases/resolve_url"
	"github.com/exanubes/url-shortener/internal/infrastructure/api/http"
	"github.com/exanubes/url-shortener/internal/infrastructure/clock"
	encoding "github.com/exanubes/url-shortener/internal/infrastructure/encoding/base_62"
	"github.com/exanubes/url-shortener/internal/infrastructure/persistence/dynamodb"
)

const feistel_key uint64 = 0x8c3f19d2e4a761b5
const epoch = "2026-02-21T11:04:57.497"

func main() {
	ctx := context.Background()
	client, err := dynamodb.NewLocalStackClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	table := dynamodb.NewRepository(client)
	encoder := encoding.New()
	clock := clock.NewClock()
	hash := fnv.New64a()
	hash.Write([]byte("local_http_server"))

	epoch_date, _ := time.Parse("2006-01-02T15:04:05.000", epoch)
	token_generator := shortcode.NewSnowflakeGenerator(hash, epoch_date, clock)
	scrambler := shortcode.NewFeistel(feistel_key)
	policy_factory := createshorturl.NewRetryPolicyFactory(3)
	expiration_factory := expiration.NewFactory()
	shortcodes_service := shortcode.NewService(token_generator, scrambler, encoder)
	create_short_url_use_case := createshorturl.New(table, shortcodes_service, policy_factory, expiration_factory)
	visit_url_use_case := resolveurl.New(table, table)

	driver := http.NewHttpDriver(create_short_url_use_case, visit_url_use_case)

	if err := driver.Run(ctx, http.DefaultConfig); err != nil {
		log.Fatal(err)
	}
}
