package usecase

import (
	"context"
	"errors"

	"github.com/exanubes/url-shortener/internal/domain"
)

type CreateShortUrl struct {
	persistence          domain.PersistenceProvider
	short_code_generator domain.ShortCodeGenerator
}

func NewCreateShortUrl(provider domain.PersistenceProvider, short_code_generator domain.ShortCodeGenerator) *CreateShortUrl {
	return &CreateShortUrl{
		persistence:          provider,
		short_code_generator: short_code_generator,
	}
}

func (usecase *CreateShortUrl) Execute(ctx context.Context, url string) (string, error) {
	// TODO: Replace with retry policy
	retries := 0
	for retries < 3 {
		short_code, err := usecase.short_code_generator.Generate()
		if err != nil {
			return "", err
		}

		if err := usecase.persistence.Save(ctx, domain.Url{
			Long:  url,
			Short: short_code.String(),
		}); err != nil {
			if !errors.Is(err, domain.ErrShortCodeCollision) {
				return "", err
			}
			retries += 1
		} else {
			return short_code.String(), nil
		}
	}

	return "", errors.New("Failed to generate short code")
}
