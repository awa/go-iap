package playstore

import (
	"os"
	"reflect"
	"testing"
	"time"

	"code.google.com/p/goauth2/oauth"
)

func TestInit(t *testing.T) {
	expected := &oauth.Config{
		ClientId:     "dummyId",
		ClientSecret: "dummySecret",
		Scope:        "https://www.googleapis.com/auth/androidpublisher",
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://accounts.google.com/o/oauth2/token",
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

func TestInitWithConfig(t *testing.T) {
	expected := &oauth.Config{
		ClientId:     "dummyId",
		ClientSecret: "dummySecret",
		Scope:        "https://www.googleapis.com/auth/androidpublisher",
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://accounts.google.com/o/oauth2/token",
	}

	config := &oauth.Config{
		ClientId:     "dummyId",
		ClientSecret: "dummySecret",
		Scope:        "https://www.googleapis.com/auth/androidpublisher",
		AuthURL:      "https://accounts.google.com/o/oauth2/auth",
		TokenURL:     "https://accounts.google.com/o/oauth2/token",
	}
	InitWithConfig(config)
	actual := defaultConfig
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestNew(t *testing.T) {
	// Initialize config
	_config := &oauth.Config{
		ClientId:     "dummyId",
		ClientSecret: "dummySecret",
	}
	InitWithConfig(_config)

	token := &oauth.Token{
		AccessToken:  "accessToken",
		RefreshToken: "refreshToken",
		Expiry:       time.Unix(1234567890, 0).UTC(),
	}

	actual := New(token)
	val, _ := actual.httpClient.Transport.(*oauth.Transport)

	if !reflect.DeepEqual(val.Config, _config) {
		t.Errorf("got %v\nwant %v", val.Config, _config)
	}

	if !reflect.DeepEqual(val.Token, token) {
		t.Errorf("got %v\nwant %v", val.Token, token)
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
	token := &oauth.Token{
		AccessToken:  "accessToken",
		RefreshToken: "refreshToken",
		Expiry:       time.Unix(1234567890, 0).UTC(),
	}

	client := New(token)
	expected := "Get https://www.googleapis.com/androidpublisher/v2/applications/package/purchases/subscriptions/subscriptionID/tokens/purchaseToken?alt=json: OAuthError: updateToken: Unexpected HTTP status 400 Bad Request"
	_, err := client.VerifySubscription("package", "subscriptionID", "purchaseToken")

	if err.Error() != expected {
		t.Errorf("got %v", err)
	}

	// TODO Nomal scenario
}
