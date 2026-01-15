package inmemory

import (
	"context"

	"github.com/exanubes/url-shortener/internal/domain"
)

type Repository struct {
	cache map[string]domain.LinkState
}

func NewInmemoryRepository() *Repository {
	return &Repository{
		cache: make(map[string]domain.LinkState),
	}
}

func (repository *Repository) Write(ctx context.Context, link *domain.Link) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if _, exists := repository.cache[link.ShortCode().String()]; exists {
		return domain.ErrShortCodeCollision
	}

	repository.cache[link.ShortCode().String()] = link.Snapshot()
	return nil
}

func (repository *Repository) Resolve(ctx context.Context, input domain.ShortCode) (domain.Url, error) {
	if err := ctx.Err(); err != nil {
		return domain.Url{}, err
	}

	link_state, exists := repository.cache[input.String()]

	if !exists {
		return domain.Url{}, domain.ErrUrlNotFound
	}

	return link_state.Url, nil
}
