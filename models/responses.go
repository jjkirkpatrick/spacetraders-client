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
		Survey   []Survey     `json:"survey"`
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
		Mounts      []Mount     `json:"mounts"`
		Cargo       Cargo       `json:"cargo"`
		Transaction Transaction `json:"transaction"`
	} `json:"data"`
}
type RemoveMountResponse struct {
	Data struct {
		Agent       Agent       `json:"agent"`
		Mounts      []Mount     `json:"mounts"`
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

type PatchShipNacResponse struct {
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
		Produced struct {
			TradeSymbol GoodSymbol `json:"tradeSymbol"`
			Units       int        `json:"units"`
		} `json:"produced"`
		Consumed struct {
			TradeSymbol GoodSymbol `json:"tradeSymbol"`
			Units       int        `json:"units"`
		} `json:"consumed"`
	} `json:"data"`
}
