go-iap
======

[![Build Status](https://travis-ci.org/dogenzaka/go-iap.svg?branch=master)](https://travis-ci.org/dogenzaka/go-iap)
[![codecov.io](https://codecov.io/github/dogenzaka/go-iap/coverage.svg?branch=master)](https://codecov.io/github/dogenzaka/go-iap?branch=master)

go-iap verifies the purchase receipt via AppStore or GooglePlayStore.

Current API Documents:

* AppStore: [![GoDoc](https://godoc.org/github.com/dogenzaka/go-iap/appstore?status.svg)](https://godoc.org/github.com/dogenzaka/go-iap/appstore)
* GooglePlay: [![GoDoc](https://godoc.org/github.com/dogenzaka/go-iap/playstore?status.svg)](https://godoc.org/github.com/dogenzaka/go-iap/playstore)

# Dependencies
```
go get github.com/parnurzeal/gorequest
go get golang.org/x/net/context
go get golang.org/x/oauth2
go get google.golang.org/api/androidpublisher/v2
go get google.golang.org/appengine/urlfetch
go get google.golang.org/appengine/aetest
```

# Installation
```
go get github.com/dogenzaka/go-iap/appstore
go get github.com/dogenzaka/go-iap/playstore
```


# Quick Start

### In App Purchase (via App Store)

```
import(
    "github.com/dogenzaka/go-iap/appstore"
)

func main() {
	client := appstore.New()
	req := appstore.IAPRequest{
		ReceiptData: "your receipt data encoded by base64",
	}
	resp, err := client.Verify(req)
}
```

### In App Billing (via GooglePlay)

```
import(
    "golang.org/x/oauth2"

    "github.com/dogenzaka/go-iap/playstore"
)

func main() {
	// You need to prepare a public key for your Android app's in app billing
	// at https://console.developers.google.com.
	jsonKey, err := ioutil.ReadFile("jsonKey.json")
	if err != nil {
		log.Fatal(err)
	}

	client := playstore.New(jsonKey)
	resp, err := client.VerifySubscription("package", "subscriptionID", "purchaseToken")
}
```

### In App Billing (via GooglePlay/ Google App Engine)

```
import(
    "golang.org/x/oauth2"

    "github.com/dogenzaka/go-iap/playstore"
)

func foo(ctx context.Context) {
	client, _ := NewGAE(ctx)
	resp, err := client.VerifySubscription("package", "subscriptionID", "purchaseToken")
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


# License
go-iap is licensed under the MIT.
