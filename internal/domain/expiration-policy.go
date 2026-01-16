package domain

import "time"

const max_visits_limit = 100

type PolicySettings struct {
	MaxVisits int
	MaxAge    time.Duration
}

func (settings PolicySettings) HasMaxVisitsLimit() bool {
	return settings.MaxVisits > 0
}

func (settings PolicySettings) HasMaxAgeLimit() bool {
	return settings.MaxAge > 0
}

func NewPolicySettings(max_visits int, max_age time.Duration) (PolicySettings, error) {
	if max_visits < 0 {
		return PolicySettings{}, ErrExceededMinVisits
	}

	if max_visits > max_visits_limit {
		return PolicySettings{}, ErrExceededMaxVisits
	}

	if max_age < time.Minute {
		return PolicySettings{}, ErrExceededMinAge
	}
	day := 24 * time.Hour
	year := 365 * day
	if max_age > year {
		return PolicySettings{}, ErrExceededMaxAge
	}
	return PolicySettings{MaxVisits: max_visits, MaxAge: max_age}, nil
}

type ExpirationPolicy interface {
	Expired(ExpirationContext) bool
}

type ExpirationContext struct {
	CreatedAt     time.Time
	LastVisitedAt time.Time
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
