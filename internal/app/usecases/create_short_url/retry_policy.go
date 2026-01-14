package createshorturl

import (
	"errors"

	"github.com/exanubes/url-shortener/internal/domain"
)

type RetryPolicy struct {
	retries int
}

func NewRetryPolicy(retries int) *RetryPolicy {
	return &RetryPolicy{retries: retries}
}

func (policy *RetryPolicy) Next() bool {
	if policy.retries <= 0 {
		return false
	}

	policy.retries -= 1

	return true
}

func (policy *RetryPolicy) Verify(err error) bool {
	return errors.Is(err, domain.ErrShortCodeCollision)
}

type RetryPolicyFactory struct {
	retries int
}

func NewRetryPolicyFactory(retries int) RetryPolicyFactory {
	return RetryPolicyFactory{retries}
}

func (factory RetryPolicyFactory) Create() Policy {
	return NewRetryPolicy(factory.retries)
}
