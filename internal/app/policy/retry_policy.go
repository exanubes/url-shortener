package policy

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
