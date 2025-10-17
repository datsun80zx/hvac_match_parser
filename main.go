package main

import (
	"fmt"
	"log"

	"github.com/datsun80zx/hvac_match_parser/internal"
)

func main() {

	csvFileEquip := "C:/Users/mrich/Downloads/amana_equipment_list.csv"
	csvFileAHRI := "C:/Users/mrich/Downloads/amana_ac_ahri_matches.csv"

	fmt.Println("Now working on equipment list...")
	eqList, err := internal.CSVEquipReader(csvFileEquip)
	if err != nil {
		log.Printf("Something went wrong with reading csv file: %v\n", err)

	}

	fmt.Println("Now working on furnaces...")
	fmt.Printf("Total number of furnaces: %v\n\n", len(eqList.Furnaces))
	for _, item := range eqList.Furnaces[:min(3, len(eqList.Furnaces))] {
		item.NormalizedModelNumber = internal.NormalizeString(item.InputModelNumber, "furnace")
		fmt.Printf("\nOriginal Model #: %v\nNormalized Model #: %v\n", item.InputModelNumber, item.NormalizedModelNumber)
	}

	fmt.Println("Now working on indoor units...")
	fmt.Printf("Total number of indoor units: %v\n\n", len(eqList.IndoorUnits))
	for _, item := range eqList.IndoorUnits[:min(3, len(eqList.IndoorUnits))] {
		if item.AirHandler {
			fmt.Println("This is an air handler...")
			item.NormalizedModelNumber = internal.NormalizeString(item.InputModelNumber, "airhandler")
			fmt.Printf("\nOriginal Model #: %v\nNormalized Model #: %v\n\n", item.InputModelNumber, item.NormalizedModelNumber)
		} else {
			fmt.Println("This is an evap coil...")
			item.NormalizedModelNumber = internal.NormalizeString(item.InputModelNumber, "coil")
			fmt.Printf("\nOriginal Model #: %v\nNormalized Model #: %v\n\n", item.InputModelNumber, item.NormalizedModelNumber)
		}
	}

	fmt.Println("Now working on outdoor units...")
	fmt.Printf("Total number of outdoor units: %v\n\n", len(eqList.OutdoorUnits))
	for _, item := range eqList.OutdoorUnits[:min(3, len(eqList.OutdoorUnits))] {
		item.NormalizedModelNumber = internal.NormalizeString(item.InputModelNumber, "outdoor")
		fmt.Printf("\nOriginal Model #: %v\nNormalized Model #: %v\n\n", item.InputModelNumber, item.NormalizedModelNumber)
	}

	fmt.Println("Now working on ahri list...")
	ahriList, err := internal.CSVAHRIReader(csvFileAHRI)
	if err != nil {
		log.Printf("Something went wrong with reading ahri csv file: %v", err)
	}

	fmt.Println("Now printing ahri list...")
	for _, item := range ahriList[:min(3, len(ahriList))] {
		fmt.Println("This is a Furnace...")
		f := internal.NormalizeString(item.Furnace, "furnace")
		fmt.Printf("\nOriginal Model #: %v\nNormalized Model#: %v\n\n", item.Furnace, f)

		fmt.Println("This is an indoor unit...")
		in := internal.NormalizeString(item.IndoorUnit, "indoor")
		fmt.Printf("\nOriginal Model #: %v\nNormalized Model #: %v\n\n", item.IndoorUnit, in)

		fmt.Println("This is an outdoor unit...")
		out := internal.NormalizeString(item.OutdoorUnit, "outdoor")
		fmt.Printf("\nOriginal Model #: %v\nNormalized Model #: %v\n\n", item.OutdoorUnit, out)
	}
	fmt.Printf("Total number of ahri matches: %v\n\n", len(ahriList))

}
