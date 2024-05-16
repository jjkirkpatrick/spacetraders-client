package entities

import (
	"github.com/jjkirkpatrick/spacetraders-client/client"
	"github.com/jjkirkpatrick/spacetraders-client/internal/api"
	"github.com/jjkirkpatrick/spacetraders-client/models"
)

type Contract struct {
	models.Contract
	Client *client.Client
}

func ListContracts(c *client.Client) ([]*Contract, error) {
	fetchFunc := func(meta models.Meta) ([]*Contract, models.Meta, error) {
		metaPtr := &meta
		contracts, metaPtr, err := api.ListContracts(c.Get, metaPtr)

		var convertedContracts []*Contract
		for _, modelContract := range contracts {
			convertedContract := &Contract{
				Contract: *modelContract, // Directly embed the modelContract
				Client:   c,
			}
			convertedContracts = append(convertedContracts, convertedContract)
		}

		if err != nil {
			if metaPtr == nil {
				// Use default Meta values or handle accordingly
				defaultMeta := models.Meta{Page: 1, Limit: 20, Total: 0}
				metaPtr = &defaultMeta
			}
			return convertedContracts, *metaPtr, err.AsError()
		}
		if metaPtr != nil {
			return convertedContracts, *metaPtr, nil
		} else {
			defaultMeta := models.Meta{Page: 1, Limit: 20, Total: 0}
			return convertedContracts, defaultMeta, nil
		}
	}
	return client.NewPaginator[*Contract](fetchFunc).FetchAllPages()
}

func GetContract(c *client.Client, symbol string) (*Contract, error) {
	contract, err := api.GetContract(c.Get, symbol)
	if err != nil {
		return nil, err
	}

	contractEntity := &Contract{
		Contract: *contract,
		Client:   c,
	}

	return contractEntity, nil
}

func (c *Contract) Accept() (*Agent, *Contract, error) {
	agent, contract, err := api.AcceptContract(c.Client.Post, c.Contract.ID)
	if err != nil {
		return nil, nil, err
	}

	return &Agent{Agent: *agent, Client: c.Client}, &Contract{Contract: *contract, Client: c.Client}, nil
}

func (c *Contract) DeliverCargo(shop *Ship, tradeGood models.GoodSymbol, units int) (*Contract, *models.Cargo, error) {

	contractRequest := models.DeliverContractCargoRequest{
		ShipSymbol:  shop.Symbol,
		TradeSymbol: tradeGood,
		Units:       units,
	}

	agent, cargo, err := api.DeliverContractCargo(c.Client.Post, c.Contract.ID, contractRequest)
	if err != nil {
		return nil, nil, err
	}

	return &Contract{Contract: *agent, Client: c.Client}, cargo, nil
}

func (c *Contract) Fulfill() (*models.Agent, *models.Contract, error) {
	agent, contract, err := api.FulfillContract(c.Client.Post, c.Contract.ID)
	if err != nil {
		return nil, nil, err
	}

	return agent, contract, nil
}
