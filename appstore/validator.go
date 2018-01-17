package appstore

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

const (
	// SandboxURL is the endpoint for sandbox environment.
	SandboxURL string = "https://sandbox.itunes.apple.com/verifyReceipt"
	// ProductionURL is the endpoint for production environment.
	ProductionURL string = "https://buy.itunes.apple.com/verifyReceipt"
)

// Config is a configuration to initialize client
type Config struct {
	TimeOut time.Duration
}

// IAPClient is an interface to call validation API in App Store
type IAPClient interface {
	Verify(IAPRequest, interface{}) error
}

// Client implements IAPClient
type Client struct {
	ProductionURL string
	SandboxURL    string
	TimeOut       time.Duration
}

// HandleError returns error message by status code
func HandleError(status int) error {
	var message string

	switch status {
	case 0:
		return nil

	case 21000:
		message = "The App Store could not read the JSON object you provided."

	case 21002:
		message = "The data in the receipt-data property was malformed or missing."

	case 21003:
		message = "The receipt could not be authenticated."

	case 21004:
		message = "The shared secret you provided does not match the shared secret on file for your account."

	case 21005:
		message = "The receipt server is not currently available."

	case 21007:
		message = "This receipt is from the test environment, but it was sent to the production environment for verification. Send it to the test environment instead."

	case 21008:
		message = "This receipt is from the production environment, but it was sent to the test environment for verification. Send it to the production environment instead."

	case 21010:
		message = "This receipt could not be authorized. Treat this the same as if a purchase was never made."

	default:
		if status >= 21100 && status <= 21199 {
			message = "Internal data access error."
		} else {
			message = "An unknown error occurred"
		}
	}

	return errors.New(message)
}

// New creates a client object
func New() Client {
	client := Client{
		ProductionURL: ProductionURL,
		SandboxURL:    SandboxURL,
		TimeOut:       time.Second * 5,
	}
	return client
}

// NewWithConfig creates a client with configuration
func NewWithConfig(config Config) Client {
	if config.TimeOut == 0 {
		config.TimeOut = time.Second * 5
	}

	client := Client{
		ProductionURL: ProductionURL,
		SandboxURL:    SandboxURL,
		TimeOut:       config.TimeOut,
	}

	return client
}

// Verify sends receipts and gets validation result
func (c *Client) Verify(req IAPRequest, result interface{}) error {
	client := http.Client{
		Timeout: c.TimeOut,
	}

	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(req)

	resp, err := client.Post(c.ProductionURL, "application/json; charset=utf-8", b)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(result)
	if err != nil {
		return err
	}

	// https://developer.apple.com/library/content/technotes/tn2413/_index.html#//apple_ref/doc/uid/DTS40016228-CH1-RECEIPTURL
	r, ok := result.(*HttpStatusResponse)
	if ok && r.Status == 21007 {
		b = new(bytes.Buffer)
		json.NewEncoder(b).Encode(req)
		resp, err := client.Post(c.SandboxURL, "application/json; charset=utf-8", b)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		return json.NewDecoder(resp.Body).Decode(result)
	}

	return nil
}
