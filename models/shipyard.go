package models

type ShipType string

const (
	ShipProbe             ShipType = "SHIP_PROBE"
	ShipMiningDrone       ShipType = "SHIP_MINING_DRONE"
	ShipSiphonDrone       ShipType = "SHIP_SIPHON_DRONE"
	ShipInterceptor       ShipType = "SHIP_INTERCEPTOR"
	ShipLightHauler       ShipType = "SHIP_LIGHT_HAULER"
	ShipCommandFrigate    ShipType = "SHIP_COMMAND_FRIGATE"
	ShipExplorer          ShipType = "SHIP_EXPLORER"
	ShipHeavyFreighter    ShipType = "SHIP_HEAVY_FREIGHTER"
	ShipLightShuttle      ShipType = "SHIP_LIGHT_SHUTTLE"
	ShipOreHound          ShipType = "SHIP_ORE_HOUND"
	ShipRefiningFreighter ShipType = "SHIP_REFINING_FREIGHTER"
	ShipSurveyor          ShipType = "SHIP_SURVEYOR"
)

type Shipyard struct {
	Symbol        string        `json:"symbol" required:"true"`
	ShipTypes     []ShipType    `json:"shipTypes" required:"true"`
	Transactions  []Transaction `json:"transactions"`
	Ships         []Ship        `json:"ships"`
	Modifications int           `json:"monificationFee"`
}

type Transaction struct {
	WaypointSymbol string   `json:"waypointSymbol" required:"true"`
	ShipSymbol     string   `json:"shipSymbol" required:"true" deprecated:"true"`
	ShipType       ShipType `json:"shipType" required:"true"`
	Price          int      `json:"price" required:"true"`
	AgentSymbol    string   `json:"agentSymbol" required:"true"`
	Timestamp      string   `json:"timestamp" required:"true" format:"date-time"`
}

type ShipyardShip struct {
	Type             ShipType `json:"type" required:"true"`
	Name             string   `json:"name" required:"true"`
	Description      string   `json:"description" required:"true"`
	Supply           string   `json:"supply" required:"true"`
	Activity         string   `json:"activity"`
	PurchasePrice    int      `json:"purchasePrice" required:"true"`
	Frame            Frame    `json:"frame" required:"true"`
	Reactor          Reactor  `json:"reactor" required:"true"`
	Engine           Engine   `json:"engine" required:"true"`
	Modules          []Module `json:"modules" required:"true"`
	Mounts           []Mount  `json:"mounts" required:"true"`
	Crew             Crew     `json:"crew" required:"true"`
	ModificationsFee int      `json:"modificationsFee" required:"true"`
}

type Frame struct {
	Symbol         string  `json:"symbol" required:"true"`
	Name           string  `json:"name" required:"true"`
	Description    string  `json:"description" required:"true"`
	Condition      float64 `json:"condition" required:"true"`
	Integrity      float64 `json:"integrity" required:"true"`
	ModuleSlots    int     `json:"moduleSlots" required:"true"`
	MountingPoints int     `json:"mountingPoints" required:"true"`
	FuelCapacity   int     `json:"fuelCapacity" required:"true"`
}

type Reactor struct {
	Symbol      string  `json:"symbol" required:"true"`
	Name        string  `json:"name" required:"true"`
	Description string  `json:"description" required:"true"`
	Condition   float64 `json:"condition" required:"true"`
	Integrity   float64 `json:"integrity" required:"true"`
	PowerOutput int     `json:"powerOutput" required:"true"`
}

type Engine struct {
	Symbol      string  `json:"symbol" required:"true"`
	Name        string  `json:"name" required:"true"`
	Description string  `json:"description" required:"true"`
	Condition   float64 `json:"condition" required:"true"`
	Integrity   float64 `json:"integrity" required:"true"`
	Speed       int     `json:"speed" required:"true"`
}

type Module struct {
	Symbol      string `json:"symbol" required:"true"`
	Capacity    int    `json:"capacity" required:"true"`
	Range       int    `json:"range" required:"true"`
	Name        string `json:"name" required:"true"`
	Description string `json:"description" required:"true"`
}

type Mount struct {
	Symbol      string   `json:"symbol" required:"true"`
	Name        string   `json:"name" required:"true"`
	Description string   `json:"description"`
	Strength    int      `json:"strength" required:"true"`
	Deposits    []string `json:"deposits"`
}

type Crew struct {
	Required int `json:"required" required:"true"`
	Capacity int `json:"capacity" required:"true"`
}
