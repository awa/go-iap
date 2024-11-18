//go:generate mockgen  -destination=mocks/appstore.go -package=mocks github.com/awa/go-iap/appstore IAPClient

package appstore

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	// SandboxURL is the endpoint for sandbox environment.
	SandboxURL string = "https://sandbox.itunes.apple.com/verifyReceipt"
	// ProductionURL is the endpoint for production environment.
	ProductionURL string = "https://buy.itunes.apple.com/verifyReceipt"
	// ContentType is the request content-type for apple store.
	ContentType string = "application/json; charset=utf-8"
)

// IAPClient is an interface to call validation API in App Store
type IAPClient interface {
	Verify(ctx context.Context, reqBody IAPRequest, resp interface{}) error
	VerifyWithStatus(ctx context.Context, reqBody IAPRequest, resp interface{}) (int, error)
	ParseNotificationV2(tokenStr string, result *jwt.Token) error
	ParseNotificationV2WithClaim(tokenStr string, result jwt.Claims) error
}

// Client implements IAPClient
type Client struct {
	ProductionURL string
	SandboxURL    string
	httpCli       *http.Client
}

// list of errors
var (
	ErrAppStoreServer = errors.New("AppStore server error")

	ErrInvalidJSON            = errors.New("The App Store could not read the JSON object you provided.")
	ErrInvalidReceiptData     = errors.New("The data in the receipt-data property was malformed or missing.")
	ErrReceiptUnauthenticated = errors.New("The receipt could not be authenticated.")
	ErrInvalidSharedSecret    = errors.New("The shared secret you provided does not match the shared secret on file for your account.")
	ErrServerUnavailable      = errors.New("The receipt server is not currently available.")
	ErrReceiptIsForTest       = errors.New("This receipt is from the test environment, but it was sent to the production environment for verification. Send it to the test environment instead.")
	ErrReceiptIsForProduction = errors.New("This receipt is from the production environment, but it was sent to the test environment for verification. Send it to the production environment instead.")
	ErrReceiptUnauthorized    = errors.New("This receipt could not be authorized. Treat this the same as if a purchase was never made.")

	ErrInternalDataAccessError = errors.New("Internal data access error.")
	ErrUnknown                 = errors.New("An unknown error occurred")
)

// HandleError returns error message by status code
func HandleError(status int) error {
	var e error
	switch status {
	case 0:
		return nil
	case 21000:
		e = ErrInvalidJSON
	case 21002:
		e = ErrInvalidReceiptData
	case 21003:
		e = ErrReceiptUnauthenticated
	case 21004:
		e = ErrInvalidSharedSecret
	case 21005:
		e = ErrServerUnavailable
	case 21007:
		e = ErrReceiptIsForTest
	case 21008:
		e = ErrReceiptIsForProduction
	case 21009:
		e = ErrInternalDataAccessError
	case 21010:
		e = ErrReceiptUnauthorized
	default:
		if status >= 21100 && status <= 21199 {
			e = ErrInternalDataAccessError
		} else {
			e = ErrUnknown
		}
	}

	return fmt.Errorf("status %d: %w", status, e)
}

// New creates a client object
func New() *Client {
	client := &Client{
		ProductionURL: ProductionURL,
		SandboxURL:    SandboxURL,
		httpCli: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
	return client
}

// NewWithClient creates a client with a custom http client.
func NewWithClient(client *http.Client) *Client {
	return &Client{
		ProductionURL: ProductionURL,
		SandboxURL:    SandboxURL,
		httpCli:       client,
	}
}

// Verify sends receipts and gets validation result
func (c *Client) Verify(ctx context.Context, reqBody IAPRequest, result interface{}) error {
	_, err := c.verify(ctx, reqBody, result)
	return err
}

// VerifyWithStatus sends receipts and gets validation result with status code
// If the Apple verification receipt server is unhealthy and responds with an HTTP status code in the 5xx range, that status code will be returned.
func (c *Client) VerifyWithStatus(ctx context.Context, reqBody IAPRequest, result interface{}) (int, error) {
	return c.verify(ctx, reqBody, result)
}

func (c *Client) verify(ctx context.Context, reqBody IAPRequest, result interface{}) (int, error) {
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(reqBody); err != nil {
		return 0, err
	}

	req, err := http.NewRequest("POST", c.ProductionURL, b)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Content-Type", ContentType)
	req = req.WithContext(ctx)
	resp, err := c.httpCli.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 500 {
		return resp.StatusCode, fmt.Errorf("Received http status code %d from the App Store: %w", resp.StatusCode, ErrAppStoreServer)
	}
	return c.parseResponse(resp, result, ctx, reqBody)
}

func (c *Client) parseResponse(resp *http.Response, result interface{}, ctx context.Context, reqBody IAPRequest) (int, error) {
	// Read the body now so that we can unmarshal it twice
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	err = json.Unmarshal(buf, &result)
	if err != nil {
		return 0, err
	}

	// https://developer.apple.com/library/content/technotes/tn2413/_index.html#//apple_ref/doc/uid/DTS40016228-CH1-RECEIPTURL
	var r StatusResponse
	err = json.Unmarshal(buf, &r)
	if err != nil {
		return 0, err
	}
	if r.Status == 21007 {
		b := new(bytes.Buffer)
		if err := json.NewEncoder(b).Encode(reqBody); err != nil {
			return 0, err
		}

		req, err := http.NewRequest("POST", c.SandboxURL, b)
		if err != nil {
			return 0, err
		}
		req.Header.Set("Content-Type", ContentType)
		req = req.WithContext(ctx)
		resp, err := c.httpCli.Do(req)
		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 500 {
			return resp.StatusCode, fmt.Errorf("Received http status code %d from the App Store Sandbox: %w", resp.StatusCode, ErrAppStoreServer)
		}

		return r.Status, json.NewDecoder(resp.Body).Decode(result)
	}

	return r.Status, nil
}

// ParseNotificationV2 parse notification from App Store Server
func (c *Client) ParseNotificationV2(tokenStr string, result *jwt.Token) error {
	cert := Cert{}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return cert.ExtractPublicKeyFromToken(tokenStr)
	})
	if token != nil {
		*result = *token
	}
	return err
}

// ParseNotificationV2WithClaim parse notification from App Store Server
func (c *Client) ParseNotificationV2WithClaim(tokenStr string, result jwt.Claims) error {
	cert := Cert{}

	_, err := jwt.ParseWithClaims(tokenStr, result, func(token *jwt.Token) (interface{}, error) {
		return cert.ExtractPublicKeyFromToken(tokenStr)
	})
	return err
}
