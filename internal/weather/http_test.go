package weather

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type stubZipcodeResolver struct {
	city string
	err  error
}

func (s stubZipcodeResolver) ResolveCity(context.Context, string) (string, error) {
	return s.city, s.err
}

type stubTemperatureProvider struct {
	tempC float64
	err   error
}

func (s stubTemperatureProvider) CurrentTempC(context.Context, string) (float64, error) {
	return s.tempC, s.err
}

func TestHandleWeatherSuccess(t *testing.T) {
	service := NewService(
		stubZipcodeResolver{city: "Sao Paulo"},
		stubTemperatureProvider{tempC: 28.5},
	)
	req := httptest.NewRequest(http.MethodGet, "/weather/01001000", nil)
	rec := httptest.NewRecorder()

	HandleWeather(service).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var body TemperatureResponse
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.TempC != 28.5 || body.TempF != 83.3 || body.TempK != 301.65 {
		t.Fatalf("body = %+v, want converted temperatures", body)
	}
}

func TestHandleWeatherTemperaturesAlias(t *testing.T) {
	service := NewService(
		stubZipcodeResolver{city: "Sao Paulo"},
		stubTemperatureProvider{tempC: 28.5},
	)
	req := httptest.NewRequest(http.MethodGet, "/temperatures/01001000", nil)
	rec := httptest.NewRecorder()

	HandleWeather(service).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestHandleWeatherInvalidZipcode(t *testing.T) {
	service := NewService(stubZipcodeResolver{}, stubTemperatureProvider{})
	req := httptest.NewRequest(http.MethodGet, "/weather/abc", nil)
	rec := httptest.NewRecorder()

	HandleWeather(service).ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnprocessableEntity)
	}
	if rec.Body.String() != "invalid zipcode" {
		t.Fatalf("body = %q, want invalid zipcode", rec.Body.String())
	}
}

func TestHandleWeatherZipcodeNotFound(t *testing.T) {
	service := NewService(
		stubZipcodeResolver{err: ErrZipcodeNotFound},
		stubTemperatureProvider{},
	)
	req := httptest.NewRequest(http.MethodGet, "/weather/99999999", nil)
	rec := httptest.NewRecorder()

	HandleWeather(service).ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
	if rec.Body.String() != "can not find zipcode" {
		t.Fatalf("body = %q, want can not find zipcode", rec.Body.String())
	}
}

func TestHandleWeatherProviderError(t *testing.T) {
	service := NewService(
		stubZipcodeResolver{city: "Sao Paulo"},
		stubTemperatureProvider{err: errors.New("provider down")},
	)
	req := httptest.NewRequest(http.MethodGet, "/weather?zipcode=01001000", nil)
	rec := httptest.NewRecorder()

	HandleWeather(service).ServeHTTP(rec, req)

	if rec.Code != http.StatusBadGateway {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadGateway)
	}
}

func TestHandleHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
	rec := httptest.NewRecorder()

	HandleHealth(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}
