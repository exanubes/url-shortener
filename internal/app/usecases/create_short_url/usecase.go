package createshorturl

import (
	"context"
	"errors"

	"github.com/exanubes/url-shortener/internal/app/services/shortcode"
	"github.com/exanubes/url-shortener/internal/domain"
)

type CreateShortUrl struct {
	writer             UrlWriter
	short_code_service shortcode.Service
	policy_factory     RetryPolicyFactory
}

func New(writer UrlWriter, short_code_service shortcode.Service, policy_factory RetryPolicyFactory) *CreateShortUrl {
	return &CreateShortUrl{
		writer:             writer,
		short_code_service: short_code_service,
		policy_factory:     policy_factory,
	}
}

func (usecase *CreateShortUrl) Execute(ctx context.Context, url domain.Url) (domain.ShortCode, error) {
	retry_policy := usecase.policy_factory.Create()
	for retry_policy.Next() {
		short_code, err := usecase.short_code_service.Generate()
		if err != nil {
			return domain.ShortCode{}, err
		}

		if err := usecase.writer.Write(ctx, short_code, url); err != nil {
			if retry_policy.Verify(err) {
				continue
			}

			return domain.ShortCode{}, err
		} else {
			return short_code, nil
		}
	}

	return domain.ShortCode{}, errors.New("Failed to generate a unique short code")
}
