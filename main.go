package main

import (
	"context"
	"log"

	"github.com/exanubes/url-shortener/internal/app/services/expiration"
	"github.com/exanubes/url-shortener/internal/app/services/shortcode"
	createshorturl "github.com/exanubes/url-shortener/internal/app/usecases/create_short_url"
	visitshorturl "github.com/exanubes/url-shortener/internal/app/usecases/visit_short_url"
	encoding "github.com/exanubes/url-shortener/internal/infrastructure/encoding/base_62"
	"github.com/exanubes/url-shortener/internal/infrastructure/http"
	"github.com/exanubes/url-shortener/internal/infrastructure/persistence/inmemory"
)

func main() {
	provider := inmemory.NewInmemoryRepository()
	encoder := encoding.New()
	token_generator := shortcode.NewGenerator(int64(7))
	policy_factory := createshorturl.NewRetryPolicyFactory(3)
	expiration_factory := expiration.NewFactory()
	shortcodes_service := shortcode.NewService(token_generator, encoder)
	create_short_url_use_case := createshorturl.New(provider, shortcodes_service, policy_factory, expiration_factory)
	visit_url_use_case := visitshorturl.New(provider)

	driver := http.NewHttpDriver(create_short_url_use_case, visit_url_use_case)
	ctx := context.Background()
	if err := driver.Run(ctx, http.DefaultConfig); err != nil {
		log.Fatal(err)
	}
}
