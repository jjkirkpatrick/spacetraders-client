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
	Symbol    string `json:"symbol" `
	ShipTypes []struct {
		Type ShipType `json:"type" `
	} `json:"shipTypes" `
	Transactions  []Transaction `json:"transactions"`
	Ships         []Ship        `json:"ships"`
	Modifications int           `json:"monificationFee"`
}

type Transaction struct {
	WaypointSymbol string `json:"waypointSymbol"`
	ShipSymbol     string `json:"shipSymbol"`
	TradeSymbol    string `json:"tradeSymbol"`
	Type           string `json:"type" enum:"PURCHASE,SELL"`
	Units          int    `json:"units"`
	PricePerUnit   int    `json:"pricePerUnit"`
	TotalPrice     int    `json:"totalPrice"`
	Timestamp      string `json:"timestamp"`
}

type ShipyardShip struct {
	Type             ShipType `json:"type" `
	Name             string   `json:"name" `
	Description      string   `json:"description" `
	Supply           string   `json:"supply" `
	Activity         string   `json:"activity"`
	PurchasePrice    int      `json:"purchasePrice" `
	Frame            Frame    `json:"frame" `
	Reactor          Reactor  `json:"reactor" `
	Engine           Engine   `json:"engine" `
	Modules          []Module `json:"modules" `
	Mounts           []Mount  `json:"mounts" `
	Crew             Crew     `json:"crew" `
	ModificationsFee int      `json:"modificationsFee" `
}

type Frame struct {
	Symbol         string  `json:"symbol" `
	Name           string  `json:"name" `
	Description    string  `json:"description" `
	Condition      float64 `json:"condition" `
	Integrity      float64 `json:"integrity" `
	ModuleSlots    int     `json:"moduleSlots" `
	MountingPoints int     `json:"mountingPoints" `
	FuelCapacity   int     `json:"fuelCapacity" `
}

type Reactor struct {
	Symbol      string  `json:"symbol" `
	Name        string  `json:"name" `
	Description string  `json:"description" `
	Condition   float64 `json:"condition" `
	Integrity   float64 `json:"integrity" `
	PowerOutput int     `json:"powerOutput" `
}

type Engine struct {
	Symbol      string  `json:"symbol" `
	Name        string  `json:"name" `
	Description string  `json:"description" `
	Condition   float64 `json:"condition" `
	Integrity   float64 `json:"integrity" `
	Speed       int     `json:"speed" `
}

type Module struct {
	Symbol      string `json:"symbol" `
	Capacity    int    `json:"capacity" `
	Range       int    `json:"range" `
	Name        string `json:"name" `
	Description string `json:"description" `
}

type Mount struct {
	Symbol      string   `json:"symbol" `
	Name        string   `json:"name" `
	Description string   `json:"description"`
	Strength    int      `json:"strength" `
	Deposits    []string `json:"deposits"`
}

type Crew struct {
	Required int `json:"required" `
	Capacity int `json:"capacity" `
}
