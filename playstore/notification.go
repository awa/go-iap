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

// DeveloperNotification is sent by a Pub/Sub topic.
// Detailed description is following.
// https://developer.android.com/google/play/billing/rtdn-reference#json_specification
type DeveloperNotification struct {
	Version                    string                     `json:"version"`
	PackageName                string                     `json:"packageName"`
	EventTimeMillis            int64                      `json:"eventTimeMillis"`
	SubscriptionNotification   SubscriptionNotification   `json:"subscriptionNotification,omitempty"`
	OneTimeProductNotification OneTimeProductNotification `json:"oneTimeProductNotification,omitempty"`
	TestNotification           TestNotification           `json:"testNotification,omitempty"`
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

// TestNotification is the test publish that are sent only through the Google Play Developer Console
type TestNotification struct {
	Version string `json:"version"`
}
