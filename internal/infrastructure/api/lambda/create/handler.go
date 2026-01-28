package create

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	createshorturl "github.com/exanubes/url-shortener/internal/app/usecases/create_short_url"
	"github.com/exanubes/url-shortener/internal/infrastructure/api"
)

var invalid_encoding_response = events.APIGatewayV2HTTPResponse{
	StatusCode: http.StatusBadRequest, Body: "Invalid encoding",
}

var invalid_payload_response = events.APIGatewayV2HTTPResponse{
	StatusCode: http.StatusBadRequest, Body: "Invalid payload",
}

type Response struct {
	Message string `json:"message"`
}
type CreateUrlHandler struct {
	usecase createshorturl.UseCase
}

func NewHandler(usecase createshorturl.UseCase) *CreateUrlHandler {
	return &CreateUrlHandler{
		usecase: usecase,
	}
}

func (handler CreateUrlHandler) Handle(ctx context.Context, req events.APIGatewayV2HTTPRequest) events.APIGatewayV2HTTPResponse {
	var body = decode_request_body(req.IsBase64Encoded, req.Body)
	if body == nil {
		return invalid_encoding_response
	}
	var payload api.CreateUrlRequest

	if err := json.Unmarshal(body, &payload); err != nil {
		return invalid_payload_response
	}

	command, err := api.ToCreateLinkCommand(payload)
	if err != nil {
		return invalid_payload_response
	}

	link, err := handler.usecase.Execute(ctx, command)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest, Body: err.Error(),
		}
	}

	response, err := json.Marshal(api.CreateUrlResponse{
		ShortUrl: link.ShortCode().String(),
	})

	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError, Body: err.Error(),
		}
	}
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(response),
	}
}

func decode_request_body(is_encoded bool, body string) []byte {
	if is_encoded {
		decoded, err := base64.StdEncoding.DecodeString(body)
		if err != nil {
			return nil
		}
		return decoded

	} else {
		return []byte(body)
	}

}
