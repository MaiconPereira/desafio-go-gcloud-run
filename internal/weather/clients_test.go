package weather

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestViaCEPClientResolveCity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/01001000/json/" {
			t.Fatalf("path = %q, want /01001000/json/", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"localidade":"Sao Paulo"}`))
	}))
	defer server.Close()

	client := NewViaCEPClient(server.Client(), server.URL)
	city, err := client.ResolveCity(context.Background(), "01001000")
	if err != nil {
		t.Fatal(err)
	}
	if city != "Sao Paulo" {
		t.Fatalf("city = %q, want Sao Paulo", city)
	}
}

func TestViaCEPClientNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"erro":"true"}`))
	}))
	defer server.Close()

	client := NewViaCEPClient(server.Client(), server.URL)
	_, err := client.ResolveCity(context.Background(), "99999999")
	if err != ErrZipcodeNotFound {
		t.Fatalf("err = %v, want ErrZipcodeNotFound", err)
	}
}

func TestWeatherAPIClientCurrentTempC(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("key") != "secret" {
			t.Fatalf("key query = %q, want secret", r.URL.Query().Get("key"))
		}
		if r.URL.Query().Get("q") != "Sao Paulo" {
			t.Fatalf("q query = %q, want Sao Paulo", r.URL.Query().Get("q"))
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"current":{"temp_c":21.7}}`))
	}))
	defer server.Close()

	client := NewWeatherAPIClient(server.Client(), server.URL, "secret")
	tempC, err := client.CurrentTempC(context.Background(), "Sao Paulo")
	if err != nil {
		t.Fatal(err)
	}
	if tempC != 21.7 {
		t.Fatalf("tempC = %v, want 21.7", tempC)
	}
}
