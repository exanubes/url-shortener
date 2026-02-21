package shortcode_test

import (
	"hash/fnv"
	"math/big"
	"testing"
	"time"

	"github.com/exanubes/url-shortener/internal/app/services/shortcode"
)

type clock struct{}

func (clock) Now() time.Time {
	str := "2026-02-21T11:04:57.497"
	date, _ := time.Parse("2006-01-02T15:04:05.000", str)
	return date
}

func TestSnowflakeIdGenerator(t *testing.T) {
	hasher := fnv.New64a()
	hasher.Write([]byte("testing"))
	generator := shortcode.NewSnowflakeGenerator(hasher, clock{})

	token, err := generator.Generate()
	expected := big.NewInt(7430930526361071616)

	if err != nil {
		t.Fatalf("ERROR: %s", err.Error())
	}

	if token.Value().Cmp(expected) != 0 {
		t.Fatalf("Expected %d, received %d", expected.Uint64(), token.Value().Uint64())
	}

}
