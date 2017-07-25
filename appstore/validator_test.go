package appstore

import (
	"errors"
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
	client.TimeOut = time.Millisecond * 100

	req := IAPRequest{
		ReceiptData: "dummy data",
	}
	result := &IAPResponse{}
	err := client.Verify(req, result)
	if err == nil {
		t.Errorf("error should be occurred because of timeout")
	}

	client = New()
	expected := &IAPResponse{
		Status: 21002,
	}
	client.Verify(req, result)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("got %v\nwant %v", result, expected)
	}
}
