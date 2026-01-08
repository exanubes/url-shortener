package inmemory

import (
	"context"
	"fmt"

	"github.com/exanubes/url-shortener/internal/domain"
)

type Repository struct {
	id_counter int
	cache      map[int]domain.Url
}

func NewInmemoryRepository() *Repository {
	return &Repository{
		id_counter: 11_157,
	}
}

func (repository *Repository) Save(ctx context.Context, input domain.Url) error {
	repository.cache[input.ID] = input
	return nil
}

func (repository *Repository) Get(ctx context.Context, id int) domain.GetUrlOutput {
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
