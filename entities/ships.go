package entities

import (
	"fmt"

	"github.com/jjkirkpatrick/spacetraders-client/client"
	"github.com/jjkirkpatrick/spacetraders-client/internal/api"
	"github.com/jjkirkpatrick/spacetraders-client/internal/models"
)

type Ship struct {
	models.Ship
	client *client.Client
}

func ListShips(c *client.Client) ([]*Ship, error) {
	fetchFunc := func(meta models.Meta) ([]*Ship, models.Meta, error) {
		metaPtr := &meta

		// Check if ships are in cache
		ships, metaPtr, err := api.ListShips(c.Get, metaPtr)

		var convertedShips []*Ship
		for _, modelShip := range ships {
			convertedShip := &Ship{
				Ship:   *modelShip, // Directly embed the modelShip
				client: c,
			}
			convertedShips = append(convertedShips, convertedShip)
		}

		if err != nil {
			if metaPtr == nil {
				// Use default Meta values or handle accordingly
				defaultMeta := models.Meta{Page: 1, Limit: 25, Total: 0}
				metaPtr = &defaultMeta
			}
			return convertedShips, *metaPtr, err.AsError()
		}
		if metaPtr != nil {
			// Store ships in cache
			return convertedShips, *metaPtr, nil
		} else {
			defaultMeta := models.Meta{Page: 1, Limit: 25, Total: 0}
			return convertedShips, defaultMeta, nil
		}
	}
	return client.NewPaginator[*Ship](fetchFunc).FetchAllPages()
}

func GetShip(c *client.Client, symbol string) (*Ship, error) {
	ship, err := api.GetShip(c.Get, symbol)
	if err != nil {
		return nil, err
	}

	shipEntity := &Ship{
		Ship:   *ship,
		client: c,
	}

	return shipEntity, nil
}

func PurchaseShip(c *client.Client, shipType string, waypoint string) (*models.Agent, *Ship, *models.Transaction, error) {
	purchaseShipRequest := &models.PurchaseShipRequest{
		ShipType:       models.ShipType(shipType),
		WaypointSymbol: waypoint,
	}

	response, err := api.PurchaseShip(c.Post, purchaseShipRequest)
	if err != nil {
		return nil, nil, nil, err.AsError()
	}

	shipEntity := &Ship{
		Ship:   response.Data.Ship,
		client: c,
	}

	c.CacheClient.Delete("all_ships")

	return &response.Data.Agent, shipEntity, &response.Data.Transaction, nil
}

func (s *Ship) Orbit() error {
	nav, err := api.OrbitShip(s.client.Post, s.Symbol)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return err.AsError()
	}

	s.Nav = *nav

	return nil
}

func (s *Ship) Dock() error {
	nav, err := api.DockShip(s.client.Post, s.Symbol)
	if err != nil {
		return err.AsError()
	}

	s.Nav = *nav

	return nil
}

func (s *Ship) FetchCargo() (*models.Cargo, error) {
	cargo, err := api.GetShipCargo(s.client.Get, s.Symbol)
	if err != nil {
		return nil, err
	}

	s.Cargo = *cargo

	return cargo, nil
}

func (s *Ship) Refine(produce string) (*models.Produced, *models.Consumed, error) {
	refineRequest := &models.RefineRequest{
		Produce: produce,
	}

	response, err := api.ShipRefine(s.client.Post, s.Symbol, refineRequest)
	if err != nil {
		return nil, nil, err.AsError()
	}

	s.Cargo = response.Data.Cargo
	s.Cooldown = response.Data.Cooldown

	return &response.Data.Produced, &response.Data.Consumed, nil
}

func (s *Ship) Chart() (*models.Chart, *models.Waypoint, error) {
	nav, err := api.CreateChart(s.client.Post, s.Symbol)
	if err != nil {
		return nil, nil, err.AsError()
	}

	return &nav.Data.Chart, &nav.Data.Waypoint, nil
}

func (s *Ship) FetchCooldown() error {
	cooldown, err := api.GetShipCooldown(s.client.Get, s.Symbol)
	if err != nil {
		return err
	}

	s.Cooldown = *cooldown

	return nil
}

func (s *Ship) Survey() ([]models.Survey, error) {
	response, err := api.CreateSurvey(s.client.Post, s.Symbol)
	if err != nil {
		return nil, err.AsError()
	}

	s.Cooldown = response.Data.Cooldown

	return response.Data.Survey, nil
}

func (s *Ship) Extract() (*models.Extraction, error) {
	response, err := api.ExtractResources(s.client.Post, s.Symbol)
	if err != nil {
		return nil, err.AsError()
	}

	s.Cargo = response.Data.Cargo
	s.Cooldown = response.Data.Cooldown

	return &response.Data.Extraction, nil
}

func (s *Ship) Siphon() (*models.Extraction, error) {
	response, err := api.SiphonResources(s.client.Post, s.Symbol)
	if err != nil {
		return nil, err.AsError()
	}

	s.Cargo = response.Data.Cargo
	s.Cooldown = response.Data.Cooldown

	return &response.Data.Extraction, nil
}

func (s *Ship) ExtractWithSurvey(survey models.Survey) (*models.Extraction, error) {
	extractWithSurveyRequest := &models.ExtractWithSurveyRequest{
		Survey: survey,
	}

	response, err := api.ExtractResourcesWithSurvey(s.client.Post, s.Symbol, extractWithSurveyRequest)
	if err != nil {
		return nil, err.AsError()
	}

	s.Cargo = response.Data.Cargo
	s.Cooldown = response.Data.Cooldown

	return &response.Data.Extraction, nil
}

func (s *Ship) Jettison(goodSymbol models.GoodSymbol, units int) (*models.Cargo, error) {
	jettisonRequest := &models.JettisonRequest{
		Symbol: goodSymbol,
		Units:  units,
	}

	response, err := api.JettisonCargo(s.client.Post, s.Symbol, jettisonRequest)
	if err != nil {
		return nil, err.AsError()
	}

	s.Cargo = response.Data.Cargo

	return &response.Data.Cargo, nil
}

func (s *Ship) Jump(systemSymbol string) (*models.ShipNav, *models.ShipCooldown, *models.Transaction, *models.Agent, error) {
	jumpRequest := &models.JumpShipRequest{
		WaypointSymbol: systemSymbol,
	}

	response, err := api.JumpShip(s.client.Post, s.Symbol, jumpRequest)
	if err != nil {
		return nil, nil, nil, nil, err.AsError()
	}

	s.Nav = response.Data.Nav
	s.Cooldown = response.Data.Cooldown

	return &response.Data.Nav, &response.Data.Cooldown, &response.Data.Transaction, &response.Data.Agent, nil
}

func (s *Ship) Navigate(waypointSymbol string) (*models.FuelDetails, *models.ShipNav, []models.Event, error) {
	navigateRequest := &models.NavigateRequest{
		WaypointSymbol: waypointSymbol,
	}

	response, err := api.NavigateShip(s.client.Post, s.Symbol, navigateRequest)
	if err != nil {
		return nil, nil, nil, err.AsError()
	}

	s.Fuel = response.Data.Fuel
	s.Nav = response.Data.Nav

	return &response.Data.Fuel, &response.Data.Nav, response.Data.Events, nil
}

func (s *Ship) SetFlightMode(flightmode models.FlightMode) error {
	flightModeRequest := &models.NavUpdateRequest{
		FlightMode: flightmode,
	}

	response, err := api.PatchShipNav(s.client.Patch, s.Symbol, flightModeRequest)
	if err != nil {
		return err.AsError()
	}

	s.Nav.FlightMode = response.Data.FlightMode
	s.Nav.Status = response.Data.Status
	s.Nav.Route = response.Data.Route
	s.Nav.SystemSymbol = response.Data.SystemSymbol
	s.Nav.WaypointSymbol = response.Data.WaypointSymbol

	return nil
}

func (s *Ship) FetchNavigationStatus() (*models.ShipNav, error) {
	response, err := api.GetShipNav(s.client.Get, s.Symbol)
	if err != nil {
		return nil, err.AsError()
	}

	return response, nil
}

func (s *Ship) Warp(waypointSymbol string) (*models.FuelDetails, *models.ShipNav, error) {
	warpRequest := &models.WarpRequest{
		WaypointSymbol: waypointSymbol,
	}

	response, err := api.WarpShip(s.client.Post, s.Symbol, warpRequest)
	if err != nil {
		return nil, nil, err.AsError()
	}

	s.Fuel = response.Data.Fuel
	s.Nav = response.Data.Nav

	return &response.Data.Fuel, &response.Data.Nav, nil
}

func (s *Ship) SellCargo(goodSymbol models.GoodSymbol, units int) (*models.Agent, *models.Cargo, *models.Transaction, error) {
	sellRequest := &models.SellCargoRequest{
		Symbol: goodSymbol,
		Units:  units,
	}

	response, err := api.SellCargo(s.client.Post, s.Symbol, sellRequest)
	if err != nil {
		return nil, nil, nil, err.AsError()
	}

	s.Cargo = response.Data.Cargo

	return &response.Data.Agent, &response.Data.Cargo, &response.Data.Transaction, nil
}

func (s *Ship) ScanSystems() (*models.ShipCooldown, []models.System, error) {
	response, err := api.ScanSystems(s.client.Post, s.Symbol)
	if err != nil {
		return nil, nil, err.AsError()
	}

	s.Cooldown = response.Data.Cooldown

	return &response.Data.Cooldown, response.Data.Systems, nil
}

func (s *Ship) ScanWaypoints() (*models.ShipCooldown, []models.Waypoint, error) {
	response, err := api.ScanWaypoints(s.client.Post, s.Symbol)
	if err != nil {
		return nil, nil, err.AsError()
	}

	s.Cooldown = response.Data.Cooldown

	return &response.Data.Cooldown, response.Data.Waypoints, nil
}

func (s *Ship) Refuel(amount int, fromCargo bool) (*models.Agent, *models.FuelDetails, *models.Transaction, error) {
	refuelRequest := &models.RefuelShipRequest{
		Units:     amount,
		FromCargo: fromCargo,
	}
	response, err := api.RefuelShip(s.client.Post, s.Symbol, refuelRequest)
	if err != nil {
		return nil, nil, nil, err.AsError()
	}

	s.Fuel = response.Data.Fuel

	return &response.Data.Agent, &response.Data.Fuel, &response.Data.Transaction, nil
}

func (s *Ship) PurchaseCargo(goodSymbol models.GoodSymbol, units int) (*models.Agent, *models.Cargo, *models.Transaction, error) {
	purchaseRequest := &models.PurchaseCargoRequest{
		Symbol: goodSymbol,
		Units:  units,
	}

	response, err := api.PurchaseCargo(s.client.Post, s.Symbol, purchaseRequest)
	if err != nil {
		return nil, nil, nil, err.AsError()
	}

	s.Cargo = response.Data.Cargo

	return &response.Data.Agent, &response.Data.Cargo, &response.Data.Transaction, nil
}

func (s *Ship) TransferCargo(goodSymbol models.GoodSymbol, units int, shipSymbol string) (*models.Cargo, error) {
	transferRequest := &models.TransferCargoRequest{
		TradeSymbol: goodSymbol,
		Units:       units,
		ShipSymbol:  shipSymbol,
	}

	response, err := api.TransferCargo(s.client.Post, s.Symbol, transferRequest)
	if err != nil {
		return nil, err.AsError()
	}

	s.Cargo = response.Data.Cargo

	return &response.Data.Cargo, nil
}

func (s *Ship) NegotiateContract() (*models.Contract, error) {

	response, err := api.NegotiateContract(s.client.Post, s.Symbol)
	if err != nil {
		return nil, err.AsError()
	}

	return &response.Data.Contract, nil
}

func (s *Ship) GetMounts() (*models.MountSymbol, string, string, int, []string, models.ShipRequirements, error) {
	response, err := api.GetMounts(s.client.Get, s.Symbol)
	if err != nil {
		return nil, "", "", 0, nil, models.ShipRequirements{}, err.AsError()
	}

	return &response.Data.Symbol, response.Data.Name, response.Data.Description, response.Data.Strength, response.Data.Depsits, response.Data.Requirements, nil
}

func (s *Ship) InstallMount(mountSymbol models.MountSymbol) (*models.Agent, []models.Mount, *models.Cargo, *models.Transaction, error) {
	installRequest := &models.InstallMountRequest{
		Symbol: mountSymbol,
	}

	response, err := api.InstallMount(s.client.Post, s.Symbol, installRequest)
	if err != nil {
		return nil, nil, nil, nil, err.AsError()
	}

	s.Cargo = response.Data.Cargo

	return &response.Data.Agent, response.Data.Mounts, &response.Data.Cargo, &response.Data.Transaction, nil
}

func (s *Ship) RemoveMount(mountSymbol models.MountSymbol) (*models.Agent, []models.Mount, *models.Cargo, *models.Transaction, error) {
	removeRequest := &models.RemoveMountRequest{
		Symbol: mountSymbol,
	}

	response, err := api.RemoveMount(s.client.Post, s.Symbol, removeRequest)
	if err != nil {
		return nil, nil, nil, nil, err.AsError()
	}

	s.Cargo = response.Data.Cargo

	return &response.Data.Agent, response.Data.Mounts, &response.Data.Cargo, &response.Data.Transaction, nil
}

func (s *Ship) GetScrapPrice() (*models.Transaction, error) {
	response, err := api.GetScrapShip(s.client.Get, s.Symbol)
	if err != nil {
		return nil, err.AsError()
	}

	return &response.Data.Transaction, nil
}

func (s *Ship) ScrapShip() (*models.Transaction, error) {
	response, err := api.ScrapShip(s.client.Post, s.Symbol)
	if err != nil {
		return nil, err.AsError()
	}

	s.client.CacheClient.Delete("all_ships")

	return &response.Data.Transaction, nil
}

func (s *Ship) GetRepairPrice() (*models.Transaction, error) {
	response, err := api.GetRepairShip(s.client.Get, s.Symbol)
	if err != nil {
		return nil, err.AsError()
	}

	return &response.Data.Transaction, nil
}

func (s *Ship) RepairShip() (*models.Ship, *models.Transaction, error) {
	response, err := api.RepairShip(s.client.Post, s.Symbol)
	if err != nil {
		return nil, nil, err.AsError()
	}

	s.Ship = response.Data.Ship

	return &response.Data.Ship, &response.Data.Transaction, nil
}
