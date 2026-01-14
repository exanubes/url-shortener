package domain

import "net/url"

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
