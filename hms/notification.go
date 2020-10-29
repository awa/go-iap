package hms

// SubscriptionNotification Request parameters when a developer server is called by HMS API.
//
// https://developer.huawei.com/consumer/en/doc/HMSCore-References-V5/api-notifications-about-subscription-events-0000001050706084-V5
type SubscriptionNotification struct {
	// Notification message, which is a JSON string. For details, please refer to statusUpdateNotification.
	StatusUpdateNotification string `json:"statusUpdateNotification"`

	// Signature string for the StatusUpdateNotification parameter. The signature algorithm is SHA256withRSA.
	//
	// After your server receives the signature string, you need to use the public payment key to verify the signature of StatusUpdateNotification in JSON format.
	// For details, please refer to https://developer.huawei.com/consumer/en/doc/HMSCore-Guides-V5/verifying-signature-returned-result-0000001050033088-V5
	//
	// For details about how to obtain the public key, please refer to https://developer.huawei.com/consumer/en/doc/HMSCore-Guides-V5/query-payment-info-0000001050166299-V5
	NotifycationSignature string `json:"notifycationSignature"`
}

// StatusUpdateNotification JSON content when unmarshal NotificationRequest.StatusUpdateNotification
// https://developer.huawei.com/consumer/en/doc/HMSCore-References-V5/api-notifications-about-subscription-events-0000001050706084-V5#EN-US_TOPIC_0000001050706084__section18290165220716
type StatusUpdateNotification struct {
	// Environment for sending a notification. Value could be one of either:
	//    "PROD": general production environment
	//    "SandBox": sandbox testing environment
	Environment string `json:"environment"`

	// Notification event type. For details, please refer to const NotificationTypeInitialBuy etc.
	NotificationType int64 `json:"notificationType"`

	// Subscription ID
	SubscriptionID string `json:"subscriptionId"`

	// Timestamp, which is passed only when notificationType is CANCEL(1).
	CancellationDate int64 `json:"cancellationDate,omitempty"`

	// Order ID used for payment during subscription renewal.
	OrderID string `json:"orderId"`

	// PurchaseToken of the latest receipt, which is passed only when notificationType is INITIAL_BUY(0), RENEWAL(2), or INTERACTIVE_RENEWAL(3) and the renewal is successful.
	LatestReceipt string `json:"latestReceipt,omitempty"`

	// Latest receipt, which is a JSON string. This parameter is left empty when notificationType is CANCEL(1).
	// For details about the parameters contained, please refer to https://developer.huawei.com/consumer/en/doc/HMSCore-References-V5/server-data-model-0000001050986133-V5#EN-US_TOPIC_0000001050986133__section264617465219
	LatestReceiptInfo string `json:"latestReceiptInfo,omitempty"`

	// Signature string for the LatestReceiptInfo parameter. The signature algorithm is SHA256withRSA.
	//
	// After your server receives the signature string, you need to use the public payment key to verify the signature of LatestReceiptInfo in JSON format.
	// For details, please refer to https://developer.huawei.com/consumer/en/doc/HMSCore-Guides-V5/verifying-signature-returned-result-0000001050033088-V5
	//
	// For details about how to obtain the public key, please refer to https://developer.huawei.com/consumer/en/doc/HMSCore-Guides-V5/query-payment-info-0000001050166299-V5
	LatestReceiptInfoSignature string `json:"latestReceiptInfoSignature,omitempty"`

	// Token of the latest expired receipt. This parameter has a value only when NotificationType is RENEWAL(2) or INTERACTIVE_RENEWAL(3).
	LatestExpiredReceipt string `json:"latestExpiredReceipt,omitempty"`

	// Latest expired receipt, which is a JSON string. This parameter has a value only when NotificationType is RENEWAL(2) or INTERACTIVE_RENEWAL(3).
	LatestExpiredReceiptInfo string `json:"latestExpiredReceiptInfo,omitempty"`

	// Signature string for the LatestExpiredReceiptInfo parameter. The signature algorithm is SHA256withRSA.
	//
	// After your server receives the signature string, you need to use the public payment key to verify the signature of LatestExpiredReceiptInfo in JSON format.
	// For details, please refer to https://developer.huawei.com/consumer/en/doc/HMSCore-Guides-V5/verifying-signature-returned-result-0000001050033088-V5
	//
	// For details about how to obtain the public key, please refer to https://developer.huawei.com/consumer/en/doc/HMSCore-Guides-V5/query-payment-info-0000001050166299-V5
	LatestExpiredReceiptInfoSignature string `json:"latestExpiredReceiptInfoSignature,omitempty"`

	// Renewal status. Value could be one of either:
	//    1: The subscription renewal is normal.
	//    0: The user has canceled subscription renewal.
	AutoRenewStatus int64 `json:"autoRenewStatus"`

	// Refund order ID. This parameter has a value only when NotificationType is CANCEL(1).
	RefundPayOrderID string `json:"refundPayOrderId,omitempty"`

	// Product ID.
	ProductID string `json:"productId"`

	// App ID.
	ApplicationID string `json:"applicationId,omitempty"`

	// Expiration reason. This parameter has a value only when NotificationType is RENEWAL(2) or INTERACTIVE_RENEWAL(3), and the renewal is successful.
	ExpirationIntent int64 `json:"expirationIntent,omitempty"`
}

// Constants for StatusUpdateNotification.NotificationType
// https://developer.huawei.com/consumer/en/doc/HMSCore-References-V5/api-notifications-about-subscription-events-0000001050706084-V5#EN-US_TOPIC_0000001050706084__section18290165220716
const (
	NotificationTypeInitialBuy           int64 = 0
	NotificationTypeCancel               int64 = 1
	NotificationTypeRenewal              int64 = 2
	NotificationTypeInteractiveRenewal   int64 = 3
	NotificationTypeNewRenewalPref       int64 = 4
	NotificationTypeRenewalStopped       int64 = 5
	NotificationTypeRenewalRestored      int64 = 6
	NotificationTypeRenewalRecurring     int64 = 7
	NotificationTypeInGracePeriod        int64 = 8
	NotificationTypeOnHold               int64 = 9
	NotificationTypePaused               int64 = 10
	NotificationTypePausePlanChanged     int64 = 11
	NotificationTypePriceChangeConfirmed int64 = 12
	NotificationTypeDeferred             int64 = 13
)
