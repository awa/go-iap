package hms

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// HMS OAuth url
const tokenURL = "https://oauth-login.cloud.huawei.com/oauth2/v3/token"

// AccessToken expires grace period in seconds.
// The actural ExpiredAt will be substracted with this number to avoid boundray problems.
const accessTokenExpiresGracePeriod = 60

// global variable to store API AccessToken.
// All clients within an instance share one AccessToken grantee scalebility and to avoid rate limit.
var applicationAccessTokens = make(map[[16]byte]ApplicationAccessToken)

// lock when writing to applicationAccessTokens map
var applicationAccessTokensLock sync.Mutex

// ApplicationAccessToken model, received from HMS OAuth API
// https://developer.huawei.com/consumer/en/doc/HMSCore-Guides/open-platform-oauth-0000001050123437#EN-US_TOPIC_0000001050123437__section12493191334711
type ApplicationAccessToken struct {
	// App-level access token.
	AccessToken string `json:"access_token"`

	// Remaining validity period of an access token, in seconds.
	ExpiresIn int64 `json:"expires_in"`
	// This value is always Bearer, indicating the type of the returned access token.
	// TokenType    string	`json:"token_type"`

	// Save the timestamp when AccessToken is obtained
	ExpiredAt int64 `json:"-"`

	// Request header string
	HeaderString string `json:"-"`
}

// Client implements VerifySignature, VerifyOrder and VerifySubscription methods
type Client struct {
	clientID            string
	clientSecret        string
	clientIDSecretHash  [16]byte
	httpCli             *http.Client
	orderSiteURL        string // site URL to request order information
	subscriptionSiteURL string // site URL to request subscription information
}

// New returns client with credentials.
// Required client_id and client_secret which could be acquired from the HMS API Console.
// When user accountFlag is not equals to 1, orderSiteURL/subscriptionSiteURL are the site URLs that will be used to connect to HMS IAP API services.
// If orderSiteURL or subscriptionSiteURL are not set, default to AppTouch Germany site.
//
// Please refer https://developer.huawei.com/consumer/en/doc/start/api-console-guide
// and https://developer.huawei.com/consumer/en/doc/HMSCore-References/api-common-statement-0000001050986127 for details.
func New(clientID, clientSecret, orderSiteURL, subscriptionSiteURL string) *Client {
	// Set default order / subscription iap site to AppTouch Germany if it is not provided
	if !strings.HasPrefix(orderSiteURL, "http") {
		orderSiteURL = "https://orders-at-dre.iap.dbankcloud.com"
	}
	if !strings.HasPrefix(subscriptionSiteURL, "http") {
		subscriptionSiteURL = "https://subscr-at-dre.iap.dbankcloud.com"
	}

	// Create http client
	return &Client{
		clientID:           clientID,
		clientSecret:       clientSecret,
		clientIDSecretHash: md5.Sum([]byte(clientID + clientSecret)),
		httpCli: &http.Client{
			Timeout: 10 * time.Second,
		},
		orderSiteURL:        orderSiteURL,
		subscriptionSiteURL: subscriptionSiteURL,
	}
}

// GetApplicationAccessTokenHeader obtain OAuth AccessToken from HMS
//
// Source code originated from https://github.com/HMS-Core/hms-iap-serverdemo/blob/92241f97fed1b68ddeb7cb37ea4ca6e6d33d2a87/demo/atdemo.go#L37
func (c *Client) GetApplicationAccessTokenHeader() (string, error) {
	// To complie with the rate limit (1000/5min as of July 24th, 2020)
	// new AccessTokens are requested only when it is expired.
	// Please refer https://developer.huawei.com/consumer/en/doc/HMSCore-Guides/open-platform-oauth-0000001050123437 for detailes
	if applicationAccessTokens[c.clientIDSecretHash].ExpiredAt > time.Now().Unix() {
		return applicationAccessTokens[c.clientIDSecretHash].HeaderString, nil
	}

	urlValue := url.Values{"grant_type": {"client_credentials"}, "client_secret": {c.clientSecret}, "client_id": {c.clientID}}
	resp, err := c.httpCli.PostForm(tokenURL, urlValue)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var atResponse ApplicationAccessToken
	err = json.Unmarshal(bodyBytes, &atResponse)
	if err != nil {
		return "", err
	}
	if atResponse.AccessToken != "" {
		// update expire time
		atResponse.ExpiredAt = atResponse.ExpiresIn + time.Now().Unix() - accessTokenExpiresGracePeriod
		// parse request header string
		atResponse.HeaderString = fmt.Sprintf(
			"Basic %s",
			base64.StdEncoding.EncodeToString([]byte(
				fmt.Sprintf("APPAT:%s",
					atResponse.AccessToken,
				),
			)),
		)
		// save AccessToken info to global variable
		applicationAccessTokensLock.Lock()
		applicationAccessTokens[c.clientIDSecretHash] = atResponse
		applicationAccessTokensLock.Unlock()
		return atResponse.HeaderString, nil
	}
	return "", errors.New("Get token fail, " + string(bodyBytes))
}

// Returns root order URL by flag, prefixing with "https://"
func (c *Client) getRootOrderURLByFlag(flag int64) string {
	switch flag {
	case 1:
		return "https://orders-drcn.iap.cloud.huawei.com.cn"
	case 2:
		return "https://orders-dre.iap.cloud.huawei.eu"
	case 3:
		return "https://orders-dra.iap.cloud.huawei.asia"
	case 4:
		return "https://orders-drru.iap.cloud.huawei.ru"
	}
	return c.orderSiteURL
}

// Returns root subscription URL by flag, prefixing with "https://"
func (c *Client) getRootSubscriptionURLByFlag(flag int64) string {
	switch flag {
	case 1:
		return "https://subscr-drcn.iap.cloud.huawei.com.cn"
	case 2:
		return "https://subscr-dre.iap.cloud.huawei.eu"
	case 3:
		return "https://subscr-dra.iap.cloud.huawei.asia"
	case 4:
		return "https://subscr-drru.iap.cloud.huawei.ru"
	}
	return c.subscriptionSiteURL
}

// get error based on result code returned from api
func (c *Client) getResponseErrorByCode(code string) error {
	switch code {
	case "0":
		return nil
	case "5":
		return ErrorResponseInvalidParameter
	case "6":
		return ErrorResponseCritical
	case "8":
		return ErrorResponseProductNotBelongToUser
	case "9":
		return ErrorResponseConsumedProduct
	case "11":
		return ErrorResponseAbnormalUserAccount
	default:
		return ErrorResponseUnknown
	}
}

// Errors

// ErrorResponseUnknown error placeholder for undocumented errors
var ErrorResponseUnknown error = errors.New("Unknown error from API response")

// ErrorResponseSignatureVerifyFailed failed to verify dataSignature against the response json string.
// https://developer.huawei.com/consumer/en/doc/HMSCore-Guides/verifying-signature-returned-result-0000001050033088
// var ErrorResponseSignatureVerifyFailed error = errors.New("Failed to verify dataSignature against the response json string")

// ErrorResponseInvalidParameter The parameter passed to the API is invalid.
// This error may also indicate that an agreement is not signed or parameters are not set correctly for the in-app purchase settlement in HUAWEI IAP, or the required permission is not in the list.
//
// Check whether the parameter passed to the API is correctly set. If so, check whether required settings in HUAWEI IAP are correctly configured.
// https://developer.huawei.com/consumer/en/doc/HMSCore-References/server-error-code-0000001050166248
var ErrorResponseInvalidParameter error = errors.New("The parameter passed to the API is invalid")

// ErrorResponseCritical A critical error occurs during API operations.
//
// Rectify the fault based on the error information in the response. If the fault persists, contact Huawei technical support.
// https://developer.huawei.com/consumer/en/doc/HMSCore-References/server-error-code-0000001050166248
var ErrorResponseCritical error = errors.New("A critical error occurs during API operations")

// ErrorResponseProductNotBelongToUser A user failed to consume or confirm a product because the user does not own the product.
//
// https://developer.huawei.com/consumer/en/doc/HMSCore-References/server-error-code-0000001050166248
var ErrorResponseProductNotBelongToUser error = errors.New("A user failed to consume or confirm a product because the user does not own the product")

// ErrorResponseConsumedProduct The product cannot be consumed or confirmed because it has been consumed or confirmed.
//
// https://developer.huawei.com/consumer/en/doc/HMSCore-References/server-error-code-0000001050166248
var ErrorResponseConsumedProduct error = errors.New("The product cannot be consumed or confirmed because it has been consumed or confirmed")

// ErrorResponseAbnormalUserAccount The user account is abnormal, for example, the user has been deregistered.
//
// https://developer.huawei.com/consumer/en/doc/HMSCore-References/server-error-code-0000001050166248
var ErrorResponseAbnormalUserAccount error = errors.New("The user account is abnormal, for example, the user has been deregistered")
