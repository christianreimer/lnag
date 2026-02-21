package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/creimer/lnag/internal/data"
	"github.com/creimer/lnag/internal/formatter"
	"github.com/creimer/lnag/internal/matcher"
	"github.com/creimer/lnag/internal/units"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "  lnag <number> --unit <unit>\n")
	fmt.Fprintf(os.Stderr, "  lnag <number> --dimension <dimension>\n")
	os.Exit(1)
}

func main() {
	var number string
	var unitFlag, dimFlag string

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--unit":
			i++
			if i >= len(args) {
				usage()
			}
			unitFlag = args[i]
		case "--dimension":
			i++
			if i >= len(args) {
				usage()
			}
			dimFlag = args[i]
		default:
			if number != "" {
				usage()
			}
			number = args[i]
		}
	}

	if number == "" {
		usage()
	}

	value, err := strconv.ParseFloat(number, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %q is not a valid number\n", number)
		os.Exit(1)
	}

	if (unitFlag == "") == (dimFlag == "") {
		fmt.Fprintf(os.Stderr, "Error: exactly one of --unit or --dimension must be provided\n")
		os.Exit(1)
	}

	store, err := data.NewConceptStore()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading concepts: %v\n", err)
		os.Exit(1)
	}

	if unitFlag != "" {
		baseValue, dimension, err := units.Convert(value, unitFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		result, err := matcher.FindUnitMatch(baseValue, dimension, store)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(formatter.FormatUnitResult(result, value, unitFlag))
	} else {
		result, err := matcher.FindDimensionMatch(value, dimFlag, store)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println(formatter.FormatDimensionResult(result))
	}
}
