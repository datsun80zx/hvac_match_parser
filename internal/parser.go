package internal

import "strings"

/* To Normalize Model #'s:
Step 1: determine type of equipment
Step 2: truncate length of string depending on type of equipment
*/

func NormalizeString(m string, equipmentType string) string {
	// Max Equipment Model # Lengths:
	const ahl = 11
	const othEq = 11

	switch equipmentType {
	case "airhandler":
		if len(m) >= ahl {
			return m[:ahl]
		}
	case "coil":
		if len(m) > 0 && strings.ToLower(m)[0] != 'c' {
			if len(m) > 2 {
				return m[2:]
			}
		}
		if len(m) >= othEq {
			return m[:othEq]
		}

	default:
		if len(m) >= othEq {
			return m[:othEq]
		}
	}
	return m
}
