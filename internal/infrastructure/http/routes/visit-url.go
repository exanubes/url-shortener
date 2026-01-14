package routes

import (
	"context"
	"net/http"
	"time"

	visitshorturl "github.com/exanubes/url-shortener/internal/app/usecases/visit_short_url"
	"github.com/exanubes/url-shortener/internal/domain"
)

type VisitUrlRoute struct {
	usecase         visitshorturl.UseCase
	request_timeout time.Duration
}

func NewVisitUrlRoute(request_timeout time.Duration, usecase visitshorturl.UseCase) *VisitUrlRoute {
	return &VisitUrlRoute{usecase, request_timeout}
}

func (route *VisitUrlRoute) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(request.Context(), route.request_timeout)
	defer cancel()
	short_url, err := domain.NewShortCodeFromParam(request.PathValue("short_url"))

	if err != nil {
		write_error(response, http.StatusBadRequest, "INVALID_SHORT_CODE", err.Error())
		return
	}
	result, err := route.usecase.Execute(ctx, short_url)

	if err != nil {
		if err == domain.ErrUrlNotFound {
			write_error(response, http.StatusNotFound, "URL_NOT_FOUND", err.Error())
			return
		}

		write_error(response, http.StatusInternalServerError, "", err.Error())
		return
	}

	http.Redirect(response, request, result.String(), http.StatusMovedPermanently)
}
