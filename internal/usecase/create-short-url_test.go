package usecase

import (
	"context"
	"testing"

	"github.com/exanubes/url-shortener/internal/app/policy"
	"github.com/exanubes/url-shortener/internal/domain"
	"github.com/exanubes/url-shortener/internal/infrastructure/persistence/inmemory"
)

type mock_shortcode_factory struct{}

func (mock_shortcode_factory) Generate() (domain.ShortCode, error) {
	return domain.NewShortCode("2TY", 7, "0")
}

func TestCreateShortUrl(t *testing.T) {
	provider := inmemory.NewInmemoryRepository()

	usecase := NewCreateShortUrl(provider, mock_shortcode_factory{}, policy.NewRetryPolicyFactory(3))

	expected, _ := domain.NewShortCode("2TY", 7, "0")
	long_url, _ := domain.NewUrl("https://exanubes.com")
	result, err := usecase.Execute(context.TODO(), long_url)

	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}

	if result.String() != expected.String() {
		t.Fatalf("Expected '%s', received: '%s'", expected, result)
	}

	res, _ := provider.Get(context.TODO(), expected)

	if res.String() != long_url.String() {
		t.Fatalf("Expected: '%s', received: '%s'", long_url, res.String())
	}

}

type mock_provider struct {
	called_counter int
}

func (p *mock_provider) Save(_ context.Context, _ domain.Url, _ domain.ShortCode) error {
	p.called_counter += 1
	return domain.ErrShortCodeCollision
}

func (*mock_provider) Get(_ context.Context, _ domain.ShortCode) (domain.Url, error) {
	return domain.Url{}, nil
}

func TestRetryFlow(t *testing.T) {
	provider := &mock_provider{}
	usecase := NewCreateShortUrl(provider, mock_shortcode_factory{}, policy.NewRetryPolicyFactory(3))
	expected := 3
	long_url, _ := domain.NewUrl("https://exanubes.com")
	_, err := usecase.Execute(context.TODO(), long_url)

	if provider.called_counter != expected {
		t.Fatalf("Expected %d retries, received %d", expected, provider.called_counter)
	}

	if err == nil {
		t.Fatal("Expected to receive error, received nil")
	}
}
