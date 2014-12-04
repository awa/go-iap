go-iap
======

go-iap verifies the purchase receipt via AppStore or GooglePlayStore


# Installation
```
go get github.com/dogenzaka/go-iap
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
	resp, err := client.Verify(&req)
}
```

# ToDo
- [x] App Store Client
- [ ] Google Play Store Client


# Support
iOS7 or above


# License
Gorv is licensed under the MIT.
