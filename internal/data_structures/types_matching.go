package data_structures

type AHRIRecord struct {
	AHRINumber  string
	OutdoorUnit string
	IndoorUnit  string
	Furnace     string
}

type ComponentKey struct {
	Brand       string
	Furnace     Equipment
	IndoorUnit  Equipment
	OutdoorUnit Equipment
	SystemType  string
}
