package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

const testSeatClaimID = "6c42a3f4-fcbd-4aa8-a593-e5410db47616"

var (
	_ SeatManager = (*StoreClient)(nil)
	_ SeatManager = (*StoreClientWithSandboxFallback)(nil)
)

type seatOperation struct {
	name     string
	method   string
	path     string
	body     string
	response string
	want     any
	call     func(context.Context, SeatManager) (any, error)
}

func seatOperations() []seatOperation {
	return []seatOperation{
		{name: "get", method: "GET", path: "/groups/v1/transactions/original/seats", body: "", response: `{"availableSeatCount":2,"totalSeatCount":5,"assignedSeatClaimIds":["6c42a3f4-fcbd-4aa8-a593-e5410db47616"]}`,
			want: &GetSeatInfoResponse{AvailableSeatCount: 2, TotalSeatCount: 5, AssignedSeatClaimIds: []string{testSeatClaimID}},
			call: func(ctx context.Context, c SeatManager) (any, error) { return c.GetSeatInfo(ctx, "original") }},
		{name: "claim without customer", method: "POST", path: "/groups/v1/transactions/original/seats", body: `{"requestIdentifier":"claim-request"}`, response: `{"seatClaimId":"6c42a3f4-fcbd-4aa8-a593-e5410db47616"}`,
			want: &ClaimSeatResponse{SeatClaimId: testSeatClaimID},
			call: func(ctx context.Context, c SeatManager) (any, error) {
				return c.ClaimSeat(ctx, "original", ClaimSeatRequest{RequestIdentifier: "claim-request"})
			}},
		{name: "claim with customer", method: "POST", path: "/groups/v1/transactions/original/seats", body: `{"requestIdentifier":"claim-request","appTransactionId":"customer"}`, response: `{"seatClaimId":"6c42a3f4-fcbd-4aa8-a593-e5410db47616"}`,
			want: &ClaimSeatResponse{SeatClaimId: testSeatClaimID},
			call: func(ctx context.Context, c SeatManager) (any, error) {
				return c.ClaimSeat(ctx, "original", ClaimSeatRequest{RequestIdentifier: "claim-request", AppTransactionId: "customer"})
			}},
		{name: "unclaim", method: "DELETE", path: "/groups/v1/transactions/original/seats/" + testSeatClaimID, response: "", want: http.StatusOK,
			call: func(ctx context.Context, c SeatManager) (any, error) {
				return c.UnclaimSeat(ctx, "original", testSeatClaimID)
			}},
		{name: "associate", method: "POST", path: "/groups/v1/transactions/original/seats/" + testSeatClaimID + "/associate", body: `{"appTransactionId":"customer"}`, response: "", want: http.StatusOK,
			call: func(ctx context.Context, c SeatManager) (any, error) {
				return c.AssociateSeat(ctx, "original", testSeatClaimID, AssociateSeatRequest{AppTransactionId: "customer"})
			}},
	}
}

func TestSeatEndpoints(t *testing.T) {
	t.Parallel()
	for _, op := range seatOperations() {
		t.Run(op.name, func(t *testing.T) {
			t.Parallel()
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != op.method || r.URL.Path != op.path {
					t.Errorf("request = %s %s, want %s %s", r.Method, r.URL.Path, op.method, op.path)
				}
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Error(err)
				}
				if strings.TrimSpace(string(body)) != op.body {
					t.Errorf("body = %s, want %s", body, op.body)
				}
				if r.Header.Get("Content-Type") != "application/json" || !strings.HasPrefix(r.Header.Get("Authorization"), "Bearer ") {
					t.Error("missing JSON or bearer headers")
				}
				_, _ = io.WriteString(w, op.response)
			}))
			t.Cleanup(srv.Close)
			c := NewStoreClient(&StoreConfig{KeyContent: testPrivateKeyPEM(t), HostDebug: srv.URL, Sandbox: true})
			got, err := op.call(t.Context(), c)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, op.want) {
				t.Errorf("response = %+v, want %+v", got, op.want)
			}
		})
	}
}

func TestSeatEndpointErrors(t *testing.T) {
	t.Parallel()
	for _, op := range seatOperations() {
		for _, tt := range []struct {
			status, code int
			want         error
		}{
			{400, 4000241, BundleTransactionIdNotSupportedError},
			{400, 4000251, NoAvailableSeatsError},
			{400, 4000273, TransactionNotDevManagedGroupError},
			{403, 4030039, SubscriptionNotActiveError},
			{404, 4040023, AppTransactionIdNotFoundError},
			{404, 4040025, SeatClaimIdNotFoundError},
			{429, 4290000, RateLimitExceededError},
			{500, 5000001, GeneralInternalRetryableError},
		} {
			t.Run(op.name+"/"+tt.want.Error(), func(t *testing.T) {
				t.Parallel()
				srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Retry-After", "60")
					w.WriteHeader(tt.status)
					_ = json.NewEncoder(w).Encode(map[string]any{"errorCode": tt.code, "errorMessage": "provider error"})
				}))
				t.Cleanup(srv.Close)
				c := NewStoreClient(&StoreConfig{KeyContent: testPrivateKeyPEM(t), HostDebug: srv.URL})
				_, err := op.call(t.Context(), c)
				if !errors.Is(err, tt.want) {
					t.Fatalf("error = %v, want %v", err, tt.want)
				}
				if tt.status == 429 {
					apiErr, ok := errors.AsType[*Error](err)
					if !ok || apiErr.RetryAfter() != 60 {
						t.Errorf("rate limit error = %v", err)
					}
				}
			})
		}
		t.Run(op.name+"/uncoded status", func(t *testing.T) {
			var req *http.Request
			c := newGroupTestClient(t, http.StatusUnauthorized, "", &req)
			_, err := op.call(t.Context(), c)
			statusErr, ok := errors.AsType[*httpStatusError](err)
			if !ok || statusErr.StatusCode() != 401 {
				t.Errorf("error = %v, want HTTP 401", err)
			}
		})
		t.Run(op.name+"/canceled", func(t *testing.T) {
			ctx, cancel := context.WithCancel(t.Context())
			cancel()
			c := NewStoreClient(&StoreConfig{KeyContent: testPrivateKeyPEM(t), HostDebug: "http://127.0.0.1:1"})
			_, err := op.call(ctx, c)
			if !errors.Is(err, context.Canceled) {
				t.Errorf("error = %v, want context.Canceled", err)
			}
		})
	}
}

func TestSeatResponsesRejectMalformedJSON(t *testing.T) {
	t.Parallel()
	for _, op := range seatOperations()[:3] {
		for _, body := range []string{`{"invalid":`, `null`, `[]`, `{"seatClaimId":123,"totalSeatCount":"five"}`} {
			t.Run(op.name+"/"+body, func(t *testing.T) {
				var req *http.Request
				c := newGroupTestClient(t, http.StatusOK, body, &req)
				_, err := op.call(t.Context(), c)
				if err == nil {
					t.Fatal("expected response decoding error")
				}
			})
		}
	}
}

func TestSeatEndpointsRejectMissingInputs(t *testing.T) {
	t.Parallel()
	c := NewStoreClient(&StoreConfig{}) // Invalid inputs must fail before token generation.
	for _, tt := range []struct {
		name, field string
		call        func() error
	}{
		{"get transaction", "originalTransactionId", func() error { _, err := c.GetSeatInfo(t.Context(), ""); return err }},
		{"claim transaction", "originalTransactionId", func() error {
			_, err := c.ClaimSeat(t.Context(), "", ClaimSeatRequest{RequestIdentifier: "request"})
			return err
		}},
		{"claim request", "requestIdentifier", func() error { _, err := c.ClaimSeat(t.Context(), "original", ClaimSeatRequest{}); return err }},
		{"long request", "requestIdentifier", func() error {
			_, err := c.ClaimSeat(t.Context(), "original", ClaimSeatRequest{RequestIdentifier: strings.Repeat("a", 129)})
			return err
		}},
		{"unclaim transaction", "originalTransactionId", func() error { _, err := c.UnclaimSeat(t.Context(), "", testSeatClaimID); return err }},
		{"unclaim seat", "seatClaimId", func() error { _, err := c.UnclaimSeat(t.Context(), "original", ""); return err }},
		{"associate transaction", "originalTransactionId", func() error {
			_, err := c.AssociateSeat(t.Context(), "", testSeatClaimID, AssociateSeatRequest{AppTransactionId: "customer"})
			return err
		}},
		{"associate seat", "seatClaimId", func() error {
			_, err := c.AssociateSeat(t.Context(), "original", "", AssociateSeatRequest{AppTransactionId: "customer"})
			return err
		}},
		{"associate customer", "appTransactionId", func() error {
			_, err := c.AssociateSeat(t.Context(), "original", testSeatClaimID, AssociateSeatRequest{})
			return err
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.call(); err == nil || !strings.Contains(err.Error(), tt.field) {
				t.Errorf("error = %v, want validation for %s", err, tt.field)
			}
		})
	}
}

func TestSeatEndpointEscapesPathParameters(t *testing.T) {
	t.Parallel()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if want := "/groups/v1/transactions/original%2Fid%3Fquery/seats/seat%2Fid%3Fquery/associate"; r.URL.EscapedPath() != want || r.URL.RawQuery != "" {
			t.Errorf("URL = %s, want path %s", r.URL, want)
		}
	}))
	t.Cleanup(srv.Close)
	c := NewStoreClient(&StoreConfig{KeyContent: testPrivateKeyPEM(t), HostDebug: srv.URL})
	if _, err := c.AssociateSeat(t.Context(), "original/id?query", "seat/id?query", AssociateSeatRequest{AppTransactionId: "customer"}); err != nil {
		t.Fatal(err)
	}
}

func TestFallbackSeatEndpointsUseSandboxOnce(t *testing.T) {
	t.Parallel()
	for _, op := range seatOperations() {
		t.Run(op.name, func(t *testing.T) {
			for _, status := range []int{200, 401} {
				t.Run(http.StatusText(status), func(t *testing.T) {
					calls := 0
					srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						calls++
						w.WriteHeader(status)
						if status == 200 {
							_, _ = io.WriteString(w, op.response)
						}
					}))
					t.Cleanup(srv.Close)
					c := &StoreClientWithSandboxFallback{
						productionCli: &stubClient{}, // Any production call panics.
						sandboxCli:    NewStoreClient(&StoreConfig{KeyContent: testPrivateKeyPEM(t), HostDebug: srv.URL, Sandbox: true}),
					}
					got, err := op.call(t.Context(), c)
					if status == 200 && (err != nil || !reflect.DeepEqual(got, op.want)) {
						t.Errorf("response = %+v, error = %v", got, err)
					}
					if status != 200 && err == nil {
						t.Error("expected Sandbox error")
					}
					if calls != 1 {
						t.Errorf("Sandbox calls = %d, want 1", calls)
					}
				})
			}
		})
	}
}
