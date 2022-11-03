package api

import "github.com/awa/go-iap/appstore"

// OrderLookupResponse https://developer.apple.com/documentation/appstoreserverapi/orderlookupresponse
type OrderLookupResponse struct {
	Status             int      `json:"status"`
	SignedTransactions []string `json:"signedTransactions"`
}

// HistoryResponse https://developer.apple.com/documentation/appstoreserverapi/historyresponse
type HistoryResponse struct {
	AppAppleId         int      `json:"appAppleId"`
	BundleId           string   `json:"bundleId"`
	Environment        string   `json:"environment"`
	HasMore            bool     `json:"hasMore"`
	Revision           string   `json:"revision"`
	SignedTransactions []string `json:"signedTransactions"`
}

// RefundLookupResponse https://developer.apple.com/documentation/appstoreserverapi/refundlookupresponse
type RefundLookupResponse struct {
	HasMore            bool     `json:"hasMore"`
	Revision           string   `json:"revision"`
	SignedTransactions []string `json:"signedTransactions"`
}

// StatusResponse https://developer.apple.com/documentation/appstoreserverapi/get_all_subscription_statuses
type StatusResponse struct {
	Environment string                            `json:"environment"`
	AppAppleId  int                               `json:"appAppleId"`
	BundleId    string                            `json:"bundleId"`
	Data        []SubscriptionGroupIdentifierItem `json:"data"`
}

type SubscriptionGroupIdentifierItem struct {
	SubscriptionGroupIdentifier string                 `json:"subscriptionGroupIdentifier"`
	LastTransactions            []LastTransactionsItem `json:"lastTransactions"`
}

type LastTransactionsItem struct {
	OriginalTransactionId string `json:"originalTransactionId"`
	Status                int    `json:"status"`
	SignedRenewalInfo     string `json:"signedRenewalInfo"`
	SignedTransactionInfo string `json:"signedTransactionInfo"`
}

// ConsumptionRequestBody https://developer.apple.com/documentation/appstoreserverapi/consumptionrequest
type ConsumptionRequestBody struct {
	AccountTenure            int    `json:"accountTenure"`
	AppAccountToken          string `json:"appAccountToken"`
	ConsumptionStatus        int    `json:"consumptionStatus"`
	CustomerConsented        bool   `json:"customerConsented"`
	DeliveryStatus           int    `json:"deliveryStatus"`
	LifetimeDollarsPurchased int    `json:"lifetimeDollarsPurchased"`
	LifetimeDollarsRefunded  int    `json:"lifetimeDollarsRefunded"`
	Platform                 int    `json:"platform"`
	PlayTime                 int    `json:"playTime"`
	SampleContentProvided    bool   `json:"sampleContentProvided"`
	UserStatus               int    `json:"userStatus"`
}

type JWSRenewalInfoDecodedPayload struct {
}

// JWSDecodedHeader https://developer.apple.com/documentation/appstoreserverapi/jwsdecodedheader
type JWSDecodedHeader struct {
	Alg string   `json:"alg,omitempty"`
	Kid string   `json:"kid,omitempty"`
	X5C []string `json:"x5c,omitempty"`
}

// JWSTransaction https://developer.apple.com/documentation/appstoreserverapi/jwstransaction
type JWSTransaction struct {
	TransactionID               string `json:"transactionId,omitempty"`
	OriginalTransactionId       string `json:"originalTransactionId,omitempty"`
	WebOrderLineItemId          string `json:"webOrderLineItemId,omitempty"`
	BundleID                    string `json:"bundleId,omitempty"`
	ProductID                   string `json:"productId,omitempty"`
	SubscriptionGroupIdentifier string `json:"subscriptionGroupIdentifier,omitempty"`
	PurchaseDate                int64  `json:"purchaseDate,omitempty"`
	OriginalPurchaseDate        int64  `json:"originalPurchaseDate,omitempty"`
	ExpiresDate                 int64  `json:"expiresDate,omitempty"`
	Quantity                    int64  `json:"quantity,omitempty"`
	Type                        string `json:"type,omitempty"`
	AppAccountToken             string `json:"appAccountToken,omitempty"`
	InAppOwnershipType          string `json:"inAppOwnershipType,omitempty"`
	SignedDate                  int64  `json:"signedDate,omitempty"`
	OfferType                   int64  `json:"offerType,omitempty"`
	OfferIdentifier             string `json:"offerIdentifier,omitempty"`
	RevocationDate              int64  `json:"revocationDate,omitempty"`
	RevocationReason            int    `json:"revocationReason,omitempty"`
	IsUpgraded                  bool   `json:"isUpgraded,omitempty"`
}

func (J JWSTransaction) Valid() error {
	return nil
}

// https://developer.apple.com/documentation/appstoreserverapi/extendreasoncode
type ExtendReasonCode int

const (
	UndeclaredExtendReasonCode = iota
	CustomerSatisfaction
	OtherReasons
	ServiceIssueOrOutage
)

// ExtendRenewalDateRequest https://developer.apple.com/documentation/appstoreserverapi/extendrenewaldaterequest
type ExtendRenewalDateRequest struct {
	ExtendByDays      int              `json:"extendByDays"`
	ExtendReasonCode  ExtendReasonCode `json:"extendReasonCode"`
	RequestIdentifier string           `json:"requestIdentifier"`
}

// NotificationHistoryRequest https://developer.apple.com/documentation/appstoreserverapi/notificationhistoryrequest
type NotificationHistoryRequest struct {
	StartDate             int64                       `json:"startDate"`
	EndDate               int64                       `json:"endDate"`
	OriginalTransactionId string                      `json:"originalTransactionId,omitempty"`
	NotificationType      appstore.NotificationTypeV2 `json:"notificationType,omitempty"`
	NotificationSubtype   appstore.SubtypeV2          `json:"notificationSubtype,omitempty"`
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
}

// https://developer.apple.com/documentation/appstoreserverapi/firstsendattemptresult
type FirstSendAttemptResult string

const (
	FirstSendAttemptResultSuccess            FirstSendAttemptResult = "SUCCESS"
	FirstSendAttemptResultCircularRedirect   FirstSendAttemptResult = "CIRCULAR_REDIRECT"
	FirstSendAttemptResultInvalidResponse    FirstSendAttemptResult = "INVALID_RESPONSE"
	FirstSendAttemptResultNoResponse         FirstSendAttemptResult = "NO_RESPONSE"
	FirstSendAttemptResultOther              FirstSendAttemptResult = "OTHER"
	FirstSendAttemptResultPrematureClose     FirstSendAttemptResult = "PREMATURE_CLOSE"
	FirstSendAttemptResultSocketIssue        FirstSendAttemptResult = "SOCKET_ISSUE"
	FirstSendAttemptResultTimedOut           FirstSendAttemptResult = "TIMED_OUT"
	FirstSendAttemptResultTlsIssue           FirstSendAttemptResult = "TLS_ISSUE"
	FirstSendAttemptResultUnsupportedCharset FirstSendAttemptResult = "UNSUPPORTED_CHARSET"
)

// SendTestNotificationResponse https://developer.apple.com/documentation/appstoreserverapi/sendtestnotificationresponse
type SendTestNotificationResponse struct {
	TestNotificationToken string `json:"testNotificationToken"`
}
