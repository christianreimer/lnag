# Concept Library & CLI Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a CLI tool that takes a number + unit and outputs a human-readable comparison using a two-tier concept library (items + distances).

**Architecture:** Embedded JSON data files define items (small countable things) and distances (large references). A matcher scores all (item, distance) pairings by ratio quality and picks the best. A formatter turns the result into a readable sentence. The CLI wires it together.

**Tech Stack:** Go 1.24, `go:embed`, stdlib only (no external deps).

---

### Task 1: Project Scaffolding

**Files:**
- Create: `go.mod`
- Create: `internal/data/items.json`
- Create: `internal/data/distances.json`

**Step 1: Initialize go module**

Run: `go mod init github.com/creimer/numberviz`
Expected: `go.mod` created

**Step 2: Create items.json**

Create `internal/data/items.json`:
```json
[
  {"name": "iPhone", "length_m": 0.1461, "plural": "iPhones"},
  {"name": "penny", "length_m": 0.01905, "plural": "pennies"},
  {"name": "credit card", "length_m": 0.0856, "plural": "credit cards"},
  {"name": "baseball bat", "length_m": 1.067, "plural": "baseball bats"},
  {"name": "school bus", "length_m": 10.67, "plural": "school buses"},
  {"name": "football field", "length_m": 91.44, "plural": "football fields"},
  {"name": "Olympic swimming pool", "length_m": 50.0, "plural": "Olympic swimming pools"}
]
```

**Step 3: Create distances.json**

Create `internal/data/distances.json`:
```json
[
  {"name": "height of the Statue of Liberty", "distance_m": 93},
  {"name": "height of the Eiffel Tower", "distance_m": 330},
  {"name": "height of Mount Everest", "distance_m": 8849},
  {"name": "length of the Grand Canyon", "distance_m": 446000},
  {"name": "distance from New York to Los Angeles", "distance_m": 3944000},
  {"name": "diameter of the Earth", "distance_m": 12742000},
  {"name": "distance from the Earth to the Moon", "distance_m": 384400000},
  {"name": "distance from the Earth to the Sun", "distance_m": 149597870700}
]
```

**Step 4: Commit**

```bash
git add go.mod internal/data/items.json internal/data/distances.json
git commit -m "Initialize go module and add concept data files"
```

---

### Task 2: Data Loading Package

**Files:**
- Create: `internal/data/embed.go`
- Create: `internal/data/embed_test.go`

**Step 1: Write the failing tests**

Create `internal/data/embed_test.go`:
```go
package data

import "testing"

func TestLoadItems(t *testing.T) {
	items, err := LoadItems()
	if err != nil {
		t.Fatalf("LoadItems() error: %v", err)
	}
	if len(items) == 0 {
		t.Fatal("LoadItems() returned no items")
	}
	// Spot-check a known item
	found := false
	for _, item := range items {
		if item.Name == "iPhone" {
			found = true
			if item.LengthM != 0.1461 {
				t.Errorf("iPhone length = %f, want 0.1461", item.LengthM)
			}
			if item.Plural != "iPhones" {
				t.Errorf("iPhone plural = %q, want \"iPhones\"", item.Plural)
			}
		}
	}
	if !found {
		t.Error("iPhone not found in items")
	}
}

func TestLoadDistances(t *testing.T) {
	distances, err := LoadDistances()
	if err != nil {
		t.Fatalf("LoadDistances() error: %v", err)
	}
	if len(distances) == 0 {
		t.Fatal("LoadDistances() returned no distances")
	}
	// Spot-check a known distance
	found := false
	for _, d := range distances {
		if d.Name == "height of the Eiffel Tower" {
			found = true
			if d.DistanceM != 330 {
				t.Errorf("Eiffel Tower distance = %f, want 330", d.DistanceM)
			}
		}
	}
	if !found {
		t.Error("Eiffel Tower not found in distances")
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./internal/data/ -v`
Expected: FAIL — `LoadItems` and `LoadDistances` undefined

**Step 3: Write minimal implementation**

Create `internal/data/embed.go`:
```go
package data

import (
	_ "embed"
	"encoding/json"
)

//go:embed items.json
var itemsJSON []byte

//go:embed distances.json
var distancesJSON []byte

type Item struct {
	Name    string  `json:"name"`
	LengthM float64 `json:"length_m"`
	Plural  string  `json:"plural"`
}

type Distance struct {
	Name      string  `json:"name"`
	DistanceM float64 `json:"distance_m"`
}

func LoadItems() ([]Item, error) {
	var items []Item
	err := json.Unmarshal(itemsJSON, &items)
	return items, err
}

func LoadDistances() ([]Distance, error) {
	var distances []Distance
	err := json.Unmarshal(distancesJSON, &distances)
	return distances, err
}
```

**Step 4: Run tests to verify they pass**

Run: `go test ./internal/data/ -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/data/embed.go internal/data/embed_test.go
git commit -m "Add data loading package with embedded JSON"
```

---

### Task 3: Unit Conversion

**Files:**
- Create: `internal/units/units.go`
- Create: `internal/units/units_test.go`

**Step 1: Write the failing tests**

Create `internal/units/units_test.go`:
```go
package units

import (
	"math"
	"testing"
)

func TestToMeters(t *testing.T) {
	tests := []struct {
		value float64
		unit  string
		want  float64
	}{
		{1, "m", 1},
		{1, "meters", 1},
		{1, "km", 1000},
		{1, "kilometers", 1000},
		{1, "ft", 0.3048},
		{1, "feet", 0.3048},
		{1, "mi", 1609.344},
		{1, "miles", 1609.344},
		{5, "km", 5000},
		{0.5, "mi", 804.672},
	}

	for _, tt := range tests {
		t.Run(tt.unit, func(t *testing.T) {
			got, err := ToMeters(tt.value, tt.unit)
			if err != nil {
				t.Fatalf("ToMeters(%f, %q) error: %v", tt.value, tt.unit, err)
			}
			if math.Abs(got-tt.want) > 0.001 {
				t.Errorf("ToMeters(%f, %q) = %f, want %f", tt.value, tt.unit, got, tt.want)
			}
		})
	}
}

func TestToMetersUnknownUnit(t *testing.T) {
	_, err := ToMeters(1, "cubits")
	if err == nil {
		t.Error("ToMeters with unknown unit should return error")
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./internal/units/ -v`
Expected: FAIL — `ToMeters` undefined

**Step 3: Write minimal implementation**

Create `internal/units/units.go`:
```go
package units

import "fmt"

var factors = map[string]float64{
	"m":          1,
	"meters":     1,
	"km":         1000,
	"kilometers": 1000,
	"ft":         0.3048,
	"feet":       0.3048,
	"mi":         1609.344,
	"miles":      1609.344,
}

func ToMeters(value float64, unit string) (float64, error) {
	factor, ok := factors[unit]
	if !ok {
		return 0, fmt.Errorf("unknown unit: %q", unit)
	}
	return value * factor, nil
}
```

**Step 4: Run tests to verify they pass**

Run: `go test ./internal/units/ -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/units/units.go internal/units/units_test.go
git commit -m "Add unit conversion package"
```

---

### Task 4: Matcher

**Files:**
- Create: `internal/matcher/matcher.go`
- Create: `internal/matcher/matcher_test.go`

**Step 1: Write the failing tests**

Create `internal/matcher/matcher_test.go`:
```go
package matcher

import (
	"math"
	"testing"

	"github.com/creimer/numberviz/internal/data"
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
	// 3.0 should score better (lower) than 7.3
	if ScoreRatio(3.0) >= ScoreRatio(7.3) {
		t.Error("expected ScoreRatio(3.0) < ScoreRatio(7.3)")
	}
}

func TestFindBestMatch(t *testing.T) {
	items := []data.Item{
		{Name: "penny", LengthM: 0.01905, Plural: "pennies"},
		{Name: "school bus", LengthM: 10.67, Plural: "school buses"},
	}
	distances := []data.Distance{
		{Name: "height of the Eiffel Tower", DistanceM: 330},
		{Name: "distance from the Earth to the Moon", DistanceM: 384400000},
	}

	result, err := FindBestMatch(330, items, distances)
	if err != nil {
		t.Fatalf("FindBestMatch() error: %v", err)
	}

	// 330m is exactly the Eiffel Tower, so ratio should be ~1
	if math.Abs(result.Ratio-1.0) > 0.01 {
		t.Errorf("ratio = %f, want ~1.0", result.Ratio)
	}
	if result.Distance.Name != "height of the Eiffel Tower" {
		t.Errorf("distance = %q, want Eiffel Tower", result.Distance.Name)
	}
	if result.Count < 1 {
		t.Errorf("count = %f, want >= 1", result.Count)
	}
}

func TestFindBestMatchNoValidPairings(t *testing.T) {
	items := []data.Item{
		{Name: "bus", LengthM: 10.67, Plural: "buses"},
	}
	distances := []data.Distance{
		{Name: "to the Moon", DistanceM: 384400000},
	}

	// Extremely tiny input — no valid pairings (count < 1)
	_, err := FindBestMatch(0.0001, items, distances)
	if err == nil {
		t.Error("expected error for no valid pairings")
	}
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./internal/matcher/ -v`
Expected: FAIL — `ScoreRatio`, `FindBestMatch` undefined

**Step 3: Write minimal implementation**

Create `internal/matcher/matcher.go`:
```go
package matcher

import (
	"fmt"
	"math"

	"github.com/creimer/numberviz/internal/data"
)

var niceNumbers = []float64{0.5, 1, 2, 3, 5, 10, 20, 50, 100, 500, 1000}

type Result struct {
	Item     data.Item
	Distance data.Distance
	Count    float64
	Ratio    float64
}

// ScoreRatio returns how "far" a ratio is from the nearest nice number.
// Lower is better.
func ScoreRatio(ratio float64) float64 {
	best := math.MaxFloat64
	for _, n := range niceNumbers {
		// Use log-space distance so 2x and 0.5x are equally far from 1x
		dist := math.Abs(math.Log10(ratio) - math.Log10(n))
		if dist < best {
			best = dist
		}
	}
	return best
}

func FindBestMatch(meters float64, items []data.Item, distances []data.Distance) (Result, error) {
	var best Result
	bestScore := math.MaxFloat64
	found := false

	for _, item := range items {
		count := meters / item.LengthM
		if count < 1 {
			continue
		}
		for _, dist := range distances {
			ratio := meters / dist.DistanceM
			if ratio < 0.01 || ratio > 100000 {
				continue
			}

			score := ScoreRatio(ratio)
			if score < bestScore {
				bestScore = score
				best = Result{
					Item:     item,
					Distance: dist,
					Count:    count,
					Ratio:    ratio,
				}
				found = true
			}
		}
	}

	if !found {
		return Result{}, fmt.Errorf("no valid comparison found for %f meters", meters)
	}
	return best, nil
}
```

**Step 4: Run tests to verify they pass**

Run: `go test ./internal/matcher/ -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/matcher/matcher.go internal/matcher/matcher_test.go
git commit -m "Add matcher package with ratio scoring and best-match selection"
```

---

### Task 5: Formatter

**Files:**
- Create: `internal/formatter/formatter.go`
- Create: `internal/formatter/formatter_test.go`

**Step 1: Write the failing tests**

Create `internal/formatter/formatter_test.go`:
```go
package formatter

import (
	"strings"
	"testing"

	"github.com/creimer/numberviz/internal/data"
	"github.com/creimer/numberviz/internal/matcher"
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
		{1.0, ""},  // ratio ≈ 1 doesn't need a label
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

func TestFormatResult(t *testing.T) {
	t.Run("ratio about 1", func(t *testing.T) {
		r := matcher.Result{
			Item:     data.Item{Name: "penny", LengthM: 0.01905, Plural: "pennies"},
			Distance: data.Distance{Name: "height of the Eiffel Tower", DistanceM: 330},
			Count:    17322.8,
			Ratio:    1.0,
		}
		out := FormatResult(r)
		if !strings.Contains(out, "Eiffel Tower") {
			t.Errorf("output should mention Eiffel Tower: %q", out)
		}
		if !strings.Contains(out, "pennies") {
			t.Errorf("output should mention pennies: %q", out)
		}
	})

	t.Run("ratio less than 1", func(t *testing.T) {
		r := matcher.Result{
			Item:     data.Item{Name: "bus", LengthM: 10.67, Plural: "buses"},
			Distance: data.Distance{Name: "height of Mount Everest", DistanceM: 8849},
			Count:    415.5,
			Ratio:    0.5,
		}
		out := FormatResult(r)
		if !strings.Contains(out, "half") {
			t.Errorf("output should contain 'half': %q", out)
		}
	})

	t.Run("ratio greater than 1", func(t *testing.T) {
		r := matcher.Result{
			Item:     data.Item{Name: "iPhone", LengthM: 0.1461, Plural: "iPhones"},
			Distance: data.Distance{Name: "height of the Eiffel Tower", DistanceM: 330},
			Count:    68445.6,
			Ratio:    10.0,
		}
		out := FormatResult(r)
		if !strings.Contains(out, "10x") {
			t.Errorf("output should contain '10x': %q", out)
		}
	})
}
```

**Step 2: Run tests to verify they fail**

Run: `go test ./internal/formatter/ -v`
Expected: FAIL — `HumanizeRatio`, `HumanizeCount`, `FormatResult` undefined

**Step 3: Write minimal implementation**

Create `internal/formatter/formatter.go`:
```go
package formatter

import (
	"fmt"
	"math"
	"strings"

	"github.com/creimer/numberviz/internal/matcher"
)

func HumanizeRatio(ratio float64) string {
	switch {
	case math.Abs(ratio-0.25) < 0.01:
		return "a quarter"
	case math.Abs(ratio-0.5) < 0.01:
		return "half"
	case math.Abs(ratio-1.0) < 0.1:
		return ""
	case ratio == math.Floor(ratio):
		return fmt.Sprintf("%dx", int(ratio))
	default:
		return fmt.Sprintf("%.1fx", ratio)
	}
}

func HumanizeCount(count float64) string {
	n := int(math.Round(count))
	if n < 0 {
		return fmt.Sprintf("-%s", HumanizeCount(float64(-n)))
	}

	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}

	var parts []string
	for len(s) > 3 {
		parts = append([]string{s[len(s)-3:]}, parts...)
		s = s[:len(s)-3]
	}
	parts = append([]string{s}, parts...)
	return strings.Join(parts, ",")
}

func FormatResult(r matcher.Result) string {
	count := HumanizeCount(r.Count)
	ratioStr := HumanizeRatio(r.Ratio)

	switch {
	case ratioStr == "": // ratio ≈ 1
		return fmt.Sprintf("That's roughly the same as the %s — about %s %s lined up.",
			r.Distance.Name, count, r.Item.Plural)
	case r.Ratio < 1:
		return fmt.Sprintf("That's about %s of the %s. You'd need to line up %s %s to cover it.",
			ratioStr, r.Distance.Name, count, r.Item.Plural)
	default:
		return fmt.Sprintf("That's about %s the %s. Imagine lining up %s %s!",
			ratioStr, r.Distance.Name, count, r.Item.Plural)
	}
}
```

**Step 4: Run tests to verify they pass**

Run: `go test ./internal/formatter/ -v`
Expected: PASS

**Step 5: Commit**

```bash
git add internal/formatter/formatter.go internal/formatter/formatter_test.go
git commit -m "Add formatter package with ratio humanization and templates"
```

---

### Task 6: CLI Entry Point

**Files:**
- Create: `cmd/numberviz/main.go`

**Step 1: Write the CLI**

Create `cmd/numberviz/main.go`:
```go
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/creimer/numberviz/internal/data"
	"github.com/creimer/numberviz/internal/formatter"
	"github.com/creimer/numberviz/internal/matcher"
	"github.com/creimer/numberviz/internal/units"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Usage: numberviz <value> <unit>\n")
		fmt.Fprintf(os.Stderr, "  Units: m, meters, km, kilometers, ft, feet, mi, miles\n")
		os.Exit(1)
	}

	value, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %q is not a valid number\n", os.Args[1])
		os.Exit(1)
	}

	meters, err := units.ToMeters(value, os.Args[2])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	items, err := data.LoadItems()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading items: %v\n", err)
		os.Exit(1)
	}

	distances, err := data.LoadDistances()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading distances: %v\n", err)
		os.Exit(1)
	}

	result, err := matcher.FindBestMatch(meters, items, distances)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(formatter.FormatResult(result))
}
```

**Step 2: Build and manually test**

Run: `go build -o numberviz ./cmd/numberviz && ./numberviz 330 meters`
Expected: Output mentioning the Eiffel Tower with ratio ~1

Run: `./numberviz 50000 km`
Expected: A comparison output with a large distance

Run: `./numberviz`
Expected: Usage message on stderr, exit code 1

**Step 3: Run all tests**

Run: `go test ./... -v`
Expected: All tests PASS

**Step 4: Commit**

```bash
git add cmd/numberviz/main.go
git commit -m "Add CLI entry point"
```

---

### Task 7: Update CLAUDE.md

**Files:**
- Modify: `CLAUDE.md`

**Step 1: Add build and test commands to CLAUDE.md**

Add the following sections:

```markdown
## Build & Test

- Build: `go build -o numberviz ./cmd/numberviz`
- Run all tests: `go test ./...`
- Run single package tests: `go test ./internal/matcher/ -v`
- Run single test: `go test ./internal/matcher/ -run TestScoreRatio -v`

## Architecture

Two-tier concept model for number visualization:

- `internal/data/` — Embedded JSON concept library (items + distances) loaded via `go:embed`
- `internal/units/` — Unit conversion (meters, km, feet, miles → meters)
- `internal/matcher/` — Scores all (item, distance) pairings by ratio quality, picks the best match
- `internal/formatter/` — Humanizes ratios and counts, selects output template
- `cmd/numberviz/` — CLI entry point
```

**Step 2: Commit**

```bash
git add CLAUDE.md
git commit -m "Update CLAUDE.md with build commands and architecture"
```
