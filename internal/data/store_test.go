package data

import (
	"math"
	"testing"
)

func TestNewConceptStore(t *testing.T) {
	store, err := NewConceptStore()
	if err != nil {
		t.Fatalf("NewConceptStore() error: %v", err)
	}
	if len(store.All) == 0 {
		t.Fatal("store.All is empty")
	}
	if len(store.ByDimension) == 0 {
		t.Fatal("store.ByDimension is empty")
	}
}

func TestConceptStoreHasAllDimensions(t *testing.T) {
	store, err := NewConceptStore()
	if err != nil {
		t.Fatalf("NewConceptStore() error: %v", err)
	}
	for _, dim := range []string{"length", "height", "width", "weight", "volume", "area", "duration"} {
		idx, ok := store.ByDimension[dim]
		if !ok {
			t.Errorf("missing dimension index for %q", dim)
			continue
		}
		if len(idx.Entries) == 0 {
			t.Errorf("dimension index for %q is empty", dim)
		}
	}
}

func TestDimensionIndexSorted(t *testing.T) {
	store, err := NewConceptStore()
	if err != nil {
		t.Fatalf("NewConceptStore() error: %v", err)
	}
	for dim, idx := range store.ByDimension {
		for i := 1; i < len(idx.Entries); i++ {
			if idx.Entries[i].Value < idx.Entries[i-1].Value {
				t.Errorf("dimension %q: entries not sorted at index %d: %f < %f",
					dim, i, idx.Entries[i].Value, idx.Entries[i-1].Value)
			}
		}
	}
}

func TestFindClosest(t *testing.T) {
	v1, v2, v3 := 1.0, 10.0, 100.0
	c1 := Concept{Name: "Small"}
	c2 := Concept{Name: "Medium"}
	c3 := Concept{Name: "Large"}

	idx := &DimensionIndex{
		Entries: []IndexEntry{
			{Concept: &c1, Value: v1},
			{Concept: &c2, Value: v2},
			{Concept: &c3, Value: v3},
		},
	}

	tests := []struct {
		target   float64
		wantName string
	}{
		{1.0, "Small"},
		{10.0, "Medium"},
		{100.0, "Large"},
		{5.0, "Small"},      // closer to 1 in abs diff (4 vs 5)
		{6.0, "Medium"},     // closer to 10 (4 vs 5)
		{50.0, "Medium"},    // closer to 10 in abs diff (40 vs 50)
		{55.0, "Medium"},    // equidistant to 10 and 100 (45 vs 45), picks lower
		{56.0, "Large"},     // closer to 100 (44 vs 46)
		{0.5, "Small"},      // below range
		{200.0, "Large"},    // above range
	}

	for _, tt := range tests {
		entry := idx.FindClosest(tt.target)
		if entry == nil {
			t.Errorf("FindClosest(%f) returned nil", tt.target)
			continue
		}
		if entry.Concept.Name != tt.wantName {
			t.Errorf("FindClosest(%f) = %q, want %q", tt.target, entry.Concept.Name, tt.wantName)
		}
	}
}

func TestFindClosestEmpty(t *testing.T) {
	idx := &DimensionIndex{Entries: nil}
	if entry := idx.FindClosest(5.0); entry != nil {
		t.Error("FindClosest on empty index should return nil")
	}
}

func TestFindClosestSingleEntry(t *testing.T) {
	c := Concept{Name: "Only"}
	idx := &DimensionIndex{
		Entries: []IndexEntry{{Concept: &c, Value: 42.0}},
	}
	entry := idx.FindClosest(100.0)
	if entry == nil || entry.Concept.Name != "Only" {
		t.Error("FindClosest on single-entry index should return that entry")
	}
}

func TestConceptStoreDurationIndex(t *testing.T) {
	store, err := NewConceptStore()
	if err != nil {
		t.Fatalf("NewConceptStore() error: %v", err)
	}
	idx := store.ByDimension["duration"]
	if idx == nil {
		t.Fatal("no duration index")
	}
	if len(idx.Entries) < 50 {
		t.Errorf("expected at least 50 duration entries, got %d", len(idx.Entries))
	}
	// Check that marathon is in the duration index
	found := false
	for _, e := range idx.Entries {
		if e.Concept.Name == "Marathon world record (men)" {
			found = true
			if math.Abs(e.Value-7084) > 0.01 {
				t.Errorf("Marathon value = %f, want 7084", e.Value)
			}
		}
	}
	if !found {
		t.Error("Marathon world record not found in duration index")
	}
}
