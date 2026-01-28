package http

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	createshorturl "github.com/exanubes/url-shortener/internal/app/usecases/create_short_url"
	"github.com/exanubes/url-shortener/internal/infrastructure/api"
)

type CreateUrlRoute struct {
	usecase         createshorturl.UseCase
	request_timeout time.Duration
}

func new_create_url_route(request_timeout time.Duration, usecase createshorturl.UseCase) *CreateUrlRoute {
	return &CreateUrlRoute{usecase, request_timeout}
}

func (route *CreateUrlRoute) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(request.Context(), route.request_timeout)
	defer cancel()
	var payload api.CreateUrlRequest
	if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
		WriteError(response, http.StatusBadRequest, "INVALID_PAYLOAD", err.Error())
		return
	}
	command, err := api.ToCreateLinkCommand(payload)
	if err != nil {
		WriteError(response, http.StatusBadRequest, "INVALID_PAYLOAD", err.Error())
	}

	link, err := route.usecase.Execute(ctx, command)

	if err != nil {
		WriteError(response, http.StatusInternalServerError, "", err.Error())
		return
	}

	response.Header().Set("Content-Type", "application/json")

	json.NewEncoder(response).Encode(api.CreateUrlResponse{
		ShortUrl: link.ShortCode().String(),
	})
}
