package resolve

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	visitshorturl "github.com/exanubes/url-shortener/internal/app/usecases/visit_short_url"
	"github.com/exanubes/url-shortener/internal/domain"
)

type Response struct {
	Message string `json:"message"`
}
type ResolveUrlHandler struct {
	usecase visitshorturl.UseCase
}

func NewHandler(usecase visitshorturl.UseCase) *ResolveUrlHandler {
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

	url, err := handler.usecase.Execute(ctx, sc)

	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest, Body: err.Error(),
		}
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusTemporaryRedirect,
		Headers: map[string]string{
			"Location": url.String(),
		},
	}
}
