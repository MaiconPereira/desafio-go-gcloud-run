package weather

import "math"

type TemperatureResponse struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func ConvertTemperature(celsius float64) TemperatureResponse {
	return TemperatureResponse{
		TempC: round2(celsius),
		TempF: round2(celsius*1.8 + 32),
		TempK: round2(celsius + 273.15),
	}
}

func round2(value float64) float64 {
	return math.Round(value*100) / 100
}
