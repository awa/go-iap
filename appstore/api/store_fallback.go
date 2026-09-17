package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// NewStoreClientWithSandboxFallback creates a client that looks a transaction
// up in production and retries against Sandbox when production cannot answer
// for it, as App Review and TestFlight purchases need. config.Sandbox is
// ignored: both hosts are built.
func NewStoreClientWithSandboxFallback(config *StoreConfig) *StoreClientWithSandboxFallback {
	prodConf, sandboxConf := *config, *config
	prodConf.Sandbox, sandboxConf.Sandbox = false, true
	return &StoreClientWithSandboxFallback{
		productionCli: NewStoreClient(&prodConf),
		sandboxCli:    NewStoreClient(&sandboxConf),
	}
}

type StoreClientWithSandboxFallback struct {
	productionCli StoreAPIClient
	sandboxCli    StoreAPIClient
}

// hasFallback reports whether Sandbox may answer a production failure: the
// transaction lives in the other environment, or production refuses the bundle
// ID because the app has never shipped there. It hangs off this client so
// StoreClient, which has no Sandbox host to ask, cannot reach it.
func (c *StoreClientWithSandboxFallback) hasFallback(err error) bool {
	if errors.Is(err, TransactionIdNotFoundError) {
		return true
	}
	var statusErr *httpStatusError
	return errors.As(err, &statusErr) && statusErr.StatusCode() == http.StatusUnauthorized
}

func (c *StoreClientWithSandboxFallback) GetALLSubscriptionStatuses(ctx context.Context, originalTransactionId string, query *url.Values) (*StatusResponse, error) {
	rsp, productionErr := c.productionCli.GetALLSubscriptionStatuses(ctx, originalTransactionId, query)
	if !c.hasFallback(productionErr) {
		return rsp, productionErr
	}
	rsp, err := c.sandboxCli.GetALLSubscriptionStatuses(ctx, originalTransactionId, query)
	if err == nil {
		return rsp, nil
	}
	return rsp, fmt.Errorf("production: %v; sandbox: %w", productionErr, err)
}

func (c *StoreClientWithSandboxFallback) GetTransactionInfo(ctx context.Context, transactionId string) (*TransactionInfoResponse, error) {
	rsp, productionErr := c.productionCli.GetTransactionInfo(ctx, transactionId)
	if !c.hasFallback(productionErr) {
		return rsp, productionErr
	}
	rsp, err := c.sandboxCli.GetTransactionInfo(ctx, transactionId)
	if err == nil {
		return rsp, nil
	}
	return rsp, fmt.Errorf("production: %v; sandbox: %w", productionErr, err)
}

func (c *StoreClientWithSandboxFallback) Verify(ctx context.Context, transactionId string) (*TransactionInfoResponse, error) {
	return c.GetTransactionInfo(ctx, transactionId)
}

func (c *StoreClientWithSandboxFallback) GetTransactionHistory(ctx context.Context, transactionId string, query *url.Values) ([]*HistoryResponse, error) {
	responses, productionErr := c.productionCli.GetTransactionHistory(ctx, transactionId, query)
	if !c.hasFallback(productionErr) {
		return responses, productionErr
	}
	responses, err := c.sandboxCli.GetTransactionHistory(ctx, transactionId, query)
	if err == nil {
		return responses, nil
	}
	return responses, fmt.Errorf("production: %v; sandbox: %w", productionErr, err)
}

func (c *StoreClientWithSandboxFallback) GetRefundHistory(ctx context.Context, originalTransactionId string) ([]*RefundLookupResponse, error) {
	responses, productionErr := c.productionCli.GetRefundHistory(ctx, originalTransactionId)
	if !c.hasFallback(productionErr) {
		return responses, productionErr
	}
	responses, err := c.sandboxCli.GetRefundHistory(ctx, originalTransactionId)
	if err == nil {
		return responses, nil
	}
	return responses, fmt.Errorf("production: %v; sandbox: %w", productionErr, err)
}

func (c *StoreClientWithSandboxFallback) GetAppTransactionInfo(ctx context.Context, transactionId string) (*AppTransactionInfoResponse, error) {
	rsp, productionErr := c.productionCli.GetAppTransactionInfo(ctx, transactionId)
	if !c.hasFallback(productionErr) {
		return rsp, productionErr
	}
	rsp, err := c.sandboxCli.GetAppTransactionInfo(ctx, transactionId)
	if err == nil {
		return rsp, nil
	}
	return rsp, fmt.Errorf("production: %v; sandbox: %w", productionErr, err)
}

// The rest go to production only. Writes are not repeated elsewhere, and the
// remaining reads fail in ways Sandbox cannot answer.

func (c *StoreClientWithSandboxFallback) GetSubscriptionRenewalDataStatus(ctx context.Context, productId, requestIdentifier string) (int, *MassExtendRenewalDateStatusResponse, error) {
	return c.productionCli.GetSubscriptionRenewalDataStatus(ctx, productId, requestIdentifier)
}

func (c *StoreClientWithSandboxFallback) LookupOrderID(ctx context.Context, orderId string) (*OrderLookupResponse, error) {
	return c.productionCli.LookupOrderID(ctx, orderId)
}

func (c *StoreClientWithSandboxFallback) FinishTransaction(ctx context.Context, transactionId string) (int, error) {
	return c.productionCli.FinishTransaction(ctx, transactionId)
}

func (c *StoreClientWithSandboxFallback) ExtendSubscriptionRenewalDate(ctx context.Context, originalTransactionId string, body ExtendRenewalDateRequest) (int, error) {
	return c.productionCli.ExtendSubscriptionRenewalDate(ctx, originalTransactionId, body)
}

func (c *StoreClientWithSandboxFallback) ExtendSubscriptionRenewalDateForAll(ctx context.Context, body MassExtendRenewalDateRequest) (int, error) {
	return c.productionCli.ExtendSubscriptionRenewalDateForAll(ctx, body)
}

func (c *StoreClientWithSandboxFallback) ParseSignedTransactions(transactions []string) ([]*JWSTransaction, error) {
	return c.productionCli.ParseSignedTransactions(transactions)
}

func (c *StoreClientWithSandboxFallback) ParseJWSEncodeString(jwsEncode string) (interface{}, error) {
	return c.productionCli.ParseJWSEncodeString(jwsEncode)
}

func (c *StoreClientWithSandboxFallback) ParseSignedTransaction(transaction string) (*JWSTransaction, error) {
	return c.productionCli.ParseSignedTransaction(transaction)
}

func (c *StoreClientWithSandboxFallback) GetAllNotificationHistory(ctx context.Context, body NotificationHistoryRequest, duration time.Duration) ([]NotificationHistoryResponseItem, error) {
	return c.productionCli.GetAllNotificationHistory(ctx, body, duration)
}

func (c *StoreClientWithSandboxFallback) GetNotificationHistory(ctx context.Context, body NotificationHistoryRequest, paginationToken string) (*NotificationHistoryResponses, error) {
	return c.productionCli.GetNotificationHistory(ctx, body, paginationToken)
}

func (c *StoreClientWithSandboxFallback) GetTestNotificationStatus(ctx context.Context, testNotificationToken string) (int, []byte, error) {
	return c.productionCli.GetTestNotificationStatus(ctx, testNotificationToken)
}

func (c *StoreClientWithSandboxFallback) SendRequestTestNotification(ctx context.Context) (int, []byte, error) {
	return c.productionCli.SendRequestTestNotification(ctx)
}

func (c *StoreClientWithSandboxFallback) SendConsumptionInfo(ctx context.Context, originalTransactionId string, body ConsumptionRequestBody) (int, error) {
	return c.productionCli.SendConsumptionInfo(ctx, originalTransactionId, body)
}

func (c *StoreClientWithSandboxFallback) SendConsumptionInfoV2(ctx context.Context, transactionId string, body ConsumptionRequest) (int, error) {
	return c.productionCli.SendConsumptionInfoV2(ctx, transactionId, body)
}

func (c *StoreClientWithSandboxFallback) SetAppAccountToken(ctx context.Context, originalTransactionId string, body UpdateAppAccountTokenRequest) (int, error) {
	return c.productionCli.SetAppAccountToken(ctx, originalTransactionId, body)
}

func (c *StoreClientWithSandboxFallback) Do(ctx context.Context, method string, url string, body io.Reader) (int, []byte, error) {
	return c.productionCli.Do(ctx, method, url, body)
}
