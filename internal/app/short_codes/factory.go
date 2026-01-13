package shortcodes

import "github.com/exanubes/url-shortener/internal/domain"

type ShortcodesGenerator struct {
	tokens  domain.TokenSpaceGenerator
	encoder domain.Encoder
}

func New(generator domain.TokenSpaceGenerator, encoder domain.Encoder) *ShortcodesGenerator {
	return &ShortcodesGenerator{
		tokens:  generator,
		encoder: encoder,
	}
}

func (generator *ShortcodesGenerator) Generate() (domain.ShortCode, error) {
	token, err := generator.tokens.Generate()

	if err != nil {
		return domain.ShortCode{}, err
	}

	short_url := generator.encoder.Encode(token)

	return domain.NewShortCode(short_url, int(token.Size()), "0")
}
