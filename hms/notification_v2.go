package hms

// SubscriptionNotificationV2 Request parameters when a developer server is called by HMS API.
//
// https://developer.huawei.com/consumer/en/doc/HMSCore-References/api-notifications-about-subscription-events-v2-0000001385268541
type SubscriptionNotificationV2 struct {
	// Notification service version, which is set to v2.
	Version string `json:"version,omitempty"`

	//Notification type. The value can be:
	//ORDER: order
	//SUBSCRIPTION: subscription
	EventType string `json:"eventType,omitempty"`

	//Timestamp of the time when a notification is sent (in UTC), which is the number of milliseconds from 00:00:00 on January 1, 1970 to the time when the notification is sent.
	NotifyTime int64 `json:"notifyTime,omitempty"`

	//App ID.
	ApplicationID string `json:"applicationId,omitempty"`
	// Content of an order notification, which is returned when eventType is ORDER.
	OrderNotification OrderNotification `json:"orderNotification,omitempty"`

	//Content of a subscription notification, which is returned when eventType is SUBSCRIPTION.
	SubNotification SubNotification `json:"subNotification,omitempty"`
}

// OrderNotification JSON content when unmarshal NotificationRequest.OrderNotification
// https://developer.huawei.com/consumer/en/doc/HMSCore-References/api-notifications-about-subscription-events-v2-0000001385268541
type OrderNotification struct {
	//Notification service version, which is set to v2.
	Version string `json:"version,omitempty"`

	// Notification type. The value can be:
	//1: successful payment
	//2: successful refund
	NotificationType int64 `json:"notificationType"`

	// Subscription token, which matches a unique subscription ID.
	PurchaseToken string `json:"purchaseToken"`

	// Product ID.
	ProductID string `json:"productId"`
}

// SubNotification JSON content when unmarshal NotificationRequest.SubNotification
// https://developer.huawei.com/consumer/en/doc/HMSCore-References/api-notifications-about-subscription-events-v2-0000001385268541
type SubNotification struct {
	//Notification service version, which is set to v2.
	Version string `json:"version,omitempty"`

	//Notification message, in JSON format. For details, please refer to statusUpdateNotification.
	StatusUpdateNotification string `json:"statusUpdateNotification"`

	// Signature string of the statusUpdateNotification field. Find the signature algorithm from the value of signatureAlgorithm.
	//After your server receives the signature string, you need to use the IAP public key to verify the signature of statusUpdateNotification (in JSON format). For details, please refer to Verifying the Signature in the Returned Result.
	//For details about how to obtain the public key, please refer to Querying IAP Information.
	NotificationSignature string `json:"notificationSignature"`

	//Signature algorithm.
	SignatureAlgorithm string `json:"signatureAlgorithm"`
}
