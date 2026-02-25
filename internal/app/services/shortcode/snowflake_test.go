package shortcode_test

import (
	"hash/fnv"
	"math/big"
	"testing"
	"time"

	"github.com/exanubes/url-shortener/internal/app/services/shortcode"
)

func TestSnowflakeIdGenerator(t *testing.T) {
	hasher := fnv.New64a()
	hasher.Write([]byte("testing"))
	generator := shortcode.NewSnowflakeGenerator(hasher, time.Time{}, clock{})

	token, err := generator.Generate()
	expected := big.NewInt(7430930526361071616)

	if err != nil {
		t.Fatalf("ERROR: %s", err.Error())
	}

	if token.Value().Cmp(expected) != 0 {
		t.Fatalf("Expected %d, received %d", expected.Uint64(), token.Value().Uint64())
	}

}
