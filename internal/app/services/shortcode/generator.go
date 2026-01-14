package shortcode

import (
	"crypto/rand"
	"math/big"
)

type Base62TokenSpace struct {
	size int64
}

var _ TokenSpaceGenerator = (*Base62TokenSpace)(nil)

func NewGenerator(size int64) *Base62TokenSpace {
	return &Base62TokenSpace{
		size: size,
	}
}

func (generator Base62TokenSpace) Generate() (Token, error) {
	max_value := new(big.Int).Exp(big.NewInt(62), big.NewInt(generator.size), nil)
	result, err := rand.Int(rand.Reader, max_value)

	if err != nil {
		return Token{}, err
	}

	return NewToken(result, generator.size)
}
