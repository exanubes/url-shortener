package resolve

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	resolveurl "github.com/exanubes/url-shortener/internal/app/usecases/resolve_url"
	"github.com/exanubes/url-shortener/internal/domain"
)

var CACHE_MAX_AGE = 86400 * time.Second // 1 day

type Response struct {
	Message string `json:"message"`
}
type ResolveUrlHandler struct {
	usecase resolveurl.UseCase
}

func NewHandler(usecase resolveurl.UseCase) *ResolveUrlHandler {
	return &ResolveUrlHandler{
		usecase: usecase,
	}
}

func (handler ResolveUrlHandler) Handle(ctx context.Context, req events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse {
	short_code, exists := req.PathParameters["short_code"]

	if !exists {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       CreateErrorBody(http.StatusBadRequest, "INVALID_SHORT_CODE", "short code not included in path"),
		}
	}

	sc, err := domain.NewShortCodeFromParam(short_code)

	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       CreateErrorBody(http.StatusBadRequest, "INVALID_SHORT_CODE", err.Error()),
		}
	}

	output, err := handler.usecase.Execute(ctx, sc)

	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest, Body: err.Error(),
		}
	}

	var cache_control string

	if output.Status.Consumed || output.Status.ExpiresAt.IsZero() {
		cache_control = "no-store"
	} else {
		now := time.Now()
		duration := output.Status.ExpiresAt.Sub(now)
		if duration <= 0 {
			cache_control = "no-store"
		} else {
			max_age := min(output.Status.ExpiresAt.Sub(now), CACHE_MAX_AGE).Seconds()
			cache_control = fmt.Sprintf(
				"public, max-age=0, s-maxage=%d",
				int(max_age),
			)
		}
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusTemporaryRedirect,
		Headers: map[string]string{
			"Location":      output.Url.String(),
			"Cache-Control": cache_control,
		},
	}
}
