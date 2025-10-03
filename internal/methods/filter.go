package methods

import (
	"strings"

	"github.com/datsun80zx/hvac_match_parser/internal/objects"
)

// matchesEquipment checks if all specified columns contain equipment from our list
func matchesEquipment(record []string, equipmentSet objects.EquipmentSet, columns []int) bool {
	for _, colIdx := range columns {
		if colIdx < len(record) {
			model := strings.TrimSpace(record[colIdx])
			if model != "" && !equipmentSet[strings.ToUpper(model)] {
				return false
			}
		}
	}
	return true
}

// FilterMatches filters a slice of matches based on equipment set
func FilterMatches(matches []objects.Match, equipmentSet objects.EquipmentSet, columns []int) []objects.Match {
	filtered := make([]objects.Match, 0, len(matches))

	for _, match := range matches {
		if matchesEquipment(match.Data, equipmentSet, columns) {
			filtered = append(filtered, match)
		}
	}

	return filtered
}
