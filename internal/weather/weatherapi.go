package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type WeatherAPIClient struct {
	httpClient *http.Client
	baseURL    string
	apiKey     string
}

func NewWeatherAPIClient(httpClient *http.Client, baseURL, apiKey string) *WeatherAPIClient {
	return &WeatherAPIClient{
		httpClient: httpClient,
		baseURL:    baseURL,
		apiKey:     apiKey,
	}
}

func (c *WeatherAPIClient) CurrentTempC(ctx context.Context, city string) (float64, error) {
	if c.apiKey == "" {
		return 0, fmt.Errorf("WEATHER_API_KEY is required")
	}

	endpoint, err := url.Parse(c.baseURL)
	if err != nil {
		return 0, err
	}
	query := endpoint.Query()
	query.Set("key", c.apiKey)
	query.Set("q", city)
	query.Set("aqi", "no")
	endpoint.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint.String(), nil)
	if err != nil {
		return 0, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode > 299 {
		return 0, fmt.Errorf("weather api returned status %d", res.StatusCode)
	}

	var payload struct {
		Current struct {
			TempC float64 `json:"temp_c"`
		} `json:"current"`
	}
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return 0, err
	}

	return payload.Current.TempC, nil
}
