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

func (repository *Repository) Save(ctx context.Context, input domain.Url) error {
	if _, exists := repository.cache[input.Short]; exists {
		return domain.ErrShortCodeCollision
	}

	repository.cache[input.Short] = input
	return nil
}

func (repository *Repository) Get(ctx context.Context, input string) (domain.Url, error) {
	url, exists := repository.cache[input]

	if !exists {
		return domain.Url{}, domain.ErrUrlNotFound
	}

	return url, nil
}
