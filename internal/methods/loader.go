package methods

import (
	"encoding/csv"
	"os"
	"strings"

	"github.com/datsun80zx/hvac_match_parser/internal/objects"
)

// LoadEquipmentList reads the equipment CSV and returns a set of model numbers
func LoadEquipmentList(filename string, column int) (objects.EquipmentSet, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	equipmentSet := make(objects.EquipmentSet)

	// Skip header row and process data
	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}
		if column < len(record) {
			model := strings.TrimSpace(record[column])
			if model != "" {
				equipmentSet[strings.ToUpper(model)] = true
			}
		}
	}

	return equipmentSet, nil
}

// LoadExampleHeaders reads the header row from the example file
func LoadExampleHeaders(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	return headers, nil
}

// LoadAndFilterAHRIMatches loads AHRI matches and filters by equipment list
func LoadAndFilterAHRIMatches(filename string, equipmentSet objects.EquipmentSet, columns []int) ([]objects.Match, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var matches []objects.Match

	// Skip header row and process data
	for i, record := range records {
		if i == 0 {
			continue // Skip header
		}

		// Filter records using the filter function
		if matchesEquipment(record, equipmentSet, columns) {
			matches = append(matches, objects.Match{Data: record})
		}
	}

	return matches, nil
}
