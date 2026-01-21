package expiration

import "github.com/exanubes/url-shortener/internal/domain"

type ExpirationFactory struct{}

func NewFactory() ExpirationFactory {
	return ExpirationFactory{}
}

func (factory ExpirationFactory) Create(settings domain.PolicySettings) []domain.PolicySpec {
	var policies []domain.PolicySpec

	if settings.HasMaxAgeLimit() {
		policies = append(policies, domain.PolicySpec{
			Kind:   domain.PolicyKind_MaxAge,
			Params: map[string]any{"duration": settings.MaxAge},
		})
	}

	if settings.IsSingleUse() {
		policies = append(policies, domain.PolicySpec{
			Kind: domain.PolicyKind_SingleUse,
		})
	}

	return policies
}
