package resolveurl

import (
	"context"

	"github.com/exanubes/url-shortener/internal/domain"
)

type LinkResolver interface {
	Resolve(context.Context, domain.ShortCode) (*domain.Link, error)
}

type LinkConsumer interface {
	Consume(context.Context, domain.ShortCode) error
}

type UseCase interface {
	Execute(context.Context, domain.ShortCode) (domain.Url, error)
}

type EventPublisher interface {
	Publish(context.Context, domain.Event) error
}
