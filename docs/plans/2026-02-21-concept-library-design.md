# NumberViz Concept Library & Matching Logic Design

## Overview

Build a CLI tool that visualizes large or small numbers by comparing them to physical concepts. This first iteration focuses on the length/height dimension using a two-tier concept model.

## Input

```
numberviz <value> <unit>
```

Supported units: `m`, `meters`, `km`, `kilometers`, `ft`, `feet`, `mi`, `miles`. All converted to meters internally.

## Data Model (Two-Tier)

Two embedded JSON files (`go:embed`):

**`data/items.json`** — Small, countable things with a known length:
```json
[
  {"name": "iPhone", "length_m": 0.1461, "plural": "iPhones"},
  {"name": "penny", "length_m": 0.01905, "plural": "pennies"},
  {"name": "school bus", "length_m": 10.67, "plural": "school buses"}
]
```

**`data/distances.json`** — Large reference distances:
```json
[
  {"name": "height of the Eiffel Tower", "distance_m": 330},
  {"name": "distance from New York to Los Angeles", "distance_m": 3944000},
  {"name": "distance from the Earth to the Moon", "distance_m": 384400000}
]
```

## Matching Logic

1. **Parse input** — convert value + unit to meters.
2. **Score all pairings** — for every (item, distance) pair compute:
   - `count = input_meters / item.length_m`
   - `ratio = input_meters / distance.distance_m`
3. **Filter** — discard pairings where count < 1 or ratio < 0.01 or ratio > 100,000.
4. **Rank** — score by proximity to a "nice" number from `[0.5, 1, 2, 3, 5, 10, 20, 50, 100, 500, 1000]`.
5. **Pick best** — return the top-scoring pairing.

The matcher is a pure function: `(value_in_meters) -> (item, distance, count, ratio)`.

## Output Templates

Selected based on ratio:

- **ratio < 1:** `"That's about {ratio} of the {distance}. You'd need to line up {count} {items} to cover it."`
- **ratio ≈ 1:** `"That's roughly the same as the {distance} — about {count} {items} lined up."`
- **ratio > 1:** `"That's about {ratio}x the {distance}. Imagine lining up {count} {items}!"`

Ratios are humanized: 0.5 → "half", 0.25 → "a quarter", whole numbers as-is, others rounded to one decimal.

## Package Structure

```
numberviz/
├── cmd/
│   └── numberviz/
│       └── main.go          # CLI entry point, arg parsing
├── internal/
│   ├── data/
│   │   ├── items.json
│   │   ├── distances.json
│   │   └── embed.go         # go:embed + loading/parsing
│   ├── matcher/
│   │   └── matcher.go       # scoring/ranking logic
│   └── formatter/
│       └── formatter.go     # template selection + humanized output
├── go.mod
└── CLAUDE.md
```

## Technology

- Go (project preference)
- Test-first development (write unit tests before implementation)
- Embedded JSON via `go:embed`
