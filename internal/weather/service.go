package weather

import (
	"context"
	"errors"
	"regexp"
)

var (
	ErrInvalidZipcode  = errors.New("invalid zipcode")
	ErrZipcodeNotFound = errors.New("can not find zipcode")

	zipcodePattern = regexp.MustCompile(`^\d{8}$`)
)

type ZipcodeResolver interface {
	ResolveCity(ctx context.Context, zipcode string) (string, error)
}

type TemperatureProvider interface {
	CurrentTempC(ctx context.Context, city string) (float64, error)
}

type Service struct {
	zipcodes ZipcodeResolver
	weather  TemperatureProvider
}

func NewService(zipcodes ZipcodeResolver, weather TemperatureProvider) *Service {
	return &Service{zipcodes: zipcodes, weather: weather}
}

func (s *Service) WeatherByZipcode(ctx context.Context, zipcode string) (TemperatureResponse, error) {
	if !zipcodePattern.MatchString(zipcode) {
		return TemperatureResponse{}, ErrInvalidZipcode
	}

	city, err := s.zipcodes.ResolveCity(ctx, zipcode)
	if err != nil {
		return TemperatureResponse{}, err
	}

	tempC, err := s.weather.CurrentTempC(ctx, city)
	if err != nil {
		return TemperatureResponse{}, err
	}

	return ConvertTemperature(tempC), nil
}
