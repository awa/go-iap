package storekit

import (
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"net/url"
)

// HistoryResponse https://developer.apple.com/documentation/appstoreserverapi/historyresponse
type HistoryResponse struct {
	AppAppleId         int64    `json:"appAppleId"`         // The app’s identifier in the App Store.
	BundleId           string   `json:"bundleId"`           // The bundle identifier of the app.
	Environment        string   `json:"environment"`        // The server environment in which you’re making the request, sandbox or production.
	HasMore            bool     `json:"hasMore"`            // A Boolean value that indicates whether App Store has more transactions than are returned in this request.
	Revision           string   `json:"revision"`           // A token you use in a query to request the next set transactions from the Get Transaction History endpoint.
	SignedTransactions []string `json:"signedTransactions"` // An array of in-app purchase transactions for the customer, signed by Apple, in JSON Web Signature format.
	DecodeTransactions []JWSTransactionDecodedPayload
}

// JWSTransactionDecodedPayload https://developer.apple.com/documentation/appstoreserverapi/jwstransactiondecodedpayload
type JWSTransactionDecodedPayload struct {
	jwt.StandardClaims
	BundleId                    string `json:"bundleId"`
	ExpiresDate                 int64  `json:"expiresDate"`
	InAppOwnershipType          string `json:"inAppOwnershipType"`
	OfferType                   int    `json:"offerType"`
	OriginalPurchaseDate        int64  `json:"originalPurchaseDate"`
	OriginalTransactionId       string `json:"originalTransactionId"`
	ProductId                   string `json:"productId"`
	PurchaseDate                int64  `json:"purchaseDate"`
	Quantity                    int    `json:"quantity"`
	SignedDate                  int64  `json:"signedDate"`
	SubscriptionGroupIdentifier string `json:"subscriptionGroupIdentifier"`
	TransactionId               string `json:"transactionId"`
	Type                        string `json:"type"`
	WebOrderLineItemId          string `json:"webOrderLineItemId"`
}

// GetTransactionHistory Get a customer’s transaction history, including all of their in-app purchases in your app.
func (client *Client) GetTransactionHistory(originalTransactionId, revision string) (resp *HistoryResponse, err error) {
	u := ProductionURL + "/history/" + originalTransactionId
	if client.Sandbox {
		u = SandboxURL + "/history/" + originalTransactionId
	}
	// add query
	query := url.Values{}
	if len(revision) > 0 {
		query.Set("revision", revision)
	}

	req, _ := http.NewRequest(http.MethodGet, u+"?"+query.Encode(), nil)
	resp = new(HistoryResponse)
	if err = client.Do(req, resp); err != nil {
		return
	}
	// decode transaction
	resp.DecodeTransactions = make([]JWSTransactionDecodedPayload, 0, len(resp.SignedTransactions))
	for _, raw := range resp.SignedTransactions {
		var tmp = new(JWSTransactionDecodedPayload)
		_, _ = jwt.ParseWithClaims(raw, tmp, func(token *jwt.Token) (interface{}, error) {
			return client.PrivateKey, nil
		})
		resp.DecodeTransactions = append(resp.DecodeTransactions, *tmp)
	}

	return
}
