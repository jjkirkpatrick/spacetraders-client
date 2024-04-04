package models

// Faction represents the faction details
type Faction struct {
	Symbol       string         `json:"symbol"`
	Name         string         `json:"name"`
	Description  string         `json:"description"`
	Headquarters string         `json:"headquarters"`
	Traits       []FactionTrait `json:"traits"`
	IsRecruiting bool           `json:"isRecruiting"`
}

// FactionTrait represents a trait of a faction
type FactionTrait struct {
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
