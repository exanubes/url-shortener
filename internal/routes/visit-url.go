package routes

import (
	"net/http"

	"github.com/exanubes/url-shortener/internal/domain"
)

type VisitUrlRoute struct {
	usecase domain.ForVisitingUrls
}

func NewVisitUrlRoute(usecase domain.ForVisitingUrls) *VisitUrlRoute {
	return &VisitUrlRoute{
		usecase,
	}
}

func (route *VisitUrlRoute) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	short_url := request.PathValue("short_url")

	result, err := route.usecase.Execute(short_url)

	if err != nil {
		response.WriteHeader(500)
		return
	}

	http.Redirect(response, request, result, http.StatusMovedPermanently)
}
