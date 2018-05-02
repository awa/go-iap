package appstore

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestHandleError(t *testing.T) {
	var expected, actual error

	// status 0
	expected = nil
	actual = HandleError(0)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}

	// status 21000
	expected = errors.New("The App Store could not read the JSON object you provided.")
	actual = HandleError(21000)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}

	// status 21002
	expected = errors.New("The data in the receipt-data property was malformed or missing.")
	actual = HandleError(21002)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}

	// status 21003
	expected = errors.New("The receipt could not be authenticated.")
	actual = HandleError(21003)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}

	// status 21004
	expected = errors.New("The shared secret you provided does not match the shared secret on file for your account.")
	actual = HandleError(21004)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}

	// status 21005
	expected = errors.New("The receipt server is not currently available.")
	actual = HandleError(21005)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}

	// status 21007
	expected = errors.New("This receipt is from the test environment, but it was sent to the production environment for verification. Send it to the test environment instead.")
	actual = HandleError(21007)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}

	// status 21008
	expected = errors.New("This receipt is from the production environment, but it was sent to the test environment for verification. Send it to the production environment instead.")
	actual = HandleError(21008)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}

	// status 21010
	expected = errors.New("This receipt could not be authorized. Treat this the same as if a purchase was never made.")
	actual = HandleError(21010)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}

	// status 21100 - 21199
	expected = errors.New("Internal data access error.")
	actual = HandleError(21155)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}

	// status unknown
	expected = errors.New("An unknown error occurred")
	actual = HandleError(100)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestNew(t *testing.T) {
	expected := Client{
		ProductionURL: ProductionURL,
		SandboxURL:    SandboxURL,
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	actual := New()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestNewWithEnvironment(t *testing.T) {
	expected := Client{
		ProductionURL: ProductionURL,
		SandboxURL:    SandboxURL,
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	os.Setenv("IAP_ENVIRONMENT", "production")
	actual := New()
	os.Clearenv()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestSetHTTPClient(t *testing.T) {
	expected := Client{
		ProductionURL: ProductionURL,
		SandboxURL:    SandboxURL,
		HTTPClient: &http.Client{
			Timeout: time.Second * 2,
		},
	}

	actual := New()
	actual.SetHTTPClient(&http.Client{
		Timeout: time.Second * 2,
	})
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestVerifyTimeout(t *testing.T) {
	client := New()
	client.SetHTTPClient(&http.Client{
		Timeout: time.Millisecond,
	})

	req := IAPRequest{
		ReceiptData: "dummy data",
	}
	result := &IAPResponse{}
	err := client.Verify(req, result)
	if err == nil {
		t.Errorf("error should be occurred because of timeout")
	}
}

func TestVerifyBadURL(t *testing.T) {
	client := New()
	client.ProductionURL = "127.0.0.1"

	req := IAPRequest{
		ReceiptData: "dummy data",
	}
	result := &IAPResponse{}
	err := client.Verify(req, result)
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
	client.SetHTTPClient(&http.Client{
		Timeout: time.Second * 100,
	})
	client.SandboxURL = "localhost"

	for i, tc := range testCases {
		defer tc.testServer.Close()
		client.ProductionURL = tc.testServer.URL
		if tc.sandboxServ != nil {
			client.SandboxURL = tc.sandboxServ.URL
		}

		err := client.Verify(req, result)
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
	client.SetHTTPClient(&http.Client{
		Timeout: time.Second * 100,
	})
	client.SandboxURL = "localhost"

	for i, tc := range testCases {
		defer tc.testServer.Close()
		client.ProductionURL = tc.testServer.URL

		err := client.Verify(req, result)
		if err == nil {
			t.Errorf("Test case %d - expected error to be not nil since the sandbox is not responding", i)
		}
	}
}

func TestCannotReadBody(t *testing.T) {
	client := New()
	testResponse := http.Response{Body: ioutil.NopCloser(errReader(0))}

	if client.parseResponse(&testResponse, IAPResponse{}, IAPRequest{}) == nil {
		t.Errorf("expected redirectToSandbox to fail to read the body")
	}
}

func TestCannotUnmarshalBody(t *testing.T) {
	client := New()
	testResponse := http.Response{Body: ioutil.NopCloser(strings.NewReader(`{"status": true}`))}

	if client.parseResponse(&testResponse, StatusResponse{}, IAPRequest{}) == nil {
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
			w.Write([]byte(response))
			return
		} else {
			w.Write([]byte(`unsupported request`))
		}

		w.WriteHeader(statusCode)
	})
}
