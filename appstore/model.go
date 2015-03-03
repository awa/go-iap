package appstore

type (
	// The IAPRequest type has the request parameter
	IAPRequest struct {
		ReceiptData string `json:"receipt-data"`
		Password    string `json:"password,omitempty"`
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
		Quantity                  string `json:"quantity"`
		ProductID                 string `json:"product_id"`
		TransactionID             string `json:"transaction_id"`
		OriginalTransactionID     string `json:"original_transaction_id"`
		IsTrialPeriod             string `json:"is_trial_period"`
		AppItemID                 string `json:"app_item_id"`
		VersionExternalIdentifier string `json:"version_external_identifier"`
		WebOrderLineItemID        string `json:"web_order_line_item_id"`
		PurchaseDate
		OriginalPurchaseDate
		ExpiresDate
		CancellationDate
	}

	// The Receipt type has whole data of receipt
	Receipt struct {
		ReceiptType                string  `json:"receipt_type"`
		AdamID                     int64   `json:"adam_id"`
		AppItemID                  int64   `json:"app_item_id"`
		BundleID                   string  `json:"bundle_id"`
		ApplicationVersion         string  `json:"application_version"`
		DownloadID                 int64   `json:"download_id"`
		OriginalApplicationVersion string  `json:"original_application_version"`
		InApp                      []InApp `json:"in_app"`
		RequestDate
		OriginalPurchaseDate
	}

	// The IAPResponse type has the response properties
	IAPResponse struct {
		Status            int     `json:"status"`
		Environment       string  `json:"environment"`
		Receipt           Receipt `json:"receipt"`
		LatestReceiptInfo []InApp `json:"latest_receipt_info"`
		LatestReceipt     string  `json:"latest_receipt"`
	}
)
