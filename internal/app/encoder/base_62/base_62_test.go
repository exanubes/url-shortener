package base62

import (
	"math/big"
	"testing"

	"github.com/exanubes/url-shortener/internal/domain"
)

func TestBase62Encoding(t *testing.T) {
	encoding := Base62Encoder{}
	token, _ := domain.NewToken(big.NewInt(11_157), 7)
	result := encoding.Encode(token)
	expected := "2TX"

	if result != expected {
		t.Fatalf("Nope! Expected '%s', received: '%s'", expected, result)
	}
}
