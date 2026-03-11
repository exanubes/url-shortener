package dynamodb

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/exanubes/url-shortener/internal/domain"
	"github.com/exanubes/url-shortener/internal/infrastructure/persistence/dynamodb/internal"
)

type Repository struct {
	client *client
	clock  domain.Clock
}

func NewRepository(client *client, clock domain.Clock) *Repository {
	return &Repository{client: client, clock: clock}
}

func (repository *Repository) Write(ctx context.Context, link *domain.Link) error {
	snapshot := link.Snapshot()

	queries := repository.client.Queries()
	row := internal.LinkRow{
		Shortcode: snapshot.Shortcode.String(),
		Url:       snapshot.Url.String(),
		CreatedAt: snapshot.CreatedAt,
		Version:   "1.0.0",
	}

	policy_dtos, err := internal.SerializePolicies(snapshot.PolicySpecs)
	if err != nil {
		return err
	}

	row.PolicySpecs = policy_dtos

	err = queries.CreateLink(ctx, row)

	if err != nil {
		var exception *types.ConditionalCheckFailedException
		if errors.As(err, &exception) {
			return domain.ErrShortCodeCollision
		}

		return err
	}

	return nil
}

func (repository *Repository) Resolve(ctx context.Context, input domain.ShortCode) (*domain.Link, error) {
	queries := repository.client.Queries()

	row, err := queries.GetLink(ctx, input.String())

	if err != nil {
		return nil, err
	}
	url, err := domain.NewUrl(row.Url)

	if err != nil {
		return nil, err
	}

	short_code, err := domain.NewShortCodeFromParam(row.Shortcode)

	if err != nil {
		return nil, err
	}
	policy_specs, err := internal.DeserializePolicies(row.PolicySpecs)

	if err != nil {
		return nil, err
	}

	state := domain.LinkState{
		Url:         url,
		Shortcode:   short_code,
		PolicySpecs: policy_specs,
		CreatedAt:   row.CreatedAt,
		ConsumedAt:  row.ConsumedAt,
	}

	return domain.RehydrateLink(state), nil
}

// Consume single-use link, do not use with multi-use links
func (repository *Repository) Consume(ctx context.Context, input domain.ShortCode) error {
	err := repository.client.Queries().ConsumeSingleUseLink(ctx, internal.ConsumeSingleUseLinkParams{Shortcode: input.String(), ConsumedAt: repository.clock.Now()})

	var exception *types.ConditionalCheckFailedException
	if errors.As(err, &exception) {
		return domain.ErrLinkConsumed
	}

	return err
}

func (repository *Repository) Visit(ctx context.Context, event domain.LinkVisited) error {
	return repository.client.Queries().LogLinkVisit(ctx, internal.LogLinkVisitParams{Shortcode: event.ShortCode, VisitedAt: event.VisitedAt})
}
