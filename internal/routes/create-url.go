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
		http.Error(response, "invalid payload", http.StatusBadRequest)
		return
	}

	// TODO: url validator
	if payload.Url == "" {
		http.Error(response, "invalid payload", http.StatusBadRequest)
		return
	}

	result, err := route.usecase.Execute(ctx, payload.Url, policy.NewRetryPolicy(3))

	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	output := PostResponseBody{
		ShortUrl: result,
	}

	response.Header().Set("Content-Type", "application/json")

	json.NewEncoder(response).Encode(output)

}

type PostRequestBody struct {
	Url string `json:"url"`
}

type PostResponseBody struct {
	ShortUrl string `json:"short_url"`
}
