package usecase

import (
	"context"

	"github.com/exanubes/url-shortener/internal/domain"
)

type VisitShortUrl struct {
	persistence domain.PersistenceProvider
}

func NewVisitShortUrl(provider domain.PersistenceProvider) *VisitShortUrl {
	return &VisitShortUrl{
		persistence: provider,
	}
}

func (usecase *VisitShortUrl) Execute(ctx context.Context, short_url string) (string, error) {
	result, err := usecase.persistence.Get(ctx, short_url)

	if err != nil {
		return "", err
	}

	return result.Long, nil
}
