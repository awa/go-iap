package playstore

import (
	"errors"
	"os"
	"reflect"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

func TestInit(t *testing.T) {
	expected := &oauth2.Config{
		ClientID:     "dummyId",
		ClientSecret: "dummySecret",
		Scopes:       []string{"https://www.googleapis.com/auth/androidpublisher"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}
	os.Setenv("IAB_CLIENT_ID", "dummyId")
	os.Setenv("IAB_CLIENT_SECRET", "dummySecret")
	Init()
	os.Clearenv()
	actual := defaultConfig
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestInitWithoutClientSecret(t *testing.T) {
	expected := errors.New("Client Secret Key is required")

	os.Setenv("IAB_CLIENT_ID", "dummyId")
	actual := Init()
	os.Clearenv()

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestInitWithConfig(t *testing.T) {
	expected := &oauth2.Config{
		ClientID:     "dummyId",
		ClientSecret: "dummySecret",
		Scopes:       []string{"https://www.googleapis.com/auth/androidpublisher"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}

	config := &oauth2.Config{
		ClientID:     "dummyId",
		ClientSecret: "dummySecret",
	}

	InitWithConfig(config)
	actual := defaultConfig
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestInitWithConfigErrors(t *testing.T) {
	expected := errors.New("Client ID is required")

	config := &oauth2.Config{
		Scopes: []string{"https://www.googleapis.com/auth/androidpublisher"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}
	actual := InitWithConfig(config)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}

	expected = errors.New("Client Secret Key is required")
	config = &oauth2.Config{
		ClientID: "dummyId",
		Scopes:   []string{"https://www.googleapis.com/auth/androidpublisher"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}
	actual = InitWithConfig(config)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestNew(t *testing.T) {
	// Initialize config
	_config := &oauth2.Config{
		ClientID:     "dummyId",
		ClientSecret: "dummySecret",
		RedirectURL:  "REDIRECT_URL",
		Scopes:       []string{"scope1", "scope2"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "http://example.com/auth",
			TokenURL: "http://example.com/token",
		},
	}
	InitWithConfig(_config)

	_token := &oauth2.Token{
		AccessToken:  "accessToken",
		RefreshToken: "refreshToken",
		Expiry:       time.Unix(2234567890, 0).UTC(),
	}

	actual := New(_token)
	val := actual.httpClient.Transport.(*oauth2.Transport)
	token, _ := val.Source.Token()
	if !reflect.DeepEqual(token, _token) {
		t.Errorf("got %v\nwant %v", token, _token)
	}
}

func TestSetTimeout(t *testing.T) {
	timeout := time.Second * 3
	SetTimeout(timeout)

	if defaultTimeout != timeout {
		t.Errorf("got %#v\nwant %#v", defaultTimeout, timeout)
	}
}

func TestVerifySubscription(t *testing.T) {
	Init()

	// Exception scenario
	token := &oauth2.Token{
		AccessToken:  "accessToken",
		RefreshToken: "refreshToken",
		Expiry:       time.Unix(2234567890, 0).UTC(),
	}

	client := New(token)
	expected := "googleapi: Error 401: Invalid Credentials, authError"
	_, err := client.VerifySubscription("package", "subscriptionID", "purchaseToken")

	if err.Error() != expected {
		t.Errorf("got %v", err)
	}

	// TODO Normal scenario
}

func TestVerifySubscriptionAndroidPublisherError(t *testing.T) {
	client := Client{nil}
	expected := errors.New("client is nil")
	_, actual := client.VerifySubscription("package", "subscriptionID", "purchaseToken")

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestVerifyProduct(t *testing.T) {
	Init()

	// Exception scenario
	token := &oauth2.Token{
		AccessToken:  "accessToken",
		RefreshToken: "refreshToken",
		Expiry:       time.Unix(2234567890, 0).UTC(),
	}

	client := New(token)
	expected := "googleapi: Error 401: Invalid Credentials, authError"
	_, err := client.VerifyProduct("package", "productID", "purchaseToken")

	if err.Error() != expected {
		t.Errorf("got %v", err)
	}

	// TODO Normal scenario
}

func TestVerifyProductAndroidPublisherError(t *testing.T) {
	client := Client{nil}
	expected := errors.New("client is nil")
	_, actual := client.VerifyProduct("package", "productID", "purchaseToken")

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}
