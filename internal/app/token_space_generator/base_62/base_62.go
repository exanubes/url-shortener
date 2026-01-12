package base62

import (
	"crypto/rand"
	"math/big"

	"github.com/exanubes/url-shortener/internal/domain"
)

type Base62TokenSpace struct {
	size int64
}

var _ domain.TokenSpaceGenerator = (*Base62TokenSpace)(nil)

func New(size int64) *Base62TokenSpace {
	return &Base62TokenSpace{
		size: size,
	}
}

func (generator Base62TokenSpace) Generate() (domain.Token, error) {
	max_value := new(big.Int).Exp(big.NewInt(62), big.NewInt(generator.size), nil)
	result, err := rand.Int(rand.Reader, max_value)

	if err != nil {
		return domain.Token{}, err
	}

	return domain.NewToken(result, generator.size)
}
