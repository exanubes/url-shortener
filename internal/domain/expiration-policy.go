package domain

import (
	"time"
)

const max_visits_limit = 100

type PolicySettings struct {
	MaxVisits int
	MaxAge    time.Duration
	Usage     LinkUsage
}

func (settings PolicySettings) HasMaxVisitsLimit() bool {
	return settings.MaxVisits > 0
}

func (settings PolicySettings) HasMaxAgeLimit() bool {
	return settings.MaxAge > 0
}

func (settings PolicySettings) IsSingleUse() bool {
	return settings.Usage == LinkUsage_Single
}

func NewPolicySettings(max_visits int, max_age time.Duration, usage LinkUsage) (PolicySettings, error) {
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
	return PolicySettings{MaxVisits: max_visits, MaxAge: max_age, Usage: usage}, nil
}

type ExpirationPolicy interface {
	Expired(ExpirationContext) bool
}

type ExpirationContext struct {
	CreatedAt     time.Time
	LastVisitedAt time.Time
	VisitCount    int
	Now           time.Time
	Status        LinkStatus
}

type MaxLinkAgeExpirationPolicy struct {
	age time.Duration
}

func NewMaxLinkAgeExpirationPolicy(age time.Duration) (MaxLinkAgeExpirationPolicy, error) {
	if age < time.Minute {
		return MaxLinkAgeExpirationPolicy{}, ErrExceededMinAge
	}

	return MaxLinkAgeExpirationPolicy{age}, nil
}

func (policy MaxLinkAgeExpirationPolicy) Expired(context ExpirationContext) bool {
	expires_at := context.CreatedAt.Add(policy.age)

	return context.Now.After(expires_at)

}

type OneTimeLinkExpirationPolicy struct {
}

func NewOneTimeLinkExpirationPolicy() OneTimeLinkExpirationPolicy {

	return OneTimeLinkExpirationPolicy{}
}

func (OneTimeLinkExpirationPolicy) Expired(context ExpirationContext) bool {
	return context.Status != LinkStatus_New
}

type ChainExpirationPolicy struct {
	policies []ExpirationPolicy
}

func NewChainExpirationPolicy(policies []ExpirationPolicy) (ChainExpirationPolicy, error) {
	if len(policies) == 0 {
		return ChainExpirationPolicy{}, ErrUndefinedExpirationPolicy
	}

	return ChainExpirationPolicy{policies}, nil
}

func (policy ChainExpirationPolicy) Expired(context ExpirationContext) bool {
	for _, p := range policy.policies {
		if p.Expired(context) {
			return true
		}
	}

	return false
}
