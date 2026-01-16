package createshorturl

import (
	"context"
	"errors"
	"time"

	"github.com/exanubes/url-shortener/internal/app/services/shortcode"
	"github.com/exanubes/url-shortener/internal/domain"
)

type CreateShortUrl struct {
	writer             LinkWriter
	short_code_service shortcode.Service
	policy_factory     RetryPolicyFactory
}

func New(writer LinkWriter, short_code_service shortcode.Service, policy_factory RetryPolicyFactory) *CreateShortUrl {
	return &CreateShortUrl{
		writer:             writer,
		short_code_service: short_code_service,
		policy_factory:     policy_factory,
	}
}

func (usecase *CreateShortUrl) Execute(ctx context.Context, url domain.Url) (*domain.Link, error) {
	retry_policy := usecase.policy_factory.Create()
	// day := 24 * time.Hour
	expiration_policy := domain.NewMaxLinkAgeExpirationPolicy(30 * time.Second)
	for retry_policy.Next() {
		short_code, err := usecase.short_code_service.Generate()
		if err != nil {
			return nil, err
		}
		link := domain.CreateLink(url, short_code, expiration_policy, time.Now())

		if err := usecase.writer.Write(ctx, link); err != nil {
			if retry_policy.Verify(err) {
				continue
			}

			return nil, err
		} else {
			return link, nil
		}
	}

	return nil, errors.New("Failed to generate a unique short code")
}
