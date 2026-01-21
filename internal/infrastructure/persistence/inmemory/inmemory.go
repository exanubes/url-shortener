package inmemory

import (
	"context"
	"time"

	"github.com/exanubes/url-shortener/internal/domain"
)

type Repository struct {
	links  map[string]domain.LinkState
	visits map[string][]analytics
}

func NewInmemoryRepository() *Repository {
	return &Repository{
		links:  make(map[string]domain.LinkState),
		visits: make(map[string][]analytics),
	}
}

func (repository *Repository) Write(ctx context.Context, link *domain.Link) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if _, exists := repository.links[link.ShortCode().String()]; exists {
		return domain.ErrShortCodeCollision
	}

	repository.links[link.ShortCode().String()] = link.Snapshot()
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

	if link_state.ConsumedAt.IsZero() {
		link_state.ConsumedAt = time.Now()
		repository.links[input.String()] = link_state
		return nil
	}

	return domain.ErrLinkExpired
}

func (repository *Repository) Visit(ctx context.Context, key domain.ShortCode, date time.Time) error {
	visit := analytics{shortcode: key.String(), visited_at: date, ip_address: "0.0.0.0"}
	repository.visits[key.String()] = append(repository.visits[key.String()], visit)

	return nil
}

func (repository *Repository) get(key string) (domain.LinkState, error) {
	link_state, exists := repository.links[key]

	if !exists {
		return domain.LinkState{}, domain.ErrUrlNotFound
	}

	return link_state, nil
}

type analytics struct {
	shortcode  string
	visited_at time.Time
	ip_address string
}
