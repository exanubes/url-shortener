package inmemory

import (
	"context"
	"fmt"

	"github.com/exanubes/url-shortener/internal/domain"
)

type Repository struct {
	id_counter uint64
	cache      map[uint64]domain.Url
}

func NewInmemoryRepository() *Repository {
	return &Repository{
		id_counter: 11_157,
		cache:      make(map[uint64]domain.Url),
	}
}

func (repository *Repository) Save(ctx context.Context, input domain.Url) error {
	repository.cache[input.ID] = input
	return nil
}

func (repository *Repository) Get(ctx context.Context, id uint64) domain.GetUrlOutput {
	url, exists := repository.cache[id]
	if !exists {
		return domain.GetUrlOutput{
			Err: fmt.Errorf("Url does not exist"),
		}
	}
	return domain.GetUrlOutput{
		Data: url,
	}
}

func (repository *Repository) GenerateID(ctx context.Context) domain.GenerateIDOutput {
	repository.id_counter += 1
	return domain.GenerateIDOutput{
		Data: repository.id_counter,
	}
}
