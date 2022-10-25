package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	HostSandBox    = "https://api.storekit-sandbox.itunes.apple.com"
	HostProduction = "https://api.storekit.itunes.apple.com"

	PathLookUp                        = "/inApps/v1/lookup/{orderId}"
	PathTransactionHistory            = "/inApps/v1/history/{originalTransactionId}"
	PathRefundHistory                 = "/inApps/v2/refund/lookup/{originalTransactionId}"
	PathGetALLSubscriptionStatus      = "/inApps/v1/subscriptions/{originalTransactionId}"
	PathConsumptionInfo               = "/inApps/v1/transactions/consumption/{originalTransactionId}"
	PathExtendSubscriptionRenewalDate = "/inApps/v1/subscriptions/extend/{originalTransactionId}"
	PathGetNotificationHistory        = "/inApps/v1/notifications/history"
)

type StoreConfig struct {
	KeyContent []byte // Loads a .p8 certificate
	KeyID      string // Your private key ID from App Store Connect (Ex: 2X9R4HXF34)
	BundleID   string // Your app’s bundle ID
	Issuer     string // Your issuer ID from the Keys page in App Store Connect (Ex: "57246542-96fe-1a63-e053-0824d011072a")
	Sandbox    bool   // default is Production
}

type StoreClient struct {
	Token *Token
	cert  *Cert
}

// NewStoreClient create appstore server api client
func NewStoreClient(config *StoreConfig) *StoreClient {
	token := &Token{}
	token.WithConfig(config)

	client := &StoreClient{
		Token: token,
		cert:  &Cert{},
	}
	return client
}

// GetALLSubscriptionStatuses https://developer.apple.com/documentation/appstoreserverapi/get_all_subscription_statuses
func (a *StoreClient) GetALLSubscriptionStatuses(originalTransactionId string) (rsp *StatusResponse, err error) {
	URL := HostProduction + PathGetALLSubscriptionStatus
	if a.Token.Sandbox {
		URL = HostSandBox + PathGetALLSubscriptionStatus
	}
	URL = strings.Replace(URL, "{originalTransactionId}", originalTransactionId, -1)
	statusCode, body, err := a.Do(http.MethodGet, URL, nil)
	if err != nil {
		return
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("appstore api: %v return status code %v", URL, statusCode)
	}

	err = json.Unmarshal(body, &rsp)
	if err != nil {
		return nil, err
	}

	return
}

// LookupOrderID https://developer.apple.com/documentation/appstoreserverapi/look_up_order_id
func (a *StoreClient) LookupOrderID(invoiceOrderId string) (rsp *OrderLookupResponse, err error) {
	URL := HostProduction + PathLookUp
	if a.Token.Sandbox {
		URL = HostSandBox + PathLookUp
	}
	URL = strings.Replace(URL, "{orderId}", invoiceOrderId, -1)
	statusCode, body, err := a.Do(http.MethodGet, URL, nil)
	if err != nil {
		return
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("appstore api: %v return status code %v", URL, statusCode)
	}

	err = json.Unmarshal(body, &rsp)
	if err != nil {
		return nil, err
	}

	return
}

// GetTransactionHistory https://developer.apple.com/documentation/appstoreserverapi/get_transaction_history
func (a *StoreClient) GetTransactionHistory(originalTransactionId string, query *url.Values) (responses []*HistoryResponse, err error) {
	URL := HostProduction + PathTransactionHistory
	if a.Token.Sandbox {
		URL = HostSandBox + PathTransactionHistory
	}
	URL = strings.Replace(URL, "{originalTransactionId}", originalTransactionId, -1)

	if query == nil {
		query = &url.Values{}
	}

	for {
		rsp := HistoryResponse{}

		statusCode, body, errOmit := a.Do(http.MethodGet, URL+"?"+query.Encode(), nil)
		if errOmit != nil {
			return nil, errOmit
		}

		if statusCode != http.StatusOK {
			return nil, fmt.Errorf("appstore api: %v return status code %v", URL, statusCode)
		}

		err = json.Unmarshal(body, &rsp)
		if err != nil {
			return nil, err
		}

		responses = append(responses, &rsp)
		if !rsp.HasMore {
			break
		}

		if rsp.HasMore && rsp.Revision != "" {
			query.Set("revision", rsp.Revision)
		}

		time.Sleep(10 * time.Millisecond)
	}

	return
}

// GetRefundHistory https://developer.apple.com/documentation/appstoreserverapi/get_refund_history
func (a *StoreClient) GetRefundHistory(originalTransactionId string) (responses []*RefundLookupResponse, err error) {
	baseURL := HostProduction + PathRefundHistory
	if a.Token.Sandbox {
		baseURL = HostSandBox + PathRefundHistory
	}
	baseURL = strings.Replace(baseURL, "{originalTransactionId}", originalTransactionId, -1)

	URL := baseURL
	for {
		rsp := RefundLookupResponse{}

		statusCode, body, errOmit := a.Do(http.MethodGet, URL, nil)
		if errOmit != nil {
			return nil, errOmit
		}

		if statusCode != http.StatusOK {
			return nil, fmt.Errorf("appstore api: %v return status code %v", URL, statusCode)
		}

		err = json.Unmarshal(body, &rsp)
		if err != nil {
			return nil, err
		}

		responses = append(responses, &rsp)
		if !rsp.HasMore {
			break
		}

		data := url.Values{}
		if rsp.HasMore && rsp.Revision != "" {
			data.Set("revision", rsp.Revision)
			URL = baseURL + "?" + data.Encode()
		}

		time.Sleep(10 * time.Millisecond)
	}
	return
}

// SendConsumptionInfo https://developer.apple.com/documentation/appstoreserverapi/send_consumption_information
func (a *StoreClient) SendConsumptionInfo(originalTransactionId string, body ConsumptionRequestBody) (statusCode int, err error) {
	URL := HostProduction + PathConsumptionInfo
	if a.Token.Sandbox {
		URL = HostSandBox + PathConsumptionInfo
	}
	URL = strings.Replace(URL, "{originalTransactionId}", originalTransactionId, -1)

	bodyBuf := new(bytes.Buffer)
	err = json.NewEncoder(bodyBuf).Encode(body)
	if err != nil {
		return 0, err
	}

	statusCode, _, err = a.Do(http.MethodPut, URL, bodyBuf)
	if err != nil {
		return statusCode, err
	}
	return statusCode, nil
}

// ExtendSubscriptionRenewalDate https://developer.apple.com/documentation/appstoreserverapi/extend_a_subscription_renewal_date
func (a *StoreClient) ExtendSubscriptionRenewalDate(originalTransactionId string, body ExtendRenewalDateRequest) (statusCode int, err error) {
	URL := HostProduction + PathExtendSubscriptionRenewalDate
	if a.Token.Sandbox {
		URL = HostSandBox + PathExtendSubscriptionRenewalDate
	}
	URL = strings.Replace(URL, "{originalTransactionId}", originalTransactionId, -1)

	bodyBuf := new(bytes.Buffer)
	err = json.NewEncoder(bodyBuf).Encode(body)
	if err != nil {
		return 0, err
	}

	statusCode, _, err = a.Do(http.MethodPut, URL, bodyBuf)
	if err != nil {
		return statusCode, err
	}
	return statusCode, nil
}

// GetNotificationHistory https://developer.apple.com/documentation/appstoreserverapi/get_notification_history
// Note: Notification history is available starting on June 6, 2022. Use a startDate of June 6, 2022 or later in your request.
func (a *StoreClient) GetNotificationHistory(body NotificationHistoryRequest) (responses []NotificationHistoryResponseItem, err error) {
	baseURL := HostProduction + PathGetNotificationHistory
	if a.Token.Sandbox {
		baseURL = HostSandBox + PathGetNotificationHistory
	}

	bodyBuf := new(bytes.Buffer)
	err = json.NewEncoder(bodyBuf).Encode(body)
	if err != nil {
		return nil, err
	}

	URL := baseURL
	for {
		rsp := NotificationHistoryResponses{}
		rsp.NotificationHistory = make([]NotificationHistoryResponseItem, 0)

		statusCode, rspBody, errOmit := a.Do(http.MethodPost, URL, bodyBuf)
		if errOmit != nil {
			return nil, errOmit
		}

		if statusCode != http.StatusOK {
			return nil, fmt.Errorf("appstore api: %v return status code %v", URL, statusCode)
		}

		err = json.Unmarshal(rspBody, &rsp)
		if err != nil {
			return nil, err
		}

		responses = append(responses, rsp.NotificationHistory...)
		if !rsp.HasMore {
			break
		}

		data := url.Values{}
		if rsp.HasMore && rsp.PaginationToken != "" {
			data.Set("paginationToken", rsp.PaginationToken)
			URL = baseURL + "?" + data.Encode()
		}

		time.Sleep(10 * time.Millisecond)
	}

	return responses, nil
}

// ParseSignedTransactions parse the jws singed transactions
// Per doc: https://datatracker.ietf.org/doc/html/rfc7515#section-4.1.6
func (a *StoreClient) ParseSignedTransactions(transactions []string) ([]*JWSTransaction, error) {
	result := make([]*JWSTransaction, 0)
	for _, v := range transactions {
		trans, err := a.parseSignedTransaction(v)
		if err == nil && trans != nil {
			result = append(result, trans)
		}
	}

	return result, nil
}

func (a *StoreClient) parseSignedTransaction(transaction string) (*JWSTransaction, error) {
	tran := &JWSTransaction{}

	rootCertStr, err := a.cert.extractCertByIndex(transaction, 2)
	if err != nil {
		return nil, err
	}
	intermediaCertStr, err := a.cert.extractCertByIndex(transaction, 1)
	if err != nil {
		return nil, err
	}
	if err = a.cert.verifyCert(rootCertStr, intermediaCertStr); err != nil {
		return nil, err
	}

	_, err = jwt.ParseWithClaims(transaction, tran, func(token *jwt.Token) (interface{}, error) {
		return a.cert.extractPublicKeyFromToken(transaction)
	})
	if err != nil {
		return nil, err
	}

	return tran, nil
}

// Per doc: https://developer.apple.com/documentation/appstoreserverapi#topics
func (a *StoreClient) Do(method string, url string, body io.Reader) (int, []byte, error) {
	authToken, err := a.Token.GenerateIfExpired()
	if err != nil {
		return 0, nil, fmt.Errorf("appstore generate token err %w", err)
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return 0, nil, fmt.Errorf("appstore new http request err %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("User-Agent", "App Store Client")

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("appstore http client do err %w", err)
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, fmt.Errorf("appstore read http body err %w", err)
	}

	return resp.StatusCode, bytes, err
}
