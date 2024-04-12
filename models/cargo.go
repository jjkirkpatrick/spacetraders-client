package models

type DeliverContractCargoRequest struct {
	ShipSymbol  string     `json:"shipSymbol"`
	TradeSymbol GoodSymbol `json:"tradeSymbol"`
	Units       int        `json:"units"`
}

type Cargo struct {
	Capacity  int         `json:"capacity"`
	Units     int         `json:"units"`
	Inventory []Inventory `json:"inventory"`
}

type Inventory struct {
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Units       int    `json:"units"`
}
