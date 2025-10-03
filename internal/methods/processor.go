package methods

import (
	"fmt"

	"github.com/datsun80zx/hvac_match_parser/config"
)

// ProcessMatches orchestrates the entire matching pipeline
func ProcessMatches(cfg *config.Config) error {
	// Step 1: Load equipment list
	equipmentSet, err := LoadEquipmentList(cfg.EquipmentFile, cfg.EquipmentColumn)
	if err != nil {
		return fmt.Errorf("failed to load equipment list: %w", err)
	}
	fmt.Printf("Loaded %d equipment models\n", len(equipmentSet))

	// Step 2: Load example format headers
	headers, err := LoadExampleHeaders(cfg.ExampleFile)
	if err != nil {
		return fmt.Errorf("failed to load example headers: %w", err)
	}
	fmt.Printf("Using format with %d columns\n", len(headers))

	// Step 3: Load and filter AHRI matches
	matches, err := LoadAndFilterAHRIMatches(cfg.AHRIMatchesFile, equipmentSet, cfg.AHRIColumns)
	if err != nil {
		return fmt.Errorf("failed to load AHRI matches: %w", err)
	}
	fmt.Printf("Found %d matching AHRI combinations\n", len(matches))

	// Step 4: Write output
	if err := WriteOutput(cfg.OutputFile, headers, matches); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	return nil
}
