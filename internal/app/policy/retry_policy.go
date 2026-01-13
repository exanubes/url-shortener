package policy

type RetryPolicy struct {
	retries int
	current int
}

func NewRetryPolicy(retries int) RetryPolicy {
	return RetryPolicy{retries: retries}
}

func (policy *RetryPolicy) Next() bool {
	if policy.current >= policy.retries {
		return false
	}

	policy.current += 1

	return true
}
