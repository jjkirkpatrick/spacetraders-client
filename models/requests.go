package models

type ExtractWithSurveyRequest struct {
	Signature  string    `json:"signature"`
	Symbol     string    `json:"symbol"`
	Deposits   []Deposit `json:"deposits"`
	Expiration string    `json:"expiration"  format:"date-time"`
	Size       string    `json:"size"`
}

type JettisonRequest struct {
	Symbol GoodSymbol `json:"symbol"`
	Units  int        `json:"units"`
}
type JumpShipRequest struct {
	WaypointSymbol string `json:"waypointSymbol"`
}

type NavigateRequest struct {
	WaypointSymbol string `json:"waypointSymbol"`
}
type NavUpdateRequest struct {
	FlightMode FlightMode `json:"flightMode"`
}
type WarpRequest struct {
	WaypointSymbol string `json:"waypointSymbol"`
}

type SellCargoRequest struct {
	Symbol GoodSymbol `json:"symbol"`
	Units  int        `json:"units"`
}
type RefuelShipRequest struct {
	Units     int  `json:"units"`
	FromCargo bool `json:"fromCargo"`
}
type PurchaseCargoRequest struct {
	Symbol GoodSymbol `json:"symbol"`
	Units  int        `json:"units"`
}
type TransferCargoRequest struct {
	TradeSymbol GoodSymbol `json:"tradeSymbol"`
	Units       int        `json:"units"`
	ShipSymbol  string     `json:"shipSymbol"`
}

type InstallMountRequest struct {
	Symbol MountSymbol `json:"symbol"`
}
type RemoveMountRequest struct {
	Symbol MountSymbol `json:"symbol"`
}
