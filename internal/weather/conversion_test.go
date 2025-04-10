package weather

import "testing"

func TestConvertTemperature(t *testing.T) {
	got := ConvertTemperature(28.5)

	if got.TempC != 28.5 {
		t.Fatalf("TempC = %v, want 28.5", got.TempC)
	}
	if got.TempF != 83.3 {
		t.Fatalf("TempF = %v, want 83.3", got.TempF)
	}
	if got.TempK != 301.65 {
		t.Fatalf("TempK = %v, want 301.65", got.TempK)
	}
}
