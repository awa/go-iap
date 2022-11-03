go-iap
======

![](https://img.shields.io/badge/golang-1.18+-blue.svg?style=flat)
[![Build Status](https://travis-ci.org/awa/go-iap.svg?branch=master)](https://travis-ci.org/awa/go-iap)
[![codecov.io](https://codecov.io/github/awa/go-iap/coverage.svg?branch=master)](https://codecov.io/github/awa/go-iap?branch=master)

go-iap verifies the purchase receipt via AppStore, GooglePlayStore or Amazon AppStore.

Current API Documents:

* AppStore: [![GoDoc](https://godoc.org/github.com/awa/go-iap/appstore?status.svg)](https://godoc.org/github.com/awa/go-iap/appstore)
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
