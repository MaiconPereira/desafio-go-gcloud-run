package weather

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func HandleHealth(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func HandleWeather(service *Service) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		zipcode := zipcodeFromRequest(r)

		response, err := service.WeatherByZipcode(r.Context(), zipcode)
		if err != nil {
			writeError(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	})
}

func zipcodeFromRequest(r *http.Request) string {
	for _, prefix := range []string{"/weather/", "/temperatures/"} {
		if zipcode := strings.TrimPrefix(r.URL.Path, prefix); zipcode != r.URL.Path {
			return zipcode
		}
	}

	return r.URL.Query().Get("zipcode")
}

func writeError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	switch {
	case errors.Is(err, ErrInvalidZipcode):
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte("invalid zipcode"))
	case errors.Is(err, ErrZipcodeNotFound):
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("can not find zipcode"))
	default:
		w.WriteHeader(http.StatusBadGateway)
		_, _ = w.Write([]byte("weather service unavailable"))
	}
}
