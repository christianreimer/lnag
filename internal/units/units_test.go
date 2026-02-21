package units

import (
	"math"
	"testing"
)

func TestResolve(t *testing.T) {
	tests := []struct {
		unit      string
		wantDim   string
		wantToBase float64
	}{
		{"m", "length", 1},
		{"meters", "length", 1},
		{"km", "length", 1000},
		{"ft", "length", 0.3048},
		{"feet", "length", 0.3048},
		{"mi", "length", 1609.344},
		{"miles", "length", 1609.344},
		{"kg", "weight", 1},
		{"g", "weight", 0.001},
		{"lbs", "weight", 0.453592},
		{"tons", "weight", 907.185},
		{"m3", "volume", 1},
		{"liters", "volume", 0.001},
		{"gallons", "volume", 0.003785},
		{"m2", "area", 1},
		{"acres", "area", 4046.86},
		{"hectares", "area", 10000},
		{"s", "duration", 1},
		{"sec", "duration", 1},
		{"seconds", "duration", 1},
		{"min", "duration", 60},
		{"minutes", "duration", 60},
		{"hr", "duration", 3600},
		{"hours", "duration", 3600},
		{"days", "duration", 86400},
		{"years", "duration", 31557600},
	}

	for _, tt := range tests {
		t.Run(tt.unit, func(t *testing.T) {
			info, err := Resolve(tt.unit)
			if err != nil {
				t.Fatalf("Resolve(%q) error: %v", tt.unit, err)
			}
			if info.Dimension != tt.wantDim {
				t.Errorf("Resolve(%q).Dimension = %q, want %q", tt.unit, info.Dimension, tt.wantDim)
			}
			if math.Abs(info.ToBase-tt.wantToBase) > 0.01 {
				t.Errorf("Resolve(%q).ToBase = %f, want %f", tt.unit, info.ToBase, tt.wantToBase)
			}
		})
	}
}

func TestResolveUnknownUnit(t *testing.T) {
	_, err := Resolve("cubits")
	if err == nil {
		t.Error("Resolve with unknown unit should return error")
	}
}

func TestConvert(t *testing.T) {
	tests := []struct {
		value   float64
		unit    string
		wantVal float64
		wantDim string
	}{
		{1, "km", 1000, "length"},
		{5, "m", 5, "length"},
		{1, "mi", 1609.344, "length"},
		{100, "kg", 100, "weight"},
		{1000, "g", 1, "weight"},
		{10, "liters", 0.01, "volume"},
		{1, "acres", 4046.86, "area"},
		{5, "min", 300, "duration"},
		{2, "hours", 7200, "duration"},
		{1, "days", 86400, "duration"},
		{1, "years", 31557600, "duration"},
	}

	for _, tt := range tests {
		t.Run(tt.unit, func(t *testing.T) {
			val, dim, err := Convert(tt.value, tt.unit)
			if err != nil {
				t.Fatalf("Convert(%f, %q) error: %v", tt.value, tt.unit, err)
			}
			if dim != tt.wantDim {
				t.Errorf("Convert(%f, %q) dimension = %q, want %q", tt.value, tt.unit, dim, tt.wantDim)
			}
			if math.Abs(val-tt.wantVal) > 0.01 {
				t.Errorf("Convert(%f, %q) = %f, want %f", tt.value, tt.unit, val, tt.wantVal)
			}
		})
	}
}

func TestConvertUnknownUnit(t *testing.T) {
	_, _, err := Convert(1, "cubits")
	if err == nil {
		t.Error("Convert with unknown unit should return error")
	}
}

func TestToMetersBackwardCompat(t *testing.T) {
	got, err := ToMeters(5, "km")
	if err != nil {
		t.Fatalf("ToMeters(5, km) error: %v", err)
	}
	if math.Abs(got-5000) > 0.001 {
		t.Errorf("ToMeters(5, km) = %f, want 5000", got)
	}
}
