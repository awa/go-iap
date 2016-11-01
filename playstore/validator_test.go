package playstore

import (
	"encoding/base64"
	"errors"
	"reflect"
	"testing"
	"time"

	"golang.org/x/oauth2"
)

var base64JsonKey = "ew0KICAidHlwZSI6ICJzZXJ2aWNlX2FjY291bnQiLA0KICAicHJvamVjdF9pZCI6ICJnby1pYXAiLA0KICAicHJpdmF0ZV9rZXlfaWQiOiAiZjI0MGRmZTg4Y2NlMTBmMThhZjM2ZjdhNTEyNDUyZjE0Y2RkM2EyZSIsDQogICJwcml2YXRlX2tleSI6ICItLS0tLUJFR0lOIFBSSVZBVEUgS0VZLS0tLS1cbk1JSUV2Z0lCQURBTkJna3Foa2lHOXcwQkFRRUZBQVNDQktnd2dnU2tBZ0VBQW9JQkFRQ29sbndXRmhacjlmQUNcbjNqcEozbGk5SSt4ZzUyZGl3R0pzUnBoY0U0S0xielJpM1ZWdUEvMU1EMkcwMmNOYkZ5V3NDSWg4aDMvNE1sZmxcbmtMa29wVlJlNVFlczRMY3E4RDRVOUdubjZwWXduWlI5NVBlK2Z6TVdLNGJpSWgvdHhNMzJlemxvQWwvaU5Rbi9cbmdQQ3lUc0pHNlVlNzJsdFJyblFXdkM0L24vQzNHL3lLY2ZLaHo2ODZ6NTd4ODZma28yUnEra21RaFNmY3REY3pcbng3YlNmUDU4SFdUOVhDNHlHSERqZmkwMmxnVU1YaW5Fa2VQVzB1UC9SZG8ya0dlaWhSV1JGMFFNRHlld2p1S0pcbmtaL3pCOGhyZm9RVEdmMjE4L29lQXR0UWl0ZkxkVlkrajc5Y2VxQTlBUXorbGoxQjUrSVZiRjlOQ2toaTNYYThcbndHMml6WEViQWdNQkFBRUNnZ0VBTFZnRVdnQm8yWExWc2ovSlY3THBGVDVEUnJFV3VwWGFJeHM5d1k0eHo0VUNcbmh4RFcrSGMwT3Evc2JMTWhleStYbjFUUU9RWk00aG5RVUZ1RG9hNE9LbFBabzZMeFFTaEsybUgrMWpUZlhvWVRcbnVXVExTYjUycENEaTc1R1VHdVNUTFJkcGtsTUpMUk5zOC83ZlBtWTJsTklMekRmbjFlbGhLZmhGVERHZGtmSVNcbkprNjQvd0JTRUx5ZTFQOFdYM0JLV3hZdFBKajgvc0pqTDl0SVhXcitaSjR2bGgzM1c1b25BWS9lMFJoZkZDSjJcbjBFWkZFVXkrblNUVHBUR1JmYzYyMFl2c2JyVWozcFo3dDBIRE5IN3ovbnVEeVB4UUpORTg1RVlQRDBaeU9VNnpcbmV2aW9ObmJBKzFVREdFOVg5NFVVbnJNckJ5NU5Jbk13bCt3djZ5VWFJUUtCZ1FEYnFhazVXbTRla1FJZVpsMmhcbkZ3ZU9TR0FOek9ZVmE3TWExQm9hUHBLKzVuak1xbUQxRHB4cGlFYmFGNnpwVWMyZkVIT0ZEZ21zQUZrYUw1cmNcbnRzVC84amtXRHdKSjAvQlE4QlpiekZuU1I4bHhlbFBscFNySVcremcrNWdhbDZtTXpLS2J0a29TVUsxem5uMnpcbkEydmhTTis5N1VsYTByK2d3ZCttQ0gvTVN3S0JnUURFZWVPbE1vRzM2eWxlRXFOaDdhNUlySWt0dFlCbnJuMXFcblhCSFJmWFQ5c3ROOHhPL0pYQ2s1S1JIREkzUWZJbVcrSmdJV3NwRjY1Q00rZnFRN1l2UWNwV0RGdDZvc1hRNm5cbmVhZC9RYlJ5TzlLcXY3SEo2aHJycGJrS0NvSjBBbFRKZmt4T2svTDBnM3Zkblhia242L3pVeC9JZXRsZHE4T1NcbjNlU2FBbnBNY1FLQmdRQ3diQWQ2Qk9ORXNYcGVLQ0V5N0dncEluL2pGWm9Gd2taTFdlYk5CVXllL2tRdlBQZzZcbldjM09CS0hETUJpMEcvdGxzYlRXUEh3UUpRZHJQS2pJZEJLczdrSmpNUkxKY09zbVZtM2V0TFcvYWVDa3YzYjZcbmpqbGFTbHBxS0NmMTA3RmRZRTJKZWxMcmV0aVViOHJOS0FaUkhsSjFIRXM2SXVHOW4zaWN4VjYvR1FLQmdISElcbnJUK0Vpbjg2MzFBdHR4VUZrd05mZUdwU1RMUys1cjdyNXgzTmJDMW9uUFlMRDFzcjFtdldEd1ZWeVBBbStZa3ZcbmRkSXpRL0ZKb2VlVmJBTkFnV0w5bTVlbGtCWDFKb0Z6QUwvQUM0S0VocktBSmJScnNYOTdFRGh5Y2E1QmsxekZcbm1lZC80eG9iODJZYXhUb09DTlgvODg0azV6RktRZzhTRmt2aTEzVGhBb0dCQUxRV1krV0hDL0hDU3VhaVQ1ck5cbnFiclIwM01BVTg4T3UyelZLZ1FLa0t0eDhaN0RLY2duRmNEWktQOUdPcXJFVjdhZElxdDZIckwvR1FqdVFDSGFcblRoc3dGRmVZM005MEViempGTE1WWVhlaHlFYzJidkJmcFpFSG5UM1VSNXI1ZWRYWmpGYnJySFdDMGFsYnNSR1pcbnZoQ2dBT3c4bk54SFRtMlBrMkZ2bHZyWFxuLS0tLS1FTkQgUFJJVkFURSBLRVktLS0tLVxuIiwNCiAgImNsaWVudF9lbWFpbCI6ICJnby1pYXBAZ28taWFwLmlhbS5nc2VydmljZWFjY291bnQuY29tIiwNCiAgImNsaWVudF9pZCI6ICIxMDMxMTMwNDc5ODIwNjMwMTE3MjciLA0KICAiYXV0aF91cmkiOiAiaHR0cHM6Ly9hY2NvdW50cy5nb29nbGUuY29tL28vb2F1dGgyL2F1dGgiLA0KICAidG9rZW5fdXJpIjogImh0dHBzOi8vYWNjb3VudHMuZ29vZ2xlLmNvbS9vL29hdXRoMi90b2tlbiIsDQogICJhdXRoX3Byb3ZpZGVyX3g1MDlfY2VydF91cmwiOiAiaHR0cHM6Ly93d3cuZ29vZ2xlYXBpcy5jb20vb2F1dGgyL3YxL2NlcnRzIiwNCiAgImNsaWVudF94NTA5X2NlcnRfdXJsIjogImh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL3JvYm90L3YxL21ldGFkYXRhL3g1MDkvZ28taWFwJTQwZ28taWFwLmlhbS5nc2VydmljZWFjY291bnQuY29tIg0KfQ=="

var base64dummyKey = "ew0KICAidHlwZSI6ICJzZXJ2aWNlX2FjY291bnQiLA0KICAicHJvamVjdF9pZCI6ICJnby1pYXAiLA0KICAicHJpdmF0ZV9rZXlfaWQiOiAiZHVtbXkiLA0KICAicHJpdmF0ZV9rZXkiOiAiLS0tLS1CRUdJTiBQUklWQVRFIEtFWS0tLS0tXG5NSUlFdmdJQkFEQU5CZ2txaGtpRzl3MEJBUUVGQUFTQ0JLZ3dnZ1NrQWdFQUFvSUJBUUNvbG53V0ZoWnI5ZkFDXG4zanBKM2xpOUkreGc1MmRpd0dKc1JwaGNFNEtMYnpSaTNWVnVBLzFNRDJHMDJjTmJGeVdzQ0loOGgzLzRNbGZsXG5rTGtvcFZSZTVRZXM0TGNxOEQ0VTlHbm42cFl3blpSOTVQZStmek1XSzRiaUloL3R4TTMyZXpsb0FsL2lOUW4vXG5nUEN5VHNKRzZVZTcybHRScm5RV3ZDNC9uL0MzRy95S2NmS2h6Njg2ejU3eDg2ZmtvMlJxK2ttUWhTZmN0RGN6XG54N2JTZlA1OEhXVDlYQzR5R0hEamZpMDJsZ1VNWGluRWtlUFcwdVAvUmRvMmtHZWloUldSRjBRTUR5ZXdqdUtKXG5rWi96QjhocmZvUVRHZjIxOC9vZUF0dFFpdGZMZFZZK2o3OWNlcUE5QVF6K2xqMUI1K0lWYkY5TkNraGkzWGE4XG53RzJpelhFYkFnTUJBQUVDZ2dFQUxWZ0VXZ0JvMlhMVnNqL0pWN0xwRlQ1RFJyRVd1cFhhSXhzOXdZNHh6NFVDXG5oeERXK0hjME9xL3NiTE1oZXkrWG4xVFFPUVpNNGhuUVVGdURvYTRPS2xQWm82THhRU2hLMm1IKzFqVGZYb1lUXG51V1RMU2I1MnBDRGk3NUdVR3VTVExSZHBrbE1KTFJOczgvN2ZQbVkybE5JTHpEZm4xZWxoS2ZoRlRER2RrZklTXG5KazY0L3dCU0VMeWUxUDhXWDNCS1d4WXRQSmo4L3NKakw5dElYV3IrWko0dmxoMzNXNW9uQVkvZTBSaGZGQ0oyXG4wRVpGRVV5K25TVFRwVEdSZmM2MjBZdnNiclVqM3BaN3QwSEROSDd6L251RHlQeFFKTkU4NUVZUEQwWnlPVTZ6XG5ldmlvTm5iQSsxVURHRTlYOTRVVW5yTXJCeTVOSW5Nd2wrd3Y2eVVhSVFLQmdRRGJxYWs1V200ZWtRSWVabDJoXG5Gd2VPU0dBTnpPWVZhN01hMUJvYVBwSys1bmpNcW1EMURweHBpRWJhRjZ6cFVjMmZFSE9GRGdtc0FGa2FMNXJjXG50c1QvOGprV0R3SkowL0JROEJaYnpGblNSOGx4ZWxQbHBTcklXK3pnKzVnYWw2bU16S0tidGtvU1VLMXpubjJ6XG5BMnZoU04rOTdVbGEwcitnd2QrbUNIL01Td0tCZ1FERWVlT2xNb0czNnlsZUVxTmg3YTVJcklrdHRZQm5ybjFxXG5YQkhSZlhUOXN0Tjh4Ty9KWENrNUtSSERJM1FmSW1XK0pnSVdzcEY2NUNNK2ZxUTdZdlFjcFdERnQ2b3NYUTZuXG5lYWQvUWJSeU85S3F2N0hKNmhycnBia0tDb0owQWxUSmZreE9rL0wwZzN2ZG5YYmtuNi96VXgvSWV0bGRxOE9TXG4zZVNhQW5wTWNRS0JnUUN3YkFkNkJPTkVzWHBlS0NFeTdHZ3BJbi9qRlpvRndrWkxXZWJOQlV5ZS9rUXZQUGc2XG5XYzNPQktIRE1CaTBHL3Rsc2JUV1BId1FKUWRyUEtqSWRCS3M3a0pqTVJMSmNPc21WbTNldExXL2FlQ2t2M2I2XG5qamxhU2xwcUtDZjEwN0ZkWUUySmVsTHJldGlVYjhyTktBWlJIbEoxSEVzNkl1RzluM2ljeFY2L0dRS0JnSEhJXG5yVCtFaW44NjMxQXR0eFVGa3dOZmVHcFNUTFMrNXI3cjV4M05iQzFvblBZTEQxc3IxbXZXRHdWVnlQQW0rWWt2XG5kZEl6US9GSm9lZVZiQU5BZ1dMOW01ZWxrQlgxSm9GekFML0FDNEtFaHJLQUpiUnJzWDk3RURoeWNhNUJrMXpGXG5tZWQvNHhvYjgyWWF4VG9PQ05YLzg4NGs1ekZLUWc4U0ZrdmkxM1RoQW9HQkFMUVdZK1dIQy9IQ1N1YWlUNXJOXG5xYnJSMDNNQVU4OE91MnpWS2dRS2tLdHg4WjdES2NnbkZjRFpLUDlHT3FyRVY3YWRJcXQ2SHJML0dRanVRQ0hhXG5UaHN3RkZlWTNNOTBFYnpqRkxNVllYZWh5RWMyYnZCZnBaRUhuVDNVUjVyNWVkWFpqRmJyckhXQzBhbGJzUkdaXG52aENnQU93OG5OeEhUbTJQazJGdmx2clhcbi0tLS0tRU5EIFBSSVZBVEUgS0VZLS0tLS1cbiIsDQogICJjbGllbnRfaWQiOiAiZHVtbXkiLA0KICAiYXV0aF91cmkiOiAiaHR0cHM6Ly9hY2NvdW50cy5nb29nbGUuY29tL28vb2F1dGgyL2F1dGgiLA0KICAidG9rZW5fdXJpIjogImh0dHBzOi8vYWNjb3VudHMuZ29vZ2xlLmNvbS9vL29hdXRoMi90b2tlbiIsDQogICJhdXRoX3Byb3ZpZGVyX3g1MDlfY2VydF91cmwiOiAiaHR0cHM6Ly93d3cuZ29vZ2xlYXBpcy5jb20vb2F1dGgyL3YxL2NlcnRzIiwNCiAgImNsaWVudF94NTA5X2NlcnRfdXJsIjogImh0dHBzOi8vd3d3Lmdvb2dsZWFwaXMuY29tL3JvYm90L3YxL21ldGFkYXRhL3g1MDkvZ28taWFwJTQwZ28taWFwLmlhbS5nc2VydmljZWFjY291bnQuY29tIg0KfQ=="

var jsonKey []byte
var dummyKey []byte

func init() {
	f, err := base64.StdEncoding.DecodeString(base64JsonKey)
	if err != nil {
		panic(err)
	}
	jsonKey = f
	d, err := base64.StdEncoding.DecodeString(base64dummyKey)
	if err != nil {
		panic(err)
	}
	dummyKey = d
}

func TestNew(t *testing.T) {
	t.Parallel()
	// Exception scenario
	expected := "oauth2: cannot fetch token: 401 Unauthorized\nResponse: {\n  \"error\" : \"invalid_client\",\n  \"error_description\" : \"The OAuth client was invalid.\"\n}"

	actual, _ := New(dummyKey)
	val := actual.httpClient.Transport.(*oauth2.Transport)
	token, err := val.Source.Token()
	if token != nil {
		t.Errorf("got %#v", token)
	}
	if err.Error() != expected {
		t.Errorf("got %v\nwant %v", err, expected)
	}

	// TODO Normal scenario
	actual, _ = New(jsonKey)
	val = actual.httpClient.Transport.(*oauth2.Transport)
	token, err = val.Source.Token()
	if err != nil {
		t.Errorf("got %#v", err)
	}
}

func TestSetTimeout(t *testing.T) {
	t.Parallel()
	_timeout := time.Second * 3
	SetTimeout(_timeout)

	if timeout != _timeout {
		t.Errorf("got %#v\nwant %#v", timeout, _timeout)
	}
}

func TestVerifySubscription(t *testing.T) {
	t.Parallel()
	// Exception scenario
	expected := "googleapi: Error 404: No application was found for the given package name., applicationNotFound"

	client, _ := New(jsonKey)
	_, err := client.VerifySubscription("package", "subscriptionID", "purchaseToken")

	if err.Error() != expected {
		t.Errorf("got %v\nwant %v", err, expected)
	}

	// TODO Normal scenario
}

func TestVerifySubscriptionAndroidPublisherError(t *testing.T) {
	t.Parallel()
	client := Client{nil}
	expected := errors.New("client is nil")
	_, actual := client.VerifySubscription("package", "subscriptionID", "purchaseToken")

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestVerifyProduct(t *testing.T) {
	t.Parallel()
	// Exception scenario
	expected := "googleapi: Error 404: No application was found for the given package name., applicationNotFound"

	client, _ := New(jsonKey)
	_, err := client.VerifyProduct("package", "productID", "purchaseToken")

	if err.Error() != expected {
		t.Errorf("got %v", err)
	}

	// TODO Normal scenario
}

func TestVerifyProductAndroidPublisherError(t *testing.T) {
	t.Parallel()
	client := Client{nil}
	expected := errors.New("client is nil")
	_, actual := client.VerifyProduct("package", "productID", "purchaseToken")

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}
}

func TestCancelSubscription(t *testing.T) {
	t.Parallel()
	// Exception scenario
	client := Client{nil}
	expected := errors.New("client is nil")
	actual := client.CancelSubscription("package", "productID", "purchaseToken")

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}

	client, _ = New(jsonKey)
	expectedStr := "googleapi: Error 404: No application was found for the given package name., applicationNotFound"
	actual = client.CancelSubscription("package", "productID", "purchaseToken")

	if actual.Error() != expectedStr {
		t.Errorf("got %v\nwant %v", actual, expectedStr)
	}

	// TODO Normal scenario
}

func TestRefundSubscription(t *testing.T) {
	t.Parallel()
	// Exception scenario
	client := Client{nil}
	expected := errors.New("client is nil")
	actual := client.RefundSubscription("package", "productID", "purchaseToken")

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}

	client, _ = New(jsonKey)
	expectedStr := "googleapi: Error 404: No application was found for the given package name., applicationNotFound"
	actual = client.RefundSubscription("package", "productID", "purchaseToken")

	if actual.Error() != expectedStr {
		t.Errorf("got %v\nwant %v", actual, expectedStr)
	}

	// TODO Normal scenario
}

func TestRevokeSubscription(t *testing.T) {
	t.Parallel()
	// Exception scenario
	client := Client{nil}
	expected := errors.New("client is nil")
	actual := client.RevokeSubscription("package", "productID", "purchaseToken")

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("got %v\nwant %v", actual, expected)
	}

	client, _ = New(jsonKey)
	expectedStr := "googleapi: Error 404: No application was found for the given package name., applicationNotFound"
	actual = client.RevokeSubscription("package", "productID", "purchaseToken")

	if actual.Error() != expectedStr {
		t.Errorf("got %v\nwant %v", actual, expectedStr)
	}

	// TODO Normal scenario
}

func TestVerifySignature(t *testing.T) {
	t.Parallel()
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
