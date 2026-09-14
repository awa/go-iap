package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/url"
	"time"
)

// NewStoreClientWithSandboxFallback creates a client that looks a transaction
// up in production and retries against Sandbox when production reports it
// missing, as App Review and TestFlight purchases need. config.Sandbox is
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

// Lookups Apple documents as returning TransactionIdNotFoundError. Only that
// error is retried: other failures are caller mistakes a retry would hide.

func (c *StoreClientWithSandboxFallback) GetALLSubscriptionStatuses(ctx context.Context, originalTransactionId string, query *url.Values) (*StatusResponse, error) {
	rsp, err := c.productionCli.GetALLSubscriptionStatuses(ctx, originalTransactionId, query)
	if !errors.Is(err, TransactionIdNotFoundError) {
		return rsp, err
	}
	rsp, err = c.sandboxCli.GetALLSubscriptionStatuses(ctx, originalTransactionId, query)
	return rsp, sandboxRetryErr(err)
}

func (c *StoreClientWithSandboxFallback) GetTransactionInfo(ctx context.Context, transactionId string) (*TransactionInfoResponse, error) {
	rsp, err := c.productionCli.GetTransactionInfo(ctx, transactionId)
	if !errors.Is(err, TransactionIdNotFoundError) {
		return rsp, err
	}
	rsp, err = c.sandboxCli.GetTransactionInfo(ctx, transactionId)
	return rsp, sandboxRetryErr(err)
}

func (c *StoreClientWithSandboxFallback) GetTransactionHistory(ctx context.Context, transactionId string, query *url.Values) ([]*HistoryResponse, error) {
	responses, err := c.productionCli.GetTransactionHistory(ctx, transactionId, query)
	if !errors.Is(err, TransactionIdNotFoundError) {
		return responses, err
	}
	responses, err = c.sandboxCli.GetTransactionHistory(ctx, transactionId, query)
	return responses, sandboxRetryErr(err)
}

func (c *StoreClientWithSandboxFallback) GetRefundHistory(ctx context.Context, originalTransactionId string) ([]*RefundLookupResponse, error) {
	responses, err := c.productionCli.GetRefundHistory(ctx, originalTransactionId)
	if !errors.Is(err, TransactionIdNotFoundError) {
		return responses, err
	}
	responses, err = c.sandboxCli.GetRefundHistory(ctx, originalTransactionId)
	return responses, sandboxRetryErr(err)
}

func (c *StoreClientWithSandboxFallback) GetAppTransactionInfo(ctx context.Context, transactionId string) (*AppTransactionInfoResponse, error) {
	rsp, err := c.productionCli.GetAppTransactionInfo(ctx, transactionId)
	if !errors.Is(err, TransactionIdNotFoundError) {
		return rsp, err
	}
	rsp, err = c.sandboxCli.GetAppTransactionInfo(ctx, transactionId)
	return rsp, sandboxRetryErr(err)
}

// The rest go to production only: the writes return the same error but
// repeating a write elsewhere is a separate decision, Extend and
// SetAppAccountToken report OriginalTransactionIdNotFoundError instead, and a
// retry on GetNotificationHistory would be skipped by its bulk caller.

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

// sandboxRetryErr wraps the Sandbox error, not the production one: missing in
// both environments still satisfies errors.Is, while a Sandbox outage stays
// distinguishable from a missing transaction.
func sandboxRetryErr(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("production reported the transaction missing; sandbox: %w", err)
}
