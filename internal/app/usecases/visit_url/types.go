package visiturl

import (
	"context"

	"github.com/exanubes/url-shortener/internal/domain"
)

type LinkEventStore interface {
	Visit(context.Context, domain.LinkVisited) error
}

type LinkEventParser interface {
	Parse(context.Context, string) (domain.LinkVisited, error)
}
