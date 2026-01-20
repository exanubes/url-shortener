package domain

import (
	"time"
)

type LinkUsage int

const (
	LinkUsage_Multi LinkUsage = iota
	LinkUsage_Single
)

type LinkState struct {
	Url        Url
	Shortcode  ShortCode
	Policy     ExpirationPolicy
	CreatedAt  time.Time
	ConsumedAt time.Time
	Usage      LinkUsage
}

type Link struct {
	url         Url
	shortcode   ShortCode
	policy      ExpirationPolicy
	created_at  time.Time
	consumed_at time.Time
	usage       LinkUsage
}

func RehydrateLink(state LinkState) *Link {
	return new_link(state.Url, state.Shortcode, state.Policy, state.CreatedAt, state.ConsumedAt, state.Usage)
}

func CreateLink(url Url, shortcode ShortCode, policy ExpirationPolicy, created_at time.Time, usage LinkUsage) *Link {
	return new_link(url, shortcode, policy, created_at, time.Time{}, usage)
}

func new_link(url Url, shortcode ShortCode, policy ExpirationPolicy, created_at time.Time, consumed_at time.Time, usage LinkUsage) *Link {
	return &Link{
		url:         url,
		shortcode:   shortcode,
		policy:      policy,
		created_at:  created_at,
		consumed_at: consumed_at,
		usage:       usage,
	}
}

func (link *Link) Visit(now time.Time) (Url, error) {
	expired := link.policy.Expired(ExpirationContext{
		CreatedAt:  link.created_at,
		Now:        now,
		ConsumedAt: link.consumed_at,
	})

	if expired {
		return Url{}, ErrLinkExpired
	}

	return link.url, nil
}

func (link Link) Url() Url {
	return link.url
}

func (link Link) ShortCode() ShortCode {
	return link.shortcode
}

func (link *Link) Snapshot() LinkState {
	return LinkState{
		Url:        link.url,
		Shortcode:  link.shortcode,
		Policy:     link.policy,
		CreatedAt:  link.created_at,
		ConsumedAt: link.consumed_at,
		Usage:      link.usage,
	}
}

func (link *Link) SingleUse() bool {
	return link.usage == LinkUsage_Single
}
