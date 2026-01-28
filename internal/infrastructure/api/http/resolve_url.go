package http

import (
	"context"
	"net/http"
	"time"

	visitshorturl "github.com/exanubes/url-shortener/internal/app/usecases/visit_short_url"
	"github.com/exanubes/url-shortener/internal/domain"
)

type ResolveUrlRoute struct {
	usecase         visitshorturl.UseCase
	request_timeout time.Duration
}

func new_resolve_url_route(request_timeout time.Duration, usecase visitshorturl.UseCase) *ResolveUrlRoute {
	return &ResolveUrlRoute{usecase, request_timeout}
}

func (route *ResolveUrlRoute) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	ctx, cancel := context.WithTimeout(request.Context(), route.request_timeout)
	defer cancel()
	short_url, err := domain.NewShortCodeFromParam(request.PathValue("short_code"))

	if err != nil {
		WriteError(response, http.StatusBadRequest, "INVALID_SHORT_CODE", err.Error())
		return
	}

	result, err := route.usecase.Execute(ctx, short_url)

	if err != nil {
		if err == domain.ErrUrlNotFound {
			WriteError(response, http.StatusNotFound, "URL_NOT_FOUND", err.Error())
			return
		}

		if err == domain.ErrLinkExpired {
			WriteError(response, http.StatusGone, "LINK_EXPIRED", err.Error())
			return
		}

		WriteError(response, http.StatusInternalServerError, "SERVER_ERROR", err.Error())
		return
	}

	http.Redirect(response, request, result.String(), http.StatusTemporaryRedirect)
}
