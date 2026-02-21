package shortcode

import (
	"fmt"
	"hash"
	"math/big"

	"github.com/exanubes/url-shortener/internal/domain"
)

var timestamp_mask = uint64((1 << 41) - 1)
var worker_mask = uint64((1 << 10) - 1)
var sequence_mask = uint64((1 << 12) - 1)

type Snowflake struct {
	worker_hash hash.Hash64
	clock       domain.Clock
}

var _ TokenSpaceGenerator = (*Snowflake)(nil)

func NewSnowflakeGenerator(worker_hash hash.Hash64, clock domain.Clock) *Snowflake {
	return &Snowflake{
		worker_hash: worker_hash,
		clock:       clock,
	}
}

func (generator Snowflake) Generate() (Token, error) {
	now := generator.clock.Now()
	var seq uint64 = 0

	id := ((uint64(now.Local().UnixMilli()) & timestamp_mask) << 22) |
		((generator.worker_hash.Sum64() & worker_mask) << 12) |
		(seq & sequence_mask)

	bigint := (new(big.Int)).SetUint64(id)
	return NewToken(bigint, int64(11))
}
