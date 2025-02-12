package models

type CreateChartResponse struct {
	Data struct {
		Chart    Chart    `json:"chart"`
		Waypoint Waypoint `json:"waypoint"`
	} `json:"data"`
}

type CreateSurveyResponse struct {
	Data struct {
		Cooldown ShipCooldown `json:"cooldown"`
		Surveys  []Survey     `json:"surveys"`
	} `json:"data"`
}

type ExtractionResponse struct {
	Data struct {
		Cooldown   ShipCooldown `json:"cooldown"`
		Extraction Extraction   `json:"extraction"`
		Cargo      Cargo        `json:"cargo"`
		Event      []Event      `json:"event"`
	} `json:"data"`
}

type SiphonResponse struct {
	Data struct {
		Cooldown   ShipCooldown `json:"cooldown"`
		Extraction Extraction   `json:"extraction"`
		Cargo      Cargo        `json:"cargo"`
		Event      []Event      `json:"event"`
	} `json:"data"`
}

type JettisonResponse struct {
	Data struct {
		Cargo Cargo `json:"cargo"`
	} `json:"data"`
}

type JumpShipResponse struct {
	Data struct {
		Nav         ShipNav      `json:"nav"`
		Cooldown    ShipCooldown `json:"cooldown"`
		Transaction Transaction  `json:"transaction"`
		Agent       Agent        `json:"agent"`
	} `json:"data"`
}

type NavigateResponse struct {
	Data struct {
		Fuel   FuelDetails `json:"fuel"`
		Nav    ShipNav     `json:"nav"`
		Events []Event     `json:"events"`
	} `json:"data"`
}
type NavUpdateResponse struct {
	Data struct {
		SystemSymbol   string       `json:"systemSymbol"`
		WaypointSymbol string       `json:"waypointSymbol"`
		Route          ShipNavRoute `json:"route"`
		Status         NavStatus    `json:"status"`
		FlightMode     FlightMode   `json:"flightMode"`
	} `json:"data"`
}
type GetNavStatusResponse struct {
	Data struct {
		SystemSymbol  string       `json:"systemSymbol"`
		WapointSymbol string       `json:"waypointSymbol"`
		Route         ShipNavRoute `json:"route"`
		Status        NavStatus    `json:"status"`
		FlightMode    FlightMode   `json:"flightMode"`
	} `json:"data"`
}
type WarpResponse struct {
	Data struct {
		Fuel FuelDetails `json:"fuel"`
		Nav  ShipNav     `json:"nav"`
	} `json:"data"`
}
type SellCargoResponse struct {
	Data struct {
		Agent       Agent       `json:"agent"`
		Cargo       Cargo       `json:"cargo"`
		Transaction Transaction `json:"transaction"`
	} `json:"data"`
}
type ScanSystemsResponse struct {
	Data struct {
		Cooldown ShipCooldown `json:"cooldown"`
		Systems  []System     `json:"systems"`
	} `json:"data"`
}
type ScanWaypointsResponse struct {
	Data struct {
		Cooldown  ShipCooldown `json:"cooldown"`
		Waypoints []Waypoint   `json:"waypoints"`
	} `json:"data"`
}
type ScanShipsResponse struct {
	Data struct {
		Cooldown ShipCooldown `json:"cooldown"`
		Ships    []Ship       `json:"ships"`
	} `json:"data"`
}
type RefuelShipResponse struct {
	Data struct {
		Agent       Agent       `json:"agent"`
		Fuel        FuelDetails `json:"fuel"`
		Transaction Transaction `json:"transaction"`
	} `json:"data"`
}

type PurchaseCargoResponse struct {
	Data struct {
		Agent       Agent       `json:"agent"`
		Cargo       Cargo       `json:"cargo"`
		Transaction Transaction `json:"transaction"`
	} `json:"data"`
}
type TransferCargoResponse struct {
	Data struct {
		Cargo Cargo `json:"cargo"`
	} `json:"data"`
}
type NegotiateContractResponse struct {
	Data struct {
		Contract Contract `json:"contract"`
	} `json:"data"`
}

type GetMountsResponse struct {
	Data struct {
		Symbol       MountSymbol      `json:"symbol"`
		Name         string           `json:"name"`
		Description  string           `json:"description"`
		Strength     int              `json:"strength"`
		Depsits      []string         `json:"deposits"`
		Requirements ShipRequirements `json:"requirements"`
	} `json:"data"`
}
type InstallMountResponse struct {
	Data struct {
		Agent       Agent       `json:"agent"`
		Mounts      []ShipMount `json:"mounts"`
		Cargo       Cargo       `json:"cargo"`
		Transaction Transaction `json:"transaction"`
	} `json:"data"`
}
type RemoveMountResponse struct {
	Data struct {
		Agent       Agent       `json:"agent"`
		Mounts      []ShipMount `json:"mounts"`
		Cargo       Cargo       `json:"cargo"`
		Transaction Transaction `json:"transaction"`
	} `json:"data"`
}
type GetScrapShipResponse struct {
	Data struct {
		Transaction Transaction `json:"transaction"`
	} `json:"data"`
}

type ScrapShipResponse struct {
	Data struct {
		Agent       Agent       `json:"agent"`
		Transaction Transaction `json:"transaction"`
	} `json:"data"`
}
type GetRepairShipResponse struct {
	Data struct {
		Transaction Transaction `json:"transaction"`
	} `json:"data"`
}
type RepairShipResponse struct {
	Data struct {
		Agent       Agent       `json:"agent"`
		Ship        Ship        `json:"ship"`
		Transaction Transaction `json:"transaction"`
	} `json:"data"`
}

type PatchShipNavResponse struct {
	Data struct {
		SystemSymbol   string       `json:"systemSymbol"`
		WaypointSymbol string       `json:"waypointSymbol"`
		Route          ShipNavRoute `json:"route"`
		Status         NavStatus    `json:"status"`
		FlightMode     FlightMode   `json:"flightMode"`
	} `json:"data"`
}

type PurchaseShipResponse struct {
	Data struct {
		Agent       Agent       `json:"agent"`
		Ship        Ship        `json:"ship"`
		Transaction Transaction `json:"transaction"`
	} `json:"data"`
}

type ShipRefineResponse struct {
	Data struct {
		Cargo    Cargo        `json:"cargo"`
		Cooldown ShipCooldown `json:"coolDown"`
		Produced Produced     `json:"produced"`
		Consumed Consumed     `json:"consumed"`
	} `json:"data"`
}

type ListShipsResponse struct {
	Data []*Ship `json:"data"`
	Meta Meta    `json:"meta"`
}

type ListFactionsResponse struct {
	Data []*Faction `json:"data"`
	Meta Meta       `json:"meta"`
}

type ListSystemsResponse struct {
	Symbol       string       `json:"symbol"`
	SectorSymbol string       `json:"sectorSymbol"`
	Type         WaypointType `json:"type"`
	X            int          `json:"x"`
	Y            int          `json:"y"`
	Waypoints    []Waypoint   `json:"waypoints"`
	Factions     []Faction    `json:"factions"`
}

type GetSystemResponse struct {
	Data struct {
		Symbol       string       `json:"symbol"`
		SectorSymbol string       `json:"sectorSymbol"`
		Type         WaypointType `json:"type"`
		X            int          `json:"x"`
		Y            int          `json:"y"`
		Waypoints    []Waypoint   `json:"waypoints"`
		Factions     []Faction    `json:"factions"`
	} `json:"data"`
}

type GetWaypointResponse struct {
	Data struct {
		Symbol              string           `json:"symbol"`
		Type                WaypointType     `json:"type"`
		SystemSymbol        string           `json:"systemSymbol"`
		X                   int              `json:"x"`
		Y                   int              `json:"y"`
		Orbitals            []Orbital        `json:"orbitals"`
		Orbits              string           `json:"orbits"`
		Faction             Faction          `json:"factions"`
		Traits              []WaypointTraits `json:"traits"`
		Modifiers           []Modifier       `json:"modifiers"`
		Chart               Chart            `json:"chart"`
		IsUnderConstruction bool             `json:"isUnderConstruction"`
	} `json:"data"`
}

type GetMarketResponse struct {
	Data struct {
		Symbol       string             `json:"symbol"`
		Exports      []Good             `json:"exports"`
		Imports      []Good             `json:"imports"`
		Exchange     []Good             `json:"exchange"`
		Transactions []Transaction      `json:"transactions"`
		TradeGoods   []MarketTradeGoods `json:"tradeGoods"`
	} `json:"data"`
}

type GetShipyardResponse struct {
	Data Shipyard `json:"data"`
}

type GetJumpGatesResponse struct {
	Data []JumpGate `json:"data"`
}

type GetConstructionSitesResponse struct {
	Data []ConstructionSite `json:"data"`
}

type SupplyConstructionSiteResponse struct {
	Data struct {
		Construction ConstructionSite `json:"construction"`
		Cargo        Cargo            `json:"cargo"`
	} `json:"data"`
}
