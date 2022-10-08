package api

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// Authorize Tokens For App Store Server API Request
// Doc: https://developer.apple.com/documentation/appstoreserverapi/generating_tokens_for_api_requests
var (
	ErrAuthKeyInvalidPem  = errors.New("token: AuthKey must be a valid .p8 PEM file")
	ErrAuthKeyInvalidType = errors.New("token: AuthKey must be of type ecdsa.PrivateKey")
)

// Token represents an Apple Provider Authentication Token (JSON Web Token).
type Token struct {
	sync.Mutex

	KeyContent []byte // Loads a .p8 certificate
	KeyID      string // Your private key ID from App Store Connect (Ex: 2X9R4HXF34)
	BundleID   string // Your app’s bundle ID
	Issuer     string // Your issuer ID from the Keys page in App Store Connect (Ex: "57246542-96fe-1a63-e053-0824d011072a")
	Sandbox    bool   // default is Production

	// internal variables
	AuthKey   *ecdsa.PrivateKey // .p8 private key
	ExpiredAt int64             // The token’s expiration time, in UNIX time. Tokens that expire more than 60 minutes after the time in iat are not valid (Ex: 1623086400)
	Bearer    string            // Authorized bearer token
}

// GenerateIfExpired checks to see if the token is about to expire and generates a new token.
func (t *Token) GenerateIfExpired() (string, error) {
	t.Lock()
	defer t.Unlock()

	if t.Expired() {
		err := t.Generate()
		if err != nil {
			return "", err
		}
	}

	return t.Bearer, nil
}

// Expired checks to see if the token has expired.
func (t *Token) Expired() bool {
	return time.Now().Unix() >= t.ExpiredAt
}

// Generate creates a new token.
func (t *Token) Generate() error {
	key, err := t.passKeyFromByte(t.KeyContent)
	if err != nil {
		return err
	}
	t.AuthKey = key

	issuedAt := time.Now().Unix()
	expiredAt := time.Now().Add(time.Duration(1) * time.Hour).Unix()
	jwtToken := &jwt.Token{
		Header: map[string]interface{}{
			"alg": "ES256",
			"kid": t.KeyID,
			"typ": "JWT",
		},

		Claims: jwt.MapClaims{
			"iss":   t.Issuer,
			"iat":   issuedAt,
			"exp":   expiredAt,
			"aud":   "appstoreconnect-v1",
			"nonce": uuid.New(),
			"bid":   t.BundleID,
		},
		Method: jwt.SigningMethodES256,
	}

	bearer, err := jwtToken.SignedString(t.AuthKey)
	if err != nil {
		return err
	}
	t.ExpiredAt = expiredAt
	t.Bearer = bearer

	return nil
}

// loadKeyFromFile loads a .p8 certificate from a local file and returns a *ecdsa.PrivateKey.
func (t *Token) loadKeyFromFile(filename string) (*ecdsa.PrivateKey, error) {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return t.passKeyFromByte(bytes)
}

// passKeyFromByte loads a .p8 certificate from an in memory byte array and returns an *ecdsa.PrivateKey.
func (t *Token) passKeyFromByte(bytes []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(bytes)
	if block == nil {
		return nil, ErrAuthKeyInvalidPem
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pk := key.(type) {
	case *ecdsa.PrivateKey:
		return pk, nil
	default:
		return nil, ErrAuthKeyInvalidType
	}
}
