package appstore

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
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

	// status unkown
	expected = errors.New("An unknown error ocurred")
	actual = HandleError(100)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestNew(t *testing.T) {
	expected := Client{
		URL:     "https://sandbox.itunes.apple.com/verifyReceipt",
		TimeOut: time.Second * 5,
	}

	actual := New()
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestNewWithEnvironment(t *testing.T) {
	expected := Client{
		URL:     "https://buy.itunes.apple.com/verifyReceipt",
		TimeOut: time.Second * 5,
	}

	os.Setenv("IAP_ENVIRONMENT", "production")
	actual := New()
	os.Clearenv()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestNewWithConfig(t *testing.T) {
	config := Config{
		IsProduction: true,
		TimeOut:      time.Second * 2,
	}

	expected := Client{
		URL:     "https://buy.itunes.apple.com/verifyReceipt",
		TimeOut: time.Second * 2,
	}

	actual := NewWithConfig(config)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestNewWithConfigTimeout(t *testing.T) {
	config := Config{
		IsProduction: true,
	}

	expected := Client{
		URL:     "https://buy.itunes.apple.com/verifyReceipt",
		TimeOut: time.Second * 5,
	}

	actual := NewWithConfig(config)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestVerify(t *testing.T) {
	client := New()

	expected := IAPResponse{
		Status: 21002,
	}
	req := IAPRequest{
		ReceiptData: "dummy data",
	}
	actual, _ := client.Verify(req)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestVerifyErrors(t *testing.T) {
	server, client := testTools(199, "dummy response")
	defer server.Close()

	req := IAPRequest{
		ReceiptData: "dummy data",
	}

	expected := errors.New("An error occurred in IAP - code:199")
	_, actual := client.Verify(req)
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestVerifyTimeout(t *testing.T) {
	// HTTP 100 is "continue" so it will time out
	server, client := testTools(100, "dummy response")
	defer server.Close()

	req := IAPRequest{
		ReceiptData: "dummy data",
	}

	expected := errors.New("")
	_, actual := client.Verify(req)
	if !reflect.DeepEqual(reflect.TypeOf(actual), reflect.TypeOf(expected)) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func testTools(code int, body string) (*httptest.Server, *Client) {

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, body)
	}))

	client := &Client{URL: server.URL, TimeOut: time.Second * 2}
	return server, client
}
