# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Large Number Analogy Generator (lnag) is a CLI tool that visualizes large or small numbers by comparing them to physical concepts (e.g., "500 m is about the length of 5 Soccer Fields").

Two modes:
- `--unit`: Given a measurement (e.g., 500 m), find one reference concept that produces a nice ratio
- `--dimension`: Given a count (e.g., 2000) and dimension (e.g., weight), find a small unit item and large target item for comparison

Supported dimensions: length, height, weight, volume, area, duration. Backed by a reference library of physical measurements (world_measurements.json) and durations (world_durations.json).

## Guidelines

- **Language preference:** Go > Python > TypeScript
- **Test-first development:** Always write unit tests before implementation code
- **Git commits:** Do not include mentions of Claude, Opus, AI assistants, or similar in commit messages

## Build & Test

- Build: `go build -o lnag ./cmd/lnag`
- Run all tests: `go test ./...`
- Run single package tests: `go test ./internal/matcher/ -v`
- Run single test: `go test ./internal/matcher/ -run TestScoreRatio -v`

## Architecture

- `internal/data/` — Embedded `world_measurements.json` and `world_durations.json` loaded via `go:embed`. `ConceptStore` provides pre-indexed, sorted dimension lookups with `DimensionIndex.FindClosest()` binary search. `NewConceptStore()` merges measurements + durations and builds sorted indexes for all dimensions
- `internal/units/` — Multi-dimension unit conversion (length, weight, volume, area, duration) to base SI units
- `internal/matcher/` — Two matching modes: `FindUnitMatch` (value + dimension + store → best concept) and `FindDimensionMatch` (count + dimension + store → unit item + target item pair). Both use `ScoreRatio` for nice-number proximity scoring. `FindDimensionMatch` uses O(n·k·log n) binary-search approach via `DimensionIndex.FindClosest()`
- `internal/formatter/` — `FormatUnitResult` and `FormatDimensionResult` with dimension-specific templates (duration uses "as long as" phrasing)
- `cmd/lnag/` — CLI entry point with `--unit` and `--dimension` flags
