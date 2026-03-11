package domain_test

import (
	"testing"
	"time"

	"github.com/exanubes/url-shortener/internal/domain"
)

const (
	testUrl       = "https://example.com"
	testShortCode = "abc123x"
)

var baseTime = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func mustCreateUrl(t *testing.T, url string) domain.Url {
	t.Helper()
	u, err := domain.NewUrl(url)
	if err != nil {
		t.Fatalf("Failed to create URL: %v", err)
	}
	return u
}

func mustCreateShortCode(t *testing.T, value string) domain.ShortCode {
	t.Helper()
	sc, err := domain.NewShortCodeFromParam(value)
	if err != nil {
		t.Fatalf("Failed to create ShortCode: %v", err)
	}
	return sc
}

func maxAgeSpec(ttl time.Duration) domain.PolicySpec {
	return domain.PolicySpec{
		Kind:   domain.PolicyKind_MaxAge,
		Params: domain.MaxAgeParams{TTL: ttl},
	}
}

func singleUseSpec() domain.PolicySpec {
	return domain.PolicySpec{
		Kind:   domain.PolicyKind_SingleUse,
		Params: domain.SingleUseParams{},
	}
}

func TestCreateLink_Success(t *testing.T) {
	url := mustCreateUrl(t, testUrl)
	shortcode := mustCreateShortCode(t, testShortCode)
	policies := []domain.PolicySpec{maxAgeSpec(time.Hour)}

	link := domain.CreateLink(url, shortcode, policies, baseTime)

	if link == nil {
		t.Fatalf("Expected link to be created, got nil")
	}

	snapshot := link.Snapshot()
	if snapshot.Url.String() != testUrl {
		t.Fatalf("Expected URL '%s', got '%s'", testUrl, snapshot.Url.String())
	}
	if snapshot.Shortcode.String() != testShortCode {
		t.Fatalf("Expected ShortCode '%s', got '%s'", testShortCode, snapshot.Shortcode.String())
	}
	if !snapshot.ConsumedAt.IsZero() {
		t.Fatalf("Expected ConsumedAt to be zero, got %v", snapshot.ConsumedAt)
	}
	if !snapshot.CreatedAt.Equal(baseTime) {
		t.Fatalf("Expected CreatedAt to be %v, got %v", baseTime, snapshot.CreatedAt)
	}
}

func TestCreateLink_WithNoPolicy(t *testing.T) {
	url := mustCreateUrl(t, testUrl)
	shortcode := mustCreateShortCode(t, testShortCode)
	policies := []domain.PolicySpec{}

	link := domain.CreateLink(url, shortcode, policies, baseTime)

	_, err := link.Visit(baseTime)
	if err != domain.ErrUndefinedExpirationPolicy {
		t.Fatalf("Expected ErrUndefinedExpirationPolicy, got %v", err)
	}
}

func TestRehydrateLink_Success(t *testing.T) {
	url := mustCreateUrl(t, testUrl)
	shortcode := mustCreateShortCode(t, testShortCode)
	policies := []domain.PolicySpec{maxAgeSpec(time.Hour)}
	consumedAt := baseTime.Add(30 * time.Minute)

	state := domain.LinkState{
		Url:         url,
		Shortcode:   shortcode,
		PolicySpecs: policies,
		CreatedAt:   baseTime,
		ConsumedAt:  consumedAt,
	}

	link := domain.RehydrateLink(state)

	if link == nil {
		t.Fatalf("Expected link to be rehydrated, got nil")
	}

	snapshot := link.Snapshot()
	if !snapshot.ConsumedAt.Equal(consumedAt) {
		t.Fatalf("Expected ConsumedAt to be %v, got %v", consumedAt, snapshot.ConsumedAt)
	}
	if !snapshot.CreatedAt.Equal(baseTime) {
		t.Fatalf("Expected CreatedAt to be %v, got %v", baseTime, snapshot.CreatedAt)
	}
}

func TestRehydrateLink_WithConsumedAt(t *testing.T) {
	url := mustCreateUrl(t, testUrl)
	shortcode := mustCreateShortCode(t, testShortCode)
	policies := []domain.PolicySpec{singleUseSpec()}
	consumedAt := baseTime.Add(10 * time.Minute)

	state := domain.LinkState{
		Url:         url,
		Shortcode:   shortcode,
		PolicySpecs: policies,
		CreatedAt:   baseTime,
		ConsumedAt:  consumedAt,
	}

	link := domain.RehydrateLink(state)

	snapshot := link.Snapshot()
	if snapshot.ConsumedAt.IsZero() {
		t.Fatalf("Expected ConsumedAt to be preserved, got zero time")
	}
	if !snapshot.ConsumedAt.Equal(consumedAt) {
		t.Fatalf("Expected ConsumedAt to be %v, got %v", consumedAt, snapshot.ConsumedAt)
	}
	if !link.SingleUse() {
		t.Fatalf("Expected SingleUse to be true")
	}
}

// NOTE: Link.Visit() tests
func TestVisit_Success_WithMaxAgePolicy(t *testing.T) {
	url := mustCreateUrl(t, testUrl)
	shortcode := mustCreateShortCode(t, testShortCode)
	policies := []domain.PolicySpec{maxAgeSpec(time.Hour)}

	link := domain.CreateLink(url, shortcode, policies, baseTime)

	visitTime := baseTime.Add(30 * time.Minute)
	resultUrl, err := link.Visit(visitTime)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resultUrl.String() != testUrl {
		t.Fatalf("Expected URL '%s', got '%s'", testUrl, resultUrl.String())
	}
}

func TestVisit_Expired_MaxAgePolicy(t *testing.T) {
	url := mustCreateUrl(t, testUrl)
	shortcode := mustCreateShortCode(t, testShortCode)
	policies := []domain.PolicySpec{maxAgeSpec(time.Hour)}

	link := domain.CreateLink(url, shortcode, policies, baseTime)

	visitTime := baseTime.Add(2 * time.Hour)
	_, err := link.Visit(visitTime)

	if err != domain.ErrLinkExpired {
		t.Fatalf("Expected ErrLinkExpired, got %v", err)
	}
}

func TestVisit_Success_WithSingleUsePolicy_NotConsumed(t *testing.T) {
	url := mustCreateUrl(t, testUrl)
	shortcode := mustCreateShortCode(t, testShortCode)
	policies := []domain.PolicySpec{singleUseSpec()}

	link := domain.CreateLink(url, shortcode, policies, baseTime)

	visitTime := baseTime.Add(10 * time.Minute)
	resultUrl, err := link.Visit(visitTime)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resultUrl.String() != testUrl {
		t.Fatalf("Expected URL '%s', got '%s'", testUrl, resultUrl.String())
	}
}

func TestVisit_Expired_SingleUsePolicy_Consumed(t *testing.T) {
	url := mustCreateUrl(t, testUrl)
	shortcode := mustCreateShortCode(t, testShortCode)
	policies := []domain.PolicySpec{singleUseSpec()}

	link := domain.CreateLink(url, shortcode, policies, baseTime)
	link.Consume(baseTime.Add(5 * time.Minute))

	visitTime := baseTime.Add(10 * time.Minute)
	_, err := link.Visit(visitTime)

	if err != domain.ErrLinkExpired {
		t.Fatalf("Expected ErrLinkExpired, got %v", err)
	}
}

func TestVisit_CombinedPolicies_MaxAgeExpired(t *testing.T) {
	url := mustCreateUrl(t, testUrl)
	shortcode := mustCreateShortCode(t, testShortCode)
	policies := []domain.PolicySpec{
		maxAgeSpec(time.Hour),
		singleUseSpec(),
	}

	link := domain.CreateLink(url, shortcode, policies, baseTime)

	visitTime := baseTime.Add(2 * time.Hour)
	_, err := link.Visit(visitTime)

	if err != domain.ErrLinkExpired {
		t.Fatalf("Expected ErrLinkExpired, got %v", err)
	}
}

func TestVisit_CombinedPolicies_SingleUseExpired(t *testing.T) {
	url := mustCreateUrl(t, testUrl)
	shortcode := mustCreateShortCode(t, testShortCode)
	policies := []domain.PolicySpec{
		maxAgeSpec(time.Hour),
		singleUseSpec(),
	}

	link := domain.CreateLink(url, shortcode, policies, baseTime)
	link.Consume(baseTime.Add(10 * time.Minute))

	visitTime := baseTime.Add(30 * time.Minute)
	_, err := link.Visit(visitTime)

	if err != domain.ErrLinkExpired {
		t.Fatalf("Expected ErrLinkExpired, got %v", err)
	}
}

func TestVisit_CombinedPolicies_BothValid(t *testing.T) {
	url := mustCreateUrl(t, testUrl)
	shortcode := mustCreateShortCode(t, testShortCode)
	policies := []domain.PolicySpec{
		maxAgeSpec(time.Hour),
		singleUseSpec(),
	}

	link := domain.CreateLink(url, shortcode, policies, baseTime)

	visitTime := baseTime.Add(30 * time.Minute)
	resultUrl, err := link.Visit(visitTime)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if resultUrl.String() != testUrl {
		t.Fatalf("Expected URL '%s', got '%s'", testUrl, resultUrl.String())
	}
}

func TestVisit_InvalidPolicySpec(t *testing.T) {
	url := mustCreateUrl(t, testUrl)
	shortcode := mustCreateShortCode(t, testShortCode)
	policies := []domain.PolicySpec{
		{
			Kind:   domain.PolicyKind_SingleUse,
			Params: domain.MaxAgeParams{TTL: time.Hour},
		},
	}

	link := domain.CreateLink(url, shortcode, policies, baseTime)

	_, err := link.Visit(baseTime)

	if err != domain.ErrInvalidPolicySpecParams {
		t.Fatalf("Expected ErrInvalidPolicySpecParams, got %v", err)
	}
}

// NOTE: Link.Consume()

func TestConsume_FirstCall_SetsConsumedAt(t *testing.T) {
	url := mustCreateUrl(t, testUrl)
	shortcode := mustCreateShortCode(t, testShortCode)
	policies := []domain.PolicySpec{maxAgeSpec(time.Hour)}

	link := domain.CreateLink(url, shortcode, policies, baseTime)

	consumeTime := baseTime.Add(15 * time.Minute)
	link.Consume(consumeTime)

	snapshot := link.Snapshot()
	if snapshot.ConsumedAt.IsZero() {
		t.Fatalf("Expected ConsumedAt to be set, got zero time")
	}
	if !snapshot.ConsumedAt.Equal(consumeTime) {
		t.Fatalf("Expected ConsumedAt to be %v, got %v", consumeTime, snapshot.ConsumedAt)
	}
}

func TestConsume_SecondCall_PreservesOriginalTimestamp(t *testing.T) {
	url := mustCreateUrl(t, testUrl)
	shortcode := mustCreateShortCode(t, testShortCode)
	policies := []domain.PolicySpec{maxAgeSpec(time.Hour)}

	link := domain.CreateLink(url, shortcode, policies, baseTime)

	firstConsumeTime := baseTime.Add(10 * time.Minute)
	link.Consume(firstConsumeTime)

	secondConsumeTime := baseTime.Add(20 * time.Minute)
	link.Consume(secondConsumeTime)

	snapshot := link.Snapshot()
	if !snapshot.ConsumedAt.Equal(firstConsumeTime) {
		t.Fatalf("Expected ConsumedAt to remain %v, got %v", firstConsumeTime, snapshot.ConsumedAt)
	}
}

func TestConsume_AlwaysReturnsErrLinkConsumed(t *testing.T) {
	url := mustCreateUrl(t, testUrl)
	shortcode := mustCreateShortCode(t, testShortCode)
	policies := []domain.PolicySpec{maxAgeSpec(time.Hour)}

	link := domain.CreateLink(url, shortcode, policies, baseTime)

	firstConsumeTime := baseTime.Add(10 * time.Minute)
	err := link.Consume(firstConsumeTime)
	if err != domain.ErrLinkConsumed {
		t.Fatalf("Expected first Consume() to return ErrLinkConsumed, got %v", err)
	}

	secondConsumeTime := baseTime.Add(20 * time.Minute)
	err = link.Consume(secondConsumeTime)
	if err != domain.ErrLinkConsumed {
		t.Fatalf("Expected second Consume() to return ErrLinkConsumed, got %v", err)
	}
}

func TestConsume_WithSingleUsePolicy_MarksAsExpired(t *testing.T) {
	url := mustCreateUrl(t, testUrl)
	shortcode := mustCreateShortCode(t, testShortCode)
	policies := []domain.PolicySpec{singleUseSpec()}

	link := domain.CreateLink(url, shortcode, policies, baseTime)

	consumeTime := baseTime.Add(5 * time.Minute)
	link.Consume(consumeTime)

	visitTime := baseTime.Add(10 * time.Minute)
	_, err := link.Visit(visitTime)

	if err != domain.ErrLinkExpired {
		t.Fatalf("Expected Visit() after Consume() to return ErrLinkExpired, got %v", err)
	}
}

// NOTE: Link.ExpirationStatus()

func TestExpirationStatus_NotConsumed_WithMaxAge(t *testing.T) {
	url := mustCreateUrl(t, testUrl)
	shortcode := mustCreateShortCode(t, testShortCode)
	policies := []domain.PolicySpec{maxAgeSpec(time.Hour)}

	link := domain.CreateLink(url, shortcode, policies, baseTime)

	checkTime := baseTime.Add(30 * time.Minute)
	status, err := link.ExpirationStatus(checkTime)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if status.Consumed {
		t.Fatalf("Expected Consumed to be false")
	}

	expectedExpiry := baseTime.Add(time.Hour)
	if !status.ExpiresAt.Equal(expectedExpiry) {
		t.Fatalf("Expected ExpiresAt to be %v, got %v", expectedExpiry, status.ExpiresAt)
	}
}

func TestExpirationStatus_Consumed_WithSingleUse(t *testing.T) {
	url := mustCreateUrl(t, testUrl)
	shortcode := mustCreateShortCode(t, testShortCode)
	policies := []domain.PolicySpec{singleUseSpec()}

	link := domain.CreateLink(url, shortcode, policies, baseTime)
	link.Consume(baseTime.Add(5 * time.Minute))

	checkTime := baseTime.Add(10 * time.Minute)
	status, err := link.ExpirationStatus(checkTime)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !status.Consumed {
		t.Fatalf("Expected Consumed to be true")
	}
	if !status.ExpiresAt.IsZero() {
		t.Fatalf("Expected ExpiresAt to be zero for single-use policy, got %v", status.ExpiresAt)
	}
}

func TestExpirationStatus_CombinedPolicies_ReturnsEarliestExpiry(t *testing.T) {
	url := mustCreateUrl(t, testUrl)
	shortcode := mustCreateShortCode(t, testShortCode)
	policies := []domain.PolicySpec{
		maxAgeSpec(time.Hour),
		singleUseSpec(),
	}

	link := domain.CreateLink(url, shortcode, policies, baseTime)

	checkTime := baseTime.Add(30 * time.Minute)
	status, err := link.ExpirationStatus(checkTime)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if status.Consumed {
		t.Fatalf("Expected Consumed to be false")
	}

	expectedExpiry := baseTime.Add(time.Hour)
	if !status.ExpiresAt.Equal(expectedExpiry) {
		t.Fatalf("Expected ExpiresAt to be %v, got %v", expectedExpiry, status.ExpiresAt)
	}
}

func TestExpirationStatus_InvalidPolicy(t *testing.T) {
	url := mustCreateUrl(t, testUrl)
	shortcode := mustCreateShortCode(t, testShortCode)
	policies := []domain.PolicySpec{
		{
			Kind:   domain.PolicyKind_MaxAge,
			Params: domain.SingleUseParams{},
		},
	}

	link := domain.CreateLink(url, shortcode, policies, baseTime)

	_, err := link.ExpirationStatus(baseTime)

	if err != domain.ErrInvalidPolicySpecParams {
		t.Fatalf("Expected ErrInvalidPolicySpecParams, got %v", err)
	}
}

// NOTE: ShortCode & Snapshot Tests

func TestShortCode_ReturnsCorrectValue(t *testing.T) {
	url := mustCreateUrl(t, testUrl)
	shortcode := mustCreateShortCode(t, testShortCode)
	policies := []domain.PolicySpec{maxAgeSpec(time.Hour)}

	link := domain.CreateLink(url, shortcode, policies, baseTime)

	if link.ShortCode().String() != testShortCode {
		t.Fatalf("Expected ShortCode '%s', got '%s'", testShortCode, link.ShortCode().String())
	}
}

func TestSnapshot_ReturnsCompleteState(t *testing.T) {
	url := mustCreateUrl(t, testUrl)
	shortcode := mustCreateShortCode(t, testShortCode)
	policies := []domain.PolicySpec{maxAgeSpec(time.Hour)}
	consumedAt := baseTime.Add(30 * time.Minute)

	link := domain.CreateLink(url, shortcode, policies, baseTime)
	link.Consume(consumedAt)

	snapshot := link.Snapshot()

	if snapshot.Url.String() != testUrl {
		t.Fatalf("Expected URL '%s', got '%s'", testUrl, snapshot.Url.String())
	}
	if snapshot.Shortcode.String() != testShortCode {
		t.Fatalf("Expected ShortCode '%s', got '%s'", testShortCode, snapshot.Shortcode.String())
	}
	if !snapshot.CreatedAt.Equal(baseTime) {
		t.Fatalf("Expected CreatedAt to be %v, got %v", baseTime, snapshot.CreatedAt)
	}
	if !snapshot.ConsumedAt.Equal(consumedAt) {
		t.Fatalf("Expected ConsumedAt to be %v, got %v", consumedAt, snapshot.ConsumedAt)
	}
	if len(snapshot.PolicySpecs) != len(policies) {
		t.Fatalf("Expected %d policy specs, got %d", len(policies), len(snapshot.PolicySpecs))
	}
}

func TestSnapshot_PreservesZeroConsumedAt(t *testing.T) {
	url := mustCreateUrl(t, testUrl)
	shortcode := mustCreateShortCode(t, testShortCode)
	policies := []domain.PolicySpec{maxAgeSpec(time.Hour)}

	link := domain.CreateLink(url, shortcode, policies, baseTime)

	snapshot := link.Snapshot()

	if !snapshot.ConsumedAt.IsZero() {
		t.Fatalf("Expected ConsumedAt to be zero, got %v", snapshot.ConsumedAt)
	}
}

// NOTE: Edge Cases

func TestVisit_AtExactExpirationTime(t *testing.T) {
	url := mustCreateUrl(t, testUrl)
	shortcode := mustCreateShortCode(t, testShortCode)
	policies := []domain.PolicySpec{maxAgeSpec(time.Hour)}

	link := domain.CreateLink(url, shortcode, policies, baseTime)

	visitTimeAtExpiry := baseTime.Add(time.Hour)
	resultUrl, err := link.Visit(visitTimeAtExpiry)

	if err != nil {
		t.Fatalf("Expected no error at exact expiration time, got %v", err)
	}
	if resultUrl.String() != testUrl {
		t.Fatalf("Expected URL '%s', got '%s'", testUrl, resultUrl.String())
	}

	visitTimeAfterExpiry := baseTime.Add(time.Hour).Add(time.Nanosecond)
	_, err = link.Visit(visitTimeAfterExpiry)

	if err != domain.ErrLinkExpired {
		t.Fatalf("Expected ErrLinkExpired after expiration time, got %v", err)
	}
}

func TestVisit_MultiplePolicies_AllExpired(t *testing.T) {
	url := mustCreateUrl(t, testUrl)
	shortcode := mustCreateShortCode(t, testShortCode)
	policies := []domain.PolicySpec{
		maxAgeSpec(time.Hour),
		singleUseSpec(),
	}

	link := domain.CreateLink(url, shortcode, policies, baseTime)
	link.Consume(baseTime.Add(10 * time.Minute))

	visitTime := baseTime.Add(2 * time.Hour)
	_, err := link.Visit(visitTime)

	if err != domain.ErrLinkExpired {
		t.Fatalf("Expected ErrLinkExpired when all policies expired, got %v", err)
	}
}

func TestExpirationStatus_NoConsumedAt_ShowsFalse(t *testing.T) {
	url := mustCreateUrl(t, testUrl)
	shortcode := mustCreateShortCode(t, testShortCode)
	policies := []domain.PolicySpec{maxAgeSpec(time.Hour)}

	link := domain.CreateLink(url, shortcode, policies, baseTime)

	status, err := link.ExpirationStatus(baseTime)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if status.Consumed {
		t.Fatalf("Expected Consumed to be false when ConsumedAt is zero")
	}
}
