package expiration

import "github.com/exanubes/url-shortener/internal/domain"

type Factory interface {
	Create(domain.PolicySettings) (domain.ExpirationPolicy, error)
}
