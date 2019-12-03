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
			out:  errors.New("The App Store could not read the JSON object you provided."),
		},
		{
			name: "status 21002",
			in:   21002,
			out:  errors.New("The data in the receipt-data property was malformed or missing."),
		},
		{
			name: "status 21003",
			in:   21003,
			out:  errors.New("The receipt could not be authenticated."),
		},
		{
			name: "status 21004",
			in:   21004,
			out:  errors.New("The shared secret you provided does not match the shared secret on file for your account."),
		},
		{
			name: "status 21005",
			in:   21005,
			out:  errors.New("The receipt server is not currently available."),
		},
		{
			name: "status 21007",
			in:   21007,
			out:  errors.New("This receipt is from the test environment, but it was sent to the production environment for verification. Send it to the test environment instead."),
		},
		{
			name: "status 21008",
			in:   21008,
			out:  errors.New("This receipt is from the production environment, but it was sent to the test environment for verification. Send it to the production environment instead."),
		},
		{
			name: "status 21010",
			in:   21010,
			out:  errors.New("This receipt could not be authorized. Treat this the same as if a purchase was never made."),
		},
		{
			name: "status 21100 ~ 21199",
			in:   21100,
			out:  errors.New("Internal data access error."),
		},
		{
			name: "status unknown",
			in:   100,
			out:  errors.New("An unknown error occurred"),
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			out := HandleError(v.in)

			if !reflect.DeepEqual(out, v.out) {
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
			Timeout: 10 * time.Second,
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
		testServer  *httptest.Server
		sandboxServ *httptest.Server
		expected    *IAPResponse
	}

	testCases := []testCase{
		// VerifySandboxReceipt
		{
			testServer:  httptest.NewServer(serverWithResponse(http.StatusOK, `{"status": 21007}`)),
			sandboxServ: httptest.NewServer(serverWithResponse(http.StatusOK, `{"status": 0}`)),
			expected: &IAPResponse{
				Status: 0,
			},
		},
		// VerifyBadPayload
		{
			testServer: httptest.NewServer(serverWithResponse(http.StatusOK, `{"status": 21002}`)),
			expected: &IAPResponse{
				Status: 21002,
			},
		},
		// SuccessPayload
		{
			testServer: httptest.NewServer(serverWithResponse(http.StatusBadRequest, `{"status": 0}`)),
			expected: &IAPResponse{
				Status: 0,
			},
		},
	}

	client := New()
	client.SandboxURL = "localhost"

	for i, tc := range testCases {
		defer tc.testServer.Close()
		client.ProductionURL = tc.testServer.URL
		if tc.sandboxServ != nil {
			client.SandboxURL = tc.sandboxServ.URL
		}

		ctx := context.Background()
		err := client.Verify(ctx, req, result)
		if err != nil {
			t.Errorf("Test case %d - %s", i, err.Error())
		}
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("Test case %d - got %v\nwant %v", i, result, tc.expected)
		}
	}
}

func TestErrors(t *testing.T) {
	req := IAPRequest{
		ReceiptData: "dummy data",
	}
	result := &IAPResponse{}

	type testCase struct {
		testServer *httptest.Server
	}

	testCases := []testCase{
		// VerifySandboxReceiptFailure
		{
			testServer: httptest.NewServer(serverWithResponse(http.StatusOK, `{"status": 21007}`)),
		},
		// VerifyBadResponse
		{
			testServer: httptest.NewServer(serverWithResponse(http.StatusInternalServerError, `qwerty!@#$%^`)),
		},
	}

	client := New()
	client.SandboxURL = "localhost"

	for i, tc := range testCases {
		defer tc.testServer.Close()
		client.ProductionURL = tc.testServer.URL

		ctx := context.Background()
		err := client.Verify(ctx, req, result)
		if err == nil {
			t.Errorf("Test case %d - expected error to be not nil since the sandbox is not responding", i)
		}
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
		if "POST" != r.Method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(`unsupported request`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		w.Write([]byte(response))
	})
}
