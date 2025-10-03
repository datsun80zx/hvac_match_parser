package objects

// EquipmentSet represents a collection of equipment model numbers
type EquipmentSet map[string]bool

// Match represents an AHRI match record with all its data columns
type Match struct {
	Data []string
}

// ProcessingResult holds the results of the matching process
type ProcessingResult struct {
	EquipmentCount int
	MatchCount     int
	Headers        []string
	Matches        []Match
}
