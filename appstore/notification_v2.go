package appstore

import "github.com/golang-jwt/jwt/v5"

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
	NotificationTypeV2ExternalPurchaseToken  NotificationTypeV2 = "EXTERNAL_PURCHASE_TOKEN"
	NotificationTypeV2GracePeriodExpired     NotificationTypeV2 = "GRACE_PERIOD_EXPIRED"
	NotificationTypeV2OfferRedeemed          NotificationTypeV2 = "OFFER_REDEEMED"
	NotificationTypeV2OneTimeCharge          NotificationTypeV2 = "ONE_TIME_CHARGE"
	NotificationTypeV2PriceIncrease          NotificationTypeV2 = "PRICE_INCREASE"
	NotificationTypeV2Refund                 NotificationTypeV2 = "REFUND"
	NotificationTypeV2RefundDeclined         NotificationTypeV2 = "REFUND_DECLINED"
	NotificationTypeV2RefundReversed         NotificationTypeV2 = "REFUND_REVERSED"
	NotificationTypeV2RenewalExtended        NotificationTypeV2 = "RENEWAL_EXTENDED"
	NotificationTypeV2RenewalExtension       NotificationTypeV2 = "RENEWAL_EXTENSION"
	NotificationTypeV2Revoke                 NotificationTypeV2 = "REVOKE"
	NotificationTypeV2Subscribed             NotificationTypeV2 = "SUBSCRIBED"
	NotificationTypeV2Test                   NotificationTypeV2 = "TEST"
)

// SubtypeV2 is type
type SubtypeV2 string

// list of subtypes
// https://developer.apple.com/documentation/appstoreservernotifications/subtype
const (
	SubTypeV2Accepted          = "ACCEPTED"
	SubTypeV2AutoRenewDisabled = "AUTO_RENEW_DISABLED"
	SubTypeV2AutoRenewEnabled  = "AUTO_RENEW_ENABLED"
	SubTypeV2BillingRecovery   = "BILLING_RECOVERY"
	SubTypeV2BillingRetry      = "BILLING_RETRY"
	SubTypeV2Downgrade         = "DOWNGRADE"
	SubTypeV2Failure           = "FAILURE"
	SubTypeV2GracePeriod       = "GRACE_PERIOD"
	SubTypeV2InitialBuy        = "INITIAL_BUY"
	SubTypeV2Pending           = "PENDING"
	SubTypeV2PriceIncrease     = "PRICE_INCREASE"
	SubTypeV2ProductNotForSale = "PRODUCT_NOT_FOR_SALE"
	SubTypeV2Resubscribe       = "RESUBSCRIBE"
	SubTypeV2Summary           = "SUMMARY"
	SubTypeV2Upgrade           = "UPGRADE"
	SubTypeV2Unreported        = "UNREPORTED"
	SubTypeV2Voluntary         = "VOLUNTARY"
)

type AutoRenewStatus int

const (
	Off AutoRenewStatus = iota
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

// AutoRenewableSubscriptionStatus status value is current as of the signedDate in the decoded payload, SubscriptionNotificationV2DecodedPayload.
// https://developer.apple.com/documentation/appstoreservernotifications/status
type AutoRenewableSubscriptionStatus int32

const (
	AutoRenewableSubscriptionStatusActive = iota + 1
	AutoRenewableSubscriptionStatusExpired
	AutoRenewableSubscriptionStatusBillingRetryPeriod
	AutoRenewableSubscriptionStatusBillingGracePeriod
	AutoRenewableSubscriptionStatusRevoked
)

// TransactionReason indicates the cause of a purchase transaction,
// which indicates whether it’s a customer’s purchase or a renewal for an auto-renewable subscription that the system initiates.
// https://developer.apple.com/documentation/appstoreservernotifications/transactionreason
type TransactionReason string

const (
	TransactionReasonPurchase = "PURCHASE"
	TransactionReasonRenewal  = "RENEWAL"
)

type OfferDiscountType string

const (
	OfferDiscountTypeFreeTrial  OfferDiscountType = "FREE_TRIAL"
	OfferDiscountTypePayAsYouGo OfferDiscountType = "PAY_AS_YOU_GO"
	OfferDiscountTypePayUpFront OfferDiscountType = "PAY_UP_FRONT"
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
		NotificationType    NotificationTypeV2                `json:"notificationType"`
		Subtype             SubtypeV2                         `json:"subtype"`
		NotificationUUID    string                            `json:"notificationUUID"`
		NotificationVersion string                            `json:"version"`
		SignedDate          int64                             `json:"signedDate"`
		Data                SubscriptionNotificationV2Data    `json:"data,omitempty"`
		Summary             SubscriptionNotificationV2Summary `json:"summary,omitempty"`
		jwt.RegisteredClaims
	}

	// SubscriptionNotificationV2Summary is struct
	// https://developer.apple.com/documentation/appstoreservernotifications/summary
	SubscriptionNotificationV2Summary struct {
		RequestIdentifier      string `json:"requestIdentifier"`
		Environment            string `json:"environment"`
		AppAppleId             int64  `json:"appAppleId"`
		BundleID               string `json:"bundleId"`
		ProductID              string `json:"productId"`
		StorefrontCountryCodes string `json:"storefrontCountryCodes"`
		FailedCount            int64  `json:"failedCount"`
		SucceededCount         int64  `json:"succeededCount"`
	}

	// SubscriptionNotificationV2Data is struct
	// https://developer.apple.com/documentation/appstoreservernotifications/data
	SubscriptionNotificationV2Data struct {
		AppAppleID            int                             `json:"appAppleId"`
		BundleID              string                          `json:"bundleId"`
		BundleVersion         string                          `json:"bundleVersion"`
		Environment           string                          `json:"environment"`
		SignedRenewalInfo     JWSRenewalInfo                  `json:"signedRenewalInfo"`
		SignedTransactionInfo JWSTransaction                  `json:"signedTransactionInfo"`
		Status                AutoRenewableSubscriptionStatus `json:"status"`
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
		AutoRenewProductId          string              `json:"autoRenewProductId"`
		AutoRenewStatus             AutoRenewStatus     `json:"autoRenewStatus"`
		Currency                    string              `json:"currency"`
		Environment                 Environment         `json:"environment"`
		ExpirationIntent            ExpirationIntent    `json:"expirationIntent"`
		GracePeriodExpiresDate      int64               `json:"gracePeriodExpiresDate"`
		IsInBillingRetryPeriod      bool                `json:"isInBillingRetryPeriod"`
		OfferIdentifier             string              `json:"offerIdentifier"`
		OfferType                   OfferType           `json:"offerType"`
		OriginalTransactionId       string              `json:"originalTransactionId"`
		PriceIncreaseStatus         PriceIncreaseStatus `json:"priceIncreaseStatus"`
		ProductId                   string              `json:"productId"`
		RecentSubscriptionStartDate int64               `json:"recentSubscriptionStartDate"`
		RenewalDate                 int64               `json:"renewalDate"`
		RenewalPrice                int64               `json:"renewalPrice"`
		SignedDate                  int64               `json:"signedDate"`
		jwt.RegisteredClaims
	}

	// JWSTransactionDecodedPayload contains the decoded transaction information
	// https://developer.apple.com/documentation/appstoreservernotifications/jwstransactiondecodedpayload
	JWSTransactionDecodedPayload struct {
		AppAccountToken             string            `json:"appAccountToken"`
		BundleId                    string            `json:"bundleId"`
		Currency                    string            `json:"currency,omitempty"`
		Environment                 Environment       `json:"environment"`
		ExpiresDate                 int64             `json:"expiresDate"`
		InAppOwnershipType          string            `json:"inAppOwnershipType"`
		IsUpgraded                  bool              `json:"isUpgraded"`
		OfferDiscountType           OfferDiscountType `json:"offerDiscountType"`
		OfferIdentifier             string            `json:"offerIdentifier"`
		OfferType                   OfferType         `json:"offerType"`
		OriginalPurchaseDate        int64             `json:"originalPurchaseDate"`
		OriginalTransactionId       string            `json:"originalTransactionId"`
		Price                       int64             `json:"price,omitempty"`
		ProductId                   string            `json:"productId"`
		PurchaseDate                int64             `json:"purchaseDate"`
		Quantity                    int64             `json:"quantity"`
		RevocationDate              int64             `json:"revocationDate"`
		RevocationReason            RevocationReason  `json:"revocationReason"`
		SignedDate                  int64             `json:"signedDate"`
		Storefront                  string            `json:"storefront"`
		StorefrontId                string            `json:"storefrontId"`
		SubscriptionGroupIdentifier string            `json:"subscriptionGroupIdentifier"`
		TransactionId               string            `json:"transactionId"`
		TransactionReason           TransactionReason `json:"transactionReason"`
		IAPtype                     IAPType           `json:"type"`
		WebOrderLineItemId          string            `json:"webOrderLineItemId"`
		jwt.RegisteredClaims
	}
)
