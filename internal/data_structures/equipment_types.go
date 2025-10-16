package data_structures

type Furnace struct {
	InputModelNumber      string
	NormalizedModelNumber string
	EquipmentType         string
	Brand                 string
}

type OutdoorUnit struct {
	InputModelNumber      string
	NormalizedModelNumber string
	HeatPump              bool
	Brand                 string
}

type IndoorUnit struct {
	InputModelNumber      string
	NormalizedModelNumber string
	AirHandler            bool
	Brand                 string
}

type AHRIRecord struct {
	AHRINumber  string
	OutdoorUnit string
	IndoorUnit  string
	Furnace     string
}

type ComponentKey struct {
	OutdoorUnit OutdoorUnit
	IndoorUnit  IndoorUnit
	Furnace     Furnace
}

type OutputCSV struct {
	TypeOfSystem   string
	OutdoorUnit    string
	Furnace        string
	EvaporatorCoil string
	AirHandler     string
}

type Equipment struct {
	Furnaces     []Furnace
	OutdoorUnits []OutdoorUnit
	IndoorUnits  []IndoorUnit
}
