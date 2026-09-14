package api

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"testing"
)

// The constructor returns the concrete type, so nothing else checks this.
var _ StoreAPIClient = (*StoreClientWithSandboxFallback)(nil)

// stubClient answers one lookup and panics on anything else, so an unexpected
// call fails loudly.
type stubClient struct {
	StoreAPIClient // nil on purpose: any method not stubbed below panics

	name  string
	err   error
	calls *[]string
}

func (s *stubClient) record() {
	*s.calls = append(*s.calls, s.name)
}

func (s *stubClient) GetALLSubscriptionStatuses(context.Context, string, *url.Values) (*StatusResponse, error) {
	s.record()
	if s.err != nil {
		return nil, s.err
	}
	return &StatusResponse{}, nil
}

func (s *stubClient) GetTransactionInfo(context.Context, string) (*TransactionInfoResponse, error) {
	s.record()
	if s.err != nil {
		return nil, s.err
	}
	return &TransactionInfoResponse{}, nil
}

func (s *stubClient) GetTransactionHistory(context.Context, string, *url.Values) ([]*HistoryResponse, error) {
	s.record()
	if s.err != nil {
		return nil, s.err
	}
	return []*HistoryResponse{{}}, nil
}

func (s *stubClient) GetRefundHistory(context.Context, string) ([]*RefundLookupResponse, error) {
	s.record()
	if s.err != nil {
		return nil, s.err
	}
	return []*RefundLookupResponse{{}}, nil
}

func (s *stubClient) GetAppTransactionInfo(context.Context, string) (*AppTransactionInfoResponse, error) {
	s.record()
	if s.err != nil {
		return nil, s.err
	}
	return &AppTransactionInfoResponse{}, nil
}

// lookup hides the differing signatures so every method runs the same table.
type lookup struct {
	name string
	call func(context.Context, *StoreClientWithSandboxFallback) error
}

func lookups() []lookup {
	return []lookup{
		{"GetALLSubscriptionStatuses", func(ctx context.Context, c *StoreClientWithSandboxFallback) error {
			_, err := c.GetALLSubscriptionStatuses(ctx, "1", nil)
			return err
		}},
		{"GetTransactionInfo", func(ctx context.Context, c *StoreClientWithSandboxFallback) error {
			_, err := c.GetTransactionInfo(ctx, "1")
			return err
		}},
		{"GetTransactionHistory", func(ctx context.Context, c *StoreClientWithSandboxFallback) error {
			_, err := c.GetTransactionHistory(ctx, "1", nil)
			return err
		}},
		{"GetRefundHistory", func(ctx context.Context, c *StoreClientWithSandboxFallback) error {
			_, err := c.GetRefundHistory(ctx, "1")
			return err
		}},
		{"GetAppTransactionInfo", func(ctx context.Context, c *StoreClientWithSandboxFallback) error {
			_, err := c.GetAppTransactionInfo(ctx, "1")
			return err
		}},
	}
}

func TestSandboxFallbackRetriesOnlyWhenTransactionIsMissing(t *testing.T) {
	t.Parallel()

	// A caller mistake, which must not be retried against Sandbox.
	misconfigured := fmt.Errorf("wrong bundle id: %w", AppNotFoundError)

	tests := []struct {
		name          string
		productionErr error
		sandboxErr    error
		wantCalls     []string
		wantMissing   bool // errors.Is(err, TransactionIdNotFoundError)
		wantErr       bool
	}{
		{
			name:      "production answers, sandbox is not asked",
			wantCalls: []string{"production"},
		},
		{
			name:          "production reports missing, sandbox answers",
			productionErr: TransactionIdNotFoundError,
			wantCalls:     []string{"production", "sandbox"},
		},
		{
			name:          "missing in both environments stays a missing transaction",
			productionErr: TransactionIdNotFoundError,
			sandboxErr:    TransactionIdNotFoundError,
			wantCalls:     []string{"production", "sandbox"},
			wantMissing:   true,
			wantErr:       true,
		},
		{
			name:          "sandbox failing to answer is not a missing transaction",
			productionErr: TransactionIdNotFoundError,
			sandboxErr:    misconfigured,
			wantCalls:     []string{"production", "sandbox"},
			wantErr:       true,
		},
		{
			name:          "another production error is returned as is",
			productionErr: misconfigured,
			wantCalls:     []string{"production"},
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		for _, l := range lookups() {
			t.Run(tt.name+"/"+l.name, func(t *testing.T) {
				t.Parallel()

				var calls []string
				c := &StoreClientWithSandboxFallback{
					productionCli: &stubClient{name: "production", err: tt.productionErr, calls: &calls},
					sandboxCli:    &stubClient{name: "sandbox", err: tt.sandboxErr, calls: &calls},
				}

				err := l.call(context.Background(), c)

				if (err != nil) != tt.wantErr {
					t.Fatalf("error = %v, want error = %v", err, tt.wantErr)
				}
				if got := errors.Is(err, TransactionIdNotFoundError); got != tt.wantMissing {
					t.Errorf("errors.Is(err, TransactionIdNotFoundError) = %v, want %v (err = %v)", got, tt.wantMissing, err)
				}
				if !reflect.DeepEqual(calls, tt.wantCalls) {
					t.Errorf("calls = %v, want %v", calls, tt.wantCalls)
				}
			})
		}
	}
}

// A method that does not retry must still reach production, never Sandbox.
func TestSandboxFallbackDelegatesOtherMethods(t *testing.T) {
	t.Parallel()

	var calls []string
	c := &StoreClientWithSandboxFallback{
		productionCli: &delegatingStub{calls: &calls},
		sandboxCli:    &stubClient{name: "sandbox", calls: &calls},
	}

	if _, err := c.LookupOrderID(context.Background(), "order"); err != nil {
		t.Fatalf("LookupOrderID: %v", err)
	}
	if want := []string{"production"}; !reflect.DeepEqual(calls, want) {
		t.Errorf("calls = %v, want %v", calls, want)
	}
}

type delegatingStub struct {
	StoreAPIClient
	calls *[]string
}

func (s *delegatingStub) LookupOrderID(context.Context, string) (*OrderLookupResponse, error) {
	*s.calls = append(*s.calls, "production")
	return &OrderLookupResponse{}, nil
}
