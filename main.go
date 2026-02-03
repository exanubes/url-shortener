package main

import (
	"context"
	"log"

	"github.com/exanubes/url-shortener/internal/app/services/expiration"
	"github.com/exanubes/url-shortener/internal/app/services/shortcode"
	createshorturl "github.com/exanubes/url-shortener/internal/app/usecases/create_short_url"
	resolveurl "github.com/exanubes/url-shortener/internal/app/usecases/resolve_url"
	"github.com/exanubes/url-shortener/internal/domain"
	"github.com/exanubes/url-shortener/internal/infrastructure/api/http"
	encoding "github.com/exanubes/url-shortener/internal/infrastructure/encoding/base_62"
	"github.com/exanubes/url-shortener/internal/infrastructure/event"
	"github.com/exanubes/url-shortener/internal/infrastructure/persistence/dynamodb"
)

func main() {
	ctx := context.Background()
	client, err := dynamodb.NewLocalStackClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	table := dynamodb.NewRepository(client)
	event_bus := event.NewBus(func(event domain.Event) error { return table.Visit(ctx, event.(domain.LinkVisited)) })
	encoder := encoding.New()
	token_generator := shortcode.NewGenerator(int64(7))
	policy_factory := createshorturl.NewRetryPolicyFactory(3)
	expiration_factory := expiration.NewFactory()
	shortcodes_service := shortcode.NewService(token_generator, encoder)
	create_short_url_use_case := createshorturl.New(table, shortcodes_service, policy_factory, expiration_factory)
	visit_url_use_case := resolveurl.New(table, table, event_bus)

	driver := http.NewHttpDriver(create_short_url_use_case, visit_url_use_case)
	event_bus.Start(ctx)
	if err := driver.Run(ctx, http.DefaultConfig); err != nil {
		log.Fatal(err)
	}
}
