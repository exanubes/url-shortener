package domain

import (
	"time"
)

type LinkUsage int
type LinkStatus string

const (
	LinkUsage_Multi LinkUsage = iota
	LinkUsage_Single
)

const (
	LinkStatus_New     LinkStatus = "new"
	LinkStatus_Active  LinkStatus = "active"
	LinkStatus_Expired LinkStatus = "expired"
)

func NewLinkStatus(input string) LinkStatus {
	switch input {
	case "new":
		return LinkStatus_New
	case "active":
		return LinkStatus_Active
	case "expired":
		return LinkStatus_Expired
	}

	return "unknown"
}

type LinkState struct {
	Url       Url
	Shortcode ShortCode
	Policy    ExpirationPolicy
	CreatedAt time.Time
	Visits    int
	LastVisit time.Time
	Status    LinkStatus
	Usage     LinkUsage
}

type Link struct {
	url        Url
	shortcode  ShortCode
	policy     ExpirationPolicy
	created_at time.Time
	visits     int
	last_visit time.Time
	status     LinkStatus
	usage      LinkUsage
}

func RehydrateLink(state LinkState) *Link {
	return new_link(state.Url, state.Shortcode, state.Policy, state.CreatedAt, state.Visits, state.LastVisit, state.Status, state.Usage)
}

func CreateLink(url Url, shortcode ShortCode, policy ExpirationPolicy, created_at time.Time, usage LinkUsage) *Link {
	return new_link(url, shortcode, policy, created_at, 0, time.Time{}, LinkStatus_New, usage)
}

func new_link(url Url, shortcode ShortCode, policy ExpirationPolicy, created_at time.Time, visits int, last_visit time.Time, status LinkStatus, usage LinkUsage) *Link {
	return &Link{
		url:        url,
		shortcode:  shortcode,
		policy:     policy,
		created_at: created_at,
		visits:     visits,
		last_visit: last_visit,
		status:     status,
		usage:      usage,
	}
}

func (link *Link) Visit(now time.Time) (Url, error) {
	expired := link.policy.Expired(ExpirationContext{
		CreatedAt:     link.created_at,
		LastVisitedAt: link.last_visit,
		VisitCount:    link.visits,
		Status:        link.status,
		Now:           now,
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
		Url:       link.url,
		Shortcode: link.shortcode,
		Policy:    link.policy,
		CreatedAt: link.created_at,
		Visits:    link.visits,
		LastVisit: link.last_visit,
		Status:    link.status,
		Usage:     link.usage,
	}
}

func (link *Link) SingleUse() bool {
	return link.usage == LinkUsage_Single
}
