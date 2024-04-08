package models

type PurchaseShipRequest struct {
	ShipType       ShipType `json:"shipType"`
	WaypointSymbol string   `json:"waypointSymbol"`
}

type RefineRequest struct {
	Produce string `json:"produce"`
}

// Ship represents the ship details
type Ship struct {
	Symbol       string           `json:"symbol"`
	Registration ShipRegistration `json:"registration"`
	Nav          ShipNav          `json:"nav"`
	Crew         ShipCrew         `json:"crew"`
	Frame        ShipFrame        `json:"frame"`
	Reactor      ShipReactor      `json:"reactor"`
	Engine       ShipEngine       `json:"engine"`
	Cooldown     ShipCooldown     `json:"cooldown"`
	Modules      []ShipModule     `json:"modules"`
	Mounts       []ShipMount      `json:"mounts"`
	Cargo        ShipCargo        `json:"cargo"`
	Fuel         ShipFuel         `json:"fuel"`
}

// ShipRegistration represents the registration information of a ship
type ShipRegistration struct {
	Name          string               `json:"name"`
	FactionSymbol string               `json:"factionSymbol"`
	Role          ShipRegistrationRole `json:"role"`
}

type ShipRegistrationRole string

const (
	Fabricator  ShipRegistrationRole = "FABRICATOR"
	Harvester   ShipRegistrationRole = "HARVESTER"
	Hauler      ShipRegistrationRole = "HAULER"
	Interceptor ShipRegistrationRole = "INTERCEPTOR"
	Excavator   ShipRegistrationRole = "EXCAVATOR"
	Transport   ShipRegistrationRole = "TRANSPORT"
	Repair      ShipRegistrationRole = "REPAIR"
	Surveyor    ShipRegistrationRole = "SURVEYOR"
	Command     ShipRegistrationRole = "COMMAND"
	Carrier     ShipRegistrationRole = "CARRIER"
	Patrol      ShipRegistrationRole = "PATROL"
	Satellite   ShipRegistrationRole = "SATELLITE"
	Explorer    ShipRegistrationRole = "EXPLORER"
	Refinery    ShipRegistrationRole = "REFINERY"
)

// ShipNav represents the navigation information of a ship
type ShipNav struct {
	SystemSymbol   string       `json:"systemSymbol"`
	WaypointSymbol string       `json:"waypointSymbol"`
	Route          ShipNavRoute `json:"route"`
	Status         NavStatus    `json:"status"`
	FlightMode     FlightMode   `json:"flightMode"`
}

// ShipNavRoute represents the route information of a ship's navigation
type ShipNavRoute struct {
	Destination   RouteWaypoint `json:"destination"`
	Origin        RouteWaypoint `json:"origin"`
	DepartureTime string        `json:"departureTime"`
	Arrival       string        `json:"arrival"`
}

// RouteWaypoint represents a waypoint in a ship's route
type RouteWaypoint struct {
	Symbol       string `json:"symbol"`
	Type         string `json:"type"`
	SystemSymbol string `json:"systemSymbol"`
	X            int    `json:"x"`
	Y            int    `json:"y"`
}

// ShipCrew represents the crew information of a ship
type ShipCrew struct {
	Current  int    `json:"current"`
	Required int    `json:"required"`
	Capacity int    `json:"capacity"`
	Rotation string `json:"rotation"`
	Morale   int    `json:"morale"`
	Wages    int    `json:"wages"`
}

// ShipFrame represents the frame information of a ship
type ShipFrame struct {
	Symbol         string           `json:"symbol"`
	Name           string           `json:"name"`
	Description    string           `json:"description"`
	Condition      float64          `json:"condition"`
	Integrity      float64          `json:"integrity"`
	ModuleSlots    int              `json:"moduleSlots"`
	MountingPoints int              `json:"mountingPoints"`
	FuelCapacity   int              `json:"fuelCapacity"`
	Requirements   ShipRequirements `json:"requirements"`
}

// ShipReactor represents the reactor information of a ship
type ShipReactor struct {
	Symbol       string           `json:"symbol"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	Condition    float64          `json:"condition"`
	Integrity    float64          `json:"integrity"`
	PowerOutput  int              `json:"powerOutput"`
	Requirements ShipRequirements `json:"requirements"`
}

// ShipEngine represents the engine information of a ship
type ShipEngine struct {
	Symbol       string           `json:"symbol"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	Condition    float64          `json:"condition"`
	Integrity    float64          `json:"integrity"`
	Speed        int              `json:"speed"`
	Requirements ShipRequirements `json:"requirements"`
}

// ShipCooldown represents the cooldown information of a ship
type ShipCooldown struct {
	ShipSymbol       string `json:"shipSymbol"`
	TotalSeconds     int    `json:"totalSeconds"`
	RemainingSeconds int    `json:"remainingSeconds"`
	Expiration       string `json:"expiration"`
}

// ShipModule represents a module installed in a ship
type ShipModule struct {
	Symbol       string           `json:"symbol"`
	Capacity     int              `json:"capacity"`
	Range        int              `json:"range"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	Requirements ShipRequirements `json:"requirements"`
}

// ShipMount represents a mount installed in a ship
type ShipMount struct {
	Symbol       string           `json:"symbol"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	Strength     int              `json:"strength"`
	Deposits     []string         `json:"deposits"`
	Requirements ShipRequirements `json:"requirements"`
}

// ShipCargo represents the cargo information of a ship
type ShipCargo struct {
	Capacity  int             `json:"capacity"`
	Units     int             `json:"units"`
	Inventory []ShipCargoItem `json:"inventory"`
}

// ShipCargoItem represents an item in a ship's cargo
type ShipCargoItem struct {
	Symbol      string `json:"symbol"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Units       int    `json:"units"`
}

// ShipFuel represents the fuel information of a ship
type ShipFuel struct {
	Current  int              `json:"current"`
	Capacity int              `json:"capacity"`
	Consumed ShipFuelConsumed `json:"consumed"`
}

// ShipFuelConsumed represents the fuel consumed by a ship
type ShipFuelConsumed struct {
	Amount    int    `json:"amount"`
	Timestamp string `json:"timestamp"`
}

// ShipRequirements represents the requirements for installing a component on a ship
type ShipRequirements struct {
	Power int `json:"power"`
	Crew  int `json:"crew"`
	Slots int `json:"slots"`
}

type Survey struct {
	Signature  string    `json:"signature" `
	Symbol     string    `json:"symbol" `
	Deposits   []Deposit `json:"deposits" `
	Expiration string    `json:"expiration"  format:"date-time"`
	Size       string    `json:"size" `
}

type Deposit struct {
	Symbol     string `json:"symbol" `
	Expiration string `json:"expiration"  format:"date-time"`
	Size       string `json:"size" `
}

// FuelDetails represents the details of the ship's fuel tanks
type FuelDetails struct {
	Current  int           `json:"current" `
	Capacity int           `json:"capacity" `
	Consumed *FuelConsumed `json:"consumed,omitempty"`
}

// FuelConsumed represents the fuel consumption data
type FuelConsumed struct {
	Amount    int    `json:"amount" `
	Timestamp string `json:"timestamp"  format:"date-time"`
}

type FlightMode string

const (
	FlightModeDrift   FlightMode = "DRIFT"
	FlightNodeStealth FlightMode = "STEALTH"
	FlightModeCruise  FlightMode = "CRUISE"
	FlightModeBurn    FlightMode = "BURN"
)

type NavStatus string

const (
	NavStatusDocked    NavStatus = "DOCKED"
	NavStatusInOrbit   NavStatus = "IN_ORBIT"
	NavStatusInTransit NavStatus = "IN_TRANSIT"
)

type MountSymbol string

const (
	MountGasSiphonI       MountSymbol = "MOUNT_GAS_SIPHON_I"
	MountGasSiphonII      MountSymbol = "MOUNT_GAS_SIPHON_II"
	MountGasSiphonIII     MountSymbol = "MOUNT_GAS_SIPHON_III"
	MountSurveyorI        MountSymbol = "MOUNT_SURVEYOR_I"
	MountSurveyorII       MountSymbol = "MOUNT_SURVEYOR_II"
	MountSurveyorIII      MountSymbol = "MOUNT_SURVEYOR_III"
	MountSensorArrayI     MountSymbol = "MOUNT_SENSOR_ARRAY_I"
	MountSensorArrayII    MountSymbol = "MOUNT_SENSOR_ARRAY_II"
	MountSensorArrayIII   MountSymbol = "MOUNT_SENSOR_ARRAY_III"
	MountMiningLaserI     MountSymbol = "MOUNT_MINING_LASER_I"
	MountMiningLaserII    MountSymbol = "MOUNT_MINING_LASER_II"
	MountMiningLaserIII   MountSymbol = "MOUNT_MINING_LASER_III"
	MountLaserCannonI     MountSymbol = "MOUNT_LASER_CANNON_I"
	MountMissileLauncherI MountSymbol = "MOUNT_MISSILE_LAUNCHER_I"
	MountTurretI          MountSymbol = "MOUNT_TURRET_I"
)
