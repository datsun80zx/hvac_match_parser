package data_structures

type AHRIRecord struct {
	AHRINumber  string
	OutdoorUnit string
	IndoorUnit  string
	Furnace     string
}

type ComponentKey struct {
	Brand       string
	Furnace     string
	IndoorUnit  string
	OutdoorUnit string
}
