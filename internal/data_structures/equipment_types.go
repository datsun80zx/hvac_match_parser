package data_structures

type Equipment interface {
	GetModelNumber() string
	CheckType() string
}

type Furnace struct {
	InputModelNumber      string
	NormalizedModelNumber string
	EquipmentType         string
}

func (F Furnace) GetModelNumber() string {
	return F.InputModelNumber
}

func (F Furnace) CheckType() string {
	return F.EquipmentType
}

type OutdoorUnit struct {
	InputModelNumber      string
	NormalizedModelNumber string
	HeatPump              bool
	EquipmentType         string
}

func (O OutdoorUnit) GetModelNumber() string {
	return O.InputModelNumber
}

func (O OutdoorUnit) CheckType() string {
	return O.EquipmentType
}

type IndoorUnit struct {
	InputModelNumber      string
	NormalizedModelNumber string
	AirHandler            bool
	EquipmentType         string
}

func (I IndoorUnit) GetModelNumber() string {
	return I.InputModelNumber
}

func (I IndoorUnit) CheckType() string {
	return I.EquipmentType
}

type AHRIRecord struct {
	AHRINumber  string
	OutdoorUnit OutdoorUnit
	IndoorUnit  IndoorUnit
	Furnace     Furnace
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
