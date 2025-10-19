package data_structures

type AHRIRecord struct {
	AHRINumber  string
	OutdoorUnit Equipment
	IndoorUnit  Equipment
	Furnace     Equipment
}

type ComponentKey struct {
	Brand       string
	Furnace     Equipment
	IndoorUnit  Equipment
	OutdoorUnit Equipment
	SystemType  string
}
