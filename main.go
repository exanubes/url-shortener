package main

import (
	"context"
	"log"

	encoder "github.com/exanubes/url-shortener/internal/app/encoder/base_62"
	"github.com/exanubes/url-shortener/internal/app/policy"
	shortcodes "github.com/exanubes/url-shortener/internal/app/short_codes"
	token_space "github.com/exanubes/url-shortener/internal/app/token_space_generator/base_62"
	"github.com/exanubes/url-shortener/internal/drivers"
	"github.com/exanubes/url-shortener/internal/infrastructure/persistence/inmemory"
	"github.com/exanubes/url-shortener/internal/usecase"
)

func main() {
	provider := inmemory.NewInmemoryRepository()
	encoder := encoder.New()
	token_generator := token_space.New(int64(7))
	policy_factory := policy.NewRetryPolicyFactory(3)
	shortcodes_generator := shortcodes.New(token_generator, encoder)
	create_short_url_use_case := usecase.NewCreateShortUrl(provider, shortcodes_generator, policy_factory)
	visit_url_use_case := usecase.NewVisitShortUrl(provider)

	driver := drivers.NewHttpDriver(create_short_url_use_case, visit_url_use_case)
	ctx := context.Background()
	if err := driver.Run(ctx, drivers.DefaultHttpConfig); err != nil {
		log.Fatal(err)
	}
}
