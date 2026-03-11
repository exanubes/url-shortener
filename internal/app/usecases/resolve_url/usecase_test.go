package resolveurl_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/exanubes/url-shortener/internal/app/services/expiration"
	resolveurl "github.com/exanubes/url-shortener/internal/app/usecases/resolve_url"
	"github.com/exanubes/url-shortener/internal/domain"
	"github.com/exanubes/url-shortener/internal/infrastructure/persistence/inmemory"
)

type mock_clock struct {
	now time.Time
}

func (m mock_clock) Now() time.Time {
	return m.now
}

func TestResolveUrl(t *testing.T) {
	fixed_time := time.Now()
	clock := mock_clock{now: fixed_time}

	repository := inmemory.NewInmemoryRepository(clock)

	long_url, _ := domain.NewUrl("https://exanubes.com")
	short_code, _ := domain.NewShortCode("2TX", 7, "0")

	policy_settings, _ := domain.NewPolicySettings(time.Hour, false)
	expiration_factory := expiration.NewFactory()
	policy_specs := expiration_factory.Create(policy_settings)

	link := domain.CreateLink(long_url, short_code, policy_specs, clock.Now())

	repository.Write(context.TODO(), link)

	usecase := resolveurl.New(repository, repository, clock)

	result, err := usecase.Execute(context.TODO(), short_code)

	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}

	if result.Url.String() != long_url.String() {
		t.Fatalf("Expected: %s, received: %s", long_url.String(), result.Url.String())
	}

	if result.Status.Consumed {
		t.Fatal("Expected multi-use link to not be consumed")
	}

	expected_expiry := fixed_time.Add(time.Hour)
	if !result.Status.ExpiresAt.Equal(expected_expiry) {
		t.Fatalf("Expected expiry at %s, got %s", expected_expiry, result.Status.ExpiresAt)
	}

}

func TestResolveUrl_SingleUse(t *testing.T) {
	fixed_time := time.Now()
	clock := mock_clock{now: fixed_time}
	repository := inmemory.NewInmemoryRepository(clock)

	long_url, _ := domain.NewUrl("https://exanubes.com")
	short_code, _ := domain.NewShortCode("ABC", 7, "0")

	policy_settings, _ := domain.NewPolicySettings(time.Hour, true)
	expiration_factory := expiration.NewFactory()
	policy_specs := expiration_factory.Create(policy_settings)

	link := domain.CreateLink(long_url, short_code, policy_specs, clock.Now())
	repository.Write(context.TODO(), link)

	usecase := resolveurl.New(repository, repository, clock)

	result, err := usecase.Execute(context.TODO(), short_code)

	if err != nil {
		t.Fatalf("Unexpected error: %s", err.Error())
	}

	if result.Url.String() != long_url.String() {
		t.Fatalf("Expected: %s, received: %s", long_url.String(), result.Url.String())
	}

	if !result.Status.Consumed {
		t.Fatal("Expected single-use link to be marked as consumed")
	}
}

func TestResolveUrl_ExpiredLink(t *testing.T) {
	initial_time := time.Now()
	clock := mock_clock{now: initial_time}
	repository := inmemory.NewInmemoryRepository(clock)
	consumer := &mock_consumer{}

	long_url, _ := domain.NewUrl("https://exanubes.com")
	short_code, _ := domain.NewShortCode("XYZ", 7, "0")

	policy_settings, _ := domain.NewPolicySettings(time.Hour, false)
	expiration_factory := expiration.NewFactory()
	policy_specs := expiration_factory.Create(policy_settings)

	link := domain.CreateLink(long_url, short_code, policy_specs, clock.Now())
	repository.Write(context.TODO(), link)

	expired_clock := mock_clock{now: initial_time.Add(2 * time.Hour)}
	usecase := resolveurl.New(repository, consumer, expired_clock)

	_, err := usecase.Execute(context.TODO(), short_code)

	if err == nil {
		t.Fatal("Expected error for expired link, got nil")
	}

	if !errors.Is(err, domain.ErrLinkExpired) {
		t.Fatalf("Expected ErrLinkExpired, got: %s", err.Error())
	}

	if consumer.consume_called {
		t.Fatal("Expected consumer not to be called for expired link")
	}
}

func TestResolveUrl_ConsumedSingleUseLink(t *testing.T) {
	fixed_time := time.Now()
	clock := mock_clock{now: fixed_time}
	repository := inmemory.NewInmemoryRepository(clock)

	long_url, _ := domain.NewUrl("https://exanubes.com")
	short_code, _ := domain.NewShortCode("DEF", 7, "0")

	policy_settings, _ := domain.NewPolicySettings(time.Hour, true)
	expiration_factory := expiration.NewFactory()
	policy_specs := expiration_factory.Create(policy_settings)

	link := domain.CreateLink(long_url, short_code, policy_specs, clock.Now())
	repository.Write(context.TODO(), link)

	usecase := resolveurl.New(repository, repository, clock)

	result, err := usecase.Execute(context.TODO(), short_code)
	if err != nil {
		t.Fatalf("First access should succeed, got error: %s", err.Error())
	}

	if !result.Status.Consumed {
		t.Fatal("Expected link to be consumed after first access")
	}

	_, err = usecase.Execute(context.TODO(), short_code)
	if err == nil {
		t.Fatal("Expected error for already consumed link, got nil")
	}

	if !errors.Is(err, domain.ErrLinkExpired) {
		t.Fatalf("Expected ErrLinkExpired, got: %s", err.Error())
	}
}

func TestResolveUrl_LinkNotFound(t *testing.T) {
	clock := mock_clock{now: time.Now()}
	repository := inmemory.NewInmemoryRepository(clock)

	usecase := resolveurl.New(repository, repository, clock)

	short_code, _ := domain.NewShortCode("XXX", 7, "0")
	_, err := usecase.Execute(context.TODO(), short_code)

	if err == nil {
		t.Fatal("Expected error for non-existent link, got nil")
	}

	if !errors.Is(err, domain.ErrUrlNotFound) {
		t.Fatalf("Expected ErrUrlNotFound, got: %s", err.Error())
	}
}

type mock_consumer struct {
	consume_called bool
	consume_error  error
}

func (m *mock_consumer) Consume(ctx context.Context, shortcode domain.ShortCode) error {
	m.consume_called = true
	return m.consume_error
}

func TestResolveUrl_ConsumerError(t *testing.T) {
	clock := mock_clock{now: time.Now()}
	repository := inmemory.NewInmemoryRepository(clock)

	expected_error := errors.New("consumer failed")
	consumer := &mock_consumer{consume_error: expected_error}

	long_url, _ := domain.NewUrl("https://exanubes.com")
	short_code, _ := domain.NewShortCode("GHI", 7, "0")

	policy_settings, _ := domain.NewPolicySettings(time.Hour, true)
	expiration_factory := expiration.NewFactory()
	policy_specs := expiration_factory.Create(policy_settings)

	link := domain.CreateLink(long_url, short_code, policy_specs, clock.Now())
	repository.Write(context.TODO(), link)

	usecase := resolveurl.New(repository, consumer, clock)

	_, err := usecase.Execute(context.TODO(), short_code)

	if err == nil {
		t.Fatal("Expected error when consumer fails, got nil")
	}

	if !errors.Is(err, expected_error) {
		t.Fatalf("Expected consumer error, got: %s", err.Error())
	}

	if !consumer.consume_called {
		t.Fatal("Expected consumer to be called")
	}
}
