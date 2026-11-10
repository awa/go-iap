package appstore

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

// rootPEM is generated through `openssl x509 -inform der -in AppleRootCA-G3.cer -out apple_root.pem`
const rootPEM = `
-----BEGIN CERTIFICATE-----
MIICQzCCAcmgAwIBAgIILcX8iNLFS5UwCgYIKoZIzj0EAwMwZzEbMBkGA1UEAwwS
QXBwbGUgUm9vdCBDQSAtIEczMSYwJAYDVQQLDB1BcHBsZSBDZXJ0aWZpY2F0aW9u
IEF1dGhvcml0eTETMBEGA1UECgwKQXBwbGUgSW5jLjELMAkGA1UEBhMCVVMwHhcN
MTQwNDMwMTgxOTA2WhcNMzkwNDMwMTgxOTA2WjBnMRswGQYDVQQDDBJBcHBsZSBS
b290IENBIC0gRzMxJjAkBgNVBAsMHUFwcGxlIENlcnRpZmljYXRpb24gQXV0aG9y
aXR5MRMwEQYDVQQKDApBcHBsZSBJbmMuMQswCQYDVQQGEwJVUzB2MBAGByqGSM49
AgEGBSuBBAAiA2IABJjpLz1AcqTtkyJygRMc3RCV8cWjTnHcFBbZDuWmBSp3ZHtf
TjjTuxxEtX/1H7YyYl3J6YRbTzBPEVoA/VhYDKX1DyxNB0cTddqXl5dvMVztK517
IDvYuVTZXpmkOlEKMaNCMEAwHQYDVR0OBBYEFLuw3qFYM4iapIqZ3r6966/ayySr
MA8GA1UdEwEB/wQFMAMBAf8wDgYDVR0PAQH/BAQDAgEGMAoGCCqGSM49BAMDA2gA
MGUCMQCD6cHEFl4aXTQY2e3v9GwOAEZLuN+yRhHFD/3meoyhpmvOwgPUnPWTxnS4
at+qIxUCMG1mihDK1A3UT82NQz60imOlM27jbdoXt2QfyFMm+YhidDkLF1vLUagM
6BgD56KyKA==
-----END CERTIFICATE-----
`

type Cert struct{}

// ExtractCertByIndex extracts the certificate from the token string by index.
func (c *Cert) extractCertByIndex(tokenStr string, index int) ([]byte, error) {
	if index > 2 {
		return nil, errors.New("invalid index")
	}

	tokenArr := strings.Split(tokenStr, ".")
	headerByte, err := base64.RawStdEncoding.DecodeString(tokenArr[0])
	if err != nil {
		return nil, err
	}

	type Header struct {
		Alg string   `json:"alg"`
		X5c []string `json:"x5c"`
	}
	var header Header
	err = json.Unmarshal(headerByte, &header)
	if err != nil {
		return nil, err
	}
	// App Store JWS always carries exactly three certificates: leaf, intermediate and root.
	if len(header.X5c) != 3 {
		return nil, errors.New("failed to extract cert from x5c header, possible unauthorised request detected")
	}
	certByte, err := base64.StdEncoding.DecodeString(header.X5c[index])
	if err != nil {
		return nil, err
	}

	return certByte, nil
}

// Apple App Store signing certificates carry these custom marker OIDs. Apple's
// own App Store Server Library pins them during verification. Requiring them
// ensures the signing certificate is an App Store certificate, not merely some
// other certificate that happens to chain to Apple Root CA G3.
const (
	// appStoreLeafMarkerOID is present on the App Store receipt/JWS signing leaf.
	appStoreLeafMarkerOID = "1.2.840.113635.100.6.11.1"
	// appStoreIntermediateMarkerOID is present on the Apple WWDR intermediate
	// that issues App Store signing leaves.
	appStoreIntermediateMarkerOID = "1.2.840.113635.100.6.2.1"
)

// certHasExtension reports whether cert carries an extension with the given OID.
func certHasExtension(cert *x509.Certificate, oid string) bool {
	for _, ext := range cert.Extensions {
		if ext.Id.String() == oid {
			return true
		}
	}
	return false
}

// errNotAppStoreCert is returned when a certificate chain verifies against the
// trusted root but is not an App Store signing chain.
var errNotAppStoreCert = errors.New("certificate chain is not an App Store signing chain")

// verifyCert verifies the certificate chain against the embedded Apple Root CA G3 at the current time.
func (c *Cert) verifyCert(rootCert, intermediaCert, leafCert *x509.Certificate) error {
	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(rootPEM))
	if !ok {
		return errors.New("failed to parse root certificate")
	}

	return verifyCertChain(rootCert, intermediaCert, leafCert, roots, time.Now())
}

// verifyCertChain verifies that leafCert chains to roots through intermediaCert at the given time,
// and that the verified chain is an App Store signing chain.
func verifyCertChain(rootCert, intermediaCert, leafCert *x509.Certificate, roots *x509.CertPool, now time.Time) error {
	intermedia := x509.NewCertPool()
	intermedia.AddCert(intermediaCert)

	opts := x509.VerifyOptions{
		Roots:         roots,
		Intermediates: intermedia,
		CurrentTime:   now,
	}
	_, err := rootCert.Verify(opts)
	if err != nil {
		return err
	}

	chains, err := leafCert.Verify(opts)
	if err != nil {
		return err
	}

	// Chaining to Apple Root CA G3 is necessary but not sufficient: Go's
	// x509.Verify accepts any leaf that chains to the root and carries no
	// Extended Key Usage, which includes certificates issued for other Apple
	// programs (for example an Apple Pay Payment Processing certificate, which a
	// developer can obtain for a key they control). Such a certificate could
	// otherwise be used to sign a forged transaction or notification. Require
	// Apple's App Store marker OIDs on the verified chain (leaf, intermediate,
	// root), matching Apple's own libraries.
	for _, chain := range chains {
		if len(chain) == 3 &&
			certHasExtension(chain[0], appStoreLeafMarkerOID) &&
			certHasExtension(chain[1], appStoreIntermediateMarkerOID) {
			return nil
		}
	}

	return fmt.Errorf("%w: requires a leaf with OID %s issued by an intermediate with OID %s",
		errNotAppStoreCert, appStoreLeafMarkerOID, appStoreIntermediateMarkerOID)
}

func (c *Cert) ExtractPublicKeyFromToken(token string) (*ecdsa.PublicKey, error) {
	rootCertBytes, err := c.extractCertByIndex(token, 2)
	if err != nil {
		return nil, err
	}
	rootCert, err := x509.ParseCertificate(rootCertBytes)
	if err != nil {
		return nil, fmt.Errorf("appstore failed to parse root certificate")
	}

	intermediaCertBytes, err := c.extractCertByIndex(token, 1)
	if err != nil {
		return nil, err
	}
	intermediaCert, err := x509.ParseCertificate(intermediaCertBytes)
	if err != nil {
		return nil, fmt.Errorf("appstore failed to parse intermediate certificate")
	}

	leafCertBytes, err := c.extractCertByIndex(token, 0)
	if err != nil {
		return nil, err
	}
	leafCert, err := x509.ParseCertificate(leafCertBytes)
	if err != nil {
		return nil, fmt.Errorf("appstore failed to parse leaf certificate")
	}
	if err = c.verifyCert(rootCert, intermediaCert, leafCert); err != nil {
		return nil, err
	}

	switch pk := leafCert.PublicKey.(type) {
	case *ecdsa.PublicKey:
		return pk, nil
	default:
		return nil, errors.New("appstore public key must be of type ecdsa.PublicKey")
	}
}
