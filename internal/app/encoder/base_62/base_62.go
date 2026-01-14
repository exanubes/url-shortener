package base62

import (
	"math/big"
	"slices"

	"github.com/exanubes/url-shortener/internal/domain"
)

var characters_list = []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}

type Base62Encoder struct{}

func New() Base62Encoder {
	return Base62Encoder{}
}

func (Base62Encoder) Encode(token domain.Token) string {
	val := token.Value()

	if val.Cmp(big.NewInt(0)) == 0 {
		return "0"
	}
	var digits []rune
	current := val
	quotient := big.NewInt(62)
	for current.Cmp(big.NewInt(0)) != 0 {
		temp := new(big.Int)
		temp.Mod(current, quotient)

		digits = append(digits, characters_list[temp.Int64()])
		current = current.Div(current, quotient)
	}

	slices.Reverse(digits)

	return string(digits)
}
