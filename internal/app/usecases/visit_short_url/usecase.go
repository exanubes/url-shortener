package visitshorturl

import (
	"context"
	"time"

	"github.com/exanubes/url-shortener/internal/domain"
)

type VisitShortUrl struct {
	resolver LinkResolver
	consumer LinkConsumer
}

func New(resolver LinkResolver, consumer LinkConsumer) *VisitShortUrl {
	return &VisitShortUrl{
		resolver: resolver,
		consumer: consumer,
	}
}

func (usecase *VisitShortUrl) Execute(ctx context.Context, short_url domain.ShortCode) (domain.Url, error) {
	link, err := usecase.resolver.Resolve(ctx, short_url)

	if err != nil {
		return domain.Url{}, err
	}

	url, err := link.Visit(time.Now())

	if err != nil {
		return domain.Url{}, err
	}

	if link.SingleUse() {
		if err := usecase.consumer.Consume(ctx, short_url); err != nil {
			return domain.Url{}, err
		}
	}

	return url, nil
}
