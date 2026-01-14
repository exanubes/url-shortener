package routes

import (
	"context"
	"net/http"
	"time"

	"github.com/exanubes/url-shortener/internal/domain"
)

type VisitUrlRoute struct {
	usecase         domain.ForVisitingUrls
	request_timeout time.Duration
}

func NewVisitUrlRoute(request_timeout time.Duration, usecase domain.ForVisitingUrls) *VisitUrlRoute {
	return &VisitUrlRoute{usecase, request_timeout}
}

func (route *VisitUrlRoute) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(request.Context(), route.request_timeout)
	defer cancel()
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

	http.Redirect(response, request, result.String(), http.StatusMovedPermanently)
}
