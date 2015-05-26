package playstore

import (
	"errors"
	"net/http"
	"os"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	androidpublisher "google.golang.org/api/androidpublisher/v2"
)

const (
	scope    = "https://www.googleapis.com/auth/androidpublisher"
	authURL  = "https://accounts.google.com/o/oauth2/auth"
	tokenURL = "https://accounts.google.com/o/oauth2/token"

	timeout = time.Second * 5
)

var defaultConfig *oauth2.Config
var defaultTimeout = timeout

// Init initializes the global configuration
func Init() error {
	defaultConfig = &oauth2.Config{
		Scopes: []string{scope},
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}

	clientID := os.Getenv("IAB_CLIENT_ID")
	if clientID != "" {
		defaultConfig.ClientID = clientID
	}
	if defaultConfig.ClientID == "" {
		return errors.New("Client ID is required")
	}

	clientSecret := os.Getenv("IAB_CLIENT_SECRET")
	if clientSecret != "" {
		defaultConfig.ClientSecret = clientSecret
	}
	if defaultConfig.ClientSecret == "" {
		return errors.New("Client Secret Key is required")
	}

	return nil
}

// InitWithConfig initializes the global configuration with parameters
func InitWithConfig(config *oauth2.Config) error {
	if config.ClientID == "" {
		return errors.New("Client ID is required")
	}

	if config.ClientSecret == "" {
		return errors.New("Client Secret Key is required")
	}

	if len(config.Scopes) == 0 {
		config.Scopes = []string{scope}
	}

	if config.Endpoint.AuthURL == "" {
		config.Endpoint.AuthURL = authURL
	}

	if config.Endpoint.TokenURL == "" {
		config.Endpoint.TokenURL = tokenURL
	}

	defaultConfig = config

	return nil
}

// SetTimeout sets dial timeout duration
func SetTimeout(t time.Duration) {
	defaultTimeout = t
}

// The IABClient type is an interface to verify purchase token
type IABClient interface {
	VerifySubscription(string, string, string) (*androidpublisher.SubscriptionPurchase, error)
	VerifyProduct(string, string, string) (*androidpublisher.ProductPurchase, error)
}

// The Client type implements VerifySubscription method
type Client struct {
	httpClient *http.Client
}

// New returns http client which has oauth token
func New(token *oauth2.Token) Client {
	ctx := context.WithValue(oauth2.NoContext, oauth2.HTTPClient, &http.Client{
		Timeout: defaultTimeout,
	})

	httpClient := defaultConfig.Client(ctx, token)

	return Client{httpClient}
}

// VerifySubscription Verifies subscription status
func (c *Client) VerifySubscription(
	packageName string,
	subscriptionID string,
	token string,
) (*androidpublisher.SubscriptionPurchase, error) {
	service, err := androidpublisher.New(c.httpClient)
	if err != nil {
		return nil, err
	}

	ps := androidpublisher.NewPurchasesSubscriptionsService(service)
	result, err := ps.Get(packageName, subscriptionID, token).Do()

	return result, err
}

// VerifyProduct Verifies product status
func (c *Client) VerifyProduct(
	packageName string,
	productID string,
	token string,
) (*androidpublisher.ProductPurchase, error) {
	service, err := androidpublisher.New(c.httpClient)
	if err != nil {
		return nil, err
	}

	ps := androidpublisher.NewPurchasesProductsService(service)
	result, err := ps.Get(packageName, productID, token).Do()

	return result, err
}
