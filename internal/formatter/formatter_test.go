package formatter

import (
	"testing"

	"github.com/creimer/lnag/internal/data"
	"github.com/creimer/lnag/internal/matcher"
)

func TestHumanizeRatio(t *testing.T) {
	tests := []struct {
		ratio float64
		want  string
	}{
		{0.5, "half"},
		{0.25, "a quarter"},
		{3.0, "3x"},
		{10.0, "10x"},
		{2.7, "2.7x"},
		{1.0, ""},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := HumanizeRatio(tt.ratio)
			if got != tt.want {
				t.Errorf("HumanizeRatio(%f) = %q, want %q", tt.ratio, got, tt.want)
			}
		})
	}
}

func TestHumanizeCount(t *testing.T) {
	tests := []struct {
		count float64
		want  string
	}{
		{1000, "1,000"},
		{1234567, "1,234,567"},
		{42, "42"},
		{3.7, "4"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := HumanizeCount(tt.count)
			if got != tt.want {
				t.Errorf("HumanizeCount(%f) = %q, want %q", tt.count, got, tt.want)
			}
		})
	}
}

func pf(v float64) *float64 { return &v }

func TestFormatUnitResult(t *testing.T) {
	t.Run("length ratio 5", func(t *testing.T) {
		r := matcher.UnitResult{
			Concept:   data.Concept{Name: "Soccer Field", LengthM: pf(100)},
			Ratio:     5.0,
			Dimension: "length",
		}
		got := FormatUnitResult(r, 500, "m")
		want := "500 m is about the length of 5 Soccer Fields."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("length ratio 1", func(t *testing.T) {
		r := matcher.UnitResult{
			Concept:   data.Concept{Name: "Eiffel Tower", HeightM: pf(330)},
			Ratio:     1.0,
			Dimension: "height",
		}
		got := FormatUnitResult(r, 330, "m")
		want := "330 m is about the height of 1 Eiffel Tower."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("weight", func(t *testing.T) {
		r := matcher.UnitResult{
			Concept:   data.Concept{Name: "African Elephant", WeightKg: pf(5000)},
			Ratio:     10.0,
			Dimension: "weight",
		}
		got := FormatUnitResult(r, 50000, "kg")
		want := "50,000 kg is about the weight of 10 African Elephants."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("fractional ratio", func(t *testing.T) {
		r := matcher.UnitResult{
			Concept:   data.Concept{Name: "Blue Whale", WeightKg: pf(150000)},
			Ratio:     0.5,
			Dimension: "weight",
		}
		got := FormatUnitResult(r, 75000, "kg")
		want := "75,000 kg is about half the weight of a Blue Whale."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("duration ratio >1", func(t *testing.T) {
		r := matcher.UnitResult{
			Concept:   data.Concept{Name: "Average song length (pop, 2020s)", DurationS: pf(190)},
			Ratio:     2.0,
			Dimension: "duration",
		}
		got := FormatUnitResult(r, 380, "sec")
		want := "380 sec is about as long as 2 Average song length (pop, 2020s)s."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("duration ratio 1", func(t *testing.T) {
		r := matcher.UnitResult{
			Concept:   data.Concept{Name: "Marathon world record (men)", DurationS: pf(7084)},
			Ratio:     1.0,
			Dimension: "duration",
		}
		got := FormatUnitResult(r, 7084, "sec")
		want := "7,084 sec is about as long as 1 Marathon world record (men)."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("duration ratio <1", func(t *testing.T) {
		r := matcher.UnitResult{
			Concept:   data.Concept{Name: "Average night of sleep", DurationS: pf(28800)},
			Ratio:     0.5,
			Dimension: "duration",
		}
		got := FormatUnitResult(r, 14400, "sec")
		want := "14,400 sec is about half as long as an Average night of sleep."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("proper noun height ratio 5", func(t *testing.T) {
		r := matcher.UnitResult{
			Concept:   data.Concept{Name: "Eiffel Tower", HeightM: pf(330), ProperNoun: true},
			Ratio:     5.0,
			Dimension: "height",
		}
		got := FormatUnitResult(r, 1650, "m")
		want := "1,650 m is about 5x the height of the Eiffel Tower."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("proper noun height ratio 1", func(t *testing.T) {
		r := matcher.UnitResult{
			Concept:   data.Concept{Name: "Eiffel Tower", HeightM: pf(330), ProperNoun: true},
			Ratio:     1.0,
			Dimension: "height",
		}
		got := FormatUnitResult(r, 330, "m")
		want := "330 m is about the height of the Eiffel Tower."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("proper noun height ratio <1", func(t *testing.T) {
		r := matcher.UnitResult{
			Concept:   data.Concept{Name: "Eiffel Tower", HeightM: pf(330), ProperNoun: true},
			Ratio:     0.5,
			Dimension: "height",
		}
		got := FormatUnitResult(r, 165, "m")
		want := "165 m is about half the height of the Eiffel Tower."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("proper noun duration ratio >1", func(t *testing.T) {
		r := matcher.UnitResult{
			Concept:   data.Concept{Name: "Apollo 11 total mission", DurationS: pf(691200), ProperNoun: true},
			Ratio:     2.0,
			Dimension: "duration",
		}
		got := FormatUnitResult(r, 1382400, "sec")
		want := "1,382,400 sec is about 2x as long as the Apollo 11 total mission."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("proper noun duration ratio 1", func(t *testing.T) {
		r := matcher.UnitResult{
			Concept:   data.Concept{Name: "Apollo 11 total mission", DurationS: pf(691200), ProperNoun: true},
			Ratio:     1.0,
			Dimension: "duration",
		}
		got := FormatUnitResult(r, 691200, "sec")
		want := "691,200 sec is about as long as the Apollo 11 total mission."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("proper noun duration ratio <1", func(t *testing.T) {
		r := matcher.UnitResult{
			Concept:   data.Concept{Name: "Apollo 11 total mission", DurationS: pf(691200), ProperNoun: true},
			Ratio:     0.5,
			Dimension: "duration",
		}
		got := FormatUnitResult(r, 345600, "sec")
		want := "345,600 sec is about half as long as the Apollo 11 total mission."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}

func TestFormatDimensionResult(t *testing.T) {
	t.Run("weight comparison", func(t *testing.T) {
		r := matcher.DimensionResult{
			UnitItem:   data.Concept{Name: "Watermelon", WeightKg: pf(5)},
			TargetItem: data.Concept{Name: "African Elephant", WeightKg: pf(5000)},
			Count:      2000,
			Ratio:      2.0,
			Dimension:  "weight",
		}
		got := FormatDimensionResult(r)
		want := "2,000 Watermelons would weigh about as much as 2 African Elephants."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("height comparison", func(t *testing.T) {
		r := matcher.DimensionResult{
			UnitItem:   data.Concept{Name: "Human", HeightM: pf(1.7)},
			TargetItem: data.Concept{Name: "Eiffel Tower", HeightM: pf(330)},
			Count:      100,
			Ratio:      0.5,
			Dimension:  "height",
		}
		got := FormatDimensionResult(r)
		want := "100 Humans stacked would be about half the height of an Eiffel Tower."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("length comparison ratio 1", func(t *testing.T) {
		r := matcher.DimensionResult{
			UnitItem:   data.Concept{Name: "Car", LengthM: pf(4.5)},
			TargetItem: data.Concept{Name: "Soccer Field", LengthM: pf(100)},
			Count:      100,
			Ratio:      4.5,
			Dimension:  "length",
		}
		got := FormatDimensionResult(r)
		// 4.5 is not an exact nice number, so it formats as "4.5x"
		want := "100 Cars lined up would stretch about 4.5x the length of a Soccer Field."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("duration comparison", func(t *testing.T) {
		r := matcher.DimensionResult{
			UnitItem:   data.Concept{Name: "Eye blink", DurationS: pf(0.15)},
			TargetItem: data.Concept{Name: "Average song length (pop, 2020s)", DurationS: pf(190)},
			Count:      1000,
			Ratio:      0.5,
			Dimension:  "duration",
		}
		got := FormatDimensionResult(r)
		want := "1,000 Eye blinks would last about half as long as an Average song length (pop, 2020s)."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("duration comparison ratio >1", func(t *testing.T) {
		r := matcher.DimensionResult{
			UnitItem:   data.Concept{Name: "Average song length (pop, 2020s)", DurationS: pf(190)},
			TargetItem: data.Concept{Name: "Marathon world record (men)", DurationS: pf(7084)},
			Count:      100,
			Ratio:      2.0,
			Dimension:  "duration",
		}
		got := FormatDimensionResult(r)
		want := "100 Average song length (pop, 2020s)s would last about as long as 2 Marathon world record (men)s."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("proper noun target height ratio <1", func(t *testing.T) {
		r := matcher.DimensionResult{
			UnitItem:   data.Concept{Name: "Human", HeightM: pf(1.7)},
			TargetItem: data.Concept{Name: "Eiffel Tower", HeightM: pf(330), ProperNoun: true},
			Count:      100,
			Ratio:      0.5,
			Dimension:  "height",
		}
		got := FormatDimensionResult(r)
		want := "100 Humans stacked would be about half the height of the Eiffel Tower."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("proper noun target weight ratio >1", func(t *testing.T) {
		r := matcher.DimensionResult{
			UnitItem:   data.Concept{Name: "Watermelon", WeightKg: pf(5)},
			TargetItem: data.Concept{Name: "Statue of Liberty", WeightKg: pf(204000), ProperNoun: true},
			Count:      50000,
			Ratio:      2.0,
			Dimension:  "weight",
		}
		got := FormatDimensionResult(r)
		want := "50,000 Watermelons would weigh about 2x as much as the Statue of Liberty."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("proper noun target weight ratio 1", func(t *testing.T) {
		r := matcher.DimensionResult{
			UnitItem:   data.Concept{Name: "Watermelon", WeightKg: pf(5)},
			TargetItem: data.Concept{Name: "Statue of Liberty", WeightKg: pf(204000), ProperNoun: true},
			Count:      40800,
			Ratio:      1.0,
			Dimension:  "weight",
		}
		got := FormatDimensionResult(r)
		want := "40,800 Watermelons would weigh about as much as the Statue of Liberty."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("proper noun target duration ratio >1", func(t *testing.T) {
		r := matcher.DimensionResult{
			UnitItem:   data.Concept{Name: "Eye blink", DurationS: pf(0.15)},
			TargetItem: data.Concept{Name: "Apollo 11 total mission", DurationS: pf(691200), ProperNoun: true},
			Count:      10000000,
			Ratio:      2.0,
			Dimension:  "duration",
		}
		got := FormatDimensionResult(r)
		want := "10,000,000 Eye blinks would last about 2x as long as the Apollo 11 total mission."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
