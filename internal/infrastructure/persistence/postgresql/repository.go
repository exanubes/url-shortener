package postgresql

import (
	"context"
	"encoding/json"
	"net"
	"time"

	"github.com/exanubes/url-shortener/internal/domain"
	"github.com/exanubes/url-shortener/internal/infrastructure/persistence/postgresql/sqlc"
	"github.com/sqlc-dev/pqtype"
)

type Repository struct {
	client *client
}

type policyJSON struct {
	Kind   string          `json:"kind"`
	Config json.RawMessage `json:"config"`
}

func NewPostgresqlRepository(client *client) *Repository {
	return &Repository{
		client: client,
	}
}

func (repository *Repository) Write(ctx context.Context, link *domain.Link) error {
	snapshot := link.Snapshot()

	tx, err := repository.client.db.BeginTx(ctx, nil)

	if err != nil {
		return err
	}

	defer tx.Rollback()

	queries := repository.client.Queries().WithTx(tx)

	exists, err := queries.CheckShortCodeExists(ctx, snapshot.Shortcode.String())

	if err != nil {
		return err
	}

	if exists {
		return domain.ErrShortCodeCollision
	}

	err = queries.CreateLink(ctx, sqlc.CreateLinkParams{
		ID:        snapshot.Shortcode.String(),
		Url:       snapshot.Url.String(),
		CreatedAt: snapshot.CreatedAt,
	})

	if err != nil {
		return err
	}

	for _, policy := range snapshot.PolicySpecs {
		var err error
		var config json.RawMessage

		if policy.Params != nil {
			config, err = json.Marshal(policy.Params)
			if err != nil {
				return err
			}
		} else {
			config = []byte("{}")
		}

		err = queries.CreateLinkPolicy(ctx, sqlc.CreateLinkPolicyParams{
			LinkID: snapshot.Shortcode.String(),
			Kind:   string(policy.Kind),
			Config: config,
		})

		if err != nil {
			return err
		}
	}

	return tx.Commit()

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

	short_code, err := domain.NewShortCodeFromParam(row.ID)

	if err != nil {
		return nil, err
	}

	policySpecs, err := deserializePolicySpecs(row.Policies)
	if err != nil {
		return nil, err
	}

	state := domain.LinkState{
		Url:         url,
		Shortcode:   short_code,
		PolicySpecs: policySpecs,
		CreatedAt:   row.CreatedAt,
		ConsumedAt:  row.ConsumedAt.Time,
	}

	return domain.RehydrateLink(state), nil
}

// Consume single-use link, do not use with multi-use links
func (repository *Repository) Consume(ctx context.Context, input domain.ShortCode) error {
	return repository.client.Queries().ConsumeSingleUseLink(ctx, input.String())
}

func (repository *Repository) Visit(ctx context.Context, key domain.ShortCode, date time.Time) error {
	mock_ip := pqtype.Inet{
		IPNet: net.IPNet{
			IP:   net.ParseIP("0.0.0.0"),
			Mask: net.CIDRMask(32, 32),
		},
		Valid: true,
	}
	return repository.client.Queries().LogLinkVisit(ctx, sqlc.LogLinkVisitParams{
		LinkID:    key.String(),
		VisitedAt: date,
		IpAddress: mock_ip,
	})
}

func (repository *Repository) Close() error {
	return repository.client.db.Close()
}

func deserializePolicySpecs(policiesJSON json.RawMessage) ([]domain.PolicySpec, error) {
	if len(policiesJSON) == 0 || string(policiesJSON) == "null" {
		return []domain.PolicySpec{}, nil
	}

	var policies []policyJSON
	if err := json.Unmarshal(policiesJSON, &policies); err != nil {
		return nil, err
	}

	specs := make([]domain.PolicySpec, 0, len(policies))

	for _, p := range policies {
		switch p.Kind {
		case string(domain.PolicyKind_SingleUse):
			specs = append(specs, domain.PolicySpec{
				Kind:   domain.PolicyKind_SingleUse,
				Params: nil,
			})

		case string(domain.PolicyKind_MaxAge):
			var config_map map[string]any
			if err := json.Unmarshal(p.Config, &config_map); err != nil {
				return nil, err
			}

			duration_value, ok := config_map["duration"]
			if !ok {
				return nil, domain.ErrInvalidPolicySpecParams
			}

			// NOTE: JSON unmarshals numbers as float64 (stored as float64 from JSON number when marshalling PolicySpec.Params)
			duration_nanos, ok := duration_value.(float64)
			if !ok {
				return nil, domain.ErrInvalidPolicySpecParams
			}

			specs = append(specs, domain.PolicySpec{
				Kind: domain.PolicyKind_MaxAge,
				Params: map[string]any{
					"duration": time.Duration(duration_nanos),
				},
			})

		default:
			continue
		}
	}

	return specs, nil
}
