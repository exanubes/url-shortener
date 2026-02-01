package resolveurl_test

import (
	"context"
	"testing"

	resolveurl "github.com/exanubes/url-shortener/internal/app/usecases/resolve_url"
	"github.com/exanubes/url-shortener/internal/domain"
	"github.com/exanubes/url-shortener/internal/infrastructure/persistence/inmemory"
)

func TestResolveUrl(t *testing.T) {
	provider := inmemory.NewInmemoryRepository()
	long_url, _ := domain.NewUrl("https://exanubes.com")
	short_code, _ := domain.NewShortCode("2TX", 7, "0")
	provider.Write(context.TODO(), short_code, long_url)

	usecase := resolveurl.New(provider)
	result, _ := usecase.Execute(context.TODO(), short_code)

	if result.String() != long_url.String() {
		t.Fatalf("Expected: %s, received: %s", long_url.String(), result.String())
	}
}
