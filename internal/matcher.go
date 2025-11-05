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
	comAirHandlers := []data_structures.Equipment{}
	coils := []data_structures.Equipment{}
	comCoils := []data_structures.Equipment{}
	airCons := []data_structures.Equipment{}
	comAirCons := []data_structures.Equipment{}
	heatPumps := []data_structures.Equipment{}
	comHeatPumps := []data_structures.Equipment{}

	for _, item := range list {
		if strings.Contains(item.Type, "furnace") {
			furnaces = append(furnaces, item)
		} else if strings.Contains(item.Type, "handler") {
			if strings.ContainsAny(strings.ToLower(item.NormalizedModelNumber), "ahve") {
				comAirHandlers = append(comAirHandlers, item)
			} else {
				airHandlers = append(airHandlers, item)
			}
		} else if strings.Contains(item.Type, "coil") {
			if strings.ContainsAny(strings.ToLower(item.NormalizedModelNumber), "capea") {
				comCoils = append(comCoils, item)
			} else {
				coils = append(coils, item)
			}

		} else if strings.Contains(item.Type, "ac") {
			if strings.ContainsAny(strings.ToLower(item.NormalizedModelNumber), "asxv9") {
				comAirCons = append(comAirCons, item)
			} else if strings.ContainsAny(strings.ToLower(item.NormalizedModelNumber), "axv6") {
				comAirCons = append(comAirCons, item)
			} else {
				airCons = append(airCons, item)
			}
		} else if strings.Contains(item.Type, "hp") {
			if strings.ContainsAny(strings.ToLower(item.NormalizedModelNumber), "aszv9") {
				comHeatPumps = append(comHeatPumps, item)
			} else if strings.ContainsAny(strings.ToLower(item.NormalizedModelNumber), "azv6") {
				comHeatPumps = append(comHeatPumps, item)
			} else {
				heatPumps = append(heatPumps, item)
			}
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
		for _, comHeatPump := range comHeatPumps {
			for _, comAirHandler := range comAirHandlers {
				equipCombo := data_structures.ComponentKey{
					Brand:       comHeatPump.Brand,
					IndoorUnit:  comAirHandler,
					OutdoorUnit: comHeatPump,
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
		for _, comHeatPump := range comHeatPumps {
			for _, furnace := range furnaces {
				for _, comCoil := range comCoils {
					equipCombo := data_structures.ComponentKey{
						Brand:       comHeatPump.Brand,
						IndoorUnit:  comCoil,
						Furnace:     furnace,
						OutdoorUnit: comHeatPump,
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
		for _, comAirCon := range comAirCons {
			for _, furnace := range furnaces {
				for _, comCoil := range comCoils {
					equipCombo := data_structures.ComponentKey{
						Brand:       comAirCon.Brand,
						IndoorUnit:  comCoil,
						Furnace:     furnace,
						OutdoorUnit: comAirCon,
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
		for _, comAirCon := range comAirCons {
			for _, comAirHandler := range comAirHandlers {
				equipCombo := data_structures.ComponentKey{
					Brand:       comAirCon.Brand,
					IndoorUnit:  comAirHandler,
					OutdoorUnit: comAirCon,
					SystemType:  systemTypes["central ac & air handler"],
				}
				equipConfigs = append(equipConfigs, equipCombo)
			}
		}
		return equipConfigs, nil
	case "furnace":
		for _, furnace := range furnaces {
			equipCombo := data_structures.ComponentKey{
				Brand:      furnace.Brand,
				Furnace:    furnace,
				SystemType: systemTypes["furnace"],
			}
			equipConfigs = append(equipConfigs, equipCombo)
		}
		return equipConfigs, nil
	case "central ac":
		for _, airCon := range airCons {
			for _, coil := range coils {
				equipCombo := data_structures.ComponentKey{
					Brand:       airCon.Brand,
					IndoorUnit:  coil,
					OutdoorUnit: airCon,
					SystemType:  systemTypes["central ac"],
				}
				equipConfigs = append(equipConfigs, equipCombo)
			}
		}
		for _, comAirCon := range comAirCons {
			for _, comCoil := range comCoils {
				equipCombo := data_structures.ComponentKey{
					Brand:       comAirCon.Brand,
					IndoorUnit:  comCoil,
					OutdoorUnit: comAirCon,
					SystemType:  systemTypes["central ac"],
				}
				equipConfigs = append(equipConfigs, equipCombo)
			}
		}
		return equipConfigs, nil
	}
	return nil, fmt.Errorf("there was an error processing all equipment configurations")
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

	return []string{string(variation1), string(variation2)} // this option is if we are doing both upflow and downflow
	// return []string{string(variation2)}
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
