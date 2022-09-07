package appstore

// NotificationTypeV2 is type
type NotificationTypeV2 string

// list of notificationType
// https://developer.apple.com/documentation/appstoreservernotifications/notificationtype
const (
	NotificationTypeV2ConsumptionRequest     NotificationTypeV2 = "CONSUMPTION_REQUEST"
	NotificationTypeV2DidChangeRenewalPref   NotificationTypeV2 = "DID_CHANGE_RENEWAL_PREF"
	NotificationTypeV2DidChangeRenewalStatus NotificationTypeV2 = "DID_CHANGE_RENEWAL_STATUS"
	NotificationTypeV2DidFailToRenew         NotificationTypeV2 = "DID_FAIL_TO_RENEW"
	NotificationTypeV2DidRenew               NotificationTypeV2 = "DID_RENEW"
	NotificationTypeV2Expired                NotificationTypeV2 = "EXPIRED"
	NotificationTypeV2GracePeriodExpired     NotificationTypeV2 = "GRACE_PERIOD_EXPIRED"
	NotificationTypeV2OfferRedeemed          NotificationTypeV2 = "OFFER_REDEEMED"
	NotificationTypeV2PriceIncrease          NotificationTypeV2 = "PRICE_INCREASE"
	NotificationTypeV2Refund                 NotificationTypeV2 = "REFUND"
	NotificationTypeV2RefundDeclined         NotificationTypeV2 = "REFUND_DECLINED"
	NotificationTypeV2RenewalExtended        NotificationTypeV2 = "RENEWAL_EXTENDED"
	NotificationTypeV2Revoke                 NotificationTypeV2 = "REVOKE"
	NotificationTypeV2Subscribed             NotificationTypeV2 = "SUBSCRIBED"
)

// SubtypeV2 is type
type SubtypeV2 string

// list of subtypes
// https://developer.apple.com/documentation/appstoreservernotifications/subtype
const (
	SubTypeV2InitialBuy        = "INITIAL_BUY"
	SubTypeV2Resubscribe       = "RESUBSCRIBE"
	SubTypeV2Downgrade         = "DOWNGRADE"
	SubTypeV2Upgrade           = "UPGRADE"
	SubTypeV2AutoRenewEnabled  = "AUTO_RENEW_ENABLED"
	SubTypeV2AutoRenewDisabled = "AUTO_RENEW_DISABLED"
	SubTypeV2Voluntary         = "VOLUNTARY"
	SubTypeV2BillingRetry      = "BILLING_RETRY"
	SubTypeV2PriceIncrease     = "PRICE_INCREASE"
	SubTypeV2GracePeriod       = "GRACE_PERIOD"
	SubTypeV2BillingRecovery   = "BILLING_RECOVERY"
	SubTypeV2Pending           = "PENDING"
	SubTypeV2Accepted          = "ACCEPTED"
)

type AutoRenewStatus int

const (
	Off RevocationReason = iota
	On
)

type ExpirationIntent int

const (
	CustomerCancelled ExpirationIntent = iota + 1
	BillingError
	NoPriceChangeConsent
	ProductUnavailable
)

type OfferType int

const (
	IntroductoryOffer OfferType = iota + 1
	PromotionalOffer
	SubscriptionOfferCode
)

type PriceIncreaseStatus int

const (
	CustomerNotYetConsented PriceIncreaseStatus = iota
	CustomerConsented
)

type RevocationReason int

const (
	OtherReason RevocationReason = iota
	AppIssue
)

type IAPType string

const (
	AutoRenewable IAPType = "Auto-Renewable Subscription"
	NonConsumable IAPType = "Non-Consumable"
	Consumable    IAPType = "Consumable"
	NonRenewable  IAPType = "Non-Renewing Subscription"
)

type (
	// SubscriptionNotificationV2 is struct for
	// https://developer.apple.com/documentation/appstoreservernotifications/responsebodyv2
	SubscriptionNotificationV2 struct {
		SignedPayload SubscriptionNotificationV2SignedPayload `json:"signedPayload"`
	}

	// SubscriptionNotificationV2SignedPayload is struct
	// https://developer.apple.com/documentation/appstoreservernotifications/signedpayload
	SubscriptionNotificationV2SignedPayload struct {
		SignedPayload string `json:"signedPayload"`
	}

	// SubscriptionNotificationV2DecodedPayload is struct
	// https://developer.apple.com/documentation/appstoreservernotifications/responsebodyv2decodedpayload
	SubscriptionNotificationV2DecodedPayload struct {
		NotificationType    NotificationTypeV2             `json:"notificationType"`
		Subtype             SubtypeV2                      `json:"subtype"`
		NotificationUUID    string                         `json:"notificationUUID"`
		NotificationVersion string                         `json:"version"`
		SignedDate          int64                          `json:"signedDate"`
		Data                SubscriptionNotificationV2Data `json:"data"`
	}

	// SubscriptionNotificationV2Data is struct
	// https://developer.apple.com/documentation/appstoreservernotifications/data
	SubscriptionNotificationV2Data struct {
		AppAppleID            int            `json:"appAppleId"`
		BundleID              string         `json:"bundleId"`
		BundleVersion         string         `json:"bundleVersion"`
		Environment           string         `json:"environment"`
		SignedRenewalInfo     JWSRenewalInfo `json:"signedRenewalInfo"`
		SignedTransactionInfo JWSTransaction `json:"signedTransactionInfo"`
	}

	// SubscriptionNotificationV2JWSDecodedHeader is struct
	SubscriptionNotificationV2JWSDecodedHeader struct {
		Alg string   `json:"alg"`
		Kid string   `json:"kid"`
		X5c []string `json:"x5c"`
	}

	// JWSRenewalInfo contains the Base64 encoded signed JWS payload of the renewal information
	// https://developer.apple.com/documentation/appstoreservernotifications/jwsrenewalinfo
	JWSRenewalInfo string

	// JWSTransaction contains the Base64 encoded signed JWS payload of the transaction
	// https://developer.apple.com/documentation/appstoreservernotifications/jwstransaction
	JWSTransaction string

	// JWSRenewalInfoDecodedPayload contains the decoded renewal information
	// https://developer.apple.com/documentation/appstoreservernotifications/jwsrenewalinfodecodedpayload
	JWSRenewalInfoDecodedPayload struct {
		AutoRenewProductId     string              `json:"autoRenewProductId"`
		AutoRenewStatus        AutoRenewStatus     `json:"autoRenewStatus"`
		Environment            Environment         `json:"environment"`
		ExpirationIntent       ExpirationIntent    `json:"expirationIntent"`
		GracePeriodExpiresDate int64               `json:"gracePeriodExpiresDate"`
		IsInBillingRetryPeriod bool                `json:"isInBillingRetryPeriod"`
		OfferIdentifier        string              `json:"offerIdentifier"`
		OfferType              OfferType           `json:"offerType"`
		OriginalTransactionId  string              `json:"originalTransactionId"`
		PriceIncreaseStatus    PriceIncreaseStatus `json:"priceIncreaseStatus"`
		ProductId              string              `json:"productId"`
		SignedDate             int64               `json:"signedDate"`
	}

	// JWSTransactionDecodedPayload contains the decoded transaction information
	// https://developer.apple.com/documentation/appstoreservernotifications/jwstransactiondecodedpayload
	JWSTransactionDecodedPayload struct {
		AppAccountToken             string           `json:"appAccountToken"`
		BundleId                    string           `json:"bundleId"`
		Environment                 Environment      `json:"environment"`
		ExpiresDate                 int64            `json:"expiresDate"`
		InAppOwnershipType          string           `json:"inAppOwnershipType"`
		IsUpgraded                  bool             `json:"isUpgraded"`
		OfferIdentifier             string           `json:"offerIdentifier"`
		OfferType                   OfferType        `json:"offerType"`
		OriginalPurchaseDate        int64            `json:"originalPurchaseDate"`
		OriginalTransactionId       string           `json:"originalTransactionId"`
		ProductId                   string           `json:"productId"`
		PurchaseDate                int64            `json:"purchaseDate"`
		Quantity                    int64            `json:"quantity"`
		RevocationDate              int64            `json:"revocationDate"`
		RevocationReason            RevocationReason `json:"revocationReason"`
		SignedDate                  int64            `json:"signedDate"`
		SubscriptionGroupIdentifier string           `json:"subscriptionGroupIdentifier"`
		TransactionId               string           `json:"transactionId"`
		IAPtype                     IAPType          `json:"type"`
		WebOrderLineItemId          string           `json:"webOrderLineItemId"`
	}
)
