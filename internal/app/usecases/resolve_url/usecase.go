package resolveurl

import (
	"context"
	"fmt"

	"github.com/exanubes/url-shortener/internal/domain"
)

type ResolveUrl struct {
	resolver LinkResolver
	consumer LinkConsumer
	clock    domain.Clock
}

func New(resolver LinkResolver, consumer LinkConsumer, clock domain.Clock) *ResolveUrl {
	return &ResolveUrl{
		resolver: resolver,
		consumer: consumer,
		clock:    clock,
	}
}

func (usecase *ResolveUrl) Execute(ctx context.Context, short_url domain.ShortCode) (domain.ResolveUrlCommandOutput, error) {
	link, err := usecase.resolver.Resolve(ctx, short_url)

	if err != nil {
		return domain.ResolveUrlCommandOutput{}, err
	}
	now := usecase.clock.Now()
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

	expiration_status, err := link.ExpirationStatus(usecase.clock.Now())
	if err != nil {
		fmt.Println(err.Error())
	}

	return domain.ResolveUrlCommandOutput{
		Url:    url,
		Status: expiration_status,
	}, nil
}
