package domain

import (
	"context"
)

type PersistenceProvider interface {
	Save(ctx context.Context, input Url) error
	Get(ctx context.Context, id uint64) GetUrlOutput
	GenerateID(ctx context.Context) GenerateIDOutput
}

type Codec interface {
	Encode(input uint64) string
	Decode(input string) (uint64, error)
}

type TokenSpaceGenerator interface {
	Generate() (Token, error)
}

type ForCreatingUrls interface {
	Execute(url string) (string, error)
}
type ForVisitingUrls interface {
	Execute(short_url string) (string, error)
}
