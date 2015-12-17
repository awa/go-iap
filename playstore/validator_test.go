package playstore

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/appengine/aetest"
)

type testSignature struct {
	PrivateKeyID string `json:"private_key_id"`
	PrivateKey   string `json:"private_key"`
	ClientEmail  string `json:"client_email"`
	ClientID     string `json:"client_id"`
	Type         string `json:"type"`
}

var testJSON = testSignature{
	PrivateKeyID: "dummyKeyID",
	PrivateKey:   "-----BEGIN PRIVATE KEY-----\nMIIBOQIBAAJBANXOa7wgs5KHMEVJmVo2eoRxEgeqiYF2oABPGYrebU+cQiE7Mwdy\nxv153DHME+9L9QzAj+fR4y5Rwva/fAsGAssCAwEAAQJATQwrFMtwCtC+22kvYywY\nsJuSlMKm9MmL1TCsErgfCj2rksRK1U+/ZY709tE3XJVYlZalWCeVhHTjs5p0pnk6\nYQIhAOw0FksytfIfpdfcREbful+LhFp1um5WjcVf7kQ73JDxAiEA57nJkG9pwnUd\nBCyIcElTVIAKU0+iFpd1208OnGxyT3sCIGaEBNkGXWmEytnxQ8DvAVjOmNcaGZwh\n/M4ZYLREtupBAiAsrpFkTWdqPKTcsi2Y4Tq1N39GMzvA+XGbWTIrDWo5UwIgHhp9\nEOnHuUuPCjpLfYM2vSFiYzaj8UJCImjkMtDwzbA=\n-----END PRIVATE KEY-----\n",
	ClientEmail:  "dummyEmail",
	ClientID:     "dummyClientID",
	Type:         "service_account",
}

func TestNew(t *testing.T) {
	// Exception scenario
	jsonKey, _ := json.Marshal(testJSON)
	expected := "oauth2: cannot fetch token: 400 Bad Request\nResponse: {\n  \"error\" : \"invalid_grant\"\n}"

	actual, _ := New(jsonKey)
	val := actual.httpClient.Transport.(*oauth2.Transport)
	token, err := val.Source.Token()
	if token != nil {
		t.Errorf("got %#v", token)
	}
	if err.Error() != expected {
		t.Errorf("got %v\nwant %v", err, expected)
	}

	// TODO Normal scenario
}

func TestSetTimeout(t *testing.T) {
	_timeout := time.Second * 3
	SetTimeout(_timeout)

	if timeout != _timeout {
		t.Errorf("got %#v\nwant %#v", timeout, _timeout)
	}
}

func TestVerifySubscription(t *testing.T) {
	// Exception scenario
	jsonKey, _ := json.Marshal(testJSON)

	expected := "Get https://www.googleapis.com/androidpublisher/v2/applications/package/purchases/subscriptions/subscriptionID/tokens/purchaseToken?alt=json: oauth2: cannot fetch token: 400 Bad Request\nResponse: {\n  \"error\" : \"invalid_grant\"\n}"

	client, _ := New(jsonKey)
	_, err := client.VerifySubscription("package", "subscriptionID", "purchaseToken")

	if err.Error() != expected {
		t.Errorf("got %v\nwant %v", err, expected)
	}

	// TODO Normal scenario
}

func TestVerifySubscriptionGAE(t *testing.T) {
	// Exception scenario

	expected := "googleapi: Error 401: Invalid Credentials, authError"

	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	client, _ := NewGAE(ctx)
	_, err = client.VerifySubscription("package", "subscriptionID", "purchaseToken")

	if err.Error() != expected {
		t.Errorf("got %v\nwant %v", err, expected)
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
	// Exception scenario
	jsonKey, _ := json.Marshal(testJSON)

	expected := "Get https://www.googleapis.com/androidpublisher/v2/applications/package/purchases/products/productID/tokens/purchaseToken?alt=json: oauth2: cannot fetch token: 400 Bad Request\nResponse: {\n  \"error\" : \"invalid_grant\"\n}"

	client, _ := New(jsonKey)
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

func TestCancelSubscription(t *testing.T) {
	// Exception scenario
	client := Client{nil}
	expected := errors.New("client is nil")
	actual := client.CancelSubscription("package", "productID", "purchaseToken")

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}

	jsonKey, _ := json.Marshal(testJSON)
	client, _ = New(jsonKey)
	expectedStr := "Post https://www.googleapis.com/androidpublisher/v2/applications/package/purchases/subscriptions/productID/tokens/purchaseToken:cancel?alt=json: oauth2: cannot fetch token: 400 Bad Request\nResponse: {\n  \"error\" : \"invalid_grant\"\n}"
	actual = client.CancelSubscription("package", "productID", "purchaseToken")

	if actual.Error() != expectedStr {
		t.Errorf("got %v\nwant %v", actual, expectedStr)
	}

	// TODO Normal scenario
}

func TestRefundSubscription(t *testing.T) {
	// Exception scenario
	client := Client{nil}
	expected := errors.New("client is nil")
	actual := client.RefundSubscription("package", "productID", "purchaseToken")

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}

	jsonKey, _ := json.Marshal(testJSON)
	client, _ = New(jsonKey)
	expectedStr := "Post https://www.googleapis.com/androidpublisher/v2/applications/package/purchases/subscriptions/productID/tokens/purchaseToken:refund?alt=json: oauth2: cannot fetch token: 400 Bad Request\nResponse: {\n  \"error\" : \"invalid_grant\"\n}"
	actual = client.RefundSubscription("package", "productID", "purchaseToken")

	if actual.Error() != expectedStr {
		t.Errorf("got %v\nwant %v", actual, expectedStr)
	}

	// TODO Normal scenario
}

func TestRevokeSubscription(t *testing.T) {
	// Exception scenario
	client := Client{nil}
	expected := errors.New("client is nil")
	actual := client.RevokeSubscription("package", "productID", "purchaseToken")

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}

	jsonKey, _ := json.Marshal(testJSON)
	client, _ = New(jsonKey)
	expectedStr := "Post https://www.googleapis.com/androidpublisher/v2/applications/package/purchases/subscriptions/productID/tokens/purchaseToken:revoke?alt=json: oauth2: cannot fetch token: 400 Bad Request\nResponse: {\n  \"error\" : \"invalid_grant\"\n}"
	actual = client.RevokeSubscription("package", "productID", "purchaseToken")

	if actual.Error() != expectedStr {
		t.Errorf("got %v\nwant %v", actual, expectedStr)
	}

	// TODO Normal scenario
}

func TestVerifySignature(t *testing.T) {
	receipt := `{"orderId":"GPA.xxxx-xxxx-xxxx-xxxxx","packageName":"my.package","productId":"myproduct","purchaseTime":1437564796303,"purchaseState":0,"developerPayload":"user001","purchaseToken":"some-token"}`

	// when public key format is invalid base64
	pubkey := "dummy_public_key"
	sig := "gj0N8LANKXOw4OhWkS1UZmDVUxM1UIP28F6bDzEp7BCqcVAe0DuDxmAY5wXdEgMRx/VM1Nl2crjogeV60OqCsbIaWqS/ZJwdP127aKR0jk8sbX36ssyYZ0DdZdBdCr1tBZ/eSW1GlGuD/CgVaxns0JaWecXakgoV7j+RF2AFbS4="
	expectedStr := "failed to decode public key"
	_, err := VerifySignature(pubkey, []byte(receipt), sig)
	if err.Error() != expectedStr {
		t.Errorf("got %v\nwant %v", err, expectedStr)
	}

	// when pub key is not rsa public key
	pubkey = "JTbngOdvBE0rfdOs3GeuBnPB+YEP1w/peM4VJbnVz+hN9Td25vPjAznX9YKTGQN4iDohZ07wtl+zYygIcpSCc2ozNZUs9pV0s5itayQo22aT5myJrQmkp94ZSGI2npDP4+FE6ZiF+7khl3qoE0rVZq4G2mfk5LIIyTPTSA4UvyQ="
	sig = "gj0N8LANKXOw4OhWkS1UZmDVUxM1UIP28F6bDzEp7BCqcVAe0DuDxmAY5wXdEgMRx/VM1Nl2crjogeV60OqCsbIaWqS/ZJwdP127aKR0jk8sbX36ssyYZ0DdZdBdCr1tBZ/eSW1GlGuD/CgVaxns0JaWecXakgoV7j+RF2AFbS4="
	expectedStr = "failed to parse public key"
	_, err = VerifySignature(pubkey, []byte(receipt), sig)
	if err.Error() != expectedStr {
		t.Errorf("got %v\nwant %v", err, expectedStr)
	}

	// when signature is invalid base64 format
	pubkey = "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDGvModvVUrqJ9C5fy8J77ZQ7JDC6+tf5iK8C74/3mjmcvwo4nmprCgzR/BQIEuZWJi8KX+jiJUXKXF90JPsXHkKAPq6A1SCga7kWvs/M8srMpjNS9zJdwZF+eDOR0+lJEihO04zlpAV9ybPJ3Q621y1HUeVpwdxDNLQpJTuIflnwIDAQAB"
	sig = "invalid_signature"
	expectedStr = "failed to decode signature"
	_, err = VerifySignature(pubkey, []byte(receipt), sig)
	if err.Error() != expectedStr {
		t.Errorf("got %v\nwant %v", err, expectedStr)
	}

	// when signature is invalid
	pubkey = "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDGvModvVUrqJ9C5fy8J77ZQ7JDC6+tf5iK8C74/3mjmcvwo4nmprCgzR/BQIEuZWJi8KX+jiJUXKXF90JPsXHkKAPq6A1SCga7kWvs/M8srMpjNS9zJdwZF+eDOR0+lJEihO04zlpAV9ybPJ3Q621y1HUeVpwdxDNLQpJTuIflnwIDAQAB"
	sig = "JTbngOdvBE0rfdOs3GeuBnPB+YEP1w/peM4VJbnVz+hN9Td25vPjAznX9YKTGQN4iDohZ07wtl+zYygIcpSCc2ozNZUs9pV0s5itayQo22aT5myJrQmkp94ZSGI2npDP4+FE6ZiF+7khl3qoE0rVZq4G2mfk5LIIyTPTSA4UvyQ="
	isValid, err := VerifySignature(pubkey, []byte(receipt), sig)
	if err != nil {
		t.Errorf("got %v\n", err)
	}
	if isValid {
		t.Errorf("got %v\nwant %v", isValid, false)
	}

	// when all arguments are valid
	pubkey = "MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDGvModvVUrqJ9C5fy8J77ZQ7JDC6+tf5iK8C74/3mjmcvwo4nmprCgzR/BQIEuZWJi8KX+jiJUXKXF90JPsXHkKAPq6A1SCga7kWvs/M8srMpjNS9zJdwZF+eDOR0+lJEihO04zlpAV9ybPJ3Q621y1HUeVpwdxDNLQpJTuIflnwIDAQAB"
	sig = "gj0N8LANKXOw4OhWkS1UZmDVUxM1UIP28F6bDzEp7BCqcVAe0DuDxmAY5wXdEgMRx/VM1Nl2crjogeV60OqCsbIaWqS/ZJwdP127aKR0jk8sbX36ssyYZ0DdZdBdCr1tBZ/eSW1GlGuD/CgVaxns0JaWecXakgoV7j+RF2AFbS4="
	isValid, err = VerifySignature(pubkey, []byte(receipt), sig)
	if err != nil {
		t.Errorf("got %v\n", err)
	}
	if !isValid {
		t.Errorf("got %v\nwant %v", isValid, true)
	}
}
