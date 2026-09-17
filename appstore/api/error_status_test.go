package api

import (
	"errors"
	"fmt"
	"testing"
)

func TestHTTPStatusError(t *testing.T) {
	t.Parallel()

	const url = "https://api.storekit.apple.com/inApps/v1/subscriptions/1"
	err := newHTTPStatusError(429, url)

	t.Run("keeps the message the untyped error used", func(t *testing.T) {
		want := fmt.Sprintf("appstore api: %v return status code %v", url, 429)
		if err.Error() != want {
			t.Errorf("Error() = %q, want %q", err.Error(), want)
		}
	})

	// The status must be readable outside this package, where the type cannot
	// be named.
	t.Run("exposes the status through an interface", func(t *testing.T) {
		var withStatus interface{ StatusCode() int }
		if !errors.As(fmt.Errorf("wrapping: %w", error(err)), &withStatus) {
			t.Fatal("want the status to be reachable through an interface")
		}
		if got := withStatus.StatusCode(); got != 429 {
			t.Errorf("StatusCode() = %d, want 429", got)
		}
	})

	t.Run("stays distinct from the App Store error code type", func(t *testing.T) {
		var apiErr *Error
		if errors.As(error(err), &apiErr) {
			t.Error("want a status error not to match *Error")
		}
	})
}
