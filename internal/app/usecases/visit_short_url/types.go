package visitshorturl

import (
	"context"

	"github.com/exanubes/url-shortener/internal/domain"
)

type UrlResolver interface {
	Resolve(context.Context, domain.ShortCode) (*domain.Link, error)
}

type UseCase interface {
	Execute(context.Context, domain.ShortCode) (domain.Url, error)
}
