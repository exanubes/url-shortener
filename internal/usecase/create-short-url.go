package usecase

import (
	"context"

	"github.com/exanubes/url-shortener/internal/domain"
)

type CreateShortUrl struct {
	persistence domain.PersistenceProvider
	codec       domain.Codec
}

func NewCreateShortUrl(provider domain.PersistenceProvider, codec domain.Codec) *CreateShortUrl {
	return &CreateShortUrl{
		persistence: provider,
		codec:       codec,
	}
}

func (usecase *CreateShortUrl) Execute(url string) (string, error) {
	ctx := context.Background()

	result := usecase.persistence.GenerateID(ctx)

	if result.Err != nil {
		return "", result.Err
	}

	short_url := usecase.codec.Encode(uint64(result.Data))

	if err := usecase.persistence.Save(ctx, domain.Url{
		ID:    result.Data,
		Long:  url,
		Short: short_url,
	}); err != nil {
		return "", err
	}

	return short_url, nil
}
