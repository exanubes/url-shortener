package resolve

import (
	"encoding/json"
	"net/http"

	"github.com/exanubes/url-shortener/internal/infrastructure/api"
)

func CreateErrorBody(status_code int, err_code, message string) string {
	response := api.ErrorResponse{
		Error:   http.StatusText(status_code),
		Message: message,
		Code:    err_code,
	}

	body, err := json.Marshal(response)

	if err != nil {
		return err.Error()
	}

	return string(body)
}
