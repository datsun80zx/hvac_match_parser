package internal

import (
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
				IndoorUnit:  indoorUnit,
				OutdoorUnit: outdoorUnit,
			}
			equipConfigs = append(equipConfigs, equipConfig)
		}
	}

	return equipConfigs
}
