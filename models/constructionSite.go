package models

type ConstructionSite struct {
	Symbol     string     `json:"symbol"`
	Materials  []Material `json:"materials"`
	IsComplete bool       `json:"isComplete"`
}

type Material struct {
	Good      []GoodSymbol `json:"tradeSymbol"`
	Required  int          `json:"quantity"`
	Fulfilled int          `json:"fulfilled"`
}
