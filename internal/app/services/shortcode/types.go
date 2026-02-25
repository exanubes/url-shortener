package shortcode

import "github.com/exanubes/url-shortener/internal/domain"

type Encoder interface {
	Encode(Token) string
}

type TokenSpaceGenerator interface {
	Generate() (Token, error)
}

type Service interface {
	Generate() (domain.ShortCode, error)
}

type Scrambler interface {
	Scramble(Token) (Token, error)
}
