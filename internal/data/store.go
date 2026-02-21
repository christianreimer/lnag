package data

import "sort"

type IndexEntry struct {
	Concept *Concept
	Value   float64
}

type DimensionIndex struct {
	Entries []IndexEntry // sorted ascending by Value
}

// FindClosest uses binary search to find the entry nearest to the target value.
func (di *DimensionIndex) FindClosest(target float64) *IndexEntry {
	n := len(di.Entries)
	if n == 0 {
		return nil
	}
	i := sort.Search(n, func(j int) bool {
		return di.Entries[j].Value >= target
	})
	switch {
	case i == 0:
		return &di.Entries[0]
	case i == n:
		return &di.Entries[n-1]
	default:
		lo, hi := di.Entries[i-1], di.Entries[i]
		if target-lo.Value <= hi.Value-target {
			return &di.Entries[i-1]
		}
		return &di.Entries[i]
	}
}

type ConceptStore struct {
	All         []Concept
	ByDimension map[string]*DimensionIndex
}

var dimensions = []string{"length", "height", "width", "weight", "volume", "area", "duration"}

func NewConceptStore() (*ConceptStore, error) {
	measurements, err := loadMeasurements()
	if err != nil {
		return nil, err
	}
	durations, err := loadDurations()
	if err != nil {
		return nil, err
	}

	all := append(measurements, durations...)
	byDim := make(map[string]*DimensionIndex, len(dimensions))

	for _, dim := range dimensions {
		var entries []IndexEntry
		for i := range all {
			v, ok := all[i].ValueFor(dim)
			if !ok || v == 0 {
				continue
			}
			entries = append(entries, IndexEntry{Concept: &all[i], Value: v})
		}
		sort.Slice(entries, func(a, b int) bool {
			return entries[a].Value < entries[b].Value
		})
		byDim[dim] = &DimensionIndex{Entries: entries}
	}

	return &ConceptStore{All: all, ByDimension: byDim}, nil
}
