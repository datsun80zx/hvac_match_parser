package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/datsun80zx/hvac_match_parser/internal"
	"github.com/datsun80zx/hvac_match_parser/internal/data_structures"
)

func main() {
	// Define the paths to your input CSV files
	csvFileEquip := "C:/Users/mrich/dev_work/hvac_match_parser/data/wilson_equip_list.csv"
	csvFileAHRI := "C:/Users/mrich/dev_work/hvac_match_parser/data/ahri_matches.csv"

	// Define what column headers we are expecting to see in the equipment list csv:

	equipmentFields := []string{
		"Brand",
		"Furnace",
		"Outdoor Unit (ac)",
		"Outdoor Unit (hp)",
		"Evaporator Coil",
		"Air Handler",
	}

	fmt.Printf("Reading equipment headers...\n\n")
	equipHeaders, err := internal.GetCSVHeader(csvFileEquip, equipmentFields)
	if err != nil {
		log.Fatalf("Failed to read equipment csv headers: %v", err)
	}

	for header, idx := range equipHeaders {
		fmt.Printf("Equipment Header %d: %s\n\n", idx, header)
	}

	// read and parse the equipment list csv starting after the headers have already been read:

	fmt.Printf("\nReading equipment list...\n\n")
	equipmentList, err := internal.CSVEquipReader(csvFileEquip, equipHeaders)
	if err != nil {
		log.Fatalf("failed to read equipment csv file: %v", err)
	}
	fmt.Printf("Loaded %d pieces of equipment\n\n", len(equipmentList))

	// Figure out what different brands we are working with:

	fmt.Printf("Identifying brands...\n\n")
	brandMap := internal.BrandIdentify(equipmentList)
	fmt.Printf("==== Brands (%d) ====\n\n", len(brandMap))
	for k := range brandMap {
		fmt.Printf("%s\n", k)
	}

	// Normalize equipment:
	fmt.Printf("\nNormalizing equipment model #'s...\n\n")
	for i := range equipmentList {
		equipmentList[i] = internal.NormalizeString(equipmentList[i])
	}
	fmt.Printf("Equipment normalization complete!\n\n")
	fmt.Printf("Categorizing equipment (standard vs communicating)...\n\n")
	for i := range equipmentList {
		equipmentList[i] = internal.CategorizeEquipment(equipmentList[i])
	}
	fmt.Printf("Equipment categorization complete!\n\n")

	// Optional: Add some logging to show categorization results
	standardCount := 0
	communicatingCount := 0
	for _, equip := range equipmentList {
		if equip.Category == data_structures.CategoryStandard {
			standardCount++
		} else if equip.Category == data_structures.CategoryCommunicating {
			communicatingCount++
		}
	}
	fmt.Printf("Standard equipment: %d\n", standardCount)
	fmt.Printf("Communicating equipment: %d\n\n", communicatingCount)
	fmt.Printf("First 5 pieces:\n\n")

	for i := 0; i < 5; i++ {
		fmt.Printf("Equipment type: %v\nEquipment Input Model #: %v\nEquipment Normalized Model #: %v\nEquipment brand: %v\n\n\n",
			equipmentList[i].Type,
			equipmentList[i].InputModelNumber,
			equipmentList[i].NormalizedModelNumber,
			equipmentList[i].Brand)
	}

	// Read and parse ahri certified matches:
	fmt.Printf("Reading ahri certified matches...\n\n")
	ahriList, err := internal.CSVAHRIReader(csvFileAHRI)
	if err != nil {
		log.Fatalf("Failed to read ahri csv file: %v", err)
	}
	fmt.Printf("Loaded %d ahri records\n\n", len(ahriList))

	fmt.Printf("First 5 records:\n\n")
	for i := 0; i < 5; i++ {
		fmt.Printf("Outdoor Unit: \n%v\n\nIndoor Unit: \n%v\n\nFurnace: \n%v\n\nAHRI Number: %v\n\n\n",
			ahriList[i].OutdoorUnit,
			ahriList[i].IndoorUnit,
			ahriList[i].Furnace,
			ahriList[i].AHRINumber)
	}

	// Build the ahri lookup match for equipment config certification:

	fmt.Printf("Building ahri cert lookup map...\n\n")
	ahriMap := internal.BuildAHRIMap(ahriList)
	fmt.Printf("Built ahri map with %d entries (including wildcard expansions)\n\n", len(ahriMap))
	// for key, value := range ahriMap {
	// 	fmt.Printf("Key: %v\nValue: %v\n\n\n", key, value)
	// }

	// Process through each brand and system type separately:
	fmt.Printf("Generating equipment combo's and finding matches...\n\n")

	allCertifiedMatches := make([]data_structures.OutputCSV, 0)
	systemTypes := []string{
		"central ac",
		"furnace",
		"central ac & air handler",
		"central ac & furnace",
		"heat pump & air handler",
		"heat pump & furnace",
	}
	totalCombinations := 0

	for brand := range brandMap {
		fmt.Printf("Processing brand: %s\n\n", brand)

		brandEquipment := internal.EquipmentSort(equipmentList, brand)
		fmt.Printf("   Found %d pieces of equipment for %s\n\n", len(brandEquipment), brand)

		for _, sysType := range systemTypes {
			combo, err := internal.GenerateFullSystemEquipmentConfig(brandEquipment, sysType)
			if err != nil {
				log.Printf("   Warning: Error generating %s combinations for %s: %v", sysType, brand, err)
				continue
			}

			if len(combo) == 0 {
				continue
			}

			fmt.Printf("   Generated %d combinations for %s\n\n", len(combo), sysType)
			totalCombinations += len(combo)

			fmt.Printf("Number of combo's: %d\n", len(combo))
			fmt.Printf("First 5 combo's:\n\n")
			for i := 0; i < 5; i++ {
				fmt.Printf("Combo %d:\nOutdoor Unit: \n%v\n\nIndoor Unit: \n%v\n\nFurnace: \n%v\n\nSystem Type: %v\n\n\n",
					i+1,
					combo[i].OutdoorUnit,
					combo[i].IndoorUnit,
					combo[i].Furnace,
					combo[i].SystemType)
			}

			certifiedMatches, err := internal.FindCertifiedMatches(combo, ahriMap)
			if err != nil {
				log.Printf("   Warning: Error finding matches for %s: %v", sysType, err)
				continue
			}

			if len(certifiedMatches) >= 0 {
				fmt.Printf("   + Found %d certified matches for %s\n\n", len(certifiedMatches), sysType)
				allCertifiedMatches = append(allCertifiedMatches, certifiedMatches...)
			}
		}
	}
	// Generate report on results:
	separator := strings.Repeat("=", 60)
	fmt.Printf("\n%s\n", separator)
	fmt.Printf("SUMMARY\n")
	fmt.Printf("%s\n", separator)
	fmt.Printf("Total combinations checked: %d\n", totalCombinations)
	fmt.Printf("Total certified matches found: %d\n", len(allCertifiedMatches))

	if totalCombinations > 0 {
		matchRate := float64(len(allCertifiedMatches)) / float64(totalCombinations) * 100
		fmt.Printf("Match rate: %.2f%%\n", matchRate)
	}
	// Create final csv output:
	if len(allCertifiedMatches) > 0 {
		outputFilename := "C:/Users/mrich/OneDrive/Wilson/wilson_hvac_matches/certified_hvac_matches.csv"
		fmt.Printf("\nWriting certified matches to %s...\n\n", outputFilename)

		err = internal.WriteOutputCSV(allCertifiedMatches, outputFilename)
		if err != nil {
			log.Fatalf("Failed to write output csv: %v", err)
		}
		fmt.Printf("\n✓ Complete! Certified matches have been written to %s\n", outputFilename)
		fmt.Println("\nYou can now open this file in Excel or any spreadsheet program to view your results.")
	} else {
		fmt.Println("\nNo certified matches found. No output file generated.")
	}
}
