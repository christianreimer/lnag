package matcher

import (
	"fmt"
	"math"
	"math/rand/v2"

	"github.com/creimer/lnag/internal/data"
)

const scoreThreshold = 0.1

var niceNumbers = []float64{0.5, 1, 2, 3, 5, 10, 20, 50, 100, 500, 1000}

// ScoreRatio returns how "far" a ratio is from the nearest nice number.
// Lower is better.
func ScoreRatio(ratio float64) float64 {
	best := math.MaxFloat64
	for _, n := range niceNumbers {
		dist := math.Abs(math.Log10(ratio / n))
		if dist < best {
			best = dist
		}
	}
	return best
}

type UnitResult struct {
	Concept   data.Concept
	Ratio     float64
	Dimension string
}

// compatibleDimensions returns the set of dimensions to search.
// Length and distance share the same base unit (meters), so a length query
// also searches distance concepts and vice versa.
func compatibleDimensions(dimension string) []string {
	switch dimension {
	case "length":
		return []string{"length", "distance"}
	case "distance":
		return []string{"distance", "length"}
	default:
		return []string{dimension}
	}
}

// FindUnitMatch finds the concept whose measurement in the given dimension
// produces the nicest ratio with the input value.
func FindUnitMatch(value float64, dimension string, store *data.ConceptStore) (UnitResult, error) {
	type candidate struct {
		result UnitResult
		score  float64
	}
	var candidates []candidate

	for _, dim := range compatibleDimensions(dimension) {
		idx, ok := store.ByDimension[dim]
		if !ok || len(idx.Entries) == 0 {
			continue
		}
		for _, e := range idx.Entries {
			ratio := value / e.Value
			if ratio < 0.01 || ratio > 100000 {
				continue
			}
			score := ScoreRatio(ratio)
			candidates = append(candidates, candidate{
				result: UnitResult{
					Concept:   *e.Concept,
					Ratio:     ratio,
					Dimension: dim,
				},
				score: score,
			})
		}
	}

	if len(candidates) == 0 {
		dims := compatibleDimensions(dimension)
		return UnitResult{}, fmt.Errorf("no valid comparison found for %v", dims)
	}

	bestScore := math.MaxFloat64
	for _, c := range candidates {
		if c.score < bestScore {
			bestScore = c.score
		}
	}

	var topCandidates []candidate
	for _, c := range candidates {
		if c.score <= bestScore+scoreThreshold {
			topCandidates = append(topCandidates, c)
		}
	}

	pick := topCandidates[rand.IntN(len(topCandidates))]
	return pick.result, nil
}

type DimensionResult struct {
	UnitItem   data.Concept
	TargetItem data.Concept
	Count      float64
	Ratio      float64
	Dimension  string
}

// FindDimensionMatch finds a (unitItem, targetItem) pair such that
// count * unitItem.value / targetItem.value is close to a nice number.
// Uses binary search per nice number for O(n·k·log n) complexity.
func FindDimensionMatch(count float64, dimension string, store *data.ConceptStore) (DimensionResult, error) {
	idx, ok := store.ByDimension[dimension]
	if !ok || len(idx.Entries) < 2 {
		return DimensionResult{}, fmt.Errorf("not enough concepts for dimension %q", dimension)
	}

	type candidate struct {
		result DimensionResult
		score  float64
	}
	var candidates []candidate

	for _, unitEntry := range idx.Entries {
		totalValue := count * unitEntry.Value
		for _, nice := range niceNumbers {
			idealTarget := totalValue / nice
			closest := idx.FindClosest(idealTarget)
			if closest == nil || closest.Concept == unitEntry.Concept || unitEntry.Value >= closest.Value {
				continue
			}
			ratio := totalValue / closest.Value
			if ratio < 0.01 || ratio > 100000 {
				continue
			}
			score := ScoreRatio(ratio)
			candidates = append(candidates, candidate{
				result: DimensionResult{
					UnitItem:   *unitEntry.Concept,
					TargetItem: *closest.Concept,
					Count:      count,
					Ratio:      ratio,
					Dimension:  dimension,
				},
				score: score,
			})
		}
	}

	if len(candidates) == 0 {
		return DimensionResult{}, fmt.Errorf("no valid comparison found for dimension %q", dimension)
	}

	bestScore := math.MaxFloat64
	for _, c := range candidates {
		if c.score < bestScore {
			bestScore = c.score
		}
	}

	var topCandidates []candidate
	for _, c := range candidates {
		if c.score <= bestScore+scoreThreshold {
			topCandidates = append(topCandidates, c)
		}
	}

	pick := topCandidates[rand.IntN(len(topCandidates))]
	return pick.result, nil
}
