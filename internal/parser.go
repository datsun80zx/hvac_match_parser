package internal

func NormalizeString(e Equipment) {
	const maxFL = 11
	const maxCL = 10
	model := e.GetModelNumber()

	if e.EquipmentType == "furnace" {
		if len(model) <= maxFL {
			e.NormalizedModelNumber = model
		}
		e.NormalizedModelNumber = model[:maxFL]
	}

	if e.EquipmentType == "coil" {
		if len(model) <= maxCL {
			e.NormalizedModelNumber = model
		}
		e.NormalizedModelNumber = model[:maxCL]
	}
}
