package api

type ExpirationUnit string

type ExpirationTimeDefinition struct {
	Value int            `json:"value"`
	Unit  ExpirationUnit `json:"unit"`
}

type CreateUrlRequest struct {
	Url          string                   `json:"url"`
	OneTimeLink  bool                     `json:"one_time_link"`
	ExpiresAfter ExpirationTimeDefinition `json:"expires_after"`
}

type CreateUrlResponse struct {
	ShortUrl string `json:"short_url"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}
