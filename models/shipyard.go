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
	Transactions  []Transaction  `json:"transactions"`
	Ships         []ShipyardShip `json:"ships"`
	Modifications int            `json:"monificationFee"`
}

type Transaction struct {
	WaypointSymbol string `json:"waypointSymbol"`
	ShipSymbol     string `json:"shipSymbol"`
	TradeSymbol    string `json:"tradeSymbol"`
	ShipType       string `json:"shipType"`
	Price          int    `json:"price"`
	AgentSymbol    string `json:"agentSymbol"`
	Type           string `json:"type" enum:"PURCHASE,SELL"`
	Units          int    `json:"units"`
	PricePerUnit   int    `json:"pricePerUnit"`
	TotalPrice     int    `json:"totalPrice"`
	Timestamp      string `json:"timestamp"`
}

type ShipyardShip struct {
	Type             ShipType     `json:"type" `
	Name             string       `json:"name" `
	Description      string       `json:"description" `
	Supply           string       `json:"supply" `
	Activity         string       `json:"activity"`
	PurchasePrice    int          `json:"purchasePrice" `
	Frame            ShipFrame    `json:"frame" `
	Reactor          ShipReactor  `json:"reactor" `
	Engine           ShipEngine   `json:"engine" `
	Modules          []ShipModule `json:"modules" `
	Mounts           []ShipMount  `json:"mounts" `
	Crew             ShipCrew     `json:"crew" `
	ModificationsFee int          `json:"modificationsFee" `
}
