package postgresql

import (
	"context"
	"time"

	"github.com/exanubes/url-shortener/internal/domain"
)

type Repository struct {
	client *client
}

func NewPostgresqlRepository(client *client) *Repository {
	return &Repository{
		client: client,
	}
}

func (repository *Repository) Write(ctx context.Context, link *domain.Link) error {
	return nil
}

func (repository *Repository) Resolve(ctx context.Context, input domain.ShortCode) (*domain.Link, error) {
	return nil, nil
}

// Consume single-use link, do not use with multi-use links
func (repository *Repository) Consume(ctx context.Context, input domain.ShortCode) error {
	return nil
}

func (repository *Repository) Visit(ctx context.Context, key domain.ShortCode, date time.Time) error {
	return nil
}
