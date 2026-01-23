package domain

import (
	"time"
)

type PolicyKind string

const (
	PolicyKind_SingleUse = "single_use"
	PolicyKind_MaxAge    = "max_age"
)

type PolicySpec struct {
	Kind   PolicyKind
	Params PolicyParams
}

type PolicyParams interface {
	policy_params()
}

type SingleUseParams struct{}

func (SingleUseParams) policy_params() {}

type MaxAgeParams struct {
	TTL time.Duration
}

func (MaxAgeParams) policy_params() {}

type PolicySettings struct {
	MaxAge     time.Duration
	ConsumedAt time.Time
	SingleUse  bool
}

func (settings PolicySettings) HasMaxAgeLimit() bool {
	return settings.MaxAge > 0
}

func (settings PolicySettings) IsSingleUse() bool {
	return settings.SingleUse
}

func NewPolicySettings(max_age time.Duration, single_use bool) (PolicySettings, error) {
	if max_age < time.Minute {
		return PolicySettings{}, ErrExceededMinAge
	}
	day := 24 * time.Hour
	year := 365 * day
	if max_age > year {
		return PolicySettings{}, ErrExceededMaxAge
	}

	return PolicySettings{MaxAge: max_age, SingleUse: single_use}, nil
}

type ExpirationPolicy interface {
	Expired(ExpirationContext) bool
}

type ExpirationContext struct {
	CreatedAt  time.Time
	ConsumedAt time.Time
	Now        time.Time
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
	return !context.ConsumedAt.IsZero()
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

func build_expiration_policy(spec PolicySpec) (ExpirationPolicy, error) {
	switch params := spec.Params.(type) {
	case SingleUseParams:
		if spec.Kind != PolicyKind_SingleUse {
			return nil, ErrInvalidPolicySpecParams
		}

		return NewOneTimeLinkExpirationPolicy(), nil

	case MaxAgeParams:
		if spec.Kind != PolicyKind_MaxAge {
			return nil, ErrInvalidPolicySpecParams
		}

		return NewMaxLinkAgeExpirationPolicy(params.TTL)
	}

	return nil, ErrUnsupportedPolicyKind
}
