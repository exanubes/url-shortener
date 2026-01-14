package inmemory

import (
	"context"

	"github.com/exanubes/url-shortener/internal/domain"
)

type Repository struct {
	cache map[string]domain.Url
}

func NewInmemoryRepository() *Repository {
	return &Repository{
		cache: make(map[string]domain.Url),
	}
}

func (repository *Repository) Save(ctx context.Context, url domain.Url, short_code domain.ShortCode) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if _, exists := repository.cache[short_code.String()]; exists {
		return domain.ErrShortCodeCollision
	}

	repository.cache[short_code.String()] = url
	return nil
}

func (repository *Repository) Get(ctx context.Context, input domain.ShortCode) (domain.Url, error) {
	if err := ctx.Err(); err != nil {
		return domain.Url{}, err
	}

	url, exists := repository.cache[input.String()]

	if !exists {
		return domain.Url{}, domain.ErrUrlNotFound
	}

	return url, nil
}
