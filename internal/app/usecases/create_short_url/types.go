package createshorturl

import (
	"context"

	"github.com/exanubes/url-shortener/internal/domain"
)

type UrlWriter interface {
	Write(context.Context, domain.ShortCode, domain.Url) error
}

type Policy interface {
	Next() bool
	Verify(error) bool
}

type PolicyFactory interface {
	Create() Policy
}

type UseCase interface {
	Execute(context.Context, domain.Url) (domain.ShortCode, error)
}
