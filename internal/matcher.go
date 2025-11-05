package internal

import (
	"fmt"
	"strings"

	"github.com/datsun80zx/hvac_match_parser/internal/data_structures"
)

var systemTypes = map[string]string{
	"central ac":               "central_ac",
	"central ac & air handler": "central_ac_air_handler",
	"central ac & furnace":     "central_ac_furnace",
	"heat pump & air handler":  "air_source_heat_pump_electric_heat",
	"heat pump & furnace":      "air_source_heat_pump_furnace",
	"heat pump":                "air_source_heat_pump",
	"furnace":                  "furnace",
}

/*
GenerateFullSystemEquipmentConfig generates a Cartesian product of equipment combinations.
It now separates standard and communicating equipment to ensure proper pairing.
Equipment list provided must all be from the same brand.
*/
func GenerateFullSystemEquipmentConfig(list []data_structures.Equipment, sysType string) ([]data_structures.ComponentKey, error) {
	// Create nested map: equipByTypeAndCategory[type][category][]Equipment
	equipByTypeAndCategory := make(map[string]map[string][]data_structures.Equipment)

	types := []string{"furnace", "handler", "coil", "ac", "hp"}
	categories := []string{data_structures.CategoryStandard, data_structures.CategoryCommunicating}

	// Initialize nested maps
	for _, t := range types {
		equipByTypeAndCategory[t] = make(map[string][]data_structures.Equipment)
		for _, c := range categories {
			equipByTypeAndCategory[t][c] = []data_structures.Equipment{}
		}
	}

	// Sort equipment by type and category
	for _, item := range list {
		if strings.Contains(item.Type, "furnace") {
			equipByTypeAndCategory["furnace"][item.Category] = append(
				equipByTypeAndCategory["furnace"][item.Category], item)
		} else if strings.Contains(item.Type, "handler") {
			equipByTypeAndCategory["handler"][item.Category] = append(
				equipByTypeAndCategory["handler"][item.Category], item)
		} else if strings.Contains(item.Type, "coil") {
			equipByTypeAndCategory["coil"][item.Category] = append(
				equipByTypeAndCategory["coil"][item.Category], item)
		} else if strings.Contains(item.Type, "ac") {
			equipByTypeAndCategory["ac"][item.Category] = append(
				equipByTypeAndCategory["ac"][item.Category], item)
		} else if strings.Contains(item.Type, "hp") {
			equipByTypeAndCategory["hp"][item.Category] = append(
				equipByTypeAndCategory["hp"][item.Category], item)
		} else {
			return nil, fmt.Errorf("unknown equipment type: %s", item.Type)
		}
	}

	equipConfigs := make([]data_structures.ComponentKey, 0)

	// Generate combinations for each category separately
	// This ensures communicating equipment only pairs with communicating equipment
	for _, category := range categories {
		combos, err := generateCombosForCategory(
			equipByTypeAndCategory,
			category,
			sysType,
		)
		if err != nil {
			return nil, err
		}
		equipConfigs = append(equipConfigs, combos...)
	}

	return equipConfigs, nil
}

// generateCombosForCategory creates equipment combinations within a single category
// This ensures standard equipment doesn't mix with communicating equipment
// Note: Furnaces are shared across categories since they work with both types
func generateCombosForCategory(
	equipMap map[string]map[string][]data_structures.Equipment,
	category string,
	sysType string,
) ([]data_structures.ComponentKey, error) {

	equipConfigs := make([]data_structures.ComponentKey, 0)

	// Get equipment for this category
	// Furnaces are shared - combine both standard and communicating (though typically all standard)
	furnaces := append(equipMap["furnace"][data_structures.CategoryStandard],
		equipMap["furnace"][data_structures.CategoryCommunicating]...)
	airHandlers := equipMap["handler"][category]
	coils := equipMap["coil"][category]
	airCons := equipMap["ac"][category]
	heatPumps := equipMap["hp"][category]

	switch sysType {
	case "heat pump & air handler":
		for _, hp := range heatPumps {
			for _, ah := range airHandlers {
				equipConfigs = append(equipConfigs, data_structures.ComponentKey{
					Brand:       hp.Brand,
					IndoorUnit:  ah,
					OutdoorUnit: hp,
					SystemType:  systemTypes["heat pump & air handler"],
				})
			}
		}
		return equipConfigs, nil

	case "heat pump & furnace":
		for _, hp := range heatPumps {
			for _, f := range furnaces {
				for _, c := range coils {
					equipConfigs = append(equipConfigs, data_structures.ComponentKey{
						Brand:       hp.Brand,
						IndoorUnit:  c,
						Furnace:     f,
						OutdoorUnit: hp,
						SystemType:  systemTypes["heat pump & furnace"],
					})
				}
			}
		}
		return equipConfigs, nil

	case "central ac & furnace":
		for _, ac := range airCons {
			for _, f := range furnaces {
				for _, c := range coils {
					equipConfigs = append(equipConfigs, data_structures.ComponentKey{
						Brand:       ac.Brand,
						IndoorUnit:  c,
						Furnace:     f,
						OutdoorUnit: ac,
						SystemType:  systemTypes["central ac & furnace"],
					})
				}
			}
		}
		return equipConfigs, nil

	case "central ac & air handler":
		for _, ac := range airCons {
			for _, ah := range airHandlers {
				equipConfigs = append(equipConfigs, data_structures.ComponentKey{
					Brand:       ac.Brand,
					IndoorUnit:  ah,
					OutdoorUnit: ac,
					SystemType:  systemTypes["central ac & air handler"],
				})
			}
		}
		return equipConfigs, nil

	case "furnace":
		for _, f := range furnaces {
			equipConfigs = append(equipConfigs, data_structures.ComponentKey{
				Brand:      f.Brand,
				Furnace:    f,
				SystemType: systemTypes["furnace"],
			})
		}
		return equipConfigs, nil

	case "central ac":
		for _, ac := range airCons {
			for _, c := range coils {
				equipConfigs = append(equipConfigs, data_structures.ComponentKey{
					Brand:       ac.Brand,
					IndoorUnit:  c,
					OutdoorUnit: ac,
					SystemType:  systemTypes["central ac"],
				})
			}
		}
		return equipConfigs, nil
	}

	return equipConfigs, nil
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
		furnace := NormalizeString(record.Furnace)
		indoorUnit := NormalizeString(record.IndoorUnit)
		outdoorUnit := NormalizeString(record.OutdoorUnit)
		furnaceVariations := expandFurnaceWildcard(furnace.NormalizedModelNumber)

		indoorVariations := expandIndoorUnitWildcard(indoorUnit.NormalizedModelNumber)

		// Outdoor units don't have wildcards
		outdoorVariations := []string{outdoorUnit.NormalizedModelNumber}

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
	ahriMap map[string]string,
) ([]data_structures.OutputCSV, error) {

	certifiedMatches := make([]data_structures.OutputCSV, 0)

	for _, combo := range fullSystemCombos {
		// Handle system types that don't need AHRI certification
		if combo.SystemType == systemTypes["furnace"] {
			output := data_structures.OutputCSV{
				Brand:        combo.Brand,
				Orientation:  "",
				Furnace:      combo.Furnace.InputModelNumber,
				TypeOfSystem: combo.SystemType,
			}
			certifiedMatches = append(certifiedMatches, output)
			continue
		}

		if combo.SystemType == systemTypes["central ac"] {
			// Apply filters for central ac systems
			if !isValidIndoorUnit(combo.IndoorUnit) {
				continue
			}
			if !isValidTonnageMatch(combo.OutdoorUnit, combo.IndoorUnit) {
				continue
			}

			output := data_structures.OutputCSV{
				Brand:          combo.Brand,
				Orientation:    "",
				OutdoorUnit:    combo.OutdoorUnit.InputModelNumber,
				EvaporatorCoil: combo.IndoorUnit.InputModelNumber,
				TypeOfSystem:   combo.SystemType,
			}
			certifiedMatches = append(certifiedMatches, output)
			continue
		}

		// For all other system types, apply standard filters and AHRI lookup

		// Filter out horizontal coils
		if !isValidIndoorUnit(combo.IndoorUnit) {
			continue
		}

		// Filter tonnage and cabinet mismatches for systems with coils and furnaces
		if needsCabinetValidation(combo.SystemType) {
			if !isValidCabinetAndTonnage(combo) {
				continue
			}
		}

		// Lookup AHRI certification
		ahriNumber, isCertified := FindAHRICertification(combo, ahriMap)
		if !isCertified {
			continue
		}

		output := createAHRIOutput(combo, ahriNumber)
		certifiedMatches = append(certifiedMatches, output)
	}

	return certifiedMatches, nil
}

// Helper functions for safer validation
func isValidIndoorUnit(indoor data_structures.Equipment) bool {
	if len(indoor.NormalizedModelNumber) < 2 {
		return false
	}
	// Filter out horizontal coils
	return indoor.NormalizedModelNumber[1] != 'H'
}

func isValidTonnageMatch(outdoor, indoor data_structures.Equipment) bool {
	outdoorModel := outdoor.NormalizedModelNumber
	indoorModel := indoor.NormalizedModelNumber

	// Check we have enough characters
	if len(outdoorModel) < 4 || len(indoorModel) < 7 {
		return false
	}

	// Compare tonnage (outdoor 4th and 3rd from end vs indoor positions 5-7)
	return outdoorModel[len(outdoorModel)-4:len(outdoorModel)-2] == indoorModel[5:7]
}

func isValidCabinetAndTonnage(combo data_structures.ComponentKey) bool {
	outdoor := combo.OutdoorUnit.NormalizedModelNumber
	indoor := combo.IndoorUnit.NormalizedModelNumber
	furnace := combo.Furnace.NormalizedModelNumber

	// Verify we have coils to check
	if !strings.Contains(combo.IndoorUnit.Type, "coil") {
		return true // Not applicable
	}

	// Check string lengths
	if len(outdoor) < 4 || len(indoor) < 10 || len(furnace) < 11 {
		return false
	}

	// Check tonnage match (outdoor 4th and 3rd from end)
	if outdoor[len(outdoor)-4:len(outdoor)-2] != indoor[5:7] {
		return false
	}

	// Check cabinet size match
	if indoor[9] != furnace[10] {
		return false
	}

	return true
}

func needsCabinetValidation(systemType string) bool {
	return systemType == systemTypes["central ac & furnace"] ||
		systemType == systemTypes["heat pump & furnace"]
}

func createAHRIOutput(combo data_structures.ComponentKey, ahriNumber string) data_structures.OutputCSV {
	output := data_structures.OutputCSV{
		AHRINumber:   ahriNumber,
		Brand:        combo.Brand,
		Orientation:  "",
		TypeOfSystem: combo.SystemType,
	}

	switch combo.SystemType {
	case systemTypes["central ac & air handler"]:
		output.OutdoorUnit = combo.OutdoorUnit.InputModelNumber
		output.AirHandler = combo.IndoorUnit.InputModelNumber

	case systemTypes["central ac & furnace"]:
		output.OutdoorUnit = combo.OutdoorUnit.InputModelNumber
		output.Furnace = combo.Furnace.InputModelNumber
		output.EvaporatorCoil = combo.IndoorUnit.InputModelNumber

	case systemTypes["heat pump & air handler"]:
		output.OutdoorUnit = combo.OutdoorUnit.InputModelNumber
		output.AirHandler = combo.IndoorUnit.InputModelNumber

	case systemTypes["heat pump & furnace"]:
		output.OutdoorUnit = combo.OutdoorUnit.InputModelNumber
		output.Furnace = combo.Furnace.InputModelNumber
		output.EvaporatorCoil = combo.IndoorUnit.InputModelNumber

	case systemTypes["heat pump"]:
		output.OutdoorUnit = combo.OutdoorUnit.InputModelNumber
		output.EvaporatorCoil = combo.IndoorUnit.InputModelNumber
	}

	return output
}
