package domain

import "time"

type Link struct {
	url        Url
	shortcode  ShortCode
	expiration ExpirationPolicy
	created_at time.Time
	visits     int
	last_visit *time.Time
}

func (link *Link) Visit(now time.Time) (Url, error) {

	expired := link.expiration.Expired(ExpirationContext{
		CreatedAt:     link.created_at,
		LastVisitedAt: link.last_visit,
		VisitCount:    link.visits,
		Now:           now,
	})

	if expired {
		return Url{}, ErrLinkExpired
	}

	link.visits += 1
	link.last_visit = &now

	return link.url, nil
}
