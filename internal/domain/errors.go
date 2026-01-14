package domain

import "errors"

var (
	ErrInvalidToken       = errors.New("Invalid token")
	ErrUrlNotFound        = errors.New("Url not found")
	ErrShortCodeCollision = errors.New("Shortcode already exists")
	ErrInvalidShortCode   = errors.New("Invalid short code")
	ErrExceededMaxSize    = errors.New("Exceeded maximum size")
	ErrShortCodeEmpty     = errors.New("Empty short code")
	ErrInvalidUrl         = errors.New("Invalid url")
	ErrInvalidScheme      = errors.New("Invalid scheme")
	ErrInvalidHost        = errors.New("Invalid host")
)
