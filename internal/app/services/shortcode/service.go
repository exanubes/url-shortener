package shortcode

import "github.com/exanubes/url-shortener/internal/domain"

type ShortCodeService struct {
	tokens    TokenSpaceGenerator
	encoder   Encoder
	scrambler Scrambler
}

func NewService(generator TokenSpaceGenerator, scrambler Scrambler, encoder Encoder) *ShortCodeService {
	return &ShortCodeService{
		tokens:    generator,
		encoder:   encoder,
		scrambler: scrambler,
	}
}

func (generator *ShortCodeService) Generate() (domain.ShortCode, error) {
	token, err := generator.tokens.Generate()

	if err != nil {
		return domain.ShortCode{}, err
	}

	token, err = generator.scrambler.Scramble(token)

	if err != nil {
		return domain.ShortCode{}, err
	}

	short_url := generator.encoder.Encode(token)

	return domain.NewShortCode(short_url, int(token.Size()), "0")
}
