package playstore

import (
	"errors"
	"net"
	"net/http"
	"os"
	"time"

	"code.google.com/p/goauth2/oauth"
	"code.google.com/p/google-api-go-client/androidpublisher/v2"
)

const (
	scope    = "https://www.googleapis.com/auth/androidpublisher"
	authURL  = "https://accounts.google.com/o/oauth2/auth"
	tokenURL = "https://accounts.google.com/o/oauth2/token"

	timeout = time.Second * 5
)

var defaultConfig *oauth.Config
var defaultTimeout = timeout

// Init initializes the global configuration
func Init() error {
	defaultConfig = &oauth.Config{
		Scope:    scope,
		AuthURL:  authURL,
		TokenURL: tokenURL,
	}

	clientID := os.Getenv("IAB_CLIENT_ID")
	if clientID != "" {
		defaultConfig.ClientId = clientID
	}
	if defaultConfig.ClientId == "" {
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
func InitWithConfig(config *oauth.Config) error {
	if config.ClientId == "" {
		return errors.New("Client ID is required")
	}

	if config.ClientSecret == "" {
		return errors.New("Client Secret Key is required")
	}

	if config.Scope == "" {
		config.Scope = scope
	}

	if config.AuthURL == "" {
		config.AuthURL = authURL
	}

	if config.TokenURL == "" {
		config.TokenURL = tokenURL
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
}

// The Client type implements VerifySubscription method
type Client struct {
	httpClient *http.Client
}

// New returns http client which has oauth token
func New(token *oauth.Token) Client {
	t := &oauth.Transport{
		Token:  token,
		Config: defaultConfig,
		Transport: &http.Transport{
			Dial: dialTimeout,
		},
	}

	httpClient := t.Client()
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

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, defaultTimeout)
}
