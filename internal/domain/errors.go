package domain

import "errors"

var (
	ErrInvalidToken = errors.New("Invalid token")
	UrlNotFound     = errors.New("Url not found")
)
