package data_structures

type Equipment struct {
	InputModelNumber      string
	NormalizedModelNumber string
	Brand                 string
	Type                  string
}

const (
	TypeFurnace     = "furnace"
	TypeACCondenser = "ac condenser"
	TypeHeatPump    = "heat pump"
	TypeEvapCoil    = "evaporator coil"
	TypeAirHandler  = "air handler"
)
