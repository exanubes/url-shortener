package domain

import (
	"math/big"
	"net/url"

	"github.com/exanubes/url-shortener/internal/helpers"
)

type Url struct {
	parsed *url.URL
	value  string
}

func NewUrl(input string) (Url, error) {
	result, err := url.Parse(input)

	if err != nil {
		return Url{}, ErrInvalidUrl
	}

	if result.Scheme != "http" && result.Scheme != "https" {
		return Url{}, ErrInvalidScheme
	}

	if result.Host == "" {
		return Url{}, ErrInvalidHost
	}

	return Url{
		value:  input,
		parsed: result,
	}, nil
}

func (url Url) String() string {
	return url.value
}

type ShortCode struct {
	size  int
	value string
	zero  string
}

func NewShortCodeFromParam(value string) (ShortCode, error) {
	if len(value) > 11 {
		return ShortCode{}, ErrInvalidShortCode
	}
	if !helpers.IsBase62(value) {
		return ShortCode{}, ErrInvalidShortCode
	}
	return ShortCode{
		value: value,
		size:  len(value),
		zero:  "0",
	}, nil
}

func NewShortCode(value string, size int, zero_char string) (ShortCode, error) {
	if size <= MIN_SHORT_CODE_SIZE {
		return ShortCode{}, ErrShortCodeEmpty
	}

	if size > MAX_SHORT_CODE_SIZE {
		return ShortCode{}, ErrExceededMaxSize
	}

	if len(value) > int(size) {
		return ShortCode{}, ErrInvalidShortCode
	}
	return ShortCode{
		value: value,
		size:  size,
		zero:  zero_char,
	}, nil
}

func (code ShortCode) String() string {
	return helpers.PadStart(code.value, code.size, code.zero)
}

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
