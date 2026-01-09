package domain

import "context"

type PersistenceProvider interface {
	Save(ctx context.Context, input Url) error
	Get(ctx context.Context, id int) GetUrlOutput
	GenerateID(ctx context.Context) GenerateIDOutput
}

type Encoder interface {
	Encode(input uint64) string
}

type Decoder interface {
	Decode(input string) (uint64, error)
}
