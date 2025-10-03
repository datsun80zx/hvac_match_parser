package methods

import (
	"encoding/csv"
	"os"

	"github.com/datsun80zx/hvac_match_parser/internal/objects"
)

// WriteOutput writes the filtered matches to the output file
func WriteOutput(filename string, headers []string, matches []objects.Match) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	if err := writer.Write(headers); err != nil {
		return err
	}

	// Write match data
	for _, match := range matches {
		record := prepareRecord(match.Data, len(headers))
		if err := writer.Write(record); err != nil {
			return err
		}
	}

	return nil
}

// prepareRecord ensures the record matches the header length
func prepareRecord(data []string, headerLength int) []string {
	record := make([]string, headerLength)
	for i := 0; i < headerLength && i < len(data); i++ {
		record[i] = data[i]
	}
	return record
}
