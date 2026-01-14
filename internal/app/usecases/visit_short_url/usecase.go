package visitshorturl

import (
	"context"

	"github.com/exanubes/url-shortener/internal/domain"
)

type VisitShortUrl struct {
	resolver UrlResolver
}

func New(resolver UrlResolver) *VisitShortUrl {
	return &VisitShortUrl{
		resolver: resolver,
	}
}

func (usecase *VisitShortUrl) Execute(ctx context.Context, short_url domain.ShortCode) (domain.Url, error) {
	result, err := usecase.resolver.Resolve(ctx, short_url)

	if err != nil {
		return domain.Url{}, err
	}

	return result, nil
}
