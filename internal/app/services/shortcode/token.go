package shortcode

import "math/big"

type Token struct {
	value *big.Int
	size  int64
}

func NewToken(value *big.Int, size int64) (Token, error) {

	if value == nil {
		return Token{}, ErrInvalidToken
	}

	if value.Sign() < 0 {
		return Token{}, ErrInvalidToken
	}

	max := new(big.Int).Exp(big.NewInt(62), big.NewInt(int64(size)), nil)

	if value.Cmp(max) >= 0 {
		return Token{}, ErrInvalidToken
	}

	return Token{
		value: new(big.Int).Set(value),
		size:  size,
	}, nil
}

func (token Token) Value() *big.Int {
	return new(big.Int).Set(token.value)
}

func (token Token) Size() int64 {
	return token.size
}
