package matcher

import (
	"testing"

	"github.com/creimer/lnag/internal/data"
)

func TestScoreRatio(t *testing.T) {
	tests := []struct {
		ratio    float64
		wantBest bool // should score better than a weird ratio
	}{
		{1.0, true},
		{0.5, true},
		{10.0, true},
		{7.3, false},
	}

	niceScore := ScoreRatio(1.0)
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			score := ScoreRatio(tt.ratio)
			if tt.wantBest && score > niceScore {
				t.Errorf("ScoreRatio(%f) = %f, expected <= %f (nicer ratios score lower)", tt.ratio, score, niceScore)
			}
		})
	}
}

func TestScoreRatioNiceNumbersBeatOddOnes(t *testing.T) {
	if ScoreRatio(3.0) >= ScoreRatio(7.3) {
		t.Error("expected ScoreRatio(3.0) < ScoreRatio(7.3)")
	}
}

func pf(v float64) *float64 { return &v }

func makeStore(concepts []data.Concept) *data.ConceptStore {
	dims := []string{"length", "height", "width", "weight", "volume", "area", "duration"}
	byDim := make(map[string]*data.DimensionIndex, len(dims))
	for _, dim := range dims {
		var entries []data.IndexEntry
		for i := range concepts {
			v, ok := concepts[i].ValueFor(dim)
			if !ok || v == 0 {
				continue
			}
			entries = append(entries, data.IndexEntry{Concept: &concepts[i], Value: v})
		}
		byDim[dim] = &data.DimensionIndex{Entries: entries}
	}
	return &data.ConceptStore{All: concepts, ByDimension: byDim}
}

func TestFindUnitMatch(t *testing.T) {
	concepts := []data.Concept{
		{Name: "Soccer Field", LengthM: pf(100)},
		{Name: "Eiffel Tower", HeightM: pf(330)},
		{Name: "African Elephant", WeightKg: pf(5000)},
	}
	store := makeStore(concepts)

	result, err := FindUnitMatch(500, "length", store)
	if err != nil {
		t.Fatalf("FindUnitMatch() error: %v", err)
	}
	if result.Concept.LengthM == nil {
		t.Error("expected concept with length dimension")
	}
	if result.Dimension != "length" {
		t.Errorf("dimension = %q, want length", result.Dimension)
	}
	if result.Ratio < 0.01 || result.Ratio > 100000 {
		t.Errorf("ratio = %f, out of valid range", result.Ratio)
	}
	if ScoreRatio(result.Ratio) > scoreThreshold+0.01 {
		t.Errorf("ratio %f has poor score %f", result.Ratio, ScoreRatio(result.Ratio))
	}
}

func TestFindUnitMatchWeight(t *testing.T) {
	concepts := []data.Concept{
		{Name: "Soccer Field", LengthM: pf(100)},
		{Name: "African Elephant", WeightKg: pf(5000)},
	}
	store := makeStore(concepts)

	result, err := FindUnitMatch(5000, "weight", store)
	if err != nil {
		t.Fatalf("FindUnitMatch() error: %v", err)
	}
	if result.Concept.WeightKg == nil {
		t.Error("expected concept with weight dimension")
	}
	if ScoreRatio(result.Ratio) > scoreThreshold+0.01 {
		t.Errorf("ratio %f has poor score %f", result.Ratio, ScoreRatio(result.Ratio))
	}
}

func TestFindUnitMatchNoConcepts(t *testing.T) {
	concepts := []data.Concept{
		{Name: "Soccer Field", LengthM: pf(100)},
	}
	store := makeStore(concepts)
	_, err := FindUnitMatch(500, "weight", store)
	if err == nil {
		t.Error("expected error when no concepts have the requested dimension")
	}
}

func TestFindUnitMatchDuration(t *testing.T) {
	concepts := []data.Concept{
		{Name: "Marathon world record (men)", DurationS: pf(7084)},
		{Name: "Average song length (pop, 2020s)", DurationS: pf(190)},
	}
	store := makeStore(concepts)

	result, err := FindUnitMatch(380, "duration", store)
	if err != nil {
		t.Fatalf("FindUnitMatch() error: %v", err)
	}
	if result.Concept.DurationS == nil {
		t.Error("expected concept with duration dimension")
	}
	if ScoreRatio(result.Ratio) > scoreThreshold+0.01 {
		t.Errorf("ratio %f has poor score %f", result.Ratio, ScoreRatio(result.Ratio))
	}
}

func TestFindDimensionMatch(t *testing.T) {
	concepts := []data.Concept{
		{Name: "Watermelon", WeightKg: pf(5)},
		{Name: "African Elephant", WeightKg: pf(5000)},
		{Name: "Blue Whale", WeightKg: pf(150000)},
	}
	store := makeStore(concepts)

	result, err := FindDimensionMatch(2000, "weight", store)
	if err != nil {
		t.Fatalf("FindDimensionMatch() error: %v", err)
	}
	if result.Dimension != "weight" {
		t.Errorf("dimension = %q, want weight", result.Dimension)
	}
	if result.UnitItem.WeightKg == nil {
		t.Error("unit item should have weight dimension")
	}
	if result.TargetItem.WeightKg == nil {
		t.Error("target item should have weight dimension")
	}
	if *result.UnitItem.WeightKg >= *result.TargetItem.WeightKg {
		t.Error("unit item should be smaller than target item")
	}
	if ScoreRatio(result.Ratio) > scoreThreshold+0.01 {
		t.Errorf("ratio %f has poor score %f", result.Ratio, ScoreRatio(result.Ratio))
	}
	if result.Count != 2000 {
		t.Errorf("count = %f, want 2000", result.Count)
	}
}

func TestFindDimensionMatchNoConcepts(t *testing.T) {
	concepts := []data.Concept{
		{Name: "Soccer Field", LengthM: pf(100)},
	}
	store := makeStore(concepts)
	_, err := FindDimensionMatch(100, "weight", store)
	if err == nil {
		t.Error("expected error when no concepts have the requested dimension")
	}
}

func TestFindDimensionMatchNeedsTwoConcepts(t *testing.T) {
	concepts := []data.Concept{
		{Name: "Watermelon", WeightKg: pf(5)},
	}
	store := makeStore(concepts)
	_, err := FindDimensionMatch(100, "weight", store)
	if err == nil {
		t.Error("expected error when only one concept available")
	}
}

func TestFindDimensionMatchDuration(t *testing.T) {
	concepts := []data.Concept{
		{Name: "Eye blink", DurationS: pf(0.15)},
		{Name: "Average song length (pop, 2020s)", DurationS: pf(190)},
		{Name: "Marathon world record (men)", DurationS: pf(7084)},
	}
	store := makeStore(concepts)

	result, err := FindDimensionMatch(1000, "duration", store)
	if err != nil {
		t.Fatalf("FindDimensionMatch() error: %v", err)
	}
	if result.Dimension != "duration" {
		t.Errorf("dimension = %q, want duration", result.Dimension)
	}
	if result.UnitItem.DurationS == nil {
		t.Error("unit item should have duration")
	}
	if result.TargetItem.DurationS == nil {
		t.Error("target item should have duration")
	}
}
