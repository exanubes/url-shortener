package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	createshorturl "github.com/exanubes/url-shortener/internal/app/usecases/create_short_url"
	"github.com/exanubes/url-shortener/internal/domain"
)

type CreateUrlRoute struct {
	usecase         createshorturl.UseCase
	request_timeout time.Duration
}

func NewCreateUrlRoute(request_timeout time.Duration, usecase createshorturl.UseCase) *CreateUrlRoute {
	return &CreateUrlRoute{usecase, request_timeout}
}

func (route *CreateUrlRoute) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(request.Context(), route.request_timeout)
	defer cancel()
	var payload CreateUrlRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		write_error(response, http.StatusBadRequest, "INVALID_PAYLOAD", err.Error())
		return
	}

	if payload.Url == "" {
		write_error(response, http.StatusBadRequest, "INVALID_PAYLOAD", "Url cannot be empty")
		return
	}

	url, err := domain.NewUrl(payload.Url)

	if err != nil {
		write_error(response, http.StatusBadRequest, "INVALID_PAYLOAD", err.Error())
		return
	}

	link, err := route.usecase.Execute(ctx, url)

	if err != nil {
		write_error(response, http.StatusInternalServerError, "", err.Error())
		return
	}

	response.Header().Set("Content-Type", "application/json")

	json.NewEncoder(response).Encode(CreateUrlResponse{
		ShortUrl: link.ShortCode().String(),
	})
}
