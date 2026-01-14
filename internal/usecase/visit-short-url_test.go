package usecase

import (
	"context"
	"testing"

	"github.com/exanubes/url-shortener/internal/domain"
	"github.com/exanubes/url-shortener/internal/infrastructure/persistence/inmemory"
)

func TestVisitShortUrl(t *testing.T) {
	provider := inmemory.NewInmemoryRepository()
	long_url, _ := domain.NewUrl("https://exanubes.com")
	short_code, _ := domain.NewShortCode("2TX", 7, "0")
	provider.Save(context.TODO(), long_url, short_code)

	usecase := NewVisitShortUrl(provider)
	result, _ := usecase.Execute(context.TODO(), short_code.String())

	if result.String() != long_url.String() {
		t.Fatalf("Expected: %s, received: %s", long_url.String(), result.String())
	}
}
