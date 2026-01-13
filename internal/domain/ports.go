package domain

import (
	"context"
)

type PersistenceProvider interface {
	Save(ctx context.Context, input Url) error
	Get(ctx context.Context, input string) (Url, error)
}

type Encoder interface {
	Encode(Token) string
}

type TokenSpaceGenerator interface {
	Generate() (Token, error)
}

type ForCreatingUrls interface {
	Execute(ctx context.Context, url string) (string, error)
}

type ForVisitingUrls interface {
	Execute(ctx context.Context, short_url string) (string, error)
}

type ShortCodeGenerator interface {
	Generate() (ShortCode, error)
}
