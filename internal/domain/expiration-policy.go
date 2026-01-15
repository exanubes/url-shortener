package domain

import "time"

type ExpirationPolicy interface {
	Expired(ExpirationContext) bool
}

type ExpirationContext struct {
	CreatedAt     time.Time
	LastVisitedAt *time.Time
	VisitCount    int
	Now           time.Time
}

type MaxLinkAgeExpirationPolicy struct {
	age time.Duration
}

func NewMaxLinkAgeExpirationPolicy(age time.Duration) MaxLinkAgeExpirationPolicy {
	return MaxLinkAgeExpirationPolicy{
		age: age,
	}
}

func (policy MaxLinkAgeExpirationPolicy) Expired(context ExpirationContext) bool {
	expires_at := context.CreatedAt.Add(policy.age)

	return context.Now.After(expires_at)

}
