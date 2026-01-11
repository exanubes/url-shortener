package drivers

import (
	"net/http"

	"github.com/exanubes/url-shortener/internal/app/encoder"
	"github.com/exanubes/url-shortener/internal/infrastructure/persistence/inmemory"
	"github.com/exanubes/url-shortener/internal/routes"
	"github.com/exanubes/url-shortener/internal/usecase"
)

type HttpDriver struct{}

func NewHttpDriver() *HttpDriver {
	return &HttpDriver{}
}

func (_ *HttpDriver) Run() error {
	mux := http.NewServeMux()
	provider := inmemory.NewInmemoryRepository()
	codec := encoder.New()
	create_short_url_use_case := usecase.NewCreateShortUrl(provider, codec)
	visit_url_use_case := usecase.NewVisitShortUrl(provider, codec)
	mux.Handle("POST /", routes.NewCreateUrlRoute(create_short_url_use_case))
	mux.Handle("GET /{short_url}", routes.NewVisitUrlRoute(visit_url_use_case))

	return http.ListenAndServe(":8000", mux)
}
