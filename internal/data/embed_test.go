package data

import "testing"

func TestLoadConcepts(t *testing.T) {
	concepts, err := LoadConcepts()
	if err != nil {
		t.Fatalf("LoadConcepts() error: %v", err)
	}
	if len(concepts) == 0 {
		t.Fatal("LoadConcepts() returned no concepts")
	}
	if len(concepts) < 100 {
		t.Errorf("expected at least 100 concepts, got %d", len(concepts))
	}
}

func TestLoadConceptsHasElephant(t *testing.T) {
	concepts, err := LoadConcepts()
	if err != nil {
		t.Fatalf("LoadConcepts() error: %v", err)
	}
	found := false
	for _, c := range concepts {
		if c.Name == "African Elephant" {
			found = true
			if c.WeightKg == nil || *c.WeightKg != 5000 {
				t.Errorf("African Elephant weight_kg = %v, want 5000", c.WeightKg)
			}
			if c.HeightM == nil || *c.HeightM != 3.3 {
				t.Errorf("African Elephant height_m = %v, want 3.3", c.HeightM)
			}
			if c.LengthM == nil || *c.LengthM != 6.0 {
				t.Errorf("African Elephant length_m = %v, want 6.0", c.LengthM)
			}
			if c.Category != "Animal" {
				t.Errorf("African Elephant category = %q, want Animal", c.Category)
			}
		}
	}
	if !found {
		t.Error("African Elephant not found in concepts")
	}
}

func TestConceptValueFor(t *testing.T) {
	weight := 5000.0
	height := 3.3
	duration := 60.0
	c := Concept{
		Name:      "African Elephant",
		WeightKg:  &weight,
		HeightM:   &height,
		DurationS: &duration,
	}

	tests := []struct {
		dimension string
		wantVal   float64
		wantOk    bool
	}{
		{"weight", 5000, true},
		{"height", 3.3, true},
		{"length", 0, false},
		{"volume", 0, false},
		{"duration", 60, true},
	}

	for _, tt := range tests {
		t.Run(tt.dimension, func(t *testing.T) {
			val, ok := c.ValueFor(tt.dimension)
			if ok != tt.wantOk {
				t.Errorf("ValueFor(%q) ok = %v, want %v", tt.dimension, ok, tt.wantOk)
			}
			if ok && val != tt.wantVal {
				t.Errorf("ValueFor(%q) = %f, want %f", tt.dimension, val, tt.wantVal)
			}
		})
	}
}

func TestConceptNullFields(t *testing.T) {
	concepts, err := LoadConcepts()
	if err != nil {
		t.Fatalf("LoadConcepts() error: %v", err)
	}
	// Blue Whale has null height_m
	for _, c := range concepts {
		if c.Name == "Blue Whale" {
			if c.HeightM != nil {
				t.Errorf("Blue Whale height_m should be nil, got %f", *c.HeightM)
			}
			if c.LengthM == nil {
				t.Error("Blue Whale length_m should not be nil")
			}
			return
		}
	}
	t.Error("Blue Whale not found")
}

func TestLoadConceptsHasDurations(t *testing.T) {
	concepts, err := LoadConcepts()
	if err != nil {
		t.Fatalf("LoadConcepts() error: %v", err)
	}
	found := false
	for _, c := range concepts {
		if c.Name == "Marathon world record (men)" {
			found = true
			if c.DurationS == nil {
				t.Error("Marathon world record duration_s should not be nil")
			} else if *c.DurationS != 7084 {
				t.Errorf("Marathon world record duration_s = %f, want 7084", *c.DurationS)
			}
			if c.Category != "Sports" {
				t.Errorf("Marathon world record category = %q, want Sports", c.Category)
			}
		}
	}
	if !found {
		t.Error("Marathon world record (men) not found in concepts")
	}
}

func TestLoadDurationsDoNotHaveOtherDimensions(t *testing.T) {
	durations, err := loadDurations()
	if err != nil {
		t.Fatalf("loadDurations() error: %v", err)
	}
	for _, c := range durations {
		if c.LengthM != nil || c.HeightM != nil || c.WeightKg != nil {
			t.Errorf("duration concept %q should not have physical dimensions", c.Name)
		}
		if c.DurationS == nil {
			t.Errorf("duration concept %q should have duration_s", c.Name)
		}
	}
}

func TestLoadMeasurementsDoNotHaveDuration(t *testing.T) {
	measurements, err := loadMeasurements()
	if err != nil {
		t.Fatalf("loadMeasurements() error: %v", err)
	}
	for _, c := range measurements {
		if c.DurationS != nil {
			t.Errorf("measurement concept %q should not have duration_s", c.Name)
		}
	}
}
