package units

import "fmt"

type UnitInfo struct {
	Dimension string
	ToBase    float64
}

var unitTable = map[string]UnitInfo{
	// length (base: meters)
	"m":      {"length", 1},
	"meters": {"length", 1},
	"km":     {"length", 1000},
	"ft":     {"length", 0.3048},
	"feet":   {"length", 0.3048},
	"mi":     {"length", 1609.344},
	"miles":  {"length", 1609.344},

	// weight (base: kg)
	"kg":   {"weight", 1},
	"g":    {"weight", 0.001},
	"lbs":  {"weight", 0.453592},
	"tons": {"weight", 907.185},

	// volume (base: m³)
	"m3":      {"volume", 1},
	"liters":  {"volume", 0.001},
	"gallons": {"volume", 0.003785},

	// area (base: m²)
	"m2":       {"area", 1},
	"acres":    {"area", 4046.86},
	"hectares": {"area", 10000},

	// distance (base: meters) — same physical unit as length but separate dimension
	"au": {"distance", 149597870700},
	"ly": {"distance", 9.4607e15},

	// duration (base: seconds)
	"s":       {"duration", 1},
	"sec":     {"duration", 1},
	"seconds": {"duration", 1},
	"min":     {"duration", 60},
	"minutes": {"duration", 60},
	"hr":      {"duration", 3600},
	"hours":   {"duration", 3600},
	"days":    {"duration", 86400},
	"years":   {"duration", 31557600},
}

func Resolve(unit string) (UnitInfo, error) {
	info, ok := unitTable[unit]
	if !ok {
		return UnitInfo{}, fmt.Errorf("unknown unit: %q", unit)
	}
	return info, nil
}

func Convert(value float64, unit string) (float64, string, error) {
	info, err := Resolve(unit)
	if err != nil {
		return 0, "", err
	}
	return value * info.ToBase, info.Dimension, nil
}

// ToMeters converts a value in the given unit to meters.
// Kept for backward compatibility; only works with length units.
func ToMeters(value float64, unit string) (float64, error) {
	info, err := Resolve(unit)
	if err != nil {
		return 0, err
	}
	if info.Dimension != "length" {
		return 0, fmt.Errorf("unit %q is not a length unit", unit)
	}
	return value * info.ToBase, nil
}
