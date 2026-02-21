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
		{1.0, ""},
		// Sub-1 non-integer unchanged
		{0.7, "0.7x"},
		// >1 non-integer: .1-.3 → "more than N times"
		{5.1, "more than 5 times"},
		{5.2, "more than 5 times"},
		{5.3, "more than 5 times"},
		// >1 non-integer: .4-.6 → "N and a half times"
		{5.4, "5 and a half times"},
		{5.5, "5 and a half times"},
		{5.6, "5 and a half times"},
		// >1 non-integer: .7-.9 → "almost N+1 times"
		{2.7, "almost 3 times"},
		{5.8, "almost 6 times"},
		{5.9, "almost 6 times"},
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

func TestApproxCount(t *testing.T) {
	tests := []struct {
		name  string
		ratio float64
		want  string
	}{
		{"integer", 5.0, "5"},
		{"large integer", 1000.0, "1,000"},
		{".1 decimal", 5.1, "more than 5"},
		{".2 decimal", 5.2, "more than 5"},
		{".3 decimal", 5.3, "more than 5"},
		{".4 decimal", 5.4, "5 and a half"},
		{".5 decimal", 5.5, "5 and a half"},
		{".6 decimal", 5.6, "5 and a half"},
		{".7 decimal", 5.7, "almost 6"},
		{".8 decimal", 5.8, "almost 6"},
		{".9 decimal", 5.9, "almost 6"},
		{"large approx", 1000.5, "1,000 and a half"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ApproxCount(tt.ratio)
			if got != tt.want {
				t.Errorf("ApproxCount(%f) = %q, want %q", tt.ratio, got, tt.want)
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
			Concept:   data.Concept{Name: "men's Marathon world record", DurationS: pf(7084)},
			Ratio:     1.0,
			Dimension: "duration",
		}
		got := FormatUnitResult(r, 7084, "sec")
		want := "7,084 sec is about as long as 1 men's Marathon world record."
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

	t.Run("non-integer ratio .2 uses more than", func(t *testing.T) {
		r := matcher.UnitResult{
			Concept:   data.Concept{Name: "Soccer Field", LengthM: pf(100)},
			Ratio:     5.2,
			Dimension: "length",
		}
		got := FormatUnitResult(r, 520, "m")
		want := "520 m is the length of more than 5 Soccer Fields."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("non-integer ratio .5 uses and a half", func(t *testing.T) {
		r := matcher.UnitResult{
			Concept:   data.Concept{Name: "Soccer Field", LengthM: pf(100)},
			Ratio:     5.5,
			Dimension: "length",
		}
		got := FormatUnitResult(r, 550, "m")
		want := "550 m is about the length of 5 and a half Soccer Fields."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("non-integer ratio .8 uses almost", func(t *testing.T) {
		r := matcher.UnitResult{
			Concept:   data.Concept{Name: "Soccer Field", LengthM: pf(100)},
			Ratio:     5.8,
			Dimension: "length",
		}
		got := FormatUnitResult(r, 580, "m")
		want := "580 m is the length of almost 6 Soccer Fields."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("proper noun non-integer ratio .5", func(t *testing.T) {
		r := matcher.UnitResult{
			Concept:   data.Concept{Name: "Eiffel Tower", HeightM: pf(330), ProperNoun: true},
			Ratio:     2.5,
			Dimension: "height",
		}
		got := FormatUnitResult(r, 825, "m")
		want := "825 m is about 2 and a half times the height of the Eiffel Tower."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("duration non-integer ratio .3", func(t *testing.T) {
		r := matcher.UnitResult{
			Concept:   data.Concept{Name: "Average song length (pop, 2020s)", DurationS: pf(190)},
			Ratio:     2.3,
			Dimension: "duration",
		}
		got := FormatUnitResult(r, 437, "sec")
		want := "437 sec is as long as more than 2 Average song length (pop, 2020s)s."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("distance ratio 5", func(t *testing.T) {
		r := matcher.UnitResult{
			Concept:   data.Concept{Name: "Mars", DistanceM: pf(227900000000)},
			Ratio:     5.0,
			Dimension: "distance",
		}
		got := FormatUnitResult(r, 1140000000, "km")
		want := "1,140,000,000 km is about 5x the distance to Mars."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("distance ratio 1", func(t *testing.T) {
		r := matcher.UnitResult{
			Concept:   data.Concept{Name: "Mars", DistanceM: pf(227900000000)},
			Ratio:     1.0,
			Dimension: "distance",
		}
		got := FormatUnitResult(r, 227900, "km")
		want := "227,900 km is about the distance to Mars."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("distance non-proper ratio <1", func(t *testing.T) {
		r := matcher.UnitResult{
			Concept:   data.Concept{Name: "Mars", DistanceM: pf(227900000000)},
			Ratio:     0.5,
			Dimension: "distance",
		}
		got := FormatUnitResult(r, 113950, "km")
		want := "113,950 km is about half the distance to Mars."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("distance proper noun ratio 5", func(t *testing.T) {
		r := matcher.UnitResult{
			Concept:   data.Concept{Name: "Sun", DistanceM: pf(149597870000), ProperNoun: true},
			Ratio:     5.0,
			Dimension: "distance",
		}
		got := FormatUnitResult(r, 748000000, "km")
		want := "748,000,000 km is about 5x the distance to the Sun."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("distance proper noun ratio <1", func(t *testing.T) {
		r := matcher.UnitResult{
			Concept:   data.Concept{Name: "Sun", DistanceM: pf(149597870000), ProperNoun: true},
			Ratio:     0.5,
			Dimension: "distance",
		}
		got := FormatUnitResult(r, 74800000, "km")
		want := "74,800,000 km is about half the distance to the Sun."
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

	t.Run("length comparison non-integer ratio", func(t *testing.T) {
		r := matcher.DimensionResult{
			UnitItem:   data.Concept{Name: "Car", LengthM: pf(4.5)},
			TargetItem: data.Concept{Name: "Soccer Field", LengthM: pf(100)},
			Count:      100,
			Ratio:      4.5,
			Dimension:  "length",
		}
		got := FormatDimensionResult(r)
		want := "100 Cars lined up would stretch about 4 and a half times the length of a Soccer Field."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("weight non-integer ratio .2", func(t *testing.T) {
		r := matcher.DimensionResult{
			UnitItem:   data.Concept{Name: "Watermelon", WeightKg: pf(5)},
			TargetItem: data.Concept{Name: "African Elephant", WeightKg: pf(5000)},
			Count:      2200,
			Ratio:      2.2,
			Dimension:  "weight",
		}
		got := FormatDimensionResult(r)
		want := "2,200 Watermelons would weigh as much as more than 2 African Elephants."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("weight non-integer ratio .8", func(t *testing.T) {
		r := matcher.DimensionResult{
			UnitItem:   data.Concept{Name: "Watermelon", WeightKg: pf(5)},
			TargetItem: data.Concept{Name: "African Elephant", WeightKg: pf(5000)},
			Count:      2800,
			Ratio:      2.8,
			Dimension:  "weight",
		}
		got := FormatDimensionResult(r)
		want := "2,800 Watermelons would weigh as much as almost 3 African Elephants."
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
			TargetItem: data.Concept{Name: "men's Marathon world record", DurationS: pf(7084)},
			Count:      100,
			Ratio:      2.0,
			Dimension:  "duration",
		}
		got := FormatDimensionResult(r)
		want := "100 Average song length (pop, 2020s)s would last about as long as 2 men's Marathon world records."
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

	t.Run("distance comparison", func(t *testing.T) {
		r := matcher.DimensionResult{
			UnitItem:   data.Concept{Name: "Earth", DistanceM: pf(149597870000)},
			TargetItem: data.Concept{Name: "Neptune", DistanceM: pf(4495000000000)},
			Count:      100,
			Ratio:      3.0,
			Dimension:  "distance",
		}
		got := FormatDimensionResult(r)
		want := "100x the distance to Earth would span about 3x the distance to Neptune."
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("distance proper target ratio 1", func(t *testing.T) {
		r := matcher.DimensionResult{
			UnitItem:   data.Concept{Name: "Mars", DistanceM: pf(227900000000)},
			TargetItem: data.Concept{Name: "Sun", DistanceM: pf(149597870000), ProperNoun: true},
			Count:      50,
			Ratio:      1.0,
			Dimension:  "distance",
		}
		got := FormatDimensionResult(r)
		want := "50x the distance to Mars would span about the distance to the Sun."
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
