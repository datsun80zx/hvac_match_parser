package main

import (
	"fmt"
	"log"

	"github.com/datsun80zx/hvac_match_parser/internal"
)

func main() {

	// testStrings := []string{
	// 	"asdfghjklsqldeyoualkdjfa;i",
	// 	"1234567890",
	// 	"askda;kdfl;",
	// 	"askd;afjskds",
	// }

	csvFileEquip := "C:/Users/mrich/Downloads/amana_equipment_list.csv"
	// csvFileAHRI := "C:/Users/mrich/Downloads/amana_ac_ahri_matches.csv"

	eqList, err := internal.CSVEquipReader(csvFileEquip)
	if err != nil {
		log.Printf("Something went wrong with reading csv file: %v\n", err)

	}
	fmt.Printf("Total number of furnaces: %v", len(eqList.Furnaces))
	for _, item := range eqList.Furnaces {
		item.NormalizedModelNumber = internal.NormalizeString(item.InputModelNumber, "furnace")
		fmt.Printf("\nOriginal Model #: %v\nNormalized Model #: %v", item.InputModelNumber, item.NormalizedModelNumber)
	}
	fmt.Printf("Total number of indoor units: %v", len(eqList.IndoorUnits))
	for _, item := range eqList.IndoorUnits {
		if item.AirHandler {
			item.NormalizedModelNumber = internal.NormalizeString(item.InputModelNumber, "airhandler")
			fmt.Printf("\nOriginal Model #: %v\nNormalized Model #: %v", item.InputModelNumber, item.NormalizedModelNumber)
		} else {
			item.NormalizedModelNumber = internal.NormalizeString(item.InputModelNumber, "coil")
			fmt.Printf("\nOriginal Model #: %v\nNormalized Model #: %v", item.InputModelNumber, item.NormalizedModelNumber)
		}
	}
	fmt.Printf("Total number of outdoor units: %v", len(eqList.OutdoorUnits))
	for _, item := range eqList.OutdoorUnits {
		item.NormalizedModelNumber = internal.NormalizeString(item.InputModelNumber, "outdoor")
		fmt.Printf("\nOriginal Model #: %v\nNormalized Model #: %v", item.InputModelNumber, item.NormalizedModelNumber)
	}

	// ahriList, err := internal.CSVAHRIReader(csvFileAHRI)
	// if err != nil {
	// 	log.Printf("Something went wrong with reading ahri csv file: %v", err)
	// }

	// for _, item := range ahriList {
	// 	fmt.Println(item)
	// }
	// fmt.Printf("Total number of ahri matches: %v", len(ahriList))

	// val := true
	// for idx, str := range testStrings {
	// fmt.Printf("testing list item: %v\n", idx)
	// fmt.Printf("length before normalization: %v\nstring contents: %v\n", len(str), str)
	// nStr := internal.NormalizeString(str, val)
	// fmt.Printf("length after normalization: %v\nstring contents: %v\n", len(nStr), nStr)
	// }
}
