package usecase

import (
	"context"
	"testing"

	"github.com/exanubes/url-shortener/internal/app/encoder"
	"github.com/exanubes/url-shortener/internal/infrastructure/persistence/inmemory"
)

func TestCreateShortUrl(t *testing.T) {
	provider := inmemory.NewInmemoryRepository()
	codec := encoder.New()
	usecase := NewCreateShortUrl(provider, codec)
	expected := "2TY"
	long_url := "https://exanubes.com"
	result, err := usecase.Execute(long_url)

	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}

	if result != expected {
		t.Fatalf("Expected '%s', received: '%s'", expected, result)
	}

	id, _ := codec.Decode(result)
	res := provider.Get(context.TODO(), int(id))

	if uint64(res.Data.ID) != id {
		t.Fatalf("Expected: %d, received: %d", id, res.Data.ID)
	}

	if res.Data.Long != long_url {
		t.Fatalf("Expected: '%s', received: '%s'", long_url, res.Data.Long)
	}

	if res.Data.Short != expected {
		t.Fatalf("Expected: '%s', received: '%s'", expected, res.Data.Short)
	}

}
