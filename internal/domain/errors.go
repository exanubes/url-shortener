package domain

import "errors"

var (
	ErrShortCodeCollision        = errors.New("Shortcode already exists")
	ErrUrlNotFound               = errors.New("Url not found")
	ErrExceededMaxSize           = errors.New("Exceeded maximum size")
	ErrShortCodeEmpty            = errors.New("Empty short code")
	ErrInvalidShortCode          = errors.New("Invalid short code")
	ErrInvalidUrl                = errors.New("Invalid url")
	ErrInvalidScheme             = errors.New("Invalid scheme")
	ErrInvalidHost               = errors.New("Invalid host")
	ErrLinkExpired               = errors.New("Expired link")
	ErrExceededMinVisits         = errors.New("Exceeded minimum visits")
	ErrExceededMaxVisits         = errors.New("Exceeded maximum visits")
	ErrExceededMinAge            = errors.New("Exceeded min age")
	ErrExceededMaxAge            = errors.New("Exceeded max age")
	ErrUndefinedExpirationPolicy = errors.New("Undefined expiration policy")
)
