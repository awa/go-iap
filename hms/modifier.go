package hms

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

// CancelSubscriptionRenewal Cancel a aubscription from auto-renew when expired.
// Note that this does not cancel the current subscription.
// If you want to revoke a subscription, use Client.RevokeSubscription() instead.
// Source code originated from https://github.com/HMS-Core/hms-iap-serverdemo/blob/92241f97fed1b68ddeb7cb37ea4ca6e6d33d2a87/demo/subscription.go#L54
// Document: https://developer.huawei.com/consumer/en/doc/HMSCore-References-V5/api-cancel-subscription-0000001050746115-V5
func (c *Client) CancelSubscriptionRenewal(ctx context.Context, purchaseToken, subscriptionID string, accountFlag int64) (success bool, responseMessage string, err error) {
	bodyMap := map[string]string{
		"subscriptionId": subscriptionID,
		"purchaseToken":  purchaseToken,
	}
	var resp ModifySubscriptionResponse
	success, resp, err = c.modifySubscriptionQuery(ctx, bodyMap, accountFlag, "/sub/applications/v2/purchases/stop")
	responseMessage = resp.ResponseMessage
	return
}

// ExtendSubscription extend the current subscription expiration date without chanrging the customer.
// Source code originated from https://github.com/HMS-Core/hms-iap-serverdemo/blob/92241f97fed1b68ddeb7cb37ea4ca6e6d33d2a87/demo/subscription.go#L68
// Document: https://developer.huawei.com/consumer/en/doc/HMSCore-References-V5/api-refund-subscription-fee-0000001050986131-V5
func (c *Client) ExtendSubscription(ctx context.Context, purchaseToken, subscriptionID string, currentExpirationTime, desiredExpirationTime int64, accountFlag int64) (success bool, responseMessage string, newExpirationTime int64, err error) {
	bodyMap := map[string]string{
		"subscriptionId":        subscriptionID,
		"purchaseToken":         purchaseToken,
		"currentExpirationTime": fmt.Sprintf("%v", currentExpirationTime),
		"desiredExpirationTime": fmt.Sprintf("%v", desiredExpirationTime),
	}
	var resp ModifySubscriptionResponse
	success, resp, err = c.modifySubscriptionQuery(ctx, bodyMap, accountFlag, "/sub/applications/v2/purchases/delay")
	responseMessage = resp.ResponseMessage
	newExpirationTime = resp.NewExpirationTime
	return
}

// RefundSubscription refund a subscription payment.
// Note that this does not cancel the current subscription.
// If you want to revoke a subscription, use Client.RevokeSubscription() instead.
// Source code originated from https://github.com/HMS-Core/hms-iap-serverdemo/blob/92241f97fed1b68ddeb7cb37ea4ca6e6d33d2a87/demo/subscription.go#L84
// Document: https://developer.huawei.com/consumer/en/doc/HMSCore-References-V5/api-refund-subscription-fee-0000001050986131-V5
func (c *Client) RefundSubscription(ctx context.Context, purchaseToken, subscriptionID string, accountFlag int64) (success bool, responseMessage string, err error) {
	bodyMap := map[string]string{
		"subscriptionId": subscriptionID,
		"purchaseToken":  purchaseToken,
	}
	var resp ModifySubscriptionResponse
	success, resp, err = c.modifySubscriptionQuery(ctx, bodyMap, accountFlag, "/sub/applications/v2/purchases/returnFee")
	responseMessage = resp.ResponseMessage
	return
}

// RevokeSubscription will revoke and issue a refund on a subscription immediately.
// Source code originated from https://github.com/HMS-Core/hms-iap-serverdemo/blob/92241f97fed1b68ddeb7cb37ea4ca6e6d33d2a87/demo/subscription.go#L99
// Document: https://developer.huawei.com/consumer/en/doc/HMSCore-References-V5/api-unsubscribe-0000001051066056-V5
func (c *Client) RevokeSubscription(ctx context.Context, purchaseToken, subscriptionID string, accountFlag int64) (success bool, responseMessage string, err error) {
	bodyMap := map[string]string{
		"subscriptionId": subscriptionID,
		"purchaseToken":  purchaseToken,
	}
	var resp ModifySubscriptionResponse
	success, resp, err = c.modifySubscriptionQuery(ctx, bodyMap, accountFlag, "/sub/applications/v2/purchases/withdrawal")
	responseMessage = resp.ResponseMessage
	return
}

// ModifySubscriptionResponse JSON response from {rootUrl}/sub/applications/v2/purchases/stop|delay|returnFee|withdrawal
type ModifySubscriptionResponse struct {
	ResponseCode      string `json:"responseCode"`
	ResponseMessage   string `json:"responseMessage;omitempty"`
	NewExpirationTime int64  `json:"newExpirationTime;omitempty"`
}

// public method to query {rootUrl}/sub/applications/v2/purchases/stop|delay|returnFee|withdrawal
func (c *Client) modifySubscriptionQuery(ctx context.Context, requestBodyMap map[string]string, accountFlag int64, uri string) (success bool, response ModifySubscriptionResponse, err error) {
	url := c.getRootSubscriptionURLByFlag(accountFlag) + uri

	bodyBytes, err := c.sendJSONRequest(ctx, url, requestBodyMap)
	if err != nil {
		return false, response, err
	}

	// debug
	log.Println("url:", url)
	log.Println("request:", requestBodyMap)
	log.Printf("%s", bodyBytes)

	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return false, response, err
	}

	switch response.ResponseCode {
	case "0":
		return true, response, nil
	default:
		return false, response, c.getResponseErrorByCode(response.ResponseCode)
	}
}
