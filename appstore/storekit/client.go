package storekit

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	// SandboxURL is the endpoint for sandbox environment.
	SandboxURL string = "https://api.storekit-sandbox.itunes.apple.com/inApps/v1"
	// ProductionURL is the endpoint for production environment.
	ProductionURL string = "https://api.storekit.itunes.apple.com/inApps/v1"

	// tokenExpire To get better performance from the App Store Server API, reuse the same signed token for up to 60 minutes.
	tokenExpire = 3600
)

type Client struct {
	BundleID   string            // your app bundleID
	IssuerID   string            // To generate token first, see: https://developer.apple.com/documentation/appstoreserverapi/creating_api_keys_to_use_with_the_app_store_server_api
	PrivateKey *ecdsa.PrivateKey // same as above
	token      *jwt.Token        // jwt token for requests, see: https://developer.apple.com/documentation/appstoreserverapi/generating_tokens_for_api_requests
	Sandbox    bool              // default is production

	signedLatest int64  // latest sign time
	signedToken  string // latest sign token
}

func parsePrivateKey(bytes []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, errors.New("AuthKey must be a valid .p8 PEM file")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pk := key.(type) {
	case *ecdsa.PrivateKey:
		return pk, nil
	default:
		return nil, errors.New("AuthKey must be of type ecdsa.PrivateKey")
	}
}

// New return a client for App Store Server API
func New(issuerID, keyID string, privateKey []byte, bundleId string) (*Client, error) {
	// parse privateKey
	key, err := parsePrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	token := jwt.New(jwt.SigningMethodES256)
	token.Header["kid"] = keyID

	return &Client{
		IssuerID:   issuerID,
		BundleID:   bundleId,
		PrivateKey: key,
		token:      token,
	}, nil
}

func (client *Client) setToken(req *http.Request) {
	now := time.Now().Unix()
	if now-client.signedLatest > tokenExpire {
		client.token.Claims = jwt.MapClaims{
			"iss":   client.IssuerID,
			"iat":   now,
			"exp":   now + tokenExpire,
			"aud":   "appstoreconnect-v1",
			"nonce": uuid.New().String(),
			"bid":   client.BundleID,
		}
		client.signedLatest = now
		client.signedToken, _ = client.token.SignedString(client.PrivateKey)
	}
	req.Header.Set("Authorization", "Bearer "+client.signedToken)
}

func (client *Client) Do(req *http.Request, resp interface{}) error {
	client.setToken(req)

	raw, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer raw.Body.Close()

	body, err := ioutil.ReadAll(raw.Body)
	if err != nil {
		return err
	}
	// TODO parse apple error response
	return json.Unmarshal(body, resp)
}
