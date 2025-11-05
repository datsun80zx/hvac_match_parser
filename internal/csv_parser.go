package internal

import (
	"strings"

	"github.com/datsun80zx/hvac_match_parser/internal/data_structures"
)

/* To Normalize Model #'s:
Step 1: determine type of equipment
Step 2: truncate length of string depending on type of equipment
*/

func NormalizeString(equipment data_structures.Equipment) data_structures.Equipment {
	equipmentTypes := map[string]int{
		"air handler":     11,
		"evaporator coil": 11,
		"furnace":         11,
		"condenser(ac)":   11,
		"condenser(hp)":   11,
		"default":         11,
	}

	typeLower := strings.ToLower(equipment.Type)
	input := equipment.InputModelNumber

	truncate := func(s string, maxLen int) string {
		if len(s) >= maxLen {
			return s[:maxLen]
		}
		return s
	}

	if strings.Contains(typeLower, "coil") {
		maxLength := equipmentTypes["evaporator coil"]
		if len(input) > 0 && strings.ToLower(input)[0] != 'c' && len(input) >= 2+maxLength {
			equipment.NormalizedModelNumber = input[2:(2 + maxLength)]
			return equipment
		}
		equipment.NormalizedModelNumber = truncate(input, maxLength)
		return equipment
	}

	if strings.Contains(typeLower, "air handler") {
		maxLength := equipmentTypes["air handler"]
		equipment.NormalizedModelNumber = truncate(input, maxLength)
		return equipment
	}

	maxLength := equipmentTypes["default"]
	equipment.NormalizedModelNumber = truncate(input, maxLength)
	return equipment
}

func EquipmentSort(list []data_structures.Equipment, parameter string) []data_structures.Equipment {
	outputList := []data_structures.Equipment{}

	for _, item := range list {
		if item.Brand == parameter {
			outputList = append(outputList, item)
		}
	}
	return outputList
}

func BrandIdentify(list []data_structures.Equipment) map[string]bool {
	brandMap := make(map[string]bool)

	for _, item := range list {
		brandMap[item.Brand] = true
	}

	return brandMap
}

func CategorizeEquipment(equipment data_structures.Equipment) data_structures.Equipment {
	modelLower := strings.ToLower(equipment.NormalizedModelNumber)
	typeLower := strings.ToLower(equipment.Type)

	// Check for air handlers
	if strings.Contains(typeLower, "handler") || strings.Contains(typeLower, "air handler") {
		if strings.Contains(modelLower, "ahve") {
			equipment.Category = data_structures.CategoryCommunicating
			return equipment
		}
	}

	// Check for evaporator coils
	if strings.Contains(typeLower, "coil") {
		if strings.Contains(modelLower, "capea") {
			equipment.Category = data_structures.CategoryCommunicating
			return equipment
		}
	}

	// Check for AC outdoor units (both "condenser(ac)" and "outdoor unit (ac)")
	if strings.Contains(typeLower, "ac") {
		if strings.Contains(modelLower, "axv") || strings.Contains(modelLower, "gxv") {
			equipment.Category = data_structures.CategoryCommunicating
			return equipment
		}
	}

	// Check for heat pump outdoor units (both "condenser(hp)" and "outdoor unit (hp)")
	if strings.Contains(typeLower, "hp") {
		if strings.Contains(modelLower, "aszv9") || strings.Contains(modelLower, "azv6") || strings.Contains(modelLower, "gszv9") || strings.Contains(modelLower, "gzv6") {
			equipment.Category = data_structures.CategoryCommunicating
			return equipment
		}
	}

	// Default to standard
	equipment.Category = data_structures.CategoryStandard
	return equipment
}
