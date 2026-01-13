package routes

import (
	"context"
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
	ctx := context.Background()
	short_url := request.PathValue("short_url")
	result, err := route.usecase.Execute(ctx, short_url)

	if err != nil {
		if err == domain.ErrUrlNotFound {
			http.Error(response, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(response, request, result, http.StatusMovedPermanently)
}
