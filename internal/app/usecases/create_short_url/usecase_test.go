package createshorturl_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/exanubes/url-shortener/internal/app/services/expiration"
	createshorturl "github.com/exanubes/url-shortener/internal/app/usecases/create_short_url"
	"github.com/exanubes/url-shortener/internal/domain"
	"github.com/exanubes/url-shortener/internal/infrastructure/persistence/inmemory"
)

type mock_clock struct {
	now time.Time
}

func (m mock_clock) Now() time.Time {
	return m.now
}

type mock_shortcode_factory struct {
	shortcode string
	err       error
}

func (f mock_shortcode_factory) Generate() (domain.ShortCode, error) {
	if f.err == nil {
		return domain.NewShortCode("2TY", 7, "0")
	}

	return domain.ShortCode{}, f.err
}

func TestCreateShortUrl(t *testing.T) {
	clock := mock_clock{now: time.Now()}
	provider := inmemory.NewInmemoryRepository(clock)
	expiration_factory := expiration.NewFactory()
	usecase := createshorturl.New(
		provider,
		mock_shortcode_factory{
			shortcode: "2TY",
		},
		createshorturl.NewRetryPolicyFactory(3),
		expiration_factory,
		clock,
	)

	expected_shortcode, _ := domain.NewShortCode("2TY", 7, "0")
	long_url, _ := domain.NewUrl("https://exanubes.com")

	policy_settings, _ := domain.NewPolicySettings(time.Hour, false)
	cmd := domain.CreateLinkCommand{
		Url:            long_url,
		PolicySettings: policy_settings,
	}

	result, err := usecase.Execute(context.TODO(), cmd)

	if err != nil {
		t.Fatalf("Unexpected error %s", err.Error())
	}

	if result.ShortCode().String() != expected_shortcode.String() {
		t.Fatalf("Expected '%s', received: '%s'", expected_shortcode, result.ShortCode())
	}
}

type mock_provider struct {
	called_counter int
}

func (p *mock_provider) Write(_ context.Context, _ *domain.Link) error {
	p.called_counter += 1
	return domain.ErrShortCodeCollision
}

func (*mock_provider) Resolve(_ context.Context, _ domain.ShortCode) (*domain.Link, error) {
	return nil, nil
}

func TestRetryFlow(t *testing.T) {
	clock := mock_clock{now: time.Now()}
	provider := &mock_provider{}
	expiration_factory := expiration.NewFactory()

	usecase := createshorturl.New(
		provider,
		mock_shortcode_factory{},
		createshorturl.NewRetryPolicyFactory(3),
		expiration_factory,
		clock,
	)

	expected := 3
	long_url, _ := domain.NewUrl("https://exanubes.com")

	policy_settings, _ := domain.NewPolicySettings(time.Hour, false)
	cmd := domain.CreateLinkCommand{
		Url:            long_url,
		PolicySettings: policy_settings,
	}

	_, err := usecase.Execute(context.TODO(), cmd)

	if provider.called_counter != expected {
		t.Fatalf("Expected %d retries, received %d", expected, provider.called_counter)
	}

	if err == nil {
		t.Fatal("Expected to receive error, received nil")
	}

	if !errors.Is(err, domain.ErrShortCodeCollision) {
		t.Fatal("Expected shortcode collision error")
	}
}

func TestShortCodeGenerationError(t *testing.T) {
	clock := mock_clock{now: time.Now()}
	provider := inmemory.NewInmemoryRepository(clock)
	expiration_factory := expiration.NewFactory()

	usecase := createshorturl.New(
		provider,
		mock_shortcode_factory{err: fmt.Errorf("test error")},
		createshorturl.NewRetryPolicyFactory(1),
		expiration_factory,
		clock,
	)

	long_url, _ := domain.NewUrl("https://exanubes.com")

	policy_settings, _ := domain.NewPolicySettings(time.Hour, false)
	cmd := domain.CreateLinkCommand{
		Url:            long_url,
		PolicySettings: policy_settings,
	}

	_, err := usecase.Execute(context.TODO(), cmd)

	if err == nil {
		t.Fatalf("expected error, received nil")
	}

	if err.Error() != "test error" {
		t.Fatalf("Expected 'test error', received %s", err.Error())
	}

}
