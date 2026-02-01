package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	createshorturl "github.com/exanubes/url-shortener/internal/app/usecases/create_short_url"
	resolveurl "github.com/exanubes/url-shortener/internal/app/usecases/resolve_url"
	"github.com/exanubes/url-shortener/internal/infrastructure/api"
)

type HttpConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
	RequestTimeout  time.Duration
}

var DefaultConfig = HttpConfig{
	Port:            ":8000",
	ReadTimeout:     5 * time.Second,
	WriteTimeout:    5 * time.Second,
	IdleTimeout:     60 * time.Second,
	ShutdownTimeout: 30 * time.Second,
	RequestTimeout:  5 * time.Second,
}

type HttpDriver struct {
	create_url createshorturl.UseCase
	visit_url  resolveurl.UseCase
}

func NewHttpDriver(create_url createshorturl.UseCase, visit_url resolveurl.UseCase) *HttpDriver {
	return &HttpDriver{
		create_url: create_url,
		visit_url:  visit_url,
	}
}

func (driver *HttpDriver) Run(ctx context.Context, config HttpConfig) error {
	if config.Port == "" {
		config.Port = DefaultConfig.Port
	}

	server := &http.Server{
		Addr:         config.Port,
		Handler:      driver.setup_routes(config.RequestTimeout),
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		IdleTimeout:  config.IdleTimeout,
	}

	err_channel := make(chan error, 1)

	go func() {
		fmt.Printf("Starting HTTP server on http://localhost%s\n", config.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			err_channel <- err
		}
	}()

	system_signal_channel := make(chan os.Signal, 1)
	signal.Notify(system_signal_channel, syscall.SIGTERM, syscall.SIGINT)

	select {
	case err := <-err_channel:
		return err
	case signal := <-system_signal_channel:
		fmt.Printf("Received system signal to shut down: %s\n", signal.String())
	case <-ctx.Done():
	}

	shutdown_ctx, cancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdown_ctx); err != nil {
		return fmt.Errorf("Failed to shut down the server: %w\n", err)
	}

	fmt.Println("Server exited gracefully")
	return nil
}

func (driver *HttpDriver) setup_routes(request_timeout time.Duration) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("POST /", new_create_url_route(request_timeout, driver.create_url))
	mux.Handle("GET /{short_code}", new_resolve_url_route(request_timeout, driver.visit_url))
	return mux
}

func WriteError(response http.ResponseWriter, status_code int, err_code, message string) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status_code)
	json.NewEncoder(response).Encode(api.ErrorResponse{
		Error:   http.StatusText(status_code),
		Message: message,
		Code:    err_code,
	})
}
