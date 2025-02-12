package models

type System struct {
	Symbol       string     `json:"symbol"`
	SectorSymbol string     `json:"sectorSymbol"`
	Type         string     `json:"type"`
	X            int        `json:"x"`
	Y            int        `json:"y"`
	Waypoints    []Waypoint `json:"waypoints"`
	Factions     []Faction  `json:"factions"`
}

type Waypoint struct {
	Symbol   string           `json:"symbol"`
	Type     string           `json:"type"`
	X        int              `json:"x"`
	Y        int              `json:"y"`
	Orbitals []Orbital        `json:"orbitals"`
	Traits   []WaypointTraits `json:"traits"`
}

type Orbital struct {
	Symbol string `json:"symbol"`
	Orbits string `json:"orbits,omitempty"`
}

type ModifierType string

const (
	Stripped      ModifierType = "STRIPPED"
	Unstable      ModifierType = "UNSTABLE"
	RadiationLeak ModifierType = "RADIATION_LEAK"
	CriticalLimit ModifierType = "CRITICAL_LIMIT"
	CivilUnrest   ModifierType = "CIVIL_UNREST"
)

type Modifier struct {
	Symbol      ModifierType `json:"symbol"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
}
