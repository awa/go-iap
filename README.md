go-iap
======

![](https://img.shields.io/badge/golang-1.20+-blue.svg?style=flat)
[![unit test](https://github.com/awa/go-iap/actions/workflows/unit_test.yml/badge.svg)](https://github.com/awa/go-iap/actions/workflows/unit_test.yml)

>go-iap verifies the purchase receipt via AppStore, GooglePlayStore or Amazon AppStore.

Current API Documents:

* AppStore: [![GoDoc](https://godoc.org/github.com/awa/go-iap/appstore?status.svg)](https://godoc.org/github.com/awa/go-iap/appstore)
* AppStore Server API: [![GoDoc](https://godoc.org/github.com/awa/go-iap/appstore?status.svg)](https://godoc.org/github.com/awa/go-iap/appstore/api)
* GooglePlay: [![GoDoc](https://godoc.org/github.com/awa/go-iap/playstore?status.svg)](https://godoc.org/github.com/awa/go-iap/playstore)
* Amazon AppStore: [![GoDoc](https://godoc.org/github.com/awa/go-iap/amazon?status.svg)](https://godoc.org/github.com/awa/go-iap/amazon)
* Huawei HMS: [![GoDoc](https://godoc.org/github.com/awa/go-iap/hms?status.svg)](https://godoc.org/github.com/awa/go-iap/hms)


# Installation
```
go get github.com/awa/go-iap/appstore
go get github.com/awa/go-iap/playstore
go get github.com/awa/go-iap/amazon
go get github.com/awa/go-iap/hms
```


# Quick Start

### In App Purchase (via App Store)

```go
import(
    "github.com/awa/go-iap/appstore"
)

func main() {
	client := appstore.New()
	req := appstore.IAPRequest{
		ReceiptData: "your receipt data encoded by base64",
	}
	resp := &appstore.IAPResponse{}
	ctx := context.Background()
	err := client.Verify(ctx, req, resp)
}
```

**Note**: The [verifyReceipt](https://developer.apple.com/documentation/appstorereceipts/verifyreceipt) API has been deprecated as of `5 Jun 2023`. Please use [App Store Server API](#in-app-store-server-api) instead.

### In App Billing (via GooglePlay)

```go
import(
    "github.com/awa/go-iap/playstore"
)

func main() {
	// You need to prepare a public key for your Android app's in app billing
	// at https://console.developers.google.com.
	jsonKey, err := ioutil.ReadFile("jsonKey.json")
	if err != nil {
		log.Fatal(err)
	}

	client := playstore.New(jsonKey)
	ctx := context.Background()
	resp, err := client.VerifySubscription(ctx, "package", "subscriptionID", "purchaseToken")
}
```

### In App Purchase (via Amazon App Store)

```go
import(
    "github.com/awa/go-iap/amazon"
)

func main() {
	client := amazon.New("developerSecret")

	ctx := context.Background()
	resp, err := client.Verify(ctx, "userID", "receiptID")
}
```

### In App Purchase (via Huawei Mobile Services)

```go
import(
    "github.com/awa/go-iap/hms"
)

func main() {
	// If "orderSiteURL" and/or "subscriptionSiteURL" are empty,
	// they will be default to AppTouch German.
	// Please refer to https://developer.huawei.com/consumer/en/doc/HMSCore-References-V5/api-common-statement-0000001050986127-V5 for details.
	client := hms.New("clientID", "clientSecret", "orderSiteURL", "subscriptionSiteURL")
	ctx := context.Background()
	resp, err := client.VerifySubscription(ctx, "purchaseToken", "subscriptionID", 1)
}
```

### In App Store Server API

**Note**
- The App Store Server API differentiates between a sandbox and a production environment based on the base URL:  
  - Use https://api.storekit.itunes.apple.com/ for the production environment.
  - Use https://api.storekit-sandbox.itunes.apple.com/ for the sandbox environment.
- If you're unsure about the environment, follow these steps:
  - Initiate a call to the endpoint using the production URL. If the call is successful, the transaction identifier is associated with the production environment.
  - If you encounter an error code `4040010`, indicating a `TransactionIdNotFoundError`, make a call to the endpoint using the sandbox URL.
  - If this call is successful, the transaction identifier is associated with the sandbox environment. If the call fails with the same error code, the transaction identifier doesn't exist in either environment.

- GetTransactionInfo

```go
import(
	"github.com/awa/go-iap/appstore/api"
)

//  For generate key file and download it, please refer to https://developer.apple.com/documentation/appstoreserverapi/creating_api_keys_to_use_with_the_app_store_server_api
const ACCOUNTPRIVATEKEY = `
    -----BEGIN PRIVATE KEY-----
    FAKEACCOUNTKEYBASE64FORMAT
    -----END PRIVATE KEY-----
    `
func main() {
	c := &api.StoreConfig{
		KeyContent: []byte(ACCOUNTPRIVATEKEY),  // Loads a .p8 certificate
		KeyID:      "FAKEKEYID",                // Your private key ID from App Store Connect (Ex: 2X9R4HXF34)
		BundleID:   "fake.bundle.id",           // Your app’s bundle ID
		Issuer:     "xxxxx-xx-xx-xx-xxxxxxxxxx",// Your issuer ID from the Keys page in App Store Connect (Ex: "57246542-96fe-1a63-e053-0824d011072a")
		Sandbox:    false,                      // default is Production
	}
	transactionId := "FAKETRANSACTIONID"
	a := api.NewStoreClient(c)
	ctx := context.Background()
	response, err := a.GetTransactionInfo(ctx, transactionId)

	transantion, err := a.ParseSignedTransaction(response.SignedTransactionInfo)
	if err != nil {
	    // error handling
	}

	if transaction.TransactionId == transactionId {
		// the transaction is valid
	}
}
```

- GetTransactionHistory

```go
import(
	"github.com/awa/go-iap/appstore/api"
)

//  For generate key file and download it, please refer to https://developer.apple.com/documentation/appstoreserverapi/creating_api_keys_to_use_with_the_app_store_server_api
const ACCOUNTPRIVATEKEY = `
    -----BEGIN PRIVATE KEY-----
    FAKEACCOUNTKEYBASE64FORMAT
    -----END PRIVATE KEY-----
    `
func main() {
	c := &api.StoreConfig{
		KeyContent: []byte(ACCOUNTPRIVATEKEY),  // Loads a .p8 certificate
		KeyID:      "FAKEKEYID",                // Your private key ID from App Store Connect (Ex: 2X9R4HXF34)
		BundleID:   "fake.bundle.id",           // Your app’s bundle ID
		Issuer:     "xxxxx-xx-xx-xx-xxxxxxxxxx",// Your issuer ID from the Keys page in App Store Connect (Ex: "57246542-96fe-1a63-e053-0824d011072a")
		Sandbox:    false,                      // default is Production
	}
	originalTransactionId := "FAKETRANSACTIONID"
	a := api.NewStoreClient(c)
	query := &url.Values{}
	query.Set("productType", "AUTO_RENEWABLE")
	query.Set("productType", "NON_CONSUMABLE")
	ctx := context.Background()
	responses, err := a.GetTransactionHistory(ctx, originalTransactionId, query)

	for _, response := range responses {
		transantions, err := a.ParseSignedTransactions(response.SignedTransactions)
	}
}
```
- Error handling
  - handler error per [apple store server api error](https://developer.apple.com/documentation/appstoreserverapi/error_codes) document
  - [error definition](./appstore/api/error.go)


### Parse Notification from App Store

```go
import (
	"github.com/awa/go-iap/appstore"
	"github.com/golang-jwt/jwt/v4"
)

func main() {
	tokenStr := "SignedRenewalInfo Encode String" // or SignedTransactionInfo string
	token := jwt.Token{}
	client := appstore.New()
	err := client.ParseNotificationV2(tokenStr, &token)

	claims, ok := token.Claims.(jwt.MapClaims)
	for key, val := range claims {
		fmt.Printf("Key: %v, value: %v\n", key, val) // key value of SignedRenewalInfo
	}
}
```

# ToDo
- [x] Validator for In App Purchase Receipt (AppStore)
- [x] Validator for Subscription token (GooglePlay)
- [x] Validator for Purchase Product token (GooglePlay)
- [ ] More Tests


# Support

### In App Purchase
This validator supports the receipt type for iOS7 or above.

### In App Billing
This validator uses [Version 3 API](http://developer.android.com/google/play/billing/api.html).

### In App Purchase (Amazon)
This validator uses [RVS for IAP v2.0](https://developer.amazon.com/public/apis/earn/in-app-purchasing/docs-v2/verifying-receipts-in-iap-2.0).

### In App Purchase (HMS)
This validator uses [Version 2 API](https://developer.huawei.com/consumer/en/doc/HMSCore-References-V5/api-common-statement-0000001050986127-V5).

### In App Store Server API
This validator uses [Version 1.0+](https://developer.apple.com/documentation/appstoreserverapi)

# License
go-iap is licensed under the MIT.
