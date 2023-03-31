package appstore

import "encoding/json"

type numericString string

func (n *numericString) UnmarshalJSON(b []byte) error {
	var number json.Number
	if err := json.Unmarshal(b, &number); err != nil {
		return err
	}
	*n = numericString(number.String())
	return nil
}

// Environment is alias
type Environment string

// list of Environment
const (
	Sandbox    Environment = "Sandbox"
	Production Environment = "Production"
)

type (
	// IAPRequest is struct
	// https://developer.apple.com/library/content/releasenotes/General/ValidateAppStoreReceipt/Chapters/ValidateRemotely.html
	// The IAPRequest type has the request parameter
	IAPRequest struct {
		ReceiptData string `json:"receipt-data"`
		// Only used for receipts that contain auto-renewable subscriptions.
		Password string `json:"password,omitempty"`
		// Only used for iOS7 style app receipts that contain auto-renewable or non-renewing subscriptions.
		// If value is true, response includes only the latest renewal transaction for any subscriptions.
		ExcludeOldTransactions bool `json:"exclude-old-transactions"`
	}

	// The ReceiptCreationDate type indicates the date when the app receipt was created.
	ReceiptCreationDate struct {
		CreationDate    string `json:"receipt_creation_date"`
		CreationDateMS  string `json:"receipt_creation_date_ms"`
		CreationDatePST string `json:"receipt_creation_date_pst"`
	}

	// The RequestDate type indicates the date and time that the request was sent
	RequestDate struct {
		RequestDate    string `json:"request_date"`
		RequestDateMS  string `json:"request_date_ms"`
		RequestDatePST string `json:"request_date_pst"`
	}

	// The PurchaseDate type indicates the date and time that the item was purchased
	PurchaseDate struct {
		PurchaseDate    string `json:"purchase_date"`
		PurchaseDateMS  string `json:"purchase_date_ms"`
		PurchaseDatePST string `json:"purchase_date_pst"`
	}

	// The OriginalPurchaseDate type indicates the beginning of the subscription period
	OriginalPurchaseDate struct {
		OriginalPurchaseDate    string `json:"original_purchase_date"`
		OriginalPurchaseDateMS  string `json:"original_purchase_date_ms"`
		OriginalPurchaseDatePST string `json:"original_purchase_date_pst"`
	}

	// The PreorderDate type indicates the date and time that the pre-order
	PreorderDate struct {
		PreorderDate    string `json:"preorder_date"`
		PreorderDateMS  string `json:"preorder_date_ms"`
		PreorderDatePST string `json:"preorder_date_pst"`
	}

	// The ExpiresDate type indicates the expiration date for the subscription
	ExpiresDate struct {
		ExpiresDate             string `json:"expires_date,omitempty"`
		ExpiresDateMS           string `json:"expires_date_ms,omitempty"`
		ExpiresDatePST          string `json:"expires_date_pst,omitempty"`
		ExpiresDateFormatted    string `json:"expires_date_formatted,omitempty"`
		ExpiresDateFormattedPST string `json:"expires_date_formatted_pst,omitempty"`
	}

	// The CancellationDate type indicates the time and date of the cancellation by Apple customer support
	CancellationDate struct {
		CancellationDate    string `json:"cancellation_date,omitempty"`
		CancellationDateMS  string `json:"cancellation_date_ms,omitempty"`
		CancellationDatePST string `json:"cancellation_date_pst,omitempty"`
	}

	// The GracePeriodDate type indicates the grace period date for the subscription
	GracePeriodDate struct {
		GracePeriodDate    string `json:"grace_period_expires_date,omitempty"`
		GracePeriodDateMS  string `json:"grace_period_expires_date_ms,omitempty"`
		GracePeriodDatePST string `json:"grace_period_expires_date_pst,omitempty"`
	}

	// AutoRenewStatusChangeDate type indicates the auto renew status change date
	AutoRenewStatusChangeDate struct {
		AutoRenewStatusChangeDate    string `json:"auto_renew_status_change_date"`
		AutoRenewStatusChangeDateMS  string `json:"auto_renew_status_change_date_ms"`
		AutoRenewStatusChangeDatePST string `json:"auto_renew_status_change_date_pst"`
	}

	// The InApp type has the receipt attributes
	InApp struct {
		Quantity                    string `json:"quantity"`
		ProductID                   string `json:"product_id"`
		TransactionID               string `json:"transaction_id"`
		OriginalTransactionID       string `json:"original_transaction_id"` // this field is string
		WebOrderLineItemID          string `json:"web_order_line_item_id,omitempty"`
		PromotionalOfferID          string `json:"promotional_offer_id"`
		SubscriptionGroupIdentifier string `json:"subscription_group_identifier"`
		OfferCodeRefName            string `json:"offer_code_ref_name,omitempty"`
		AppAccountToken             string `json:"app_account_token,omitempty"`

		IsTrialPeriod        string `json:"is_trial_period"`
		IsInIntroOfferPeriod string `json:"is_in_intro_offer_period,omitempty"`
		IsUpgraded           string `json:"is_upgraded,omitempty"`

		ExpiresDate

		PurchaseDate
		OriginalPurchaseDate

		CancellationDate
		CancellationReason string `json:"cancellation_reason,omitempty"`

		InAppOwnershipType string `json:"in_app_ownership_type,omitempty"`
	}

	// The Receipt type has whole data of receipt
	Receipt struct {
		ReceiptType                string        `json:"receipt_type"`
		AdamID                     int64         `json:"adam_id"`
		AppItemID                  numericString `json:"app_item_id"`
		BundleID                   string        `json:"bundle_id"`
		ApplicationVersion         string        `json:"application_version"`
		DownloadID                 int64         `json:"download_id"`
		VersionExternalIdentifier  numericString `json:"version_external_identifier"`
		OriginalApplicationVersion string        `json:"original_application_version"`
		InApp                      []InApp       `json:"in_app"`
		ReceiptCreationDate
		RequestDate
		OriginalPurchaseDate
		PreorderDate
		ExpiresDate
	}

	// PendingRenewalInfo is struct
	// A pending renewal may refer to a renewal that is scheduled in the future or a renewal that failed in the past for some reason.
	// https://developer.apple.com/documentation/appstoreservernotifications/unified_receipt/pending_renewal_info
	PendingRenewalInfo struct {
		SubscriptionExpirationIntent   string `json:"expiration_intent"`
		SubscriptionAutoRenewProductID string `json:"auto_renew_product_id"`
		SubscriptionRetryFlag          string `json:"is_in_billing_retry_period"`
		SubscriptionAutoRenewStatus    string `json:"auto_renew_status"`
		SubscriptionPriceConsentStatus string `json:"price_consent_status"`
		ProductID                      string `json:"product_id"`
		OriginalTransactionID          string `json:"original_transaction_id"`
		OfferCodeRefName               string `json:"offer_code_ref_name,omitempty"`
		PromotionalOfferID             string `json:"promotional_offer_id,omitempty"`
		PriceIncreaseStatus            string `json:"price_increase_status,omitempty"`

		GracePeriodDate
	}

	// The IAPResponse type has the response properties
	// We defined each field by the current IAP response, but some fields are not mentioned
	// in the following Apple's document;
	// https://developer.apple.com/library/ios/releasenotes/General/ValidateAppStoreReceipt/Chapters/ReceiptFields.html
	// If you get other types or fields from the IAP response, you should use the struct you defined.
	IAPResponse struct {
		Status             int                  `json:"status"`
		Environment        Environment          `json:"environment"`
		Receipt            Receipt              `json:"receipt"`
		LatestReceiptInfo  []InApp              `json:"latest_receipt_info,omitempty"`
		LatestReceipt      string               `json:"latest_receipt,omitempty"`
		PendingRenewalInfo []PendingRenewalInfo `json:"pending_renewal_info,omitempty"`
		IsRetryable        bool                 `json:"is_retryable,omitempty"`
	}

	// StatusResponse is struct
	// The HttpStatusResponse struct contains the status code returned by the store
	// Used as a workaround to detect when to hit the production appstore or sandbox appstore regardless of receipt type
	StatusResponse struct {
		Status int `json:"status"`
	}

	// IAPResponseForIOS6 is iOS 6 style receipt schema.
	IAPResponseForIOS6 struct {
		AutoRenewProductID       string         `json:"auto_renew_product_id"`
		AutoRenewStatus          int            `json:"auto_renew_status"`
		CancellationReason       string         `json:"cancellation_reason,omitempty"`
		ExpirationIntent         string         `json:"expiration_intent,omitempty"`
		IsInBillingRetryPeriod   string         `json:"is_in_billing_retry_period,omitempty"`
		Receipt                  ReceiptForIOS6 `json:"receipt"`
		LatestExpiredReceiptInfo ReceiptForIOS6 `json:"latest_expired_receipt_info"`
		LatestReceipt            string         `json:"latest_receipt"`
		LatestReceiptInfo        ReceiptForIOS6 `json:"latest_receipt_info"`
		Status                   int            `json:"status"`
	}

	// ReceiptForIOS6 is struct
	ReceiptForIOS6 struct {
		AppItemID numericString `json:"app_item_id"`
		BID       string        `json:"bid"`
		BVRS      string        `json:"bvrs"`
		CancellationDate
		ExpiresDate
		IsTrialPeriod        string `json:"is_trial_period"`
		IsInIntroOfferPeriod string `json:"is_in_intro_offer_period"`
		ItemID               string `json:"item_id"`
		ProductID            string `json:"product_id"`
		PurchaseDate
		OriginalTransactionID numericString `json:"original_transaction_id"`
		OriginalPurchaseDate
		Quantity                  string        `json:"quantity"`
		TransactionID             string        `json:"transaction_id"`
		UniqueIdentifier          string        `json:"unique_identifier"`
		UniqueVendorIdentifier    string        `json:"unique_vendor_identifier"`
		VersionExternalIdentifier numericString `json:"version_external_identifier,omitempty"`
		WebOrderLineItemID        string        `json:"web_order_line_item_id"`
	}
)
