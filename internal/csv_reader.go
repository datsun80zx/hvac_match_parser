package internal

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"

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

func CSVEquipReader(s string) (data_structures.Equipment, error) {
	file, err := os.Open(s)
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	r := csv.NewReader(file)

	_, err = r.Read()
	if err != nil {
		log.Printf("Error reading header: %v", err)
		return data_structures.Equipment{}, err
	}

	var furnaceList []data_structures.Furnace
	var outdoorList []data_structures.OutdoorUnit
	var indoorList []data_structures.IndoorUnit
	eq := data_structures.Equipment{}

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
		if record[1] != "" {
			furnaceList = append(furnaceList, data_structures.Furnace{
				InputModelNumber: record[1],
				Brand:            record[0],
			})
		}
		if record[2] != "" {
			outdoorList = append(outdoorList, data_structures.OutdoorUnit{
				InputModelNumber: record[2],
				Brand:            record[0],
				HeatPump:         false,
			})
		}
		if record[3] != "" {
			outdoorList = append(outdoorList, data_structures.OutdoorUnit{
				InputModelNumber: record[3],
				Brand:            record[0],
				HeatPump:         true,
			})
		}
		if record[4] != "" {
			indoorList = append(indoorList, data_structures.IndoorUnit{
				InputModelNumber: record[4],
				Brand:            record[0],
				AirHandler:       false,
			})
		}
		if record[5] != "" {
			indoorList = append(indoorList, data_structures.IndoorUnit{
				InputModelNumber: record[5],
				Brand:            record[0],
				AirHandler:       true,
			})
		}

	}

	eq.Furnaces = furnaceList
	eq.IndoorUnits = indoorList
	eq.OutdoorUnits = outdoorList

	return eq, nil
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
			AHRINumber:  record[0],
			OutdoorUnit: record[1],
			IndoorUnit:  record[2],
			Furnace:     record[3],
		})
	}
	return AHRIList, nil
}
