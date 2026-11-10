package api

import (
	"encoding/base64"
	"testing"
)

func TestStoreClient_ParseSignedTransaction(t *testing.T) {
	t.Parallel()

	t.Run("returns an error instead of panicking when x5c has fewer than three certificates", func(t *testing.T) {
		t.Parallel()
		// The signature is not valid; verification must fail on the x5c header
		// before the signature is checked.
		header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"ES256","x5c":["AA=="]}`))
		payload := base64.RawURLEncoding.EncodeToString([]byte(`{"transactionId":"1"}`))
		token := header + "." + payload + ".c2ln"

		client := NewStoreClient(&StoreConfig{})
		if _, err := client.ParseSignedTransaction(token); err == nil {
			t.Fatal("expected an error for an x5c with one certificate")
		}
	})
}
