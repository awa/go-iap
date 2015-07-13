package playstore

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	androidpublisher "google.golang.org/api/androidpublisher/v2"
)

const (
	defaultTimeout = time.Second * 5
)

var timeout = defaultTimeout

// SetTimeout sets dial timeout duration
func SetTimeout(t time.Duration) {
	timeout = t
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

// New returns http client which includes the credentials to access androidpublisher API.
// You should create a service account for your project at
// https://console.developers.google.com and download a JSON key file to set this argument.
func New(jsonKey []byte) (Client, error) {
	ctx := context.WithValue(oauth2.NoContext, oauth2.HTTPClient, &http.Client{
		Timeout: timeout,
	})

	conf, err := google.JWTConfigFromJSON(jsonKey, androidpublisher.AndroidpublisherScope)

	return Client{conf.Client(ctx)}, err
}

// VerifySubscription verifies subscription status
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

// VerifyProduct verifies product status
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

// CancelSubscription cancels a user's subscription purchase.
func (c *Client) CancelSubscription(packageName string, subscriptionID string, token string) error {
	service, err := androidpublisher.New(c.httpClient)
	if err != nil {
		return err
	}

	ps := androidpublisher.NewPurchasesSubscriptionsService(service)
	err = ps.Cancel(packageName, subscriptionID, token).Do()

	return err
}

// RefundSubscription refunds a user's subscription purchase, but the subscription remains valid
// until its expiration time and it will continue to recur.
func (c *Client) RefundSubscription(packageName string, subscriptionID string, token string) error {
	service, err := androidpublisher.New(c.httpClient)
	if err != nil {
		return err
	}

	ps := androidpublisher.NewPurchasesSubscriptionsService(service)
	err = ps.Refund(packageName, subscriptionID, token).Do()

	return err
}

// RevokeSubscription refunds and immediately revokes a user's subscription purchase.
// Access to the subscription will be terminated immediately and it will stop recurring.
func (c *Client) RevokeSubscription(packageName string, subscriptionID string, token string) error {
	service, err := androidpublisher.New(c.httpClient)
	if err != nil {
		return err
	}

	ps := androidpublisher.NewPurchasesSubscriptionsService(service)
	err = ps.Revoke(packageName, subscriptionID, token).Do()

	return err
}
