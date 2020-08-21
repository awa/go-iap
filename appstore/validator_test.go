package appstore

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestHandleError(t *testing.T) {
	tests := []struct {
		name string
		in   int
		out  error
	}{
		{
			name: "status 0",
			in:   0,
			out:  nil,
		},
		{
			name: "status 21000",
			in:   21000,
			out:  ErrInvalidJSON,
		},
		{
			name: "status 21002",
			in:   21002,
			out:  ErrInvalidReceiptData,
		},
		{
			name: "status 21003",
			in:   21003,
			out:  ErrReceiptUnauthenticated,
		},
		{
			name: "status 21004",
			in:   21004,
			out:  ErrInvalidSharedSecret,
		},
		{
			name: "status 21005",
			in:   21005,
			out:  ErrServerUnavailable,
		},
		{
			name: "status 21007",
			in:   21007,
			out:  ErrReceiptIsForTest,
		},
		{
			name: "status 21008",
			in:   21008,
			out:  ErrReceiptIsForProduction,
		},
		{
			name: "status 21009",
			in:   21009,
			out:  ErrInternalDataAccessError,
		},
		{
			name: "status 21010",
			in:   21010,
			out:  ErrReceiptUnauthorized,
		},
		{
			name: "status 21100 ~ 21199",
			in:   21100,
			out:  ErrInternalDataAccessError,
		},
		{
			name: "status unknown",
			in:   100,
			out:  ErrUnknown,
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			out := HandleError(v.in)

			if !errors.Is(out, v.out) {
				t.Errorf("input: %d\ngot: %v\nwant: %v\n", v.in, out, v.out)
			}
		})
	}
}

func TestNew(t *testing.T) {
	expected := &Client{
		ProductionURL: ProductionURL,
		SandboxURL:    SandboxURL,
		httpCli: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	actual := New()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestNewWithClient(t *testing.T) {
	expected := &Client{
		ProductionURL: ProductionURL,
		SandboxURL:    SandboxURL,
		httpCli: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	actual := NewWithClient(&http.Client{
		Timeout: 10 * time.Second,
	})
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestVerifyTimeout(t *testing.T) {
	client := &Client{
		ProductionURL: ProductionURL,
		SandboxURL:    SandboxURL,
		httpCli: &http.Client{
			Timeout: time.Millisecond,
		},
	}

	req := IAPRequest{
		ReceiptData: "dummy data",
	}
	result := &IAPResponse{}
	ctx := context.Background()
	err := client.Verify(ctx, req, result)
	if err == nil {
		t.Errorf("error should be occurred because of timeout")
	}
	t.Log(err)
}

func TestVerifyWithCancel(t *testing.T) {
	client := New()

	req := IAPRequest{
		ReceiptData: "dummy data",
	}
	result := &IAPResponse{}
	ctx, cancelFunc := context.WithCancel(context.Background())
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancelFunc()
	}()
	err := client.Verify(ctx, req, result)
	if err == nil {
		t.Errorf("error should be occurred because of context cancel")
	}
	t.Log(err)
}

func TestVerifyBadURL(t *testing.T) {
	client := New()
	client.ProductionURL = "127.0.0.1"

	req := IAPRequest{
		ReceiptData: "dummy data",
	}
	result := &IAPResponse{}
	ctx := context.Background()
	err := client.Verify(ctx, req, result)
	if err == nil {
		t.Errorf("error should be occurred because the server is not real")
	}
}

func TestResponses(t *testing.T) {
	req := IAPRequest{
		ReceiptData: "dummy data",
	}
	result := &IAPResponse{}

	type testCase struct {
		name        string
		testServer  *httptest.Server
		sandboxServ *httptest.Server
		expected    *IAPResponse
	}

	testCases := []testCase{
		{
			name:        "VerifySandboxReceipt",
			testServer:  httptest.NewServer(serverWithResponse(http.StatusOK, `{"status": 21007}`)),
			sandboxServ: httptest.NewServer(serverWithResponse(http.StatusOK, `{"status": 0}`)),
			expected: &IAPResponse{
				Status: 0,
			},
		},
		{
			name:       "VerifyBadPayload",
			testServer: httptest.NewServer(serverWithResponse(http.StatusOK, `{"status": 21002}`)),
			expected: &IAPResponse{
				Status: 21002,
			},
		},
		{
			name:       "SuccessPayload",
			testServer: httptest.NewServer(serverWithResponse(http.StatusBadRequest, `{"status": 0}`)),
			expected: &IAPResponse{
				Status: 0,
			},
		},
	}

	client := New()
	client.SandboxURL = "localhost"

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer tc.testServer.Close()
			client.ProductionURL = tc.testServer.URL
			if tc.sandboxServ != nil {
				client.SandboxURL = tc.sandboxServ.URL
			}

			ctx := context.Background()
			err := client.Verify(ctx, req, result)
			if err != nil {
				t.Errorf("%s", err)
			}
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("got %v\nwant %v", result, tc.expected)
			}
		})
	}
}

func TestHttpStatusErrors(t *testing.T) {
	req := IAPRequest{
		ReceiptData: "dummy data",
	}
	result := &IAPResponse{}

	type testCase struct {
		name       string
		testServer *httptest.Server
		err        error
	}

	testCases := []testCase{
		{
			name:       "status 200",
			testServer: httptest.NewServer(serverWithResponse(http.StatusOK, `{"status": 21000}`)),
			err:        nil,
		},
		{
			name:       "status 500",
			testServer: httptest.NewServer(serverWithResponse(http.StatusInternalServerError, `qwerty!@#$%^`)),
			err:        ErrAppStoreServer,
		},
	}

	client := New()
	client.SandboxURL = "localhost"

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer tc.testServer.Close()
			client.ProductionURL = tc.testServer.URL

			ctx := context.Background()
			err := client.Verify(ctx, req, result)
			if !errors.Is(err, tc.err) {
				t.Errorf("expected error to be not nil since the sandbox is not responding")
			}
		})
	}
}

func TestCannotReadBody(t *testing.T) {
	client := New()
	testResponse := http.Response{Body: ioutil.NopCloser(errReader(0))}

	ctx := context.Background()
	if client.parseResponse(&testResponse, IAPResponse{}, ctx, IAPRequest{}) == nil {
		t.Errorf("expected redirectToSandbox to fail to read the body")
	}
}

func TestCannotUnmarshalBody(t *testing.T) {
	client := New()
	testResponse := http.Response{Body: ioutil.NopCloser(strings.NewReader(`{"status": true}`))}

	ctx := context.Background()
	if client.parseResponse(&testResponse, StatusResponse{}, ctx, IAPRequest{}) == nil {
		t.Errorf("expected redirectToSandbox to fail to unmarshal the data")
	}
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}

func serverWithResponse(statusCode int, response string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if "POST" == r.Method {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(statusCode)
			w.Write([]byte(response))
		} else {
			w.Write([]byte(`unsupported request`))
		}
	})
}
