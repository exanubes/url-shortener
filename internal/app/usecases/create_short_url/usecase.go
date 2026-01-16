package createshorturl

import (
	"context"
	"errors"
	"time"

	"github.com/exanubes/url-shortener/internal/app/services/expiration"
	"github.com/exanubes/url-shortener/internal/app/services/shortcode"
	"github.com/exanubes/url-shortener/internal/domain"
)

type CreateShortUrl struct {
	writer             LinkWriter
	short_code_service shortcode.Service
	policy_factory     RetryPolicyFactory
	expiration         expiration.Factory
}

func New(writer LinkWriter, short_code_service shortcode.Service, policy_factory RetryPolicyFactory, expiration expiration.Factory) *CreateShortUrl {
	return &CreateShortUrl{
		writer:             writer,
		short_code_service: short_code_service,
		policy_factory:     policy_factory,
		expiration:         expiration,
	}
}

func (usecase *CreateShortUrl) Execute(ctx context.Context, cmd domain.CreateLinkCommand) (*domain.Link, error) {
	retry_policy := usecase.policy_factory.Create()
	expiration_policy, err := usecase.expiration.Create(cmd.PolicySettings)

	if err != nil {
		return nil, err
	}

	for retry_policy.Next() {
		short_code, err := usecase.short_code_service.Generate()
		if err != nil {
			return nil, err
		}
		link := domain.CreateLink(cmd.Url, short_code, expiration_policy, time.Now())

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
