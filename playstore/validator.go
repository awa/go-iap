package playstore

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"google.golang.org/api/option"
	"net/http"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	androidpublisher "google.golang.org/api/androidpublisher/v3"
)

//go:generate mockgen  -destination=mocks/mock_playstore.go -package=mocks github.com/awa/go-iap/playstore IABProduct,IABSubscription

// The IABProduct type is an interface for product service
type IABProduct interface {
	VerifyProduct(context.Context, string, string, string) (*androidpublisher.ProductPurchase, error)
	AcknowledgeProduct(context.Context, string, string, string, string) error
}

// The IABSubscription type is an interface  for subscription service
type IABSubscription interface {
	AcknowledgeSubscription(context.Context, string, string, string, *androidpublisher.SubscriptionPurchasesAcknowledgeRequest) error
	VerifySubscription(context.Context, string, string, string) (*androidpublisher.SubscriptionPurchase, error)
	CancelSubscription(context.Context, string, string, string) error
	RefundSubscription(context.Context, string, string, string) error
	RevokeSubscription(context.Context, string, string, string) error
}

// The Client type implements VerifySubscription method
type Client struct {
	service *androidpublisher.Service
}

// New returns http client which includes the credentials to access androidpublisher API.
// You should create a service account for your project at
// https://console.developers.google.com and download a JSON key file to set this argument.
func New(jsonKey []byte) (*Client, error) {
	c := &http.Client{Timeout: 10 * time.Second}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, c)


	conf, err := google.JWTConfigFromJSON(jsonKey, androidpublisher.AndroidpublisherScope)
	if err != nil {
		return nil, err
	}

	val := conf.Client(ctx).Transport.(*oauth2.Transport)
	_, err = val.Source.Token()
	if err != nil {
		return nil, err
	}

	service, err := androidpublisher.NewService(ctx, option.WithHTTPClient(conf.Client(ctx)))
	if err != nil {
		return nil, err
	}

	return &Client{service}, err
}

// NewWithClient returns http client which includes the custom http client.
func NewWithClient(jsonKey []byte, cli *http.Client) (*Client, error) {
	if cli == nil {
		return nil, fmt.Errorf("client is nil")
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, cli)

	conf, err := google.JWTConfigFromJSON(jsonKey, androidpublisher.AndroidpublisherScope)
	if err != nil {
		return nil, err
	}

	service, err := androidpublisher.NewService(ctx, option.WithHTTPClient(conf.Client(ctx)))
	if err != nil {
		return nil, err
	}

	return &Client{service}, err
}

// AcknowledgeSubscription acknowledges a subscription purchase.
func (c *Client) AcknowledgeSubscription(
	ctx context.Context,
	packageName string,
	subscriptionID string,
	token string,
	req *androidpublisher.SubscriptionPurchasesAcknowledgeRequest,
) error {
	ps := androidpublisher.NewPurchasesSubscriptionsService(c.service)
	err := ps.Acknowledge(packageName, subscriptionID, token, req).Context(ctx).Do()

	return err
}

// VerifySubscription verifies subscription status
func (c *Client) VerifySubscription(
	ctx context.Context,
	packageName string,
	subscriptionID string,
	token string,
) (*androidpublisher.SubscriptionPurchase, error) {
	ps := androidpublisher.NewPurchasesSubscriptionsService(c.service)
	result, err := ps.Get(packageName, subscriptionID, token).Context(ctx).Do()

	return result, err
}

// VerifyProduct verifies product status
func (c *Client) VerifyProduct(
	ctx context.Context,
	packageName string,
	productID string,
	token string,
) (*androidpublisher.ProductPurchase, error) {
	ps := androidpublisher.NewPurchasesProductsService(c.service)
	result, err := ps.Get(packageName, productID, token).Context(ctx).Do()

	return result, err
}

func (c *Client) AcknowledgeProduct(ctx context.Context, packageName, productID, token, developerPayload string) error {
	ps := androidpublisher.NewPurchasesProductsService(c.service)
	acknowledgeRequest := &androidpublisher.ProductPurchasesAcknowledgeRequest{DeveloperPayload: developerPayload}
	err := ps.Acknowledge(packageName, productID, token, acknowledgeRequest).Context(ctx).Do()

	return err
}

// CancelSubscription cancels a user's subscription purchase.
func (c *Client) CancelSubscription(ctx context.Context, packageName string, subscriptionID string, token string) error {
	ps := androidpublisher.NewPurchasesSubscriptionsService(c.service)
	err := ps.Cancel(packageName, subscriptionID, token).Context(ctx).Do()

	return err
}

// RefundSubscription refunds a user's subscription purchase, but the subscription remains valid
// until its expiration time and it will continue to recur.
func (c *Client) RefundSubscription(ctx context.Context, packageName string, subscriptionID string, token string) error {
	ps := androidpublisher.NewPurchasesSubscriptionsService(c.service)
	err := ps.Refund(packageName, subscriptionID, token).Context(ctx).Do()

	return err
}

// RevokeSubscription refunds and immediately revokes a user's subscription purchase.
// Access to the subscription will be terminated immediately and it will stop recurring.
func (c *Client) RevokeSubscription(ctx context.Context, packageName string, subscriptionID string, token string) error {
	ps := androidpublisher.NewPurchasesSubscriptionsService(c.service)
	err := ps.Revoke(packageName, subscriptionID, token).Context(ctx).Do()

	return err
}

// VerifySignature verifies in app billing signature.
// You need to prepare a public key for your Android app's in app billing
// at https://play.google.com/apps/publish/
func VerifySignature(base64EncodedPublicKey string, receipt []byte, signature string) (isValid bool, err error) {
	// prepare public key
	decodedPublicKey, err := base64.StdEncoding.DecodeString(base64EncodedPublicKey)
	if err != nil {
		return false, fmt.Errorf("failed to decode public key")
	}
	publicKeyInterface, err := x509.ParsePKIXPublicKey(decodedPublicKey)
	if err != nil {
		return false, fmt.Errorf("failed to parse public key")
	}
	publicKey, _ := publicKeyInterface.(*rsa.PublicKey)

	// generate hash value from receipt
	hasher := sha1.New()
	hasher.Write(receipt)
	hashedReceipt := hasher.Sum(nil)

	// decode signature
	decodedSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return false, fmt.Errorf("failed to decode signature")
	}

	// verify
	if err := rsa.VerifyPKCS1v15(publicKey, crypto.SHA1, hashedReceipt, decodedSignature); err != nil {
		return false, nil
	}

	return true, nil
}
