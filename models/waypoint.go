package models

type Chart struct {
	WaypointSymbol string `json:"waypointSymbol"`
	SubmittedBy    string `json:"submittedBy"`
	SubmittedOn    string `json:"submittedOn"`
}
