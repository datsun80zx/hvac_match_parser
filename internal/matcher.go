package internal

import (
	"github.com/datsun80zx/hvac_match_parser/internal/data_structures"
)

func MatchMaker(e data_structures.Equipment) data_structures.ComponentKey {
	ckList := []data_structures.ComponentKey

	for _, oUnit := range e.OutdoorUnits {
		for _, iUnit := range e.IndoorUnits {
			if iUnit.AirHandler {
				ckList = append(ckList, oUnit.NormalizedModelNumber)
				ckList = append
			}
		}
	}
}
