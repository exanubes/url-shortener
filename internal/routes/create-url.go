package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/exanubes/url-shortener/internal/app/policy"
	"github.com/exanubes/url-shortener/internal/domain"
)

type CreateUrlRoute struct {
	usecase         domain.ForCreatingUrls
	request_timeout time.Duration
}

func NewCreateUrlRoute(request_timeout time.Duration, usecase domain.ForCreatingUrls) *CreateUrlRoute {
	return &CreateUrlRoute{usecase, request_timeout}
}

func (route *CreateUrlRoute) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(request.Context(), route.request_timeout)
	defer cancel()
	var payload PostRequestBody
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

	result, err := route.usecase.Execute(ctx, url, policy.NewRetryPolicy(3))

	if err != nil {
		write_error(response, http.StatusInternalServerError, "", err.Error())
		return
	}

	response.Header().Set("Content-Type", "application/json")

	json.NewEncoder(response).Encode(PostResponseBody{
		ShortUrl: result.String(),
	})
}

type PostRequestBody struct {
	Url string `json:"url"`
}

type PostResponseBody struct {
	ShortUrl string `json:"short_url"`
}
