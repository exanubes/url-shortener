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

func (usecase *ResolveUrl) Execute(ctx context.Context, short_url domain.ShortCode) (domain.ResolveUrlCommandOutput, error) {
	link, err := usecase.resolver.Resolve(ctx, short_url)

	if err != nil {
		return domain.ResolveUrlCommandOutput{}, err
	}
	now := time.Now()
	url, err := link.Visit(now)

	if err != nil {
		return domain.ResolveUrlCommandOutput{}, err
	}

	if link.SingleUse() {
		if err := usecase.consumer.Consume(ctx, short_url); err != nil {
			return domain.ResolveUrlCommandOutput{}, err
		}

		link.Consume(now)
	}

	usecase.publisher.Publish(ctx, domain.LinkVisited{
		ShortCode: short_url.String(),
		VisitedAt: now,
	})

	expiration_status, err := link.ExpirationStatus(time.Now())
	if err != nil {
		//TODO: log
	}

	return domain.ResolveUrlCommandOutput{
		Url:    url,
		Status: expiration_status,
	}, nil
}
