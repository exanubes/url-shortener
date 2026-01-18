package analytics

import (
	"context"
	"time"

	"github.com/exanubes/url-shortener/internal/domain"
)

type LinkVisitor interface {
	Visit(context.Context, domain.ShortCode, time.Time) error
}
