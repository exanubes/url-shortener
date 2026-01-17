package expiration

import "github.com/exanubes/url-shortener/internal/domain"

type ExpirationFactory struct{}

func NewFactory() ExpirationFactory {
	return ExpirationFactory{}
}

func (factory ExpirationFactory) Create(settings domain.PolicySettings) (domain.ExpirationPolicy, error) {
	var policies []domain.ExpirationPolicy

	if settings.HasMaxAgeLimit() {
		policy, err := domain.NewMaxLinkAgeExpirationPolicy(settings.MaxAge)

		if err != nil {
			return policy, err
		}

		policies = append(policies, policy)
	}

	if settings.IsSingleUse() {
		policies = append(policies, domain.NewOneTimeLinkExpirationPolicy())
	}
	if settings.HasMaxVisitsLimit() {
		policy, err := domain.NewMaxVisitsExpirationPolicy(settings.MaxVisits)

		if err != nil {
			return policy, err
		}

		policies = append(policies, policy)
	}

	return domain.NewChainExpirationPolicy(policies)
}
