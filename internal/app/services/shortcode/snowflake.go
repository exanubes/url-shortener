package shortcode

import (
	"hash"
	"math/big"
	"time"

	"github.com/exanubes/url-shortener/internal/domain"
)

var timestamp_mask = uint64((1 << 41) - 1)
var worker_mask = uint64((1 << 10) - 1)
var sequence_mask = uint64((1 << 12) - 1)

type Snowflake struct {
	worker_hash hash.Hash64
	clock       domain.Clock
	prev_time   time.Time
	sequence    uint64
	epoch       time.Time
}

var _ TokenSpaceGenerator = (*Snowflake)(nil)

func NewSnowflakeGenerator(worker_hash hash.Hash64, epoch time.Time, clock domain.Clock) *Snowflake {
	return &Snowflake{
		worker_hash: worker_hash,
		clock:       clock,
		epoch:       epoch,
	}
}

func (generator Snowflake) Generate() (Token, error) {
	now := generator.clock.Now()
	seq := generator.determine_sequence(now)
	timestamp := now.Local().UnixMilli()
	if !generator.epoch.IsZero() {
		timestamp = now.Local().UnixMilli() - generator.epoch.Local().UnixMilli()
	}

	id := ((uint64(timestamp) & timestamp_mask) << 22) |
		((generator.worker_hash.Sum64() & worker_mask) << 12) |
		(seq & sequence_mask)

	bigint := (new(big.Int)).SetUint64(id)

	return NewToken(bigint, int64(11))
}

func (generator Snowflake) determine_sequence(now time.Time) uint64 {
	if now.Local().UnixMilli() == generator.prev_time.Local().UnixMilli() {
		generator.sequence += 1
	} else {
		generator.sequence = 0
		generator.prev_time = now
	}

	return generator.sequence
}
