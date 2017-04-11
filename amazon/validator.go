package amazon

import (
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

// Config is a configuration to initialize client
type Config struct {
	IsProduction bool
	Secret       string
	TimeOut      time.Duration
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
	Verify(string, string) (IAPResponse, error)
}

// Client implements IAPClient
type Client struct {
	URL     string
	Secret  string
	TimeOut time.Duration
}

// New creates a client object
func New(secret string) IAPClient {
	client := Client{
		URL:     getSandboxURL(),
		Secret:  secret,
		TimeOut: time.Second * 5,
	}
	if os.Getenv("IAP_ENVIRONMENT") == "production" {
		client.URL = ProductionURL
	}
	return client
}

// NewWithConfig creates a client with configuration
func NewWithConfig(config Config) Client {
	if config.TimeOut == 0 {
		config.TimeOut = time.Second * 5
	}

	client := Client{
		URL:     getSandboxURL(),
		Secret:  config.Secret,
		TimeOut: config.TimeOut,
	}
	if config.IsProduction {
		client.URL = ProductionURL
	}

	return client
}

// Verify sends receipts and gets validation result
func (c Client) Verify(userID string, receiptID string) (IAPResponse, error) {
	result := IAPResponse{}
	url := fmt.Sprintf("%v/version/1.0/verifyReceiptId/developer/%v/user/%v/receiptId/%v", c.URL, c.Secret, userID, receiptID)
	client := http.Client{
		Timeout: c.TimeOut,
	}
	resp, err := client.Get(url)
	if err != nil {
		return result, fmt.Errorf("%v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		responseError := IAPResponseError{}
		err = json.NewDecoder(resp.Body).Decode(&responseError)
		return result, errors.New(responseError.Message)
	}

	err = json.NewDecoder(resp.Body).Decode(&result)

	return result, err
}
