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
	The idea here is to generate a Cartesian product of the different system types.

It is important that the list of equipment provided to this function is all of the same brand and is a valid
system type
*/
func GenerateFullSystemEquipmentConfig(list []data_structures.Equipment, sysType string) ([]data_structures.ComponentKey, error) {
	equipConfigs := make([]data_structures.ComponentKey, 0)

	furnaces := []data_structures.Equipment{}
	airHandlers := []data_structures.Equipment{}
	coils := []data_structures.Equipment{}
	airCons := []data_structures.Equipment{}
	heatPumps := []data_structures.Equipment{}

	for _, item := range list {
		if strings.Contains(item.Type, "furnace") {
			furnaces = append(furnaces, item)
		} else if strings.Contains(item.Type, "handler") {
			airHandlers = append(airHandlers, item)
		} else if strings.Contains(item.Type, "coil") {
			coils = append(coils, item)
		} else if strings.Contains(item.Type, "ac") {
			airCons = append(airCons, item)
		} else if strings.Contains(item.Type, "hp") {
			heatPumps = append(heatPumps, item)
		} else {
			return nil, fmt.Errorf("there is an issue with sorting equipment by type in the GenerateFullSystem..func")
		}
	}

	switch sysType {
	case "heat pump & air handler":
		for _, heatPump := range heatPumps {
			for _, airHandler := range airHandlers {
				equipCombo := data_structures.ComponentKey{
					Brand:       heatPump.Brand,
					IndoorUnit:  airHandler,
					OutdoorUnit: heatPump,
					SystemType:  systemTypes["heat pump & air handler"],
				}
				equipConfigs = append(equipConfigs, equipCombo)
			}
		}
		return equipConfigs, nil
	case "heat pump & furnace":
		for _, heatPump := range heatPumps {
			for _, furnace := range furnaces {
				for _, coil := range coils {
					equipCombo := data_structures.ComponentKey{
						Brand:       heatPump.Brand,
						IndoorUnit:  coil,
						Furnace:     furnace,
						OutdoorUnit: heatPump,
						SystemType:  systemTypes["heat pump & furnace"],
					}
					equipConfigs = append(equipConfigs, equipCombo)
				}
			}
		}
		return equipConfigs, nil
	case "central ac & furnace":
		for _, airCon := range airCons {
			for _, furnace := range furnaces {
				for _, coil := range coils {
					equipCombo := data_structures.ComponentKey{
						Brand:       airCon.Brand,
						IndoorUnit:  coil,
						Furnace:     furnace,
						OutdoorUnit: airCon,
						SystemType:  systemTypes["central ac & furnace"],
					}
					equipConfigs = append(equipConfigs, equipCombo)
				}
			}
		}
		return equipConfigs, nil
	case "central ac & air handler":
		for _, airCon := range airCons {
			for _, airHandler := range airHandlers {
				equipCombo := data_structures.ComponentKey{
					Brand:       airCon.Brand,
					IndoorUnit:  airHandler,
					OutdoorUnit: airCon,
					SystemType:  systemTypes["central ac & air handler"],
				}
				equipConfigs = append(equipConfigs, equipCombo)
			}
		}
		return equipConfigs, nil
	}
	return nil, fmt.Errorf("there was an error processing all equipment configurations")
}

// func GenerateAirHandlerEquipmentConfig(e data_structures.Equipment) []data_structures.ComponentKey {
// 	equipConfigs := make([]data_structures.ComponentKey, 0)

// 	// Collect all unique brands
// 	brandMap := make(map[string]bool)

// 	for _, indoor := range e.IndoorUnits {
// 		brandMap[indoor.Brand] = true
// 	}
// 	for _, outdoor := range e.OutdoorUnits {
// 		brandMap[outdoor.Brand] = true
// 	}

// 	// Process each brand separately
// 	for brand := range brandMap {
// 		// Pre-allocate with estimated capacity
// 		brandIndoorUnits := make([]data_structures.IndoorUnit, 0, len(e.IndoorUnits)/len(brandMap))
// 		brandOutdoorUnits := make([]data_structures.OutdoorUnit, 0, len(e.OutdoorUnits)/len(brandMap))

// 		// Filter to only this brand's equipment
// 		for _, indoor := range e.IndoorUnits {
// 			if indoor.Brand == brand {
// 				brandIndoorUnits = append(brandIndoorUnits, indoor)
// 			}
// 		}

// 		for _, outdoor := range e.OutdoorUnits {
// 			if outdoor.Brand == brand {
// 				brandOutdoorUnits = append(brandOutdoorUnits, outdoor)
// 			}
// 		}

// 		// Generate combinations only within this brand
// 		for _, indoorUnit := range brandIndoorUnits {
// 			for _, outdoorUnit := range brandOutdoorUnits {
// 				equipConfig := data_structures.ComponentKey{
// 					Brand:       indoorUnit.Brand,
// 					IndoorUnit:  indoorUnit,
// 					OutdoorUnit: outdoorUnit,
// 				}
// 				equipConfigs = append(equipConfigs, equipConfig)
// 			}
// 		}
// 	}

// 	return equipConfigs
// }

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
	// fmt.Printf("Looking up: %v\n\n", key)
	ahriNumber, certified := ahriMap[key]
	return ahriNumber, certified
}

func FindCertifiedMatches(
	fullSystemCombos []data_structures.ComponentKey,
	ahriMap map[string]string,
) ([]data_structures.OutputCSV, error) {

	certifiedMatches := make([]data_structures.OutputCSV, 0)

	for _, combo := range fullSystemCombos {
		ahriNumber, isCertified := FindAHRICertification(combo, ahriMap)

		if !isCertified {
			continue
		}

		output := data_structures.OutputCSV{
			AHRINumber:   ahriNumber,
			Brand:        combo.Brand,
			Orientation:  "",
			TypeOfSystem: combo.SystemType,
		}

		switch combo.SystemType {
		case systemTypes["central ac"]:
			output.OutdoorUnit = combo.OutdoorUnit.InputModelNumber
			output.EvaporatorCoil = combo.IndoorUnit.InputModelNumber

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

		case systemTypes["furnace"]:
			output.Furnace = combo.Furnace.InputModelNumber

		default:
			return nil, fmt.Errorf("unknown system type: %s", combo.SystemType)
		}

		certifiedMatches = append(certifiedMatches, output)
	}

	return certifiedMatches, nil
}
