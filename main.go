package main

import (
	"fmt"
	"log"

	"github.com/datsun80zx/hvac_match_parser/internal"
)

func main() {
	// Define the paths to your input CSV files
	csvFileEquip := "C:/Users/mrich/Downloads/wilson_equipment _list.csv"
	csvFileAHRI := "C:/Users/mrich/Downloads/ahri_matches.csv"

	// Step 1: Read and parse the equipment list
	// This file contains all the individual components (furnaces, indoor units, outdoor units)
	// that we want to check for AHRI certification
	fmt.Println("Reading equipment list...")
	equipmentList, err := internal.CSVEquipReader(csvFileEquip)
	if err != nil {
		log.Fatalf("Failed to read equipment CSV file: %v", err)
	}
	fmt.Printf("Loaded %d furnaces, %d indoor units, %d outdoor units\n",
		len(equipmentList.Furnaces),
		len(equipmentList.IndoorUnits),
		len(equipmentList.OutdoorUnits))

	// Step 2: Read and parse the AHRI certified matches database
	// This file contains all the equipment combinations that have been officially certified
	fmt.Println("\nReading AHRI certification database...")
	ahriList, err := internal.CSVAHRIReader(csvFileAHRI)
	if err != nil {
		log.Fatalf("Failed to read AHRI CSV file: %v", err)
	}
	fmt.Printf("Loaded %d AHRI certified combinations\n", len(ahriList))

	// Step 3: Normalize all equipment model numbers
	// Normalization ensures that model numbers from different sources can be compared
	// consistently by truncating them to standard lengths and applying format rules
	fmt.Println("\nNormalizing equipment model numbers...")

	// Normalize furnace model numbers
	for i := range equipmentList.Furnaces {
		equipmentList.Furnaces[i].NormalizedModelNumber =
			internal.NormalizeString(equipmentList.Furnaces[i].InputModelNumber, "furnace")
	}

	// Normalize outdoor unit (condenser/heat pump) model numbers
	for i := range equipmentList.OutdoorUnits {
		equipmentList.OutdoorUnits[i].NormalizedModelNumber =
			internal.NormalizeString(equipmentList.OutdoorUnits[i].InputModelNumber, "outdoor")
	}

	// Normalize indoor unit model numbers (both air handlers and coils)
	// We need to check what type of indoor unit it is to apply the correct normalization
	for i := range equipmentList.IndoorUnits {
		if equipmentList.IndoorUnits[i].AirHandler {
			equipmentList.IndoorUnits[i].NormalizedModelNumber =
				internal.NormalizeString(equipmentList.IndoorUnits[i].InputModelNumber, "airhandler")
		} else {
			equipmentList.IndoorUnits[i].NormalizedModelNumber =
				internal.NormalizeString(equipmentList.IndoorUnits[i].InputModelNumber, "coil")
		}
	}

	fmt.Println("Equipment normalization complete")

	// Step 4: Normalize AHRI database model numbers
	// The AHRI data needs the same normalization so we can match it against our equipment
	fmt.Println("\nNormalizing AHRI database model numbers...")
	for i := range ahriList {
		// Normalize each component in the AHRI record
		ahriList[i].Furnace = internal.NormalizeString(ahriList[i].Furnace, "furnace")
		ahriList[i].OutdoorUnit = internal.NormalizeString(ahriList[i].OutdoorUnit, "outdoor")
		// For indoor units, we use "indoor" as a general type - the normalization
		// function will handle it appropriately based on the model number format
		ahriList[i].IndoorUnit = internal.NormalizeString(ahriList[i].IndoorUnit, "indoor")
	}
	fmt.Println("AHRI database normalization complete")

	// Step 5: Generate all possible equipment combinations
	// This creates the Cartesian product of all components to see every possible
	// way the equipment could be combined into complete HVAC systems
	fmt.Println("\nGenerating possible equipment combinations...")

	// Full systems include a furnace, indoor unit, and outdoor unit
	fullSystemCombinations := internal.GenerateFullSystemEquipmentConfig(equipmentList)
	fmt.Printf("Generated %d full system combinations (furnace + indoor + outdoor)\n",
		len(fullSystemCombinations))

	// Air handler systems don't include a furnace - just indoor and outdoor units
	airHandlerCombinations := internal.GenerateAirHandlerEquipmentConfig(equipmentList)
	fmt.Printf("Generated %d air handler combinations (indoor + outdoor)\n",
		len(airHandlerCombinations))

	totalCombinations := len(fullSystemCombinations) + len(airHandlerCombinations)
	fmt.Printf("Total combinations to check: %d\n", totalCombinations)

	// Step 6: Build the AHRI lookup map with wildcard expansion
	// This creates a fast hash map where we can instantly check if any combination
	// is certified. The wildcards in AHRI data are expanded to all their possible
	// concrete values so we can do exact string matching
	fmt.Println("\nBuilding AHRI certification lookup map...")
	ahriMap := internal.BuildAHRIMap(ahriList)
	fmt.Printf("AHRI map built with %d entries (includes wildcard expansions)\n", len(ahriMap))

	// Step 7: Find all certified matches
	// This checks each of our generated combinations against the AHRI database
	// and collects only those that are officially certified
	fmt.Println("\nSearching for certified equipment matches...")
	certifiedMatches := internal.FindCertifiedMatches(
		fullSystemCombinations,
		airHandlerCombinations,
		ahriMap,
	)

	// Report on what we found
	fmt.Printf("\nFound %d certified matches out of %d total combinations\n",
		len(certifiedMatches), totalCombinations)

	// Calculate and display the match rate as a percentage
	matchRate := float64(len(certifiedMatches)) / float64(totalCombinations) * 100
	fmt.Printf("Match rate: %.2f%%\n", matchRate)

	// Step 8: Write the certified matches to an output CSV file
	// This creates a formatted file that can be opened in Excel or used by other systems
	outputFilename := "certified_hvac_matches.csv"
	fmt.Printf("\nWriting certified matches to %s...\n", outputFilename)

	err = internal.WriteOutputCSV(certifiedMatches, outputFilename)
	if err != nil {
		log.Fatalf("Failed to write output CSV: %v", err)
	}

	// Success! Let the user know where to find their results
	fmt.Printf("\n✓ Complete! Certified matches have been written to %s\n", outputFilename)
	fmt.Println("\nYou can now open this file in Excel or any spreadsheet program to view your results.")
}
