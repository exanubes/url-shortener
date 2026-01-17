package domain

import "time"

type LinkVisited struct {
	ShortCode string
	VisitedAt time.Time
	IpAddress string
}
