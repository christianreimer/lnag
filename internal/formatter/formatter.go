package formatter

import (
	"fmt"
	"math"
	"strings"

	"github.com/creimer/lnag/internal/matcher"
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

// article returns "a" or "an" for simple English usage.
func article(name string) string {
	if len(name) == 0 {
		return "a"
	}
	first := strings.ToLower(name[:1])
	if first == "a" || first == "e" || first == "i" || first == "o" || first == "u" {
		return "an"
	}
	return "a"
}

// dimensionNoun returns the noun phrase for a dimension (e.g. "length", "weight").
func dimensionNoun(dimension string) string {
	switch dimension {
	case "length":
		return "length"
	case "height":
		return "height"
	case "width":
		return "width"
	case "weight":
		return "weight"
	case "volume":
		return "volume"
	case "area":
		return "area"
	case "duration":
		return "duration"
	default:
		return dimension
	}
}

// dimensionVerb returns the verb phrase for stacking/lining up in a dimension.
func dimensionVerb(dimension string) string {
	switch dimension {
	case "length":
		return "lined up would stretch"
	case "height":
		return "stacked would be"
	case "weight":
		return "would weigh"
	case "volume":
		return "would fill"
	case "area":
		return "would cover"
	case "duration":
		return "would last"
	default:
		return "would equal"
	}
}

// pluralize adds an "s" to simple names.
func pluralize(name string) string {
	return name + "s"
}

// FormatUnitResult formats a unit-mode result.
// Example: "500 m is about the length of 5 Soccer Fields."
func FormatUnitResult(r matcher.UnitResult, inputValue float64, unit string) string {
	ratioStr := HumanizeRatio(r.Ratio)
	countStr := HumanizeCount(r.Ratio)
	dim := dimensionNoun(r.Dimension)
	name := r.Concept.Name
	proper := r.Concept.ProperNoun
	inputStr := HumanizeCount(inputValue)

	switch r.Dimension {
	case "duration":
		switch {
		case ratioStr == "" && proper:
			return fmt.Sprintf("%s %s is about as long as the %s.", inputStr, unit, name)
		case ratioStr == "":
			return fmt.Sprintf("%s %s is about as long as 1 %s.", inputStr, unit, name)
		case r.Ratio < 1 && proper:
			return fmt.Sprintf("%s %s is about %s as long as the %s.", inputStr, unit, ratioStr, name)
		case r.Ratio < 1:
			return fmt.Sprintf("%s %s is about %s as long as %s %s.", inputStr, unit, ratioStr, article(name), name)
		case proper:
			return fmt.Sprintf("%s %s is about %s as long as the %s.", inputStr, unit, ratioStr, name)
		default:
			return fmt.Sprintf("%s %s is about as long as %s %s.", inputStr, unit, countStr, pluralize(name))
		}
	default:
		switch {
		case ratioStr == "" && proper:
			return fmt.Sprintf("%s %s is about the %s of the %s.", inputStr, unit, dim, name)
		case ratioStr == "":
			return fmt.Sprintf("%s %s is about the %s of 1 %s.", inputStr, unit, dim, name)
		case proper:
			return fmt.Sprintf("%s %s is about %s the %s of the %s.", inputStr, unit, ratioStr, dim, name)
		case r.Ratio < 1:
			return fmt.Sprintf("%s %s is about %s the %s of %s %s.", inputStr, unit, ratioStr, dim, article(name), name)
		default:
			return fmt.Sprintf("%s %s is about the %s of %s %s.", inputStr, unit, dim, countStr, pluralize(name))
		}
	}
}

// FormatDimensionResult formats a dimension-mode result.
// Example: "2,000 Watermelons would weigh about as much as 2 African Elephants."
func FormatDimensionResult(r matcher.DimensionResult) string {
	countStr := HumanizeCount(r.Count)
	unitName := pluralize(r.UnitItem.Name)
	targetName := r.TargetItem.Name
	proper := r.TargetItem.ProperNoun
	verb := dimensionVerb(r.Dimension)
	ratioStr := HumanizeRatio(r.Ratio)
	ratioCount := HumanizeCount(r.Ratio)
	art := article(targetName)
	if proper {
		art = "the"
	}

	switch r.Dimension {
	case "weight":
		switch {
		case ratioStr == "":
			return fmt.Sprintf("%s %s %s about as much as %s %s.", countStr, unitName, verb, art, targetName)
		case r.Ratio < 1:
			return fmt.Sprintf("%s %s %s about %s as much as %s %s.", countStr, unitName, verb, ratioStr, art, targetName)
		case proper:
			return fmt.Sprintf("%s %s %s about %s as much as the %s.", countStr, unitName, verb, ratioStr, targetName)
		default:
			return fmt.Sprintf("%s %s %s about as much as %s %s.", countStr, unitName, verb, ratioCount, pluralize(targetName))
		}
	case "duration":
		switch {
		case ratioStr == "":
			return fmt.Sprintf("%s %s %s about as long as %s %s.", countStr, unitName, verb, art, targetName)
		case r.Ratio < 1:
			return fmt.Sprintf("%s %s %s about %s as long as %s %s.", countStr, unitName, verb, ratioStr, art, targetName)
		case proper:
			return fmt.Sprintf("%s %s %s about %s as long as the %s.", countStr, unitName, verb, ratioStr, targetName)
		default:
			return fmt.Sprintf("%s %s %s about as long as %s %s.", countStr, unitName, verb, ratioCount, pluralize(targetName))
		}
	default: // length, height, width, area, volume
		dim := dimensionNoun(r.Dimension)
		switch {
		case ratioStr == "":
			return fmt.Sprintf("%s %s %s about the %s of %s %s.", countStr, unitName, verb, dim, art, targetName)
		case proper:
			return fmt.Sprintf("%s %s %s about %s the %s of the %s.", countStr, unitName, verb, ratioStr, dim, targetName)
		case r.Ratio < 1:
			return fmt.Sprintf("%s %s %s about %s the %s of %s %s.", countStr, unitName, verb, ratioStr, dim, article(targetName), targetName)
		default:
			return fmt.Sprintf("%s %s %s about %s the %s of %s %s.", countStr, unitName, verb, ratioStr, dim, article(targetName), targetName)
		}
	}
}
