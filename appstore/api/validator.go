package api

import (
	"context"
	"errors"
)

// IAPAPIClient is an interface to call validation API in App Store Server API
type IAPAPIClient interface {
	Verify(ctx context.Context, transactionId string) (interface{}, error)
}

type APIClient struct {
	productionCli *StoreClient
	sandboxCli    *StoreClient
}

func NewAPIClient(config StoreConfig) *APIClient {
	prodConf := config
	prodConf.Sandbox = false
	sandboxConf := config
	sandboxConf.Sandbox = true
	return &APIClient{productionCli: NewStoreClient(&prodConf), sandboxCli: NewStoreClient(&sandboxConf)}
}

func (c *APIClient) Verify(ctx context.Context, transactionId string) (interface{}, error) {
	result, err := c.productionCli.GetTransactionInfo(ctx, transactionId)
	if err != nil && errors.Is(err, TransactionIdNotFoundError) {
		result, err = c.sandboxCli.GetTransactionInfo(ctx, transactionId)
	}
	return result, err
}
