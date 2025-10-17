package internal

import (
	"strings"

	"github.com/datsun80zx/hvac_match_parser/internal/data_structures"
)

// The idea here is to generate a Cartesian product of the different system types.
func GenerateFullSystemEquipmentConfig(e data_structures.Equipment) []data_structures.ComponentKey {

	totalCombinations := len(e.Furnaces) * len(e.IndoorUnits) * len(e.OutdoorUnits)
	equipConfigs := make([]data_structures.ComponentKey, 0, totalCombinations)

	for _, furnace := range e.Furnaces {
		for _, indoorUnit := range e.IndoorUnits {
			for _, outdoorUnit := range e.OutdoorUnits {
				equipConfig := data_structures.ComponentKey{
					Brand:       furnace.Brand,
					Furnace:     furnace,
					IndoorUnit:  indoorUnit,
					OutdoorUnit: outdoorUnit,
				}
				equipConfigs = append(equipConfigs, equipConfig)
			}
		}
	}
	return equipConfigs
}

func GenerateAirHandlerEquipmentConfig(e data_structures.Equipment) []data_structures.ComponentKey {

	totalCombinations := len(e.IndoorUnits) * len(e.OutdoorUnits)
	equipConfigs := make([]data_structures.ComponentKey, 0, totalCombinations)

	for _, indoorUnit := range e.IndoorUnits {
		for _, outdoorUnit := range e.OutdoorUnits {
			equipConfig := data_structures.ComponentKey{
				Brand:       indoorUnit.Brand,
				IndoorUnit:  indoorUnit,
				OutdoorUnit: outdoorUnit,
			}
			equipConfigs = append(equipConfigs, equipConfig)
		}
	}

	return equipConfigs
}

func expandFurnaceWildcard(model string) []string {
	// If the model doesn't have a wildcard, return it as-is
	if !strings.Contains(model, "*") {
		return []string{model}
	}

	// The wildcard is always in position 1 (second character, zero-indexed)
	// Generate both possible variations
	runes := []rune(model)

	// Find the wildcard position (should be position 1, but let's be safe)
	wildcardPos := -1
	for i, char := range runes {
		if char == '*' {
			wildcardPos = i
			break
		}
	}

	if wildcardPos == -1 {
		return []string{model}
	}

	// Create both variations: one with 'R', one with 'D'
	variation1 := make([]rune, len(runes))
	copy(variation1, runes)
	variation1[wildcardPos] = 'R'

	variation2 := make([]rune, len(runes))
	copy(variation2, runes)
	variation2[wildcardPos] = 'D'

	return []string{string(variation1), string(variation2)}
}

func expandIndoorUnitWildcard(model string) []string {
	// If no wildcards, return as-is
	if !strings.Contains(model, "*") {
		return []string{model}
	}

	runes := []rune(model)

	// Find all wildcard positions
	wildcardPositions := []int{}
	for i, char := range runes {
		if char == '*' {
			wildcardPositions = append(wildcardPositions, i)
		}
	}

	// Handle the first wildcard (position 2, zero-indexed) - always 'P'
	for _, pos := range wildcardPositions {
		if pos == 2 { // Third character
			runes[pos] = 'P'
		}
	}

	// Find the second wildcard (second from last position)
	// After normalization, this should be at len(runes) - 2
	secondWildcardPos := -1
	for _, pos := range wildcardPositions {
		if pos == len(runes)-2 {
			secondWildcardPos = pos
			break
		}
	}

	// If there's no second wildcard, we're done
	if secondWildcardPos == -1 {
		return []string{string(runes)}
	}

	// Generate all four variations for the second wildcard
	possibleChars := []rune{'A', 'B', 'C', 'D'}
	results := make([]string, 0, len(possibleChars))

	for _, char := range possibleChars {
		variation := make([]rune, len(runes))
		copy(variation, runes)
		variation[secondWildcardPos] = char
		results = append(results, string(variation))
	}

	return results
}

func BuildAHRIMap(ahriList []data_structures.AHRIRecord) map[string]string {
	ahriMap := make(map[string]string)

	for _, record := range ahriList {
		// Expand each component based on its type
		furnaceVariations := expandFurnaceWildcard(record.Furnace)

		indoorVariations := expandIndoorUnitWildcard(record.IndoorUnit)

		// Outdoor units don't have wildcards
		outdoorVariations := []string{record.OutdoorUnit}

		// Create map entries for all combinations
		for _, furnace := range furnaceVariations {
			for _, indoor := range indoorVariations {
				for _, outdoor := range outdoorVariations {
					key := outdoor + "|" + indoor + "|" + furnace
					ahriMap[key] = record.AHRINumber
				}
			}
		}
	}

	return ahriMap
}

func FindAHRICertification(config data_structures.ComponentKey, ahriMap map[string]string) (string, bool) {
	// Build the lookup key from normalized model numbers
	key := config.OutdoorUnit.NormalizedModelNumber + "|" +
		config.IndoorUnit.NormalizedModelNumber + "|" +
		config.Furnace.NormalizedModelNumber

	// Look it up in the map
	ahriNumber, certified := ahriMap[key]
	return ahriNumber, certified
}

func FindCertifiedMatches(
	fullSystemCombos []data_structures.ComponentKey,
	airHandlerCombos []data_structures.ComponentKey,
	ahriMap map[string]string,
) []data_structures.OutputCSV {

	var certifiedMatches []data_structures.OutputCSV

	// Process full system combinations
	for _, combo := range fullSystemCombos {
		// Check if this combination is certified and get the AHRI number
		ahriNumber, isCertified := FindAHRICertification(combo, ahriMap)

		if isCertified {
			// Determine system type based on heat pump flag
			var systemType string
			if combo.OutdoorUnit.HeatPump {
				systemType = "air_source_heat_pump_furnace"
			} else {
				systemType = "central_ac_furnace"
			}

			// Determine indoor unit type
			var evapCoil, airHandler string
			if combo.IndoorUnit.AirHandler {
				airHandler = combo.IndoorUnit.InputModelNumber
			} else {
				evapCoil = combo.IndoorUnit.InputModelNumber
			}

			// Create output record with AHRI number included
			output := data_structures.OutputCSV{
				AHRINumber:     ahriNumber,
				Brand:          combo.Brand,
				Orientation:    "",
				TypeOfSystem:   systemType,
				OutdoorUnit:    combo.OutdoorUnit.InputModelNumber,
				Furnace:        combo.Furnace.InputModelNumber,
				EvaporatorCoil: evapCoil,
				AirHandler:     airHandler,
			}

			certifiedMatches = append(certifiedMatches, output)
		}
	}

	// Process air handler combinations
	for _, combo := range airHandlerCombos {
		ahriNumber, isCertified := FindAHRICertification(combo, ahriMap)

		if isCertified {
			var systemType string
			if combo.OutdoorUnit.HeatPump && combo.IndoorUnit.AirHandler {
				systemType = "air_source_heat_pump_electric_heat"
			} else if combo.IndoorUnit.AirHandler {
				systemType = "central_ac_electric_heat"
			}

			var evapCoil, airHandler string
			if combo.IndoorUnit.AirHandler {
				airHandler = combo.IndoorUnit.InputModelNumber
			} else {
				evapCoil = combo.IndoorUnit.InputModelNumber
			}

			output := data_structures.OutputCSV{
				AHRINumber:     ahriNumber, // And here too
				TypeOfSystem:   systemType,
				OutdoorUnit:    combo.OutdoorUnit.InputModelNumber,
				Furnace:        "",
				EvaporatorCoil: evapCoil,
				AirHandler:     airHandler,
			}

			certifiedMatches = append(certifiedMatches, output)
		}
	}

	return certifiedMatches
}
