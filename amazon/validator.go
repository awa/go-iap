package amazon

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	// SandboxURL is the endpoint for local environment.
	SandboxURL string = "http://localhost:8080/RVSSandbox"
	// ProductionURL is the endpoint for production environment.
	ProductionURL string = "https://appstore-sdk.amazon.com"
)

func getSandboxURL() string {
	url := os.Getenv("IAP_SANDBOX_URL")
	if url == "" {
		url = SandboxURL
	}
	return url
}

// The IAPResponse type has the response properties
type IAPResponse struct {
	ReceiptID       string `json:"receiptId"`
	ProductType     string `json:"productType"`
	ProductID       string `json:"productId"`
	PurchaseDate    int64  `json:"purchaseDate"`
	CancelDate      int64  `json:"cancelDate"`
	TestTransaction bool   `json:"testTransaction"`
}

// The IAPResponseError typs has error message and status.
type IAPResponseError struct {
	Message string `json:"message"`
	Status  bool   `json:"status"`
}

// IAPClient is an interface to call validation API in Amazon App Store
type IAPClient interface {
	Verify(context.Context, string, string) (IAPResponse, error)
}

// Client implements IAPClient
type Client struct {
	URL     string
	Secret  string
	httpCli *http.Client
}

// New creates a client object
func New(secret string) *Client {
	client := &Client{
		URL:    getSandboxURL(),
		Secret: secret,
		httpCli: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	if os.Getenv("IAP_ENVIRONMENT") == "production" {
		client.URL = ProductionURL
	}

	return client
}

// NewWithClient creates a client with a custom client.
func NewWithClient(secret string, cli *http.Client) *Client {
	client := &Client{
		URL:     getSandboxURL(),
		Secret:  secret,
		httpCli: cli,
	}
	if os.Getenv("IAP_ENVIRONMENT") == "production" {
		client.URL = ProductionURL
	}

	return client
}

// Verify sends receipts and gets validation result
func (c *Client) Verify(ctx context.Context, userID string, receiptID string) (IAPResponse, error) {
	result := IAPResponse{}
	url := fmt.Sprintf("%v/version/1.0/verifyReceiptId/developer/%v/user/%v/receiptId/%v", c.URL, c.Secret, userID, receiptID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return result, err
	}
	req = req.WithContext(ctx)

	resp, err := c.httpCli.Do(req)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		responseError := IAPResponseError{}
		err = json.NewDecoder(resp.Body).Decode(&responseError)
		if err != nil {
			return result, err
		}
		return result, errors.New(responseError.Message)
	}

	err = json.NewDecoder(resp.Body).Decode(&result)

	return result, err
}
