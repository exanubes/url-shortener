package createshorturl

import (
	"context"

	"github.com/exanubes/url-shortener/internal/domain"
)

type LinkWriter interface {
	Write(context.Context, *domain.Link) error
}

type Policy interface {
	Next() bool
	Verify(error) bool
}

type PolicyFactory interface {
	Create() Policy
}

type UseCase interface {
	Execute(context.Context, domain.Url) (*domain.Link, error)
}
