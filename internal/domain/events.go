package domain

import "time"

type Event interface {
	DomainEvent()
}

type LinkVisited struct {
	ShortCode string
	VisitedAt time.Time
	IpAddress string
}

func (LinkVisited) DomainEvent() {}
