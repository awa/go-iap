package amazon

// NotificationType is type
// https://developer.amazon.com/docs/in-app-purchasing/rtn-example.html
type NotificationType string

const (
	NotificationTypeSubscription = "SUBSCRIPTION_PURCHASED"
	NotificationTypeConsumable   = "CONSUMABLE_PURCHASED"
	NotificationTypeEntitlement  = "ENTITLEMENT_PURCHASED"
)

// Notification is struct for amazon notification
type Notification struct {
	Type             string `json:"Type"`
	MessageId        string `json:"MessageId"`
	TopicArn         string `json:"TopicArn"`
	Message          string `json:"Message"`
	Timestamp        string `json:"Timestamp"`
	SignatureVersion string `json:"SignatureVersion"`
	Signature        string `json:"Signature"`
	SigningCertURL   string `json:"SigningCertURL"`
	UnsubscribeURL   string `json:"UnsubscribeURL"`
}

// NotificationMessage is struct for Message field of Notification
type NotificationMessage struct {
	AppPackageName         string           `json:"appPackageName"`
	NotificationType       NotificationType `json:"notificationType"`
	AppUserId              string           `json:"appUserId"`
	ReceiptId              string           `json:"receiptId"`
	RelatedReceipts        struct{}         `json:"relatedReceipts"`
	Timestamp              int64            `json:"timestamp"`
	BetaProductTransaction bool             `json:"betaProductTransaction"`
}
