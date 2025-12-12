package api

import (
	"fmt"

	"github.com/jjkirkpatrick/spacetraders-client/models"
)

func ListShips(get GetFunc, meta *models.Meta) ([]*models.Ship, *models.Meta, *models.APIError) {
	endpoint := "/my/ships"

	var response models.ListShipsResponse

	queryParams := map[string]string{
		"page":  fmt.Sprintf("%d", meta.Page),
		"limit": fmt.Sprintf("%d", meta.Limit),
	}

	err := get(endpoint, queryParams, &response)
	if err != nil {
		return nil, nil, err
	}

	return response.Data, &response.Meta, nil
}

// PurchaseShip allows the user to purchase a new models.Ship
func PurchaseShip(post PostFunc, payload *models.PurchaseShipRequest) (*models.PurchaseShipResponse, *models.APIError) {
	endpoint := "/my/ships"

	var response models.PurchaseShipResponse

	err := post(endpoint, payload, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetShip retrieves the details of a specific models.Ship
func GetShip(get GetFunc, ShipSymbol string) (*models.Ship, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s", ShipSymbol)

	var response struct {
		Data models.Ship `json:"data"`
	}

	err := get(endpoint, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}

// GetShipCargo retrieves the cargo details of a specific models.Ship
func GetShipCargo(get GetFunc, ShipSymbol string) (*models.Cargo, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/cargo", ShipSymbol)

	var response struct {
		Data *models.Cargo `json:"data"`
	}

	err := get(endpoint, nil, &response)
	if err != nil {
		return nil, err
	}

	return response.Data, nil
}

// OrbitShip allows a models.Ship to orbit a celestial body
func OrbitShip(post PostFunc, ShipSymbol string) (*models.ShipNav, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/orbit", ShipSymbol)
	var response struct {
		Data struct {
			Nav models.ShipNav `json:"nav"`
		} `json:"data"`
	}

	err := post(endpoint, nil, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data.Nav, nil
}

// ShipRefine initiates the refining process for a models.Ship
func ShipRefine(post PostFunc, ShipSymbol string, payload *models.RefineRequest) (*models.ShipRefineResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/refine", ShipSymbol)

	var response models.ShipRefineResponse

	err := post(endpoint, payload, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// CreateChart creates a navigation chart for a models.Ship
func CreateChart(post PostFunc, ShipSymbol string) (*models.CreateChartResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/chart", ShipSymbol)

	var response models.CreateChartResponse

	err := post(endpoint, nil, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetShipCooldown retrieves the cooldown details of a specific models.Ship
func GetShipCooldown(get GetFunc, ShipSymbol string) (*models.ShipCooldown, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/cooldown", ShipSymbol)

	var response struct {
		Data models.ShipCooldown `json:"data"`
	}

	err := get(endpoint, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}

// DockShip allows a models.Ship to dock at a station or planet
func DockShip(post PostFunc, ShipSymbol string) (*models.ShipNav, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/dock", ShipSymbol)

	var response struct {
		Data struct {
			Nav models.ShipNav `json:"nav"`
		} `json:"data"`
	}

	err := post(endpoint, nil, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data.Nav, nil
}

// CreateSurvey initiates a survey process for a models.Ship
func CreateSurvey(post PostFunc, ShipSymbol string) (*models.CreateSurveyResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/survey", ShipSymbol)

	var response models.CreateSurveyResponse

	err := post(endpoint, nil, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// ExtractResources initiates the resource extraction process for a models.Ship
func ExtractResources(post PostFunc, ShipSymbol string) (*models.ExtractionResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/extract", ShipSymbol)

	var response models.ExtractionResponse

	err := post(endpoint, nil, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// SiphonResources initiates the resource siphoning process for a models.Ship
func SiphonResources(post PostFunc, ShipSymbol string) (*models.SiphonResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/siphon", ShipSymbol)

	var response models.SiphonResponse

	err := post(endpoint, nil, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// ExtractResourcesWithSurvey initiates the resource extraction process with a prior survey for a models.Ship
func ExtractResourcesWithSurvey(post PostFunc, ShipSymbol string, payload *models.ExtractWithSurveyRequest) (*models.ExtractionResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/extract/survey", ShipSymbol)

	var response models.ExtractionResponse

	err := post(endpoint, payload, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// JettisonCargo allows a models.Ship to jettison cargo into space
func JettisonCargo(post PostFunc, ShipSymbol string, payload *models.JettisonRequest) (*models.JettisonResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/jettison", ShipSymbol)

	var response models.JettisonResponse

	err := post(endpoint, payload, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// JumpShip initiates a jump for a models.Ship to another system
func JumpShip(post PostFunc, ShipSymbol string, payload *models.JumpShipRequest) (*models.JumpShipResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/jump", ShipSymbol)

	var response models.JumpShipResponse

	err := post(endpoint, payload, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// NavigateShip initiates navigation for a models.Ship to a waypoint
func NavigateShip(post PostFunc, ShipSymbol string, payload *models.NavigateRequest) (*models.NavigateResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/navigate", ShipSymbol)

	var response models.NavigateResponse
	err := post(endpoint, payload, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// PatchShipNav updates the navigation details of a models.Ship
func PatchShipNav(patch PatchFunc, ShipSymbol string, payload *models.NavUpdateRequest) (*models.PatchShipNavResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/nav", ShipSymbol)

	var response models.PatchShipNavResponse

	err := patch(endpoint, payload, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetShipNav retrieves the navigation details of a specific models.Ship
func GetShipNav(get GetFunc, ShipSymbol string) (*models.ShipNav, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/nav", ShipSymbol)

	var response struct {
		Data models.ShipNav `json:"data"`
	}

	err := get(endpoint, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response.Data, nil
}

// WarpShip initiates a warp for a models.Ship to another system
func WarpShip(post PostFunc, ShipSymbol string, payload *models.WarpRequest) (*models.WarpResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/warp", ShipSymbol)

	var response models.WarpResponse

	err := post(endpoint, payload, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// SellCargo sells cargo from a models.Ship's inventory
func SellCargo(post PostFunc, ShipSymbol string, payload *models.SellCargoRequest) (*models.SellCargoResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/sell", ShipSymbol)

	var response models.SellCargoResponse

	err := post(endpoint, payload, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// ScanSystems scans for systems within range
func ScanSystems(post PostFunc, ShipSymbol string) (*models.ScanSystemsResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/scan/systems", ShipSymbol)

	var response models.ScanSystemsResponse

	err := post(endpoint, nil, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// ScanWaypoints scans for waypoints within a system
func ScanWaypoints(post PostFunc, ShipSymbol string) (*models.ScanWaypointsResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/scan/waypoints", ShipSymbol)

	var response models.ScanWaypointsResponse

	err := post(endpoint, nil, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// ScanShips scans for models.Ships within range
func ScanShips(post PostFunc, ShipSymbol string) (*models.ScanShipsResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/scan/ships", ShipSymbol)

	var response models.ScanShipsResponse

	err := post(endpoint, nil, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// RefuelShip refuels a models.Ship
func RefuelShip(post PostFunc, ShipSymbol string, payload *models.RefuelShipRequest) (*models.RefuelShipResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/refuel", ShipSymbol)

	var response models.RefuelShipResponse
	var err *models.APIError

	if payload == nil {
		err = post(endpoint, nil, nil, &response)
	} else {
		err = post(endpoint, payload, nil, &response)
	}

	if err != nil {
		return nil, err
	}

	return &response, nil
}

// PurchaseCargo purchases cargo for a models.Ship
func PurchaseCargo(post PostFunc, ShipSymbol string, payload *models.PurchaseCargoRequest) (*models.PurchaseCargoResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/purchase", ShipSymbol)

	var response models.PurchaseCargoResponse

	err := post(endpoint, payload, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// TransferCargo transfers cargo between models.Ships or to a waypoint
func TransferCargo(post PostFunc, ShipSymbol string, payload *models.TransferCargoRequest) (*models.TransferCargoResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/transfer", ShipSymbol)

	var response models.TransferCargoResponse

	err := post(endpoint, payload, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// NegotiateContract negotiates a contract for a models.Ship
func NegotiateContract(post PostFunc, ShipSymbol string) (*models.NegotiateContractResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/negotiate/contract", ShipSymbol)

	var response models.NegotiateContractResponse

	err := post(endpoint, nil, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetMounts retrieves the mounts of a specific models.Ship
func GetMounts(get GetFunc, ShipSymbol string) (*models.GetMountsResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/mounts", ShipSymbol)

	var response models.GetMountsResponse

	err := get(endpoint, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// InstallMount installs a mount on a models.Ship
func InstallMount(post PostFunc, ShipSymbol string, payload *models.InstallMountRequest) (*models.InstallMountResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/mounts/install", ShipSymbol)

	var response models.InstallMountResponse

	err := post(endpoint, payload, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// RemoveMount removes a mount from a models.Ship
func RemoveMount(post PostFunc, ShipSymbol string, payload *models.RemoveMountRequest) (*models.RemoveMountResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/mounts/remove", ShipSymbol)

	var response models.RemoveMountResponse

	err := post(endpoint, payload, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetScrapShip retrieves the scrap value of a specific models.Ship
func GetScrapShip(get GetFunc, ShipSymbol string) (*models.GetScrapShipResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/scrap", ShipSymbol)

	var response models.GetScrapShipResponse

	err := get(endpoint, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// ScrapShip scraps a models.Ship
func ScrapShip(post PostFunc, ShipSymbol string) (*models.ScrapShipResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/scrap", ShipSymbol)

	var response models.ScrapShipResponse

	err := post(endpoint, nil, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetRepairShip retrieves the repair details of a specific models.Ship
func GetRepairShip(get GetFunc, ShipSymbol string) (*models.GetRepairShipResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/repair", ShipSymbol)

	var response models.GetRepairShipResponse

	err := get(endpoint, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// RepairShip repairs a models.Ship
func RepairShip(post PostFunc, ShipSymbol string) (*models.RepairShipResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/repair", ShipSymbol)

	var response models.RepairShipResponse

	err := post(endpoint, nil, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// GetModules retrieves the modules installed on a specific models.Ship
func GetModules(get GetFunc, ShipSymbol string) (*models.GetModulesResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/modules", ShipSymbol)

	var response models.GetModulesResponse

	err := get(endpoint, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// InstallModule installs a module on a models.Ship
func InstallModule(post PostFunc, ShipSymbol string, payload *models.InstallModuleRequest) (*models.InstallModuleResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/modules/install", ShipSymbol)

	var response models.InstallModuleResponse

	err := post(endpoint, payload, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// RemoveModule removes a module from a models.Ship
func RemoveModule(post PostFunc, ShipSymbol string, payload *models.RemoveModuleRequest) (*models.RemoveModuleResponse, *models.APIError) {
	endpoint := fmt.Sprintf("/my/ships/%s/modules/remove", ShipSymbol)

	var response models.RemoveModuleResponse

	err := post(endpoint, payload, nil, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}
