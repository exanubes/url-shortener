package visitshorturl_test

import (
	"context"
	"testing"

	visitshorturl "github.com/exanubes/url-shortener/internal/app/usecases/visit_short_url"
	"github.com/exanubes/url-shortener/internal/domain"
	"github.com/exanubes/url-shortener/internal/infrastructure/persistence/inmemory"
)

func TestVisitShortUrl(t *testing.T) {
	provider := inmemory.NewInmemoryRepository()
	long_url, _ := domain.NewUrl("https://exanubes.com")
	short_code, _ := domain.NewShortCode("2TX", 7, "0")
	provider.Write(context.TODO(), short_code, long_url)

	usecase := visitshorturl.New(provider)
	result, _ := usecase.Execute(context.TODO(), short_code)

	if result.String() != long_url.String() {
		t.Fatalf("Expected: %s, received: %s", long_url.String(), result.String())
	}
}
