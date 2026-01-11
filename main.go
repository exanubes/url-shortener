package main

import (
	"context"
	"log"

	"github.com/exanubes/url-shortener/internal/app/encoder"
	"github.com/exanubes/url-shortener/internal/drivers"
	"github.com/exanubes/url-shortener/internal/infrastructure/persistence/inmemory"
	"github.com/exanubes/url-shortener/internal/usecase"
)

func main() {
	provider := inmemory.NewInmemoryRepository()
	codec := encoder.New()
	create_short_url_use_case := usecase.NewCreateShortUrl(provider, codec)
	visit_url_use_case := usecase.NewVisitShortUrl(provider, codec)

	driver := drivers.NewHttpDriver(create_short_url_use_case, visit_url_use_case)
	ctx := context.Background()
	if err := driver.Run(ctx, drivers.DefaultHttpConfig); err != nil {
		log.Fatal(err)
	}
}
