package policy

import "github.com/exanubes/url-shortener/internal/domain"

type RetryPolicyFactory struct {
	retries int
}

func NewRetryPolicyFactory(retries int) RetryPolicyFactory {
	return RetryPolicyFactory{retries}
}

func (factory RetryPolicyFactory) Create() domain.RetryPolicy {
	return NewRetryPolicy(factory.retries)
}
