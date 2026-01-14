package shortcode

import "github.com/exanubes/url-shortener/internal/domain"

type ShortCodeService struct {
	tokens  TokenSpaceGenerator
	encoder Encoder
}

func NewService(generator TokenSpaceGenerator, encoder Encoder) *ShortCodeService {
	return &ShortCodeService{
		tokens:  generator,
		encoder: encoder,
	}
}

func (generator *ShortCodeService) Generate() (domain.ShortCode, error) {
	token, err := generator.tokens.Generate()

	if err != nil {
		return domain.ShortCode{}, err
	}

	short_url := generator.encoder.Encode(token)

	return domain.NewShortCode(short_url, int(token.Size()), "0")
}
