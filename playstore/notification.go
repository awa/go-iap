package playstore

// DeveloperNotification is sent by a Pub/Sub topic.
// Detailed description is following.
// https://developer.android.com/google/play/billing/realtime_developer_notifications.html#json_specification
type DeveloperNotification struct {
	Version                  string                   `json:"version"`
	PackageName              string                   `json:"packageName"`
	EventTimeMillis          string                   `json:"eventTimeMillis"`
	SubscriptionNotification SubscriptionNotification `json:"subscriptionNotification,omitempty"`
	TestNotification         SubscriptionNotification `json:"testNotification,omitempty"`
}

// SubscriptionNotification has subscription status as notificationType, toke and subscription id
// to confirm status by calling Google Android Publisher API.
type SubscriptionNotification struct {
	Version          string           `json:"version"`
	NotificationType NotificationType `json:"notificationType,omitempty"`
	PurchaseToken    string           `json:"purchaseToken,omitempty"`
	SubscriptionID   string           `json:"subscriptionId,omitempty"`
}

type NotificationType int

const (
	NotificationTypeRecovered NotificationType = iota + 1
	NotificationTypeRenewed
	NotificationTypeCanceled
	NotificationTypePurchased
	NotificationTypeAccountHold
	NotificationTypeGracePeriod
	NotificationTypeReactivated
)
