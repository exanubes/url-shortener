package resolveurl

import (
	"context"
	"time"

	"github.com/exanubes/url-shortener/internal/domain"
)

type ResolveUrl struct {
	resolver  LinkResolver
	consumer  LinkConsumer
	publisher EventPublisher
}

func New(resolver LinkResolver, consumer LinkConsumer, publisher EventPublisher) *ResolveUrl {
	return &ResolveUrl{
		resolver:  resolver,
		consumer:  consumer,
		publisher: publisher,
	}
}

func (usecase *ResolveUrl) Execute(ctx context.Context, short_url domain.ShortCode) (domain.Url, error) {
	link, err := usecase.resolver.Resolve(ctx, short_url)

	if err != nil {
		return domain.Url{}, err
	}
	now := time.Now()
	url, err := link.Visit(now)

	if err != nil {
		return domain.Url{}, err
	}

	if link.SingleUse() {
		if err := usecase.consumer.Consume(ctx, short_url); err != nil {
			return domain.Url{}, err
		}
	}

	usecase.publisher.Publish(ctx, domain.LinkVisited{
		ShortCode: short_url.String(),
		VisitedAt: now,
	})

	return url, nil
}
