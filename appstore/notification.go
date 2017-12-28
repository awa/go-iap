package appstore

// https://developer.apple.com/library/content/documentation/NetworkingInternet/Conceptual/StoreKitGuide/Chapters/Subscriptions.html#//apple_ref/doc/uid/TP40008267-CH7-SW16
type NotificationType string

const (
	// Initial purchase of the subscription.
	NotificationTypeInitialBuy NotificationType = "INITIAL_BUY"
	// Subscription was canceled by Apple customer support.
	NotificationTypeCancel NotificationType = "CANCEL"
	// Automatic renewal was successful for an expired subscription.
	NotificationTypeRenewal NotificationType = "RENEWAL"
	// Customer renewed a subscription interactively after it lapsed.
	NotificationTypeInteractiveRenewal NotificationType = "INTERACTIVE_RENEWAL"
	// Customer changed the plan that takes affect at the next subscription renewal. Current active plan is not affected.
	NotificationTypeDidChangeRenewalPreference NotificationType = "DID_CHANGE_RENEWAL_PREFERENCE"
)

type NotificationEnv string

const (
	NotificationEnvSandbox    NotificationEnv = "SANDBOX"
	NotificationEnvProduction NotificationEnv = "PROD"
)

type NotificationExpiresDate struct {
	ExpiresDateMS  string `json:"expires_date"`
	ExpiresDateUTC string `json:"expires_date_formatted"`
	ExpiresDatePST string `json:"expires_date_formatted_pst"`
}

type NotificationReceipt struct {
	UniqueIdentifier          string `json:"unique_identifier"`
	AppItemID                 string `json:"app_item_id"`
	Quantity                  string `json:"quantity"`
	VersionExternalIdentifier string `json:"version_external_identifier"`
	UniqueVendorIdentifier    string `json:"unique_vendor_identifier"`
	WebOrderLineItemID        string `json:"web_order_line_item_id"`
	ItemID                    string `json:"item_id"`
	ProductID                 string `json:"product_id"`
	BID                       string `json:"bid"`
	BVRS                      string `json:"bvrs"`
	TransactionID             string `json:"transaction_id"`
	OriginalTransactionID     string `json:"original_transaction_id"`
	IsTrialPeriod             string `json:"is_trial_period"`

	PurchaseDate
	OriginalPurchaseDate
	NotificationExpiresDate
	CancellationDate
}

type SubscriptionNotification struct {
	Environment      NotificationEnv  `json:"environment"`
	NotificationType NotificationType `json:"notification_type"`

	// Not show in raw notify body
	Password              string `json:"password"`
	OriginalTransactionID string `json:"original_transaction_id"`
	AutoRenewAdamID       string `json:"auto_renew_adam_id"`

	// The primary key for identifying a subscription purchase.
	// Posted only if the notification_type is CANCEL.
	WebOrderLineItemID string `json:"web_order_line_item_id"`

	// This is the same as the Subscription Expiration Intent in the receipt.
	// Posted only if notification_type is RENEWAL or INTERACTIVE_RENEWAL.
	ExpirationIntent string `json:"expiration_intent"`

	// Auto renew info
	AutoRenewStatus    string `json:"auto_renew_status"` // false or true
	AutoRenewProductID string `json:"auto_renew_product_id"`

	// Posted if the notification_type is RENEWAL or INTERACTIVE_RENEWAL, and only if the renewal is successful.
	// Posted also if the notification_type is INITIAL_BUY.
	// Not posted for notification_type CANCEL.
	LatestReceipt     string              `json:"latest_receipt"`
	LatestReceiptInfo NotificationReceipt `json:"latest_receipt_info"`

	// Posted only if the notification_type is RENEWAL or CANCEL or if renewal failed and subscription expired.
	LatestExpiredReceipt     string              `json:"latest_expired_receipt"`
	LatestExpiredReceiptInfo NotificationReceipt `json:"latest_expired_receipt_info"`

	// Posted only if the notification_type is CANCEL.
	CancellationDate
}
