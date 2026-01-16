package domain

import (
	"time"
)

type LinkState struct {
	Url       Url
	Shortcode ShortCode
	Policy    ExpirationPolicy
	CreatedAt time.Time
	Visits    int
	LastVisit time.Time
}

type Link struct {
	url        Url
	shortcode  ShortCode
	policy     ExpirationPolicy
	created_at time.Time
	visits     int
	last_visit time.Time
}

func RehydrateLink(state LinkState) *Link {
	return new_link(state.Url, state.Shortcode, state.Policy, state.CreatedAt, state.Visits, state.LastVisit)
}

func CreateLink(url Url, shortcode ShortCode, policy ExpirationPolicy, created_at time.Time) *Link {
	return new_link(url, shortcode, policy, created_at, 0, time.Time{})
}

func new_link(url Url, shortcode ShortCode, policy ExpirationPolicy, created_at time.Time, visits int, last_visit time.Time) *Link {
	return &Link{
		url:        url,
		shortcode:  shortcode,
		policy:     policy,
		created_at: created_at,
		visits:     visits,
		last_visit: last_visit,
	}
}

func (link *Link) Visit(now time.Time) (Url, error) {
	expired := link.policy.Expired(ExpirationContext{
		CreatedAt:     link.created_at,
		LastVisitedAt: link.last_visit,
		VisitCount:    link.visits,
		Now:           now,
	})

	if expired {
		return Url{}, ErrLinkExpired
	}

	link.visits += 1
	link.last_visit = now

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
		Url:       link.url,
		Shortcode: link.shortcode,
		Policy:    link.policy,
		CreatedAt: link.created_at,
		Visits:    link.visits,
		LastVisit: link.last_visit,
	}
}
