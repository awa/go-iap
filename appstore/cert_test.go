package appstore

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"math/big"
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
