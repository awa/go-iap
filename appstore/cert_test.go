package appstore

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"math/big"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func makeTestCert(t *testing.T, oids ...string) *x509.Certificate {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "test"},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(time.Hour),
	}
	for _, oid := range oids {
		tmpl.ExtraExtensions = append(tmpl.ExtraExtensions, pkix.Extension{
			Id:    parseOID(t, oid),
			Value: []byte{0x05, 0x00}, // ASN.1 NULL, like Apple's marker extensions
		})
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	return cert
}

func parseOID(t *testing.T, s string) asn1.ObjectIdentifier {
	t.Helper()
	var oid asn1.ObjectIdentifier
	part := 0
	seen := false
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			oid = append(oid, part)
			part, seen = 0, false
			continue
		}
		part = part*10 + int(s[i]-'0')
		seen = true
	}
	if seen {
		oid = append(oid, part)
	}
	return oid
}

func TestCertHasExtension(t *testing.T) {
	// App Store signing leaf carries the marker OID.
	appStoreLeaf := makeTestCert(t, appStoreLeafMarkerOID)
	if !certHasExtension(appStoreLeaf, appStoreLeafMarkerOID) {
		t.Errorf("expected App Store leaf marker OID to be detected")
	}

	// A certificate for another Apple program (e.g. Apple Pay Payment
	// Processing, OID 1.2.840.113635.100.6.32) that chains to the same root but
	// is not an App Store signing certificate must not be mistaken for one.
	applePay := makeTestCert(t, "1.2.840.113635.100.6.32")
	if certHasExtension(applePay, appStoreLeafMarkerOID) {
		t.Errorf("Apple Pay certificate must not satisfy the App Store leaf marker OID")
	}

	// Intermediate marker OID.
	appStoreIntermediate := makeTestCert(t, appStoreIntermediateMarkerOID)
	if !certHasExtension(appStoreIntermediate, appStoreIntermediateMarkerOID) {
		t.Errorf("expected App Store intermediate marker OID to be detected")
	}
	if certHasExtension(makeTestCert(t), appStoreIntermediateMarkerOID) {
		t.Errorf("certificate without extensions must not satisfy the intermediate marker OID")
	}
}

// testChainTime is a fixed verification time within the validity of every
// certificate used by the chain tests, including the genuine Apple fixtures.
var testChainTime = time.Date(2022, 6, 1, 0, 0, 0, 0, time.UTC)

// testIssuedCert is a certificate together with the key that can issue
// certificates under it.
type testIssuedCert struct {
	cert *x509.Certificate
	key  *ecdsa.PrivateKey
}

// issueTestCert issues an ECDSA P-256 certificate carrying the given marker
// OIDs. It is signed by parent, or self-signed when parent is nil.
func issueTestCert(t *testing.T, name string, parent *testIssuedCert, isCA bool, oids ...string) *testIssuedCert {
	t.Helper()
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	serial, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil {
		t.Fatal(err)
	}
	tmpl := &x509.Certificate{
		SerialNumber:          serial,
		Subject:               pkix.Name{CommonName: name},
		NotBefore:             testChainTime.Add(-time.Hour),
		NotAfter:              testChainTime.Add(time.Hour),
		BasicConstraintsValid: true,
		IsCA:                  isCA,
		KeyUsage:              x509.KeyUsageDigitalSignature,
	}
	if isCA {
		tmpl.KeyUsage = x509.KeyUsageCertSign
	}
	for _, oid := range oids {
		tmpl.ExtraExtensions = append(tmpl.ExtraExtensions, pkix.Extension{
			Id:    parseOID(t, oid),
			Value: []byte{0x05, 0x00}, // ASN.1 NULL, like Apple's marker extensions
		})
	}

	signer, signerKey := tmpl, key
	if parent != nil {
		signer, signerKey = parent.cert, parent.key
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, signer, &key.PublicKey, signerKey)
	if err != nil {
		t.Fatal(err)
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		t.Fatal(err)
	}
	return &testIssuedCert{cert: cert, key: key}
}

func certPool(certs ...*x509.Certificate) *x509.CertPool {
	pool := x509.NewCertPool()
	for _, cert := range certs {
		pool.AddCert(cert)
	}
	return pool
}

// loadTestCert reads a PEM certificate from testdata.
func loadTestCert(t *testing.T, name string) *x509.Certificate {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", name))
	if err != nil {
		t.Fatal(err)
	}
	block, _ := pem.Decode(b)
	if block == nil {
		t.Fatalf("no PEM block in %s", name)
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	return cert
}

// loadAppleRoot parses the embedded Apple Root CA G3.
func loadAppleRoot(t *testing.T) *x509.Certificate {
	t.Helper()
	block, _ := pem.Decode([]byte(rootPEM))
	if block == nil {
		t.Fatal("no PEM block in rootPEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	return cert
}

func TestVerifyCertChain(t *testing.T) {
	t.Parallel()

	t.Run("accepts a chain whose leaf and intermediate carry the App Store OIDs", func(t *testing.T) {
		t.Parallel()
		root := issueTestCert(t, "test root", nil, true)
		inter := issueTestCert(t, "test intermediate", root, true, appStoreIntermediateMarkerOID)
		leaf := issueTestCert(t, "test leaf", inter, false, appStoreLeafMarkerOID)

		err := verifyCertChain(root.cert, inter.cert, leaf.cert, certPool(root.cert), testChainTime)
		if err != nil {
			t.Fatalf("expected the chain to be accepted, got %v", err)
		}
	})

	t.Run("rejects a leaf without the App Store OID", func(t *testing.T) {
		t.Parallel()
		root := issueTestCert(t, "test root", nil, true)
		inter := issueTestCert(t, "test intermediate", root, true, appStoreIntermediateMarkerOID)
		leaf := issueTestCert(t, "test leaf", inter, false, "1.2.840.113635.100.6.32")

		err := verifyCertChain(root.cert, inter.cert, leaf.cert, certPool(root.cert), testChainTime)
		if !errors.Is(err, errNotAppStoreCert) {
			t.Fatalf("expected errNotAppStoreCert, got %v", err)
		}
	})

	t.Run("rejects an intermediate without the App Store OID", func(t *testing.T) {
		t.Parallel()
		root := issueTestCert(t, "test root", nil, true)
		inter := issueTestCert(t, "test intermediate", root, true, "1.2.840.113635.100.6.2.15")
		leaf := issueTestCert(t, "test leaf", inter, false, appStoreLeafMarkerOID)

		err := verifyCertChain(root.cert, inter.cert, leaf.cert, certPool(root.cert), testChainTime)
		if !errors.Is(err, errNotAppStoreCert) {
			t.Fatalf("expected errNotAppStoreCert, got %v", err)
		}
	})

	t.Run("rejects a leaf issued directly by the root", func(t *testing.T) {
		t.Parallel()
		// x5c[1] carries the intermediate OID, but the verified chain is
		// leaf -> root and never goes through it.
		root := issueTestCert(t, "test root", nil, true)
		inter := issueTestCert(t, "test intermediate", root, true, appStoreIntermediateMarkerOID)
		leaf := issueTestCert(t, "test leaf", root, false, appStoreLeafMarkerOID)

		err := verifyCertChain(root.cert, inter.cert, leaf.cert, certPool(root.cert), testChainTime)
		if !errors.Is(err, errNotAppStoreCert) {
			t.Fatalf("expected errNotAppStoreCert, got %v", err)
		}
	})

	t.Run("accepts the genuine App Store chain within its validity", func(t *testing.T) {
		t.Parallel()
		root := loadAppleRoot(t)
		inter := loadTestCert(t, "apple_wwdr_ca_g6.pem")
		leaf := loadTestCert(t, "app_store_leaf.pem")

		err := verifyCertChain(root, inter, leaf, certPool(root), testChainTime)
		if err != nil {
			t.Fatalf("expected the genuine chain to be accepted, got %v", err)
		}
	})

	t.Run("rejects the Apple WWDR CA - G2 chain that issues Apple Pay certificates", func(t *testing.T) {
		t.Parallel()
		// Apple WWDR CA - G2 chains to Apple Root CA G3 and has no EKU, so it
		// passes x509.Verify, but it is not an App Store signing chain.
		root := loadAppleRoot(t)
		wwdrG2 := loadTestCert(t, "apple_wwdr_ca_g2.pem")

		err := verifyCertChain(root, wwdrG2, wwdrG2, certPool(root), testChainTime)
		if !errors.Is(err, errNotAppStoreCert) {
			t.Fatalf("expected errNotAppStoreCert, got %v", err)
		}
	})
}

// makeTestJWS builds a JWS whose header carries the given x5c entries. The
// signature is not valid; it is only for tests that fail before verifying it.
func makeTestJWS(t *testing.T, x5c ...string) string {
	t.Helper()
	header, err := json.Marshal(map[string]any{"alg": "ES256", "x5c": x5c})
	if err != nil {
		t.Fatal(err)
	}
	payload := base64.RawURLEncoding.EncodeToString([]byte(`{}`))
	return base64.RawURLEncoding.EncodeToString(header) + "." + payload + ".c2ln"
}

func TestCert_extractCertByIndex(t *testing.T) {
	t.Parallel()

	t.Run("returns the requested certificate when x5c has three entries", func(t *testing.T) {
		t.Parallel()
		token := makeTestJWS(t, "AQ==", "Ag==", "Aw==")

		got, err := (&Cert{}).extractCertByIndex(token, 1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !bytes.Equal(got, []byte{0x02}) {
			t.Fatalf("expected the second certificate, got %v", got)
		}
	})

	t.Run("rejects an x5c with two entries", func(t *testing.T) {
		t.Parallel()
		token := makeTestJWS(t, "AQ==", "Ag==")

		if _, err := (&Cert{}).extractCertByIndex(token, 0); err == nil {
			t.Fatal("expected an error for an x5c with two entries")
		}
	})

	t.Run("rejects an x5c with four entries", func(t *testing.T) {
		t.Parallel()
		token := makeTestJWS(t, "AQ==", "Ag==", "Aw==", "BA==")

		if _, err := (&Cert{}).extractCertByIndex(token, 0); err == nil {
			t.Fatal("expected an error for an x5c with four entries")
		}
	})
}
