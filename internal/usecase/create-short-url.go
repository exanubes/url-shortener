package usecase

import (
	"context"
	"errors"

	"github.com/exanubes/url-shortener/internal/domain"
)

type CreateShortUrl struct {
	persistence          domain.PersistenceProvider
	short_code_generator domain.ShortCodeGenerator
	policy_factory       domain.RetryPolicyFactory
}

func NewCreateShortUrl(provider domain.PersistenceProvider, short_code_generator domain.ShortCodeGenerator, policy_factory domain.RetryPolicyFactory) *CreateShortUrl {
	return &CreateShortUrl{
		persistence:          provider,
		short_code_generator: short_code_generator,
		policy_factory:       policy_factory,
	}
}

func (usecase *CreateShortUrl) Execute(ctx context.Context, url domain.Url) (domain.ShortCode, error) {
	retry_policy := usecase.policy_factory.Create()
	for retry_policy.Next() {
		short_code, err := usecase.short_code_generator.Generate()
		if err != nil {
			return domain.ShortCode{}, err
		}

		if err := usecase.persistence.Save(ctx, url, short_code); err != nil {
			if !errors.Is(err, domain.ErrShortCodeCollision) {
				return domain.ShortCode{}, err
			}
		} else {
			return short_code, nil
		}
	}

	return domain.ShortCode{}, errors.New("Failed to generate short code")
}
