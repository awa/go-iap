package api

import (
	"context"
	"errors"
)

// IAPAPIClient is an interface to call validation API in App Store Server API
//
// Deprecated: use TransactionVerifier, which StoreAPIClient includes.
type IAPAPIClient interface {
	Verify(ctx context.Context, transactionId string) (*TransactionInfoResponse, error)
}

// Deprecated: use NewStoreClientWithSandboxFallback. It retries every lookup
// Apple documents as returning TransactionIdNotFoundError, not Verify alone,
// and satisfies StoreAPIClient.
type APIClient struct {
	productionCli *StoreClient
	sandboxCli    *StoreClient
}

// Deprecated: use NewStoreClientWithSandboxFallback.
func NewAPIClient(config StoreConfig) *APIClient {
	prodConf := config
	prodConf.Sandbox = false
	sandboxConf := config
	sandboxConf.Sandbox = true
	return &APIClient{productionCli: NewStoreClient(&prodConf), sandboxCli: NewStoreClient(&sandboxConf)}
}

// Deprecated: use the Verify of a client from NewStoreClientWithSandboxFallback.
func (c *APIClient) Verify(ctx context.Context, transactionId string) (*TransactionInfoResponse, error) {
	result, err := c.productionCli.GetTransactionInfo(ctx, transactionId)
	if err != nil && errors.Is(err, TransactionIdNotFoundError) {
		result, err = c.sandboxCli.GetTransactionInfo(ctx, transactionId)
	}
	return result, err
}
