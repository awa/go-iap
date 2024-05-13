package playstore

// https://developer.android.com/google/play/billing/rtdn-reference#sub
type SubscriptionNotificationType int

const (
	SubscriptionNotificationTypeRecovered SubscriptionNotificationType = iota + 1
	SubscriptionNotificationTypeRenewed
	SubscriptionNotificationTypeCanceled
	SubscriptionNotificationTypePurchased
	SubscriptionNotificationTypeAccountHold
	SubscriptionNotificationTypeGracePeriod
	SubscriptionNotificationTypeRestarted
	SubscriptionNotificationTypePriceChangeConfirmed
	SubscriptionNotificationTypeDeferred
	SubscriptionNotificationTypePaused
	SubscriptionNotificationTypePauseScheduleChanged
	SubscriptionNotificationTypeRevoked
	SubscriptionNotificationTypeExpired
)

// https://developer.android.com/google/play/billing/rtdn-reference#one-time
type OneTimeProductNotificationType int

const (
	OneTimeProductNotificationTypePurchased OneTimeProductNotificationType = iota + 1
	OneTimeProductNotificationTypeCanceled
)

// https://developer.android.com/google/play/billing/rtdn-reference#voided-purchase
type VoidedPurchaseProductType int

const (
	VoidedPurchaseProductTypeSubscription = iota + 1
	VoidedPurchaseProductTypeOneTime
)

type VoidedPurchaseRefundType int

const (
	VoidedPurchaseRefundTypeFullRefund VoidedPurchaseRefundType = iota + 1
	VoidedPurchaseRefundTypePartialRefund
)

// DeveloperNotification is sent by a Pub/Sub topic.
// Detailed description is following.
// https://developer.android.com/google/play/billing/rtdn-reference#json_specification
// Depreacated: use DeveloperNotificationV2 instead.
type DeveloperNotification struct {
	Version                    string                     `json:"version"`
	PackageName                string                     `json:"packageName"`
	EventTimeMillis            string                     `json:"eventTimeMillis"`
	SubscriptionNotification   SubscriptionNotification   `json:"subscriptionNotification,omitempty"`
	OneTimeProductNotification OneTimeProductNotification `json:"oneTimeProductNotification,omitempty"`
	VoidedPurchaseNotification VoidedPurchaseNotification `json:"voidedPurchaseNotification,omitempty"`
	TestNotification           TestNotification           `json:"testNotification,omitempty"`
}

// DeveloperNotificationV2 is sent by a Pub/Sub topic.
// Detailed description is following.
// https://developer.android.com/google/play/billing/rtdn-reference#json_specification
type DeveloperNotificationV2 struct {
	Version                    string                      `json:"version"`
	PackageName                string                      `json:"packageName"`
	EventTimeMillis            string                      `json:"eventTimeMillis"`
	SubscriptionNotification   *SubscriptionNotification   `json:"subscriptionNotification,omitempty"`
	OneTimeProductNotification *OneTimeProductNotification `json:"oneTimeProductNotification,omitempty"`
	VoidedPurchaseNotification *VoidedPurchaseNotification `json:"voidedPurchaseNotification,omitempty"`
	TestNotification           *TestNotification           `json:"testNotification,omitempty"`
}

// SubscriptionNotification has subscription status as notificationType, token and subscription id
// to confirm status by calling Google Android Publisher API.
type SubscriptionNotification struct {
	Version          string                       `json:"version"`
	NotificationType SubscriptionNotificationType `json:"notificationType,omitempty"`
	PurchaseToken    string                       `json:"purchaseToken,omitempty"`
	SubscriptionID   string                       `json:"subscriptionId,omitempty"`
}

// OneTimeProductNotification has one-time product status as notificationType, token and sku (product id)
// to confirm status by calling Google Android Publisher API.
type OneTimeProductNotification struct {
	Version          string                         `json:"version"`
	NotificationType OneTimeProductNotificationType `json:"notificationType,omitempty"`
	PurchaseToken    string                         `json:"purchaseToken,omitempty"`
	SKU              string                         `json:"sku,omitempty"`
}

// VoidedPurchaseNotification has token, order and product type to locate the right purchase and order.
// To learn how to get additional information about the voided purchase, check out the Google Play Voided Purchases API,
// which is a pull model that provides additional data for voided purchases between a given timestamp.
// https://developer.android.com/google/play/billing/rtdn-reference#voided-purchase
type VoidedPurchaseNotification struct {
	PurchaseToken string                    `json:"purchaseToken"`
	OrderID       string                    `json:"orderId"`
	ProductType   VoidedPurchaseProductType `json:"productType"`
	RefundType    VoidedPurchaseRefundType  `json:"refundType"`
}

// TestNotification is the test publish that are sent only through the Google Play Developer Console
type TestNotification struct {
	Version string `json:"version"`
}
