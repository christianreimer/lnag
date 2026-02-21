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
	case ratio > 1 && ratio == math.Floor(ratio):
		return fmt.Sprintf("%dx", int(ratio))
	case ratio > 1:
		floor := int(math.Floor(ratio))
		frac := ratio - float64(floor)
		switch {
		case frac < 0.4:
			return fmt.Sprintf("more than %d times", floor)
		case frac < 0.7:
			return fmt.Sprintf("%d and a half times", floor)
		default:
			return fmt.Sprintf("almost %d times", floor+1)
		}
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

// ApproxCount formats a ratio as an approximate count for display with plural nouns.
func ApproxCount(ratio float64) string {
	if ratio == math.Floor(ratio) {
		return HumanizeCount(ratio)
	}
	floor := int(math.Floor(ratio))
	frac := ratio - float64(floor)
	switch {
	case frac < 0.4:
		return fmt.Sprintf("more than %s", HumanizeCount(float64(floor)))
	case frac < 0.7:
		return fmt.Sprintf("%s and a half", HumanizeCount(float64(floor)))
	default:
		return fmt.Sprintf("almost %s", HumanizeCount(float64(floor+1)))
	}
}

// isDirectional returns true if the phrase already conveys approximation
// (e.g. "more than", "almost"), making "about" redundant.
func isDirectional(s string) bool {
	return strings.HasPrefix(s, "more than") || strings.HasPrefix(s, "almost")
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
	case "distance":
		return "distance"
	case "duration":
		return "duration"
	default:
		return dimension
	}
}

// dimensionPreposition returns the preposition used between the dimension noun
// and the concept name. Most dimensions use "of" ("the length of"),
// but distance uses "to" ("the distance to").
func dimensionPreposition(dimension string) string {
	if dimension == "distance" {
		return "to"
	}
	return "of"
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
	case "distance":
		return "would span"
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
	countStr := ApproxCount(r.Ratio)
	dim := dimensionNoun(r.Dimension)
	name := r.Concept.Name
	proper := r.Concept.ProperNoun
	inputStr := HumanizeCount(inputValue)

	about := "about "
	if isDirectional(ratioStr) || isDirectional(countStr) {
		about = ""
	}

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
			return fmt.Sprintf("%s %s is %s%s as long as the %s.", inputStr, unit, about, ratioStr, name)
		default:
			return fmt.Sprintf("%s %s is %sas long as %s %s.", inputStr, unit, about, countStr, pluralize(name))
		}
	case "distance":
		target := name
		if proper {
			target = "the " + name
		}
		switch {
		case ratioStr == "":
			return fmt.Sprintf("%s %s is about the distance to %s.", inputStr, unit, target)
		case r.Ratio < 1:
			return fmt.Sprintf("%s %s is about %s the distance to %s.", inputStr, unit, ratioStr, target)
		default:
			return fmt.Sprintf("%s %s is %s%s the distance to %s.", inputStr, unit, about, ratioStr, target)
		}
	default:
		prep := dimensionPreposition(r.Dimension)
		switch {
		case ratioStr == "" && proper:
			return fmt.Sprintf("%s %s is about the %s %s the %s.", inputStr, unit, dim, prep, name)
		case ratioStr == "":
			return fmt.Sprintf("%s %s is about the %s %s 1 %s.", inputStr, unit, dim, prep, name)
		case proper:
			return fmt.Sprintf("%s %s is %s%s the %s %s the %s.", inputStr, unit, about, ratioStr, dim, prep, name)
		case r.Ratio < 1:
			return fmt.Sprintf("%s %s is about %s the %s %s %s %s.", inputStr, unit, ratioStr, dim, prep, article(name), name)
		default:
			return fmt.Sprintf("%s %s is %sthe %s %s %s %s.", inputStr, unit, about, dim, prep, countStr, pluralize(name))
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
	ratioCount := ApproxCount(r.Ratio)
	art := article(targetName)
	if proper {
		art = "the"
	}

	about := "about "
	if isDirectional(ratioStr) || isDirectional(ratioCount) {
		about = ""
	}

	switch r.Dimension {
	case "weight":
		switch {
		case ratioStr == "":
			return fmt.Sprintf("%s %s %s about as much as %s %s.", countStr, unitName, verb, art, targetName)
		case r.Ratio < 1:
			return fmt.Sprintf("%s %s %s about %s as much as %s %s.", countStr, unitName, verb, ratioStr, art, targetName)
		case proper:
			return fmt.Sprintf("%s %s %s %s%s as much as the %s.", countStr, unitName, verb, about, ratioStr, targetName)
		default:
			return fmt.Sprintf("%s %s %s %sas much as %s %s.", countStr, unitName, verb, about, ratioCount, pluralize(targetName))
		}
	case "duration":
		switch {
		case ratioStr == "":
			return fmt.Sprintf("%s %s %s about as long as %s %s.", countStr, unitName, verb, art, targetName)
		case r.Ratio < 1:
			return fmt.Sprintf("%s %s %s about %s as long as %s %s.", countStr, unitName, verb, ratioStr, art, targetName)
		case proper:
			return fmt.Sprintf("%s %s %s %s%s as long as the %s.", countStr, unitName, verb, about, ratioStr, targetName)
		default:
			return fmt.Sprintf("%s %s %s %sas long as %s %s.", countStr, unitName, verb, about, ratioCount, pluralize(targetName))
		}
	case "distance":
		unitTarget := r.UnitItem.Name
		if r.UnitItem.ProperNoun {
			unitTarget = "the " + r.UnitItem.Name
		}
		targetRef := r.TargetItem.Name
		if proper {
			targetRef = "the " + r.TargetItem.Name
		}
		unitPhrase := fmt.Sprintf("%sx the distance to %s", countStr, unitTarget)
		switch {
		case ratioStr == "":
			return fmt.Sprintf("%s %s about the distance to %s.", unitPhrase, verb, targetRef)
		case r.Ratio < 1:
			return fmt.Sprintf("%s %s about %s the distance to %s.", unitPhrase, verb, ratioStr, targetRef)
		default:
			return fmt.Sprintf("%s %s %s%s the distance to %s.", unitPhrase, verb, about, ratioStr, targetRef)
		}
	default: // length, height, width, area, volume
		dim := dimensionNoun(r.Dimension)
		prep := dimensionPreposition(r.Dimension)
		switch {
		case ratioStr == "":
			return fmt.Sprintf("%s %s %s about the %s %s %s %s.", countStr, unitName, verb, dim, prep, art, targetName)
		case proper:
			return fmt.Sprintf("%s %s %s %s%s the %s %s the %s.", countStr, unitName, verb, about, ratioStr, dim, prep, targetName)
		case r.Ratio < 1:
			return fmt.Sprintf("%s %s %s about %s the %s %s %s %s.", countStr, unitName, verb, ratioStr, dim, prep, article(targetName), targetName)
		default:
			return fmt.Sprintf("%s %s %s %s%s the %s %s %s %s.", countStr, unitName, verb, about, ratioStr, dim, prep, article(targetName), targetName)
		}
	}
}
