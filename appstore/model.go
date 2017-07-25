package appstore

type (
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

	// The ExpiresDate type indicates the expiration date for the subscription
	ExpiresDate struct {
		ExpiresDate    string `json:"expires_date"`
		ExpiresDateMS  string `json:"expires_date_ms"`
		ExpiresDatePST string `json:"expires_date_pst"`
	}

	// The CancellationDate type indicates the time and date of the cancellation by Apple customer support
	CancellationDate struct {
		CancellationDate    string `json:"cancellation_date"`
		CancellationDateMS  string `json:"cancellation_date_ms"`
		CancellationDatePST string `json:"cancellation_date_pst"`
	}

	// The InApp type has the receipt attributes
	InApp struct {
		Quantity              string `json:"quantity"`
		ProductID             string `json:"product_id"`
		TransactionID         string `json:"transaction_id"`
		OriginalTransactionID string `json:"original_transaction_id"`
		WebOrderLineItemID    string `json:"web_order_line_item_id"`

		IsTrialPeriod string `json:"is_trial_period"`
		ExpiresDate

		PurchaseDate
		OriginalPurchaseDate

		CancellationDate
		CancellationReason string `json:"cancellation_reason"`
	}

	// The Receipt type has whole data of receipt
	Receipt struct {
		ReceiptType                string  `json:"receipt_type"`
		AdamID                     int64   `json:"adam_id"`
		AppItemID                  int64   `json:"app_item_id"`
		BundleID                   string  `json:"bundle_id"`
		ApplicationVersion         string  `json:"application_version"`
		DownloadID                 int64   `json:"download_id"`
		VersionExternalIdentifier  int64   `json:"version_external_identifier"`
		OriginalApplicationVersion string  `json:"original_application_version"`
		InApp                      []InApp `json:"in_app"`
		ReceiptCreationDate
		RequestDate
		OriginalPurchaseDate
	}

	// A pending renewal may refer to a renewal that is scheduled in the future or a renewal that failed in the past for some reason.
	PendingRenewalInfo struct {
		SubscriptionExpirationIntent   string `json:"expiration_intent"`
		SubscriptionAutoRenewProductID string `json:"auto_renew_product_id"`
		SubscriptionRetryFlag          string `json:"is_in_billing_retry_period"`
		SubscriptionAutoRenewStatus    string `json:"auto_renew_status"`
		SubscriptionPriceConsentStatus string `json:"price_consent_status"`
		ProductID                      string `json:"product_id"`
	}

	// The IAPResponse type has the response properties
	// We defined each field by the current IAP response, but some fields are not mentioned
	// in the following Apple's document;
	// https://developer.apple.com/library/ios/releasenotes/General/ValidateAppStoreReceipt/Chapters/ReceiptFields.html
	// If you get other types or fileds from the IAP response, you should use the struct you defined.
	IAPResponse struct {
		Status             int                  `json:"status"`
		Environment        string               `json:"environment"`
		Receipt            Receipt              `json:"receipt"`
		LatestReceiptInfo  []InApp              `json:"latest_receipt_info"`
		LatestReceipt      string               `json:"latest_receipt"`
		PendingRenewalInfo []PendingRenewalInfo `json:"pending_renewal_info"`
		IsRetryable        bool                 `json:"is-retryable"`
	}
)
