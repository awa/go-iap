package api

import "github.com/awa/go-iap/appstore"

// OrderLookupResponse https://developer.apple.com/documentation/appstoreserverapi/orderlookupresponse
type OrderLookupResponse struct {
	Status             int      `json:"status"`
	SignedTransactions []string `json:"signedTransactions"`
}

type Environment string

// Environment https://developer.apple.com/documentation/appstoreserverapi/environment
const (
	Sandbox    Environment = "Sandbox"
	Production Environment = "Production"
)

// HistoryResponse https://developer.apple.com/documentation/appstoreserverapi/historyresponse
type HistoryResponse struct {
	AppAppleId         int64       `json:"appAppleId"`
	BundleId           string      `json:"bundleId"`
	Environment        Environment `json:"environment"`
	HasMore            bool        `json:"hasMore"`
	Revision           string      `json:"revision"`
	SignedTransactions []string    `json:"signedTransactions"`
}

// TransactionInfoResponse https://developer.apple.com/documentation/appstoreserverapi/transactioninforesponse
type TransactionInfoResponse struct {
	SignedTransactionInfo string `json:"signedTransactionInfo"`
}

// RefundLookupResponse same as the RefundHistoryResponse https://developer.apple.com/documentation/appstoreserverapi/refundhistoryresponse
type RefundLookupResponse struct {
	HasMore            bool     `json:"hasMore"`
	Revision           string   `json:"revision"`
	SignedTransactions []string `json:"signedTransactions"`
}

// StatusResponse https://developer.apple.com/documentation/appstoreserverapi/get_all_subscription_statuses
type StatusResponse struct {
	Environment Environment                       `json:"environment"`
	AppAppleId  int64                             `json:"appAppleId"`
	BundleId    string                            `json:"bundleId"`
	Data        []SubscriptionGroupIdentifierItem `json:"data"`
}

type SubscriptionGroupIdentifierItem struct {
	SubscriptionGroupIdentifier string                 `json:"subscriptionGroupIdentifier"`
	LastTransactions            []LastTransactionsItem `json:"lastTransactions"`
}

type LastTransactionsItem struct {
	OriginalTransactionId string `json:"originalTransactionId"`
	Status                int32  `json:"status"`
	SignedRenewalInfo     string `json:"signedRenewalInfo"`
	SignedTransactionInfo string `json:"signedTransactionInfo"`
}

// MassExtendRenewalDateRequest https://developer.apple.com/documentation/appstoreserverapi/massextendrenewaldaterequest
type MassExtendRenewalDateRequest struct {
	RequestIdentifier      string   `json:"requestIdentifier"`
	ExtendByDays           int32    `json:"extendByDays"`
	ExtendReasonCode       int32    `json:"extendReasonCode"`
	ProductId              string   `json:"productId"`
	StorefrontCountryCodes []string `json:"storefrontCountryCodes"`
}

// ConsumptionRequestBody https://developer.apple.com/documentation/appstoreserverapi/consumptionrequest
type ConsumptionRequestBody struct {
	AccountTenure            int32  `json:"accountTenure"`
	AppAccountToken          string `json:"appAccountToken"`
	ConsumptionStatus        int32  `json:"consumptionStatus"`
	CustomerConsented        bool   `json:"customerConsented"`
	DeliveryStatus           int32  `json:"deliveryStatus"`
	LifetimeDollarsPurchased int32  `json:"lifetimeDollarsPurchased"`
	LifetimeDollarsRefunded  int32  `json:"lifetimeDollarsRefunded"`
	Platform                 int32  `json:"platform"`
	PlayTime                 int32  `json:"playTime"`
	SampleContentProvided    bool   `json:"sampleContentProvided"`
	UserStatus               int32  `json:"userStatus"`
}

// JWSRenewalInfoDecodedPayload https://developer.apple.com/documentation/appstoreserverapi/jwsrenewalinfodecodedpayload
type JWSRenewalInfoDecodedPayload struct {
	AutoRenewProductId          string      `json:"autoRenewProductId"`
	AutoRenewStatus             int32       `json:"autoRenewStatus"`
	Environment                 Environment `json:"environment"`
	ExpirationIntent            int32       `json:"expirationIntent"`
	GracePeriodExpiresDate      int64       `json:"gracePeriodExpiresDate"`
	IsInBillingRetryPeriod      bool        `json:"isInBillingRetryPeriod"`
	OfferIdentifier             string      `json:"offerIdentifier"`
	OfferType                   string      `json:"offerType"`
	OriginalTransactionId       string      `json:"originalTransactionId"`
	PriceIncreaseStatus         int32       `json:"priceIncreaseStatus"`
	ProductId                   string      `json:"productId"`
	RecentSubscriptionStartDate int64       `json:"recentSubscriptionStartDate"`
	RenewalDate                 int64       `json:"renewalDate"`
	SignedDate                  int64       `json:"signedDate"`
}

func (J JWSRenewalInfoDecodedPayload) Valid() error {
	return nil
}

// JWSDecodedHeader https://developer.apple.com/documentation/appstoreserverapi/jwsdecodedheader
type JWSDecodedHeader struct {
	Alg string   `json:"alg,omitempty"`
	Kid string   `json:"kid,omitempty"`
	X5C []string `json:"x5c,omitempty"`
}

// TransactionReason indicates the cause of a purchase transaction,
// https://developer.apple.com/documentation/appstoreservernotifications/transactionreason
type TransactionReason string

const (
	TransactionReasonPurchase = "PURCHASE"
	TransactionReasonRenewal  = "RENEWAL"
)

// IAPType https://developer.apple.com/documentation/appstoreserverapi/type
type IAPType string

const (
	AutoRenewable IAPType = "Auto-Renewable Subscription"
	NonConsumable IAPType = "Non-Consumable"
	Consumable    IAPType = "Consumable"
	NonRenewable  IAPType = "Non-Renewing Subscription"
)

// JWSTransaction https://developer.apple.com/documentation/appstoreserverapi/jwstransaction
type JWSTransaction struct {
	TransactionID               string            `json:"transactionId,omitempty"`
	OriginalTransactionId       string            `json:"originalTransactionId,omitempty"`
	WebOrderLineItemId          string            `json:"webOrderLineItemId,omitempty"`
	BundleID                    string            `json:"bundleId,omitempty"`
	ProductID                   string            `json:"productId,omitempty"`
	SubscriptionGroupIdentifier string            `json:"subscriptionGroupIdentifier,omitempty"`
	PurchaseDate                int64             `json:"purchaseDate,omitempty"`
	OriginalPurchaseDate        int64             `json:"originalPurchaseDate,omitempty"`
	ExpiresDate                 int64             `json:"expiresDate,omitempty"`
	Quantity                    int32             `json:"quantity,omitempty"`
	Type                        IAPType           `json:"type,omitempty"`
	AppAccountToken             string            `json:"appAccountToken,omitempty"`
	InAppOwnershipType          string            `json:"inAppOwnershipType,omitempty"`
	SignedDate                  int64             `json:"signedDate,omitempty"`
	OfferType                   int32             `json:"offerType,omitempty"`
	OfferIdentifier             string            `json:"offerIdentifier,omitempty"`
	RevocationDate              int64             `json:"revocationDate,omitempty"`
	RevocationReason            int32             `json:"revocationReason,omitempty"`
	IsUpgraded                  bool              `json:"isUpgraded,omitempty"`
	Storefront                  string            `json:"storefront,omitempty"`
	StorefrontId                string            `json:"storefrontId,omitempty"`
	TransactionReason           TransactionReason `json:"transactionReason,omitempty"`
	Environment                 Environment       `json:"environment,omitempty"`
}

func (J JWSTransaction) Valid() error {
	return nil
}

// https://developer.apple.com/documentation/appstoreserverapi/extendreasoncode
type ExtendReasonCode int32

const (
	UndeclaredExtendReasonCode = iota
	CustomerSatisfaction
	OtherReasons
	ServiceIssueOrOutage
)

// ExtendRenewalDateRequest https://developer.apple.com/documentation/appstoreserverapi/extendrenewaldaterequest
type ExtendRenewalDateRequest struct {
	ExtendByDays      int32            `json:"extendByDays"`
	ExtendReasonCode  ExtendReasonCode `json:"extendReasonCode"`
	RequestIdentifier string           `json:"requestIdentifier"`
}

// MassExtendRenewalDateStatusResponse https://developer.apple.com/documentation/appstoreserverapi/massextendrenewaldatestatusresponse
type MassExtendRenewalDateStatusResponse struct {
	RequestIdentifier string `json:"requestIdentifier"`
	Complete          bool   `json:"complete"`
	CompleteDate      int64  `json:"completeDate,omitempty"`
	FailedCount       int64  `json:"failedCount,omitempty"`
	SucceededCount    int64  `json:"succeededCount,omitempty"`
}

// NotificationHistoryRequest https://developer.apple.com/documentation/appstoreserverapi/notificationhistoryrequest
type NotificationHistoryRequest struct {
	StartDate           int64                       `json:"startDate"`
	EndDate             int64                       `json:"endDate"`
	NotificationType    appstore.NotificationTypeV2 `json:"notificationType,omitempty"`
	NotificationSubtype appstore.SubtypeV2          `json:"notificationSubtype,omitempty"`
	OnlyFailures        bool                        `json:"onlyFailures"`
	TransactionId       string                      `json:"transactionId,omitempty"`
	// Use transactionId instead.
	// Deprecated.
	OriginalTransactionId string `json:"originalTransactionId,omitempty"`
}

// NotificationHistoryResponses https://developer.apple.com/documentation/appstoreserverapi/notificationhistoryresponse
type NotificationHistoryResponses struct {
	HasMore             bool                              `json:"hasMore"`
	PaginationToken     string                            `json:"paginationToken"`
	NotificationHistory []NotificationHistoryResponseItem `json:"notificationHistory"`
}

// NotificationHistoryResponseItem https://developer.apple.com/documentation/appstoreserverapi/notificationhistoryresponseitem
type NotificationHistoryResponseItem struct {
	SignedPayload          string                 `json:"signedPayload"`
	FirstSendAttemptResult FirstSendAttemptResult `json:"firstSendAttemptResult"`
	SendAttempts           []SendAttemptItem      `json:"sendAttempts"`
}

// SendAttemptItem https://developer.apple.com/documentation/appstoreserverapi/sendattemptitem
type SendAttemptItem struct {
	AttemptDate       int64                  `json:"attemptDate"`
	SendAttemptResult FirstSendAttemptResult `json:"sendAttemptResult"`
}

// https://developer.apple.com/documentation/appstoreserverapi/firstsendattemptresult
type FirstSendAttemptResult string

const (
	FirstSendAttemptResultSuccess                     FirstSendAttemptResult = "SUCCESS"
	FirstSendAttemptResultCircularRedirect            FirstSendAttemptResult = "CIRCULAR_REDIRECT"
	FirstSendAttemptResultInvalidResponse             FirstSendAttemptResult = "INVALID_RESPONSE"
	FirstSendAttemptResultNoResponse                  FirstSendAttemptResult = "NO_RESPONSE"
	FirstSendAttemptResultOther                       FirstSendAttemptResult = "OTHER"
	FirstSendAttemptResultPrematureClose              FirstSendAttemptResult = "PREMATURE_CLOSE"
	FirstSendAttemptResultSocketIssue                 FirstSendAttemptResult = "SOCKET_ISSUE"
	FirstSendAttemptResultTimedOut                    FirstSendAttemptResult = "TIMED_OUT"
	FirstSendAttemptResultTlsIssue                    FirstSendAttemptResult = "TLS_ISSUE"
	FirstSendAttemptResultUnsupportedCharset          FirstSendAttemptResult = "UNSUPPORTED_CHARSET"
	FirstSendAttemptResultUnsupportedHTTPRESPONSECODE FirstSendAttemptResult = "UNSUCCESSFUL_HTTP_RESPONSE_CODE"
)

// SendTestNotificationResponse https://developer.apple.com/documentation/appstoreserverapi/sendtestnotificationresponse
type SendTestNotificationResponse struct {
	TestNotificationToken string `json:"testNotificationToken"`
}
