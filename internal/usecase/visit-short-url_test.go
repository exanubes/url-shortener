package usecase

import (
	"context"
	"testing"

	"github.com/exanubes/url-shortener/internal/app/encoder"
	"github.com/exanubes/url-shortener/internal/domain"
	"github.com/exanubes/url-shortener/internal/infrastructure/persistence/inmemory"
)

func TestVisitShortUrl(t *testing.T) {
	provider := inmemory.NewInmemoryRepository()
	codec := encoder.New()
	long_url := "https://exanubes.com"
	short_url := "2TX"
	id := uint64(11_157)
	provider.Save(context.TODO(), domain.Url{ID: id, Short: short_url, Long: long_url})

	usecase := NewVisitShortUrl(provider, codec)
	result, _ := usecase.Execute(short_url)

	if result != long_url {
		t.Fatalf("Expected: %s, received: %s", long_url, result)
	}
}
