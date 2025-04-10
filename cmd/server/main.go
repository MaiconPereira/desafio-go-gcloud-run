package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/esdras/goexpert-desafio-8-cloud-run/internal/weather"
)

func main() {
	cfg := config{
		port:           getenv("PORT", "8080"),
		viaCEPBaseURL:  getenv("VIACEP_BASE_URL", "https://viacep.com.br/ws"),
		weatherBaseURL: getenv("WEATHER_BASE_URL", "https://api.weatherapi.com/v1/current.json"),
		weatherAPIKey:  os.Getenv("WEATHER_API_KEY"),
	}

	httpClient := &http.Client{Timeout: 5 * time.Second}
	service := weather.NewService(
		weather.NewViaCEPClient(httpClient, cfg.viaCEPBaseURL),
		weather.NewWeatherAPIClient(
			httpClient,
			cfg.weatherBaseURL,
			cfg.weatherAPIKey,
		),
	)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", weather.HandleHealth)
	mux.HandleFunc("GET /health-check", weather.HandleHealth)
	mux.Handle("GET /weather/", weather.HandleWeather(service))
	mux.Handle("GET /weather", weather.HandleWeather(service))
	mux.Handle("GET /temperatures/", weather.HandleWeather(service))
	mux.Handle("GET /temperatures", weather.HandleWeather(service))

	server := &http.Server{
		Addr:         ":" + cfg.port,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	if err := run(server); err != nil {
		log.Fatal(err)
	}
}

type config struct {
	port           string
	viaCEPBaseURL  string
	weatherBaseURL string
	weatherAPIKey  string
}

func run(server *http.Server) error {
	stop, cleanup := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cleanup()

	go func() {
		log.Printf("server listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server failed: %v", err)
		}
	}()

	<-stop.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	log.Println("shutting down server")
	return server.Shutdown(shutdownCtx)
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
