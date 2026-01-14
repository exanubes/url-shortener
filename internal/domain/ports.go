package domain

import (
	"context"

	"github.com/exanubes/url-shortener/internal/app/policy"
)

type PersistenceProvider interface {
	Save(context.Context, Url, ShortCode) error
	Get(ctx context.Context, input string) (Url, error)
}

type Encoder interface {
	Encode(Token) string
}

type TokenSpaceGenerator interface {
	Generate() (Token, error)
}

type ForCreatingUrls interface {
	Execute(ctx context.Context, url Url, policy policy.RetryPolicy) (ShortCode, error)
}

type ForVisitingUrls interface {
	Execute(ctx context.Context, short_url string) (Url, error)
}

type ShortCodeGenerator interface {
	Generate() (ShortCode, error)
}
