package internal

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/datsun80zx/hvac_match_parser/internal/data_structures"
)

func WriteOutputCSV(matches []data_structures.OutputCSV, filename string) error {
	// Create the output file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the header row
	header := []string{
		"AHRI Number",
		"Brand",
		"Orientation",
		"Type of System",
		"Outdoor Unit",
		"Furnace",
		"Evaporator Coil",
		"Air Handler",
	}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write each match as a row
	for _, match := range matches {
		row := []string{
			match.AHRINumber,
			match.Brand,
			match.Orientation,
			match.TypeOfSystem,
			match.OutdoorUnit,
			match.Furnace,
			match.EvaporatorCoil,
			match.AirHandler,
		}
		if err := writer.Write(row); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}

	// Check for any errors that occurred during writing
	if err := writer.Error(); err != nil {
		return fmt.Errorf("csv writer error: %w", err)
	}

	return nil
}
func GetCSVHeader(filename string, reqFields []string) (map[string]int, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("There was an error with opening %s: %w", filename, err)
	}
	defer file.Close()

	r := csv.NewReader(file)

	header, err := r.Read()
	if err != nil {
		log.Printf("Error reading header: %v", err)
		return nil, err
	}

	columnIndices := make(map[string]int)
	for i, columnName := range header {
		cleanName := strings.TrimSpace(strings.ToLower(columnName))
		columnIndices[cleanName] = i
	}

	for _, colName := range reqFields {
		normColName := strings.ToLower(strings.TrimSpace(colName))
		if _, exists := columnIndices[normColName]; !exists {
			return nil, fmt.Errorf("required column '%s' not found in csv header", colName)
		}
	}
	return columnIndices, nil
}

func CSVEquipReader(filename string, headers map[string]int) ([]data_structures.Equipment, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	r := csv.NewReader(file)

	_, err = r.Read()
	if err != nil {
		log.Printf("Error reading header: %v", err)
		return []data_structures.Equipment{}, err
	}

	equipmentList := []data_structures.Equipment{}
	brandIdx := headers["brand"]

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			if pe, ok := err.(*csv.ParseError); ok {
				log.Println("Bad Column: ", pe.Column)
				log.Println("Bad Line: ", pe.Line)
				log.Println("Error reported ", pe.Err)
				if pe.Err == csv.ErrFieldCount {
					continue
				}
			}
			return nil, fmt.Errorf("error reading CSV: %w", err)
		}

		if len(record) <= brandIdx {
			log.Printf("Skipping row with insufficient columns: %v", record)
			continue
		}

		brand := record[brandIdx]

		for k, v := range headers {
			if k == "brand" {
				continue
			}

			if v >= len(record) {
				continue
			}

			if record[v] != "" {
				equipmentList = append(equipmentList, data_structures.Equipment{
					InputModelNumber: record[v],
					Brand:            brand,
					Type:             k,
				})
			}
		}
	}
	return equipmentList, nil
}

func CSVAHRIReader(s string) ([]data_structures.AHRIRecord, error) {
	file, err := os.Open(s)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	r := csv.NewReader(file)

	_, err = r.Read()
	if err != nil {
		log.Printf("Error reading header: %v", err)
		return []data_structures.AHRIRecord{}, err
	}

	var AHRIList []data_structures.AHRIRecord

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			if pe, ok := err.(*csv.ParseError); ok {
				log.Println("Bad Column: ", pe.Column)
				log.Println("Bad Line: ", pe.Line)
				log.Println("Error reported ", pe.Err)
				if pe.Err == csv.ErrFieldCount {
					continue
				}
			}
		}
		AHRIList = append(AHRIList, data_structures.AHRIRecord{
			AHRINumber: record[0],
			OutdoorUnit: data_structures.Equipment{
				InputModelNumber: record[1],
			},
			IndoorUnit: data_structures.Equipment{
				InputModelNumber: record[2],
			},
			Furnace: data_structures.Equipment{
				InputModelNumber: record[3],
			},
		})
	}
	return AHRIList, nil
}
