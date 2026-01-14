package domain

import (
	"context"
)

type PersistenceProvider interface {
	Save(context.Context, Url, ShortCode) error
	Get(context.Context, ShortCode) (Url, error)
}

type Encoder interface {
	Encode(Token) string
}

type TokenSpaceGenerator interface {
	Generate() (Token, error)
}

type RetryPolicy interface {
	Next() bool
}

type RetryPolicyFactory interface {
	Create() RetryPolicy
}

type ForCreatingUrls interface {
	Execute(ctx context.Context, url Url) (ShortCode, error)
}

type ForVisitingUrls interface {
	Execute(ctx context.Context, short_url ShortCode) (Url, error)
}

type ShortCodeGenerator interface {
	Generate() (ShortCode, error)
}
