package shortcode_test

import (
	"hash/fnv"
	"testing"
	"time"

	"github.com/exanubes/url-shortener/internal/app/services/shortcode"
)

const feistel_key uint64 = 0x8c3f19d2e4a761b5

func TestFeistelNetworkScrambler(t *testing.T) {
	hasher := fnv.New64a()
	hasher.Write([]byte("testing"))
	generator := shortcode.NewSnowflakeGenerator(hasher, time.Time{}, clock{})

	scrambler := shortcode.NewFeistel(feistel_key)

	token, err := generator.Generate()

	token, err = scrambler.Scramble(token)

	var expected uint64 = 4448212407724719056

	if err != nil {
		t.Fatalf("ERROR: %s", err.Error())
	}

	if token.Value().Uint64() != expected {
		t.Fatalf("Expected %d, received %d", expected, token.Value().Uint64())
	}

}
