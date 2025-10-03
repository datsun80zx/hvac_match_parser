package config

import (
	"flag"
	"fmt"
	"strings"
)

// Config holds the file paths and column mappings
type Config struct {
	EquipmentFile   string
	AHRIMatchesFile string
	ExampleFile     string
	OutputFile      string
	EquipmentColumn int
	AHRIColumns     []int
}

// ParseFlags parses command line flags and returns a Config
func ParseFlags() *Config {
	equipmentFile := flag.String("equipment", "equipment.csv", "Path to equipment list CSV")
	ahriFile := flag.String("ahri", "ahri_matches.csv", "Path to AHRI matches CSV")
	exampleFile := flag.String("example", "example_format.csv", "Path to example format CSV")
	outputFile := flag.String("output", "filtered_matches.csv", "Path to output CSV")
	equipmentCol := flag.Int("eq-col", 0, "Column index for model numbers in equipment file (0-based)")
	ahriCols := flag.String("ahri-cols", "0,1", "Comma-separated column indices for model numbers in AHRI file")

	flag.Parse()

	return &Config{
		EquipmentFile:   *equipmentFile,
		AHRIMatchesFile: *ahriFile,
		ExampleFile:     *exampleFile,
		OutputFile:      *outputFile,
		EquipmentColumn: *equipmentCol,
		AHRIColumns:     parseColumnIndices(*ahriCols),
	}
}

// parseColumnIndices converts a comma-separated string into a slice of integers
func parseColumnIndices(s string) []int {
	parts := strings.Split(s, ",")
	indices := make([]int, 0, len(parts))
	for _, part := range parts {
		var idx int
		fmt.Sscanf(strings.TrimSpace(part), "%d", &idx)
		indices = append(indices, idx)
	}
	return indices
}
