package inmemory

import (
	"context"
	"time"

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

func (repository *Repository) Resolve(ctx context.Context, input domain.ShortCode) (*domain.Link, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	link_state, err := repository.get(input.String())
	if err != nil {
		return nil, err
	}

	return domain.RehydrateLink(link_state), nil
}

// Consume single-use link, do not use with multi-use links
func (repository *Repository) Consume(ctx context.Context, input domain.ShortCode) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	link_state, err := repository.get(input.String())

	if err != nil {
		return err
	}

	if link_state.Usage == domain.LinkUsage_Single && link_state.Status == domain.LinkStatus_New {
		link_state.Status = domain.LinkStatus_Expired
		repository.cache[input.String()] = link_state
		return nil
	}

	return domain.ErrLinkExpired
}

func (repository *Repository) Visit(ctx context.Context, key domain.ShortCode, date time.Time) error {
	link_state, err := repository.get(key.String())

	if err != nil {
		return err
	}

	link_state.Visits += 1
	link_state.LastVisit = date

	repository.cache[key.String()] = link_state

	return nil
}

func (repository *Repository) get(key string) (domain.LinkState, error) {
	link_state, exists := repository.cache[key]

	if !exists {
		return domain.LinkState{}, domain.ErrUrlNotFound
	}

	return link_state, nil
}
