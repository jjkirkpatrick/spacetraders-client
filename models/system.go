package models

type System struct {
	Symbol       string     `json:"symbol"`
	SectorSymbol string     `json:"sectorSymbol"`
	Type         string     `json:"type"`
	X            int        `json:"x"`
	Y            int        `json:"y"`
	Waypoints    []Waypoint `json:"waypoints"`
}

type Waypoint struct {
	Symbol   string    `json:"symbol"`
	Type     string    `json:"type"`
	X        int       `json:"x"`
	Y        int       `json:"y"`
	Orbitals []Orbital `json:"orbitals"`
}

type Orbital struct {
	Symbol string `json:"symbol"`
	Orbits string `json:"orbits,omitempty"`
}
