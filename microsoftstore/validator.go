package microsoftstore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	resource = "https://onestore.microsoft.com"
)

// IAPClient is an interface to call validation API in Microsoft Store
type IAPClient interface {
	Verify(context.Context, string, string) (IAPResponse, error)
}

// Client implements IAPClient
type Client struct {
	TenantID     string
	ClientID     string
	ClientSecret string
	httpCli      *http.Client
}

// New creates a client object
func New(tenantId, clientId, secret string) *Client {
	client := &Client{
		TenantID:     tenantId,
		ClientID:     clientId,
		ClientSecret: secret,
		httpCli: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	return client
}

// Verify sends receipts and gets validation result
func (c *Client) Verify(ctx context.Context, receipt IAPRequest) (IAPResponse, error) {
	resp := IAPResponse{}
	token, err := c.getAzureADToken(ctx, c.TenantID, c.ClientID, c.ClientSecret, resource)
	if err != nil {
		return resp, err
	}

	return c.query(ctx, token, receipt)
}

// getAzureADToken obtains an Azure AD access token using client credentials flow
func (c *Client) getAzureADToken(ctx context.Context, tenantID, clientID, clientSecret, resource string) (string, error) {
	tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/token", tenantID)

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("resource", resource)

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.WithContext(ctx)

	resp, err := c.httpCli.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("failed to obtain token: %s", string(bodyBytes))
	}

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
	}
	err = json.NewDecoder(resp.Body).Decode(&tokenResponse)
	if err != nil {
		return "", err
	}
	return tokenResponse.AccessToken, nil
}

// query sends a query to Microsoft Store API
func (c *Client) query(ctx context.Context, accessToken string, receiptData IAPRequest) (IAPResponse, error) {
	queryURL := "https://collections.mp.microsoft.com/v6.0/collections/query"
	result := IAPResponse{}

	requestBody, err := json.Marshal(receiptData)
	if err != nil {
		return result, err
	}

	req, err := http.NewRequest("POST", queryURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return result, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.WithContext(ctx)

	res, err := c.httpCli.Do(req)
	if err != nil {
		return result, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		return result, fmt.Errorf("validation failed: %s", string(bodyBytes))
	}

	err = json.NewDecoder(res.Body).Decode(&result)
	return result, err
}
