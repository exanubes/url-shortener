package dto

import (
	"encoding/json"
	"net/http"
)

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

func WriteError(response http.ResponseWriter, status_code int, err_code, message string) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status_code)
	json.NewEncoder(response).Encode(ErrorResponse{
		Error:   http.StatusText(status_code),
		Message: message,
		Code:    err_code,
	})
}
