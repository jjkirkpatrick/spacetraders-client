package models

type AcceptContractResponse struct {
	Data struct {
		Agent    *Agent    `json:"agent"`
		Contract *Contract `json:"contract"`
	} `json:"data"`
}

type DeliverContractCargoResponse struct {
	Data struct {
		Contract *Contract `json:"contract"`
		Cargo    *Cargo    `json:"cargo"`
	} `json:"data"`
}

// Contract represents the contract details
type Contract struct {
	ID               string        `json:"id"`
	FactionSymbol    string        `json:"factionSymbol"`
	Type             string        `json:"type"`
	Terms            ContractTerms `json:"terms"`
	Accepted         bool          `json:"accepted"`
	Fulfilled        bool          `json:"fulfilled"`
	Expiration       string        `json:"expiration"`
	DeadlineToAccept string        `json:"deadlineToAccept"`
}

// ContractTerms represents the terms of a contract
type ContractTerms struct {
	Deadline string            `json:"deadline"`
	Payment  ContractPayment   `json:"payment"`
	Deliver  []ContractDeliver `json:"deliver"`
}

// ContractPayment represents the payment terms of a contract
type ContractPayment struct {
	OnAccepted  int `json:"onAccepted"`
	OnFulfilled int `json:"onFulfilled"`
}

// ContractDeliver represents the delivery terms of a contract
type ContractDeliver struct {
	TradeSymbol       string `json:"tradeSymbol"`
	DestinationSymbol string `json:"destinationSymbol"`
	UnitsRequired     int    `json:"unitsRequired"`
	UnitsFulfilled    int    `json:"unitsFulfilled"`
}
