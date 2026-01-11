package drivers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/exanubes/url-shortener/internal/domain"
	"github.com/exanubes/url-shortener/internal/routes"
)

type HttpConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

var DefaultHttpConfig = HttpConfig{
	Port:            ":8000",
	ReadTimeout:     5 * time.Second,
	WriteTimeout:    5 * time.Second,
	IdleTimeout:     60 * time.Second,
	ShutdownTimeout: 5 * time.Second,
}

type HttpDriver struct {
	create_url domain.ForCreatingUrls
	visit_url  domain.ForVisitingUrls
}

func NewHttpDriver(create_url domain.ForCreatingUrls, visit_url domain.ForVisitingUrls) *HttpDriver {
	return &HttpDriver{
		create_url: create_url,
		visit_url:  visit_url,
	}
}

func (driver *HttpDriver) Run(ctx context.Context, config HttpConfig) error {
	if config.Port == "" {
		config.Port = DefaultHttpConfig.Port
	}

	server := &http.Server{
		Addr:         config.Port,
		Handler:      driver.setup_routes(),
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
		return fmt.Errorf("Failed to shut down the server: %w", err)
	}

	fmt.Println("Server shutdown gracefully")
	return nil
}

func (driver *HttpDriver) setup_routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("POST /", routes.NewCreateUrlRoute(driver.create_url))
	mux.Handle("GET /{short_url}", routes.NewVisitUrlRoute(driver.visit_url))
	return mux
}
