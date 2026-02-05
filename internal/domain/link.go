package domain

import (
	"time"
)

type ExpirationStatus struct {
	Consumed  bool
	ExpiresAt time.Time
}

type LinkState struct {
	Url         Url
	Shortcode   ShortCode
	PolicySpecs []PolicySpec
	CreatedAt   time.Time
	ConsumedAt  time.Time
}

type Link struct {
	url          Url
	shortcode    ShortCode
	policy_specs []PolicySpec
	created_at   time.Time
	consumed_at  time.Time
}

func RehydrateLink(state LinkState) *Link {
	return new_link(state.Url, state.Shortcode, state.PolicySpecs, state.CreatedAt, state.ConsumedAt)
}

func CreateLink(url Url, shortcode ShortCode, policy_specs []PolicySpec, created_at time.Time) *Link {
	return new_link(url, shortcode, policy_specs, created_at, time.Time{})
}

func new_link(url Url, shortcode ShortCode, policy_specs []PolicySpec, created_at time.Time, consumed_at time.Time) *Link {
	return &Link{
		url:          url,
		shortcode:    shortcode,
		policy_specs: policy_specs,
		created_at:   created_at,
		consumed_at:  consumed_at,
	}
}

func (link *Link) Visit(now time.Time) (Url, error) {
	policy, err := link.policy()

	if err != nil {
		return Url{}, err
	}

	expired := policy.Expired(ExpirationContext{
		CreatedAt:  link.created_at,
		Now:        now,
		ConsumedAt: link.consumed_at,
	})

	if expired {
		return Url{}, ErrLinkExpired
	}

	return link.url, nil
}

func (link Link) ShortCode() ShortCode {
	return link.shortcode
}

func (link Link) Snapshot() LinkState {
	return LinkState{
		Url:         link.url,
		Shortcode:   link.shortcode,
		PolicySpecs: link.policy_specs,
		CreatedAt:   link.created_at,
		ConsumedAt:  link.consumed_at,
	}
}

func (link Link) ExpirationStatus(now time.Time) (ExpirationStatus, error) {
	expiration_policy, err := link.policy()

	if err != nil {
		return ExpirationStatus{}, err
	}

	var consumed bool

	if !link.consumed_at.IsZero() {
		consumed = true
	}

	return ExpirationStatus{
		Consumed: consumed,
		ExpiresAt: expiration_policy.ExpiresAt(ExpirationContext{
			CreatedAt:  link.created_at,
			Now:        now,
			ConsumedAt: link.consumed_at,
		}),
	}, nil
}

func (link *Link) SingleUse() bool {
	for _, spec := range link.policy_specs {
		if spec.Kind == PolicyKind_SingleUse {
			return true
		}
	}

	return false
}

func (link *Link) Consume(now time.Time) {
	link.consumed_at = now
}

func (link *Link) policy() (ExpirationPolicy, error) {
	policies := make([]ExpirationPolicy, len(link.policy_specs))

	for index, spec := range link.policy_specs {
		policy, err := build_expiration_policy(spec)
		if err != nil {
			return nil, err
		}

		policies[index] = policy
	}

	return NewChainExpirationPolicy(policies)
}
