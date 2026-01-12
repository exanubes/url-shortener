package routes

import (
	"encoding/json"
	"net/http"

	"github.com/exanubes/url-shortener/internal/domain"
)

type CreateUrlRoute struct {
	usecase domain.ForCreatingUrls
}

func NewCreateUrlRoute(usecase domain.ForCreatingUrls) *CreateUrlRoute {
	return &CreateUrlRoute{
		usecase,
	}
}

func (route *CreateUrlRoute) ServeHTTP(response http.ResponseWriter, request *http.Request) {
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

	result, err := route.usecase.Execute(payload.Url)

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
