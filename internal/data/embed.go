package data

import (
	_ "embed"
	"encoding/json"
)

//go:embed world_measurements.json
var measurementsJSON []byte

//go:embed world_durations.json
var durationsJSON []byte

type Concept struct {
	Name       string   `json:"name"`
	Category   string   `json:"category"`
	ProperNoun bool     `json:"proper_noun,omitempty"`
	LengthM    *float64 `json:"length_m,omitempty"`
	HeightM    *float64 `json:"height_m,omitempty"`
	WidthM     *float64 `json:"width_m,omitempty"`
	WeightKg   *float64 `json:"weight_kg,omitempty"`
	VolumeM3   *float64 `json:"volume_m3,omitempty"`
	AreaM2     *float64 `json:"area_m2,omitempty"`
	DistanceM  *float64 `json:"distance_m,omitempty"`
	DurationS  *float64 `json:"duration_s,omitempty"`
}

func (c Concept) ValueFor(dimension string) (float64, bool) {
	var p *float64
	switch dimension {
	case "length":
		p = c.LengthM
	case "height":
		p = c.HeightM
	case "width":
		p = c.WidthM
	case "weight":
		p = c.WeightKg
	case "volume":
		p = c.VolumeM3
	case "area":
		p = c.AreaM2
	case "distance":
		p = c.DistanceM
	case "duration":
		p = c.DurationS
	default:
		return 0, false
	}
	if p == nil {
		return 0, false
	}
	return *p, true
}

type rawDuration struct {
	Name       string  `json:"name"`
	Category   string  `json:"category"`
	ProperNoun bool    `json:"proper_noun,omitempty"`
	DurationS  float64 `json:"duration_s"`
}

func loadMeasurements() ([]Concept, error) {
	var concepts []Concept
	err := json.Unmarshal(measurementsJSON, &concepts)
	return concepts, err
}

func loadDurations() ([]Concept, error) {
	var raw []rawDuration
	if err := json.Unmarshal(durationsJSON, &raw); err != nil {
		return nil, err
	}
	concepts := make([]Concept, len(raw))
	for i, r := range raw {
		dur := r.DurationS
		concepts[i] = Concept{
			Name:       r.Name,
			Category:   r.Category,
			ProperNoun: r.ProperNoun,
			DurationS:  &dur,
		}
	}
	return concepts, nil
}

func LoadConcepts() ([]Concept, error) {
	measurements, err := loadMeasurements()
	if err != nil {
		return nil, err
	}
	durations, err := loadDurations()
	if err != nil {
		return nil, err
	}
	return append(measurements, durations...), nil
}
