package hms

import (
	"bytes"
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// VerifySignature validate inapp order or subscription data signature. Returns nil if pass.
//
// Document: https://developer.huawei.com/consumer/en/doc/development/HMSCore-Guides-V5/verifying-signature-returned-result-0000001050033088-V5
// Source code originated from https://github.com/HMS-Core/hms-iap-serverdemo/blob/92241f97fed1b68ddeb7cb37ea4ca6e6d33d2a87/demo/demo.go#L60
func VerifySignature(base64EncodedPublicKey string, data string, signature string) (err error) {
	publicKeyByte, err := base64.StdEncoding.DecodeString(base64EncodedPublicKey)
	if err != nil {
		return err
	}
	pub, err := x509.ParsePKIXPublicKey(publicKeyByte)
	if err != nil {
		return err
	}
	hashed := sha256.Sum256([]byte(data))
	signatureByte, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA256, hashed[:], signatureByte)
}

// SubscriptionVerifyResponse JSON response after requested {rootUrl}/sub/applications/v2/purchases/get
type SubscriptionVerifyResponse struct {
	ResponseCode      string `json:"responseCode"`                // Response code, if = "0" means succeed, for others see https://developer.huawei.com/consumer/en/doc/HMSCore-References-V5/server-error-code-0000001050166248-V5
	ResponseMessage   string `json:"responseMessage,omitempty"`   // Response descriptions, especially when error
	InappPurchaseData string `json:"inappPurchaseData,omitempty"` // InappPurchaseData JSON string
}

// VerifySubscription gets subscriptions info with subscriptionId and purchaseToken.
//
// Document: https://developer.huawei.com/consumer/en/doc/development/HMSCore-References-V5/api-subscription-verify-purchase-token-0000001050706080-V5
// Source code originated from https://github.com/HMS-Core/hms-iap-serverdemo/blob/92241f97fed1b68ddeb7cb37ea4ca6e6d33d2a87/demo/subscription.go#L40
func (c *Client) VerifySubscription(ctx context.Context, purchaseToken, subscriptionID string, accountFlag int64) (InAppPurchaseData, error) {
	var iap InAppPurchaseData

	dataString, err := c.GetSubscriptionDataString(ctx, purchaseToken, subscriptionID, accountFlag)
	if err != nil {
		return iap, err
	}

	if err := json.Unmarshal([]byte(dataString), &iap); err != nil {
		return iap, err
	}

	return iap, nil
}

// GetSubscriptionDataString gets subscriptions response data string.
//
// Document: https://developer.huawei.com/consumer/en/doc/development/HMSCore-References-V5/api-subscription-verify-purchase-token-0000001050706080-V5
// Source code originated from https://github.com/HMS-Core/hms-iap-serverdemo/blob/92241f97fed1b68ddeb7cb37ea4ca6e6d33d2a87/demo/subscription.go#L40
func (c *Client) GetSubscriptionDataString(ctx context.Context, purchaseToken, subscriptionID string, accountFlag int64) (string, error) {
	bodyMap := map[string]string{
		"subscriptionId": subscriptionID,
		"purchaseToken":  purchaseToken,
	}
	url := c.getRootSubscriptionURLByFlag(accountFlag) + "/sub/applications/v2/purchases/get"

	bodyBytes, err := c.sendJSONRequest(ctx, url, bodyMap)
	if err != nil {
		// log.Printf("GetSubscriptionDataString(): Encounter error: %s", err)
		return "", err
	}

	var resp SubscriptionVerifyResponse
	if err := json.Unmarshal(bodyBytes, &resp); err != nil {
		return "", err
	}
	if err := c.getResponseErrorByCode(resp.ResponseCode); err != nil {
		return "", err
	}

	return resp.InappPurchaseData, nil
}

// OrderVerifyResponse JSON response from {rootUrl}/applications/purchases/tokens/verify
type OrderVerifyResponse struct {
	ResponseCode      string `json:"responseCode"`                // Response code, if = "0" means succeed, for others see https://developer.huawei.com/consumer/en/doc/HMSCore-References-V5/server-error-code-0000001050166248-V5
	ResponseMessage   string `json:"responseMessage,omitempty"`   // Response descriptions, especially when error
	PurchaseTokenData string `json:"purchaseTokenData,omitempty"` // InappPurchaseData JSON string
	DataSignature     string `json:"dataSignature,omitempty"`     // Signature to verify PurchaseTokenData string
}

// VerifyOrder gets order (single item purchase) info with productId and purchaseToken.
//
// Note that this method does not verify the DataSignature, thus security is relied on HTTPS solely.
//
// Document: https://developer.huawei.com/consumer/en/doc/HMSCore-References-V5/api-order-verify-purchase-token-0000001050746113-V5
// Source code originated from https://github.com/HMS-Core/hms-iap-serverdemo/blob/92241f97fed1b68ddeb7cb37ea4ca6e6d33d2a87/demo/order.go#L41
func (c *Client) VerifyOrder(ctx context.Context, purchaseToken, productID string, accountFlag int64) (InAppPurchaseData, error) {
	var iap InAppPurchaseData

	dataString, _, err := c.GetOrderDataString(ctx, purchaseToken, productID, accountFlag)
	if err != nil {
		return iap, err
	}

	if err := json.Unmarshal([]byte(dataString), &iap); err != nil {
		return iap, err
	}

	return iap, nil
}

// GetOrderDataString gets order (single item purchase) response data as json string and dataSignature
//
// Document: https://developer.huawei.com/consumer/en/doc/HMSCore-References-V5/api-order-verify-purchase-token-0000001050746113-V5
// Source code originated from https://github.com/HMS-Core/hms-iap-serverdemo/blob/92241f97fed1b68ddeb7cb37ea4ca6e6d33d2a87/demo/order.go#L41
func (c *Client) GetOrderDataString(ctx context.Context, purchaseToken, productID string, accountFlag int64) (purchaseTokenData, dataSignature string, err error) {
	bodyMap := map[string]string{
		"purchaseToken": purchaseToken,
		"productId":     productID,
	}
	url := c.getRootOrderURLByFlag(accountFlag) + "/applications/purchases/tokens/verify"

	bodyBytes, err := c.sendJSONRequest(ctx, url, bodyMap)
	if err != nil {
		// log.Printf("GetOrderDataString(): Encounter error: %s", err)
		return "", "", err
	}

	var resp OrderVerifyResponse
	if err := json.Unmarshal(bodyBytes, &resp); err != nil {
		return "", "", err
	}
	if err := c.getResponseErrorByCode(resp.ResponseCode); err != nil {
		return "", "", err
	}

	return resp.PurchaseTokenData, resp.DataSignature, nil
}

// Helper function to send http json request and get response bodyBytes.
//
// Source code originated from https://github.com/HMS-Core/hms-iap-serverdemo/blob/92241f97fed1b68ddeb7cb37ea4ca6e6d33d2a87/demo/demo.go#L33
func (c *Client) sendJSONRequest(ctx context.Context, url string, bodyMap map[string]string) (bodyBytes []byte, err error) {
	bodyString, err := json.Marshal(bodyMap)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyString))
	if err != nil {
		return
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	atHeader, err := c.GetApplicationAccessTokenHeader()
	if err == nil {
		req.Header.Set("Authorization", atHeader)
	} else {
		return
	}

	resp, err := c.httpCli.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}

// GetCanceledOrRefundedPurchases gets all revoked purchases in CanceledPurchaseList{}.
// This method allow fetch over 1000 results regardles the cap implied by HMS API. Though you should still limit maxRows to a certain number to increate preformance.
//
// In case of an error, this method might return some fetch results if maxRows greater than 1000 or equals 0.
//
// Source code originated from https://github.com/HMS-Core/hms-iap-serverdemo/blob/92241f97fed1b68ddeb7cb37ea4ca6e6d33d2a87/demo/order.go#L52
// Document: https://developer.huawei.com/consumer/en/doc/HMSCore-References-V5/api-cancel-or-refund-record-0000001050746117-V5
func (c *Client) GetCanceledOrRefundedPurchases(
	// context of request
	ctx context.Context,

	// start time timestamp in milliseconds, if =0, will default to 1 month ago.
	startAt int64,

	// end time timestamp in milliseconds, if =0, will default to now.
	endAt int64,

	// rows to return. default to 1000 if maxRows>1000 or equals to 0.
	maxRows int,

	// Token returned in the last query to query the data on the next page.
	continuationToken string,

	// Query type. Ignore this parameter when continuationToken is passed. The options are as follows:
	//    0: Queries purchase information about consumables and non-consumables. This is the default value.
	//    1: Queries all purchase information about consumables, non-consumables, and subscriptions.
	productType int64,

	// Account flag to determine which API URL to use.
	accountFlag int64,
) (canceledPurchases []CanceledPurchase, newContinuationToken string, responseCode string, responseMessage string, err error) {
	// default values
	if maxRows > 1000 || maxRows < 1 {
		maxRows = 1000
	}

	switch endAt {
	case 0:
		endAt = time.Now().UnixNano() / 1000000
	case startAt:
		endAt++
	}

	bodyMap := map[string]string{
		"startAt":           fmt.Sprintf("%v", startAt),
		"endAt":             fmt.Sprintf("%v", endAt),
		"maxRows":           fmt.Sprintf("%v", maxRows),
		"continuationToken": continuationToken,
		"type":              fmt.Sprintf("%v", productType),
	}

	url := c.getRootOrderURLByFlag(accountFlag) + "/applications/v2/purchases/cancelledList"
	bodyBytes, err := c.sendJSONRequest(ctx, url, bodyMap)
	if err != nil {
		// log.Printf("GetCanceledOrRefundedPurchases(): Encounter error: %s", err)
	}

	var cpl CanceledPurchaseList // temporary variable to store api query result
	err = json.Unmarshal(bodyBytes, &cpl)
	if err != nil {
		return canceledPurchases, continuationToken, cpl.ResponseCode, cpl.ResponseMessage, err
	}
	if cpl.ResponseCode != "0" {
		return canceledPurchases, continuationToken, cpl.ResponseCode, cpl.ResponseMessage, c.getResponseErrorByCode(cpl.ResponseCode)
	}

	err = json.Unmarshal([]byte(cpl.CancelledPurchaseList), &canceledPurchases)
	if err != nil {
		return canceledPurchases, continuationToken, cpl.ResponseCode, cpl.ResponseMessage, err
	}

	return canceledPurchases, cpl.ContinuationToken, cpl.ResponseCode, cpl.ResponseMessage, nil
}
