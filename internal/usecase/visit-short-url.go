package usecase

import (
	"context"

	"github.com/exanubes/url-shortener/internal/domain"
)

type VisitShortUrl struct {
	persistence domain.PersistenceProvider
	codec       domain.Codec
}

func NewVisitShortUrl(provider domain.PersistenceProvider, codec domain.Codec) *VisitShortUrl {
	return &VisitShortUrl{
		persistence: provider,
		codec:       codec,
	}
}

func (usecase *VisitShortUrl) Execute(short_url string) (string, error) {
	ctx := context.Background()

	id, err := usecase.codec.Decode(short_url)

	if err != nil {
		return "", err
	}

	result := usecase.persistence.Get(ctx, id)

	if result.Err != nil {
		return "", result.Err
	}

	return result.Data.Long, nil
}
