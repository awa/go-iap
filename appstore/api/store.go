//go:generate mockgen -destination=../mocks/store.go -package=mocks github.com/awa/go-iap/appstore/api StoreAPIClient

package api

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	HostSandBox    = "https://api.storekit-sandbox.itunes.apple.com"
	HostProduction = "https://api.storekit.itunes.apple.com"

	PathLookUp                              = "/inApps/v1/lookup/{orderId}"
	PathTransactionHistory                  = "/inApps/v1/history/{originalTransactionId}"
	PathTransactionInfo                     = "/inApps/v1/transactions/{transactionId}"
	PathRefundHistory                       = "/inApps/v2/refund/lookup/{originalTransactionId}"
	PathGetALLSubscriptionStatus            = "/inApps/v1/subscriptions/{originalTransactionId}"
	PathConsumptionInfo                     = "/inApps/v1/transactions/consumption/{originalTransactionId}"
	PathExtendSubscriptionRenewalDate       = "/inApps/v1/subscriptions/extend/{originalTransactionId}"
	PathExtendSubscriptionRenewalDateForAll = "/inApps/v1/subscriptions/extend/mass/"
	PathGetStatusOfSubscriptionRenewalDate  = "/inApps/v1/subscriptions/extend/mass/{productId}/{requestIdentifier}"
	PathGetNotificationHistory              = "/inApps/v1/notifications/history"
	PathRequestTestNotification             = "/inApps/v1/notifications/test"
	PathGetTestNotificationStatus           = "/inApps/v1/notifications/test/{testNotificationToken}"
)

type StoreConfig struct {
	KeyContent []byte // Loads a .p8 certificate
	KeyID      string // Your private key ID from App Store Connect (Ex: 2X9R4HXF34)
	BundleID   string // Your app’s bundle ID
	Issuer     string // Your issuer ID from the Keys page in App Store Connect (Ex: "57246542-96fe-1a63-e053-0824d011072a")
	Sandbox    bool   // default is Production
}

type (
	StoreAPIClient interface {
		SubscriptionGetter
		SubscriptionExtender
		TransactionParser
		NotificationGetter
		NotificationSender
		ConsumptionSender
		Do(ctx context.Context, method string, url string, body io.Reader) (int, []byte, error)
	}

	SubscriptionGetter interface {
		GetALLSubscriptionStatuses(ctx context.Context, originalTransactionId string, query *url.Values) (rsp *StatusResponse, err error)
		GetRefundHistory(ctx context.Context, originalTransactionId string) (responses []*RefundLookupResponse, err error)
		GetSubscriptionRenewalDataStatus(ctx context.Context, productId, requestIdentifier string) (statusCode int, rsp *MassExtendRenewalDateStatusResponse, err error)
		GetTransactionHistory(ctx context.Context, originalTransactionId string, query *url.Values) (responses []*HistoryResponse, err error)
		GetTransactionInfo(ctx context.Context, transactionId string) (rsp *TransactionInfoResponse, err error)
		LookupOrderID(ctx context.Context, orderId string) (rsp *OrderLookupResponse, err error)
	}

	SubscriptionExtender interface {
		ExtendSubscriptionRenewalDate(ctx context.Context, originalTransactionId string, body ExtendRenewalDateRequest) (statusCode int, err error)
		ExtendSubscriptionRenewalDateForAll(ctx context.Context, body MassExtendRenewalDateRequest) (statusCode int, err error)
	}

	TransactionParser interface {
		ParseSignedTransactions(transactions []string) ([]*JWSTransaction, error)
		ParseJWSEncodeString(jwsEncode string) (interface{}, error)
		ParseSignedTransaction(transaction string) (*JWSTransaction, error)
	}

	NotificationGetter interface {
		GetAllNotificationHistory(ctx context.Context, body NotificationHistoryRequest, duration time.Duration) (responses []NotificationHistoryResponseItem, err error)
		GetNotificationHistory(ctx context.Context, body NotificationHistoryRequest, paginationToken string) (rsp *NotificationHistoryResponses, err error)
		GetTestNotificationStatus(ctx context.Context, testNotificationToken string) (int, []byte, error)
	}

	NotificationSender interface {
		SendRequestTestNotification(ctx context.Context) (int, []byte, error)
	}

	ConsumptionSender interface {
		SendConsumptionInfo(ctx context.Context, originalTransactionId string, body ConsumptionRequestBody) (statusCode int, err error)
	}
)

type StoreClient struct {
	Token   *Token
	httpCli *http.Client
	cert    *Cert
}

// NewStoreClient create a appstore server api client
func NewStoreClient(config *StoreConfig) *StoreClient {
	token := &Token{}
	token.WithConfig(config)

	client := &StoreClient{
		Token: token,
		cert:  &Cert{},
		httpCli: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
	return client
}

// NewStoreClientWithHTTPClient NewWithClient creates an App Store server api client with a custom http client.
func NewStoreClientWithHTTPClient(config *StoreConfig, httpClient *http.Client) *StoreClient {
	token := &Token{}
	token.WithConfig(config)

	client := &StoreClient{
		Token:   token,
		cert:    &Cert{},
		httpCli: httpClient,
	}
	return client
}

// GetALLSubscriptionStatuses https://developer.apple.com/documentation/appstoreserverapi/get_all_subscription_statuses
func (a *StoreClient) GetALLSubscriptionStatuses(ctx context.Context, originalTransactionId string, query *url.Values) (rsp *StatusResponse, err error) {
	URL := HostProduction + PathGetALLSubscriptionStatus
	if a.Token.Sandbox {
		URL = HostSandBox + PathGetALLSubscriptionStatus
	}
	URL = strings.Replace(URL, "{originalTransactionId}", originalTransactionId, -1)
	if query != nil {
		URL = URL + "?" + query.Encode()
	}
	statusCode, body, err := a.Do(ctx, http.MethodGet, URL, nil)
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
func (a *StoreClient) LookupOrderID(ctx context.Context, orderId string) (rsp *OrderLookupResponse, err error) {
	URL := HostProduction + PathLookUp
	if a.Token.Sandbox {
		URL = HostSandBox + PathLookUp
	}
	URL = strings.Replace(URL, "{orderId}", orderId, -1)
	statusCode, body, err := a.Do(ctx, http.MethodGet, URL, nil)
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
func (a *StoreClient) GetTransactionHistory(ctx context.Context, originalTransactionId string, query *url.Values) (responses []*HistoryResponse, err error) {
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

		statusCode, body, errOmit := a.Do(ctx, http.MethodGet, URL+"?"+query.Encode(), nil)
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

// GetTransactionInfo https://developer.apple.com/documentation/appstoreserverapi/get_transaction_info
func (a *StoreClient) GetTransactionInfo(ctx context.Context, transactionId string) (rsp *TransactionInfoResponse, err error) {
	URL := HostProduction + PathTransactionInfo
	if a.Token.Sandbox {
		URL = HostSandBox + PathTransactionInfo
	}
	URL = strings.Replace(URL, "{transactionId}", transactionId, -1)

	statusCode, body, err := a.Do(ctx, http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
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

// GetRefundHistory https://developer.apple.com/documentation/appstoreserverapi/get_refund_history
func (a *StoreClient) GetRefundHistory(ctx context.Context, originalTransactionId string) (responses []*RefundLookupResponse, err error) {
	baseURL := HostProduction + PathRefundHistory
	if a.Token.Sandbox {
		baseURL = HostSandBox + PathRefundHistory
	}
	baseURL = strings.Replace(baseURL, "{originalTransactionId}", originalTransactionId, -1)

	URL := baseURL
	for {
		rsp := RefundLookupResponse{}

		statusCode, body, errOmit := a.Do(ctx, http.MethodGet, URL, nil)
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
func (a *StoreClient) SendConsumptionInfo(ctx context.Context, originalTransactionId string, body ConsumptionRequestBody) (statusCode int, err error) {
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

	statusCode, _, err = a.Do(ctx, http.MethodPut, URL, bodyBuf)
	if err != nil {
		return statusCode, err
	}
	return statusCode, nil
}

// ExtendSubscriptionRenewalDate https://developer.apple.com/documentation/appstoreserverapi/extend_a_subscription_renewal_date
func (a *StoreClient) ExtendSubscriptionRenewalDate(ctx context.Context, originalTransactionId string, body ExtendRenewalDateRequest) (statusCode int, err error) {
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

	statusCode, _, err = a.Do(ctx, http.MethodPut, URL, bodyBuf)
	if err != nil {
		return statusCode, err
	}
	return statusCode, nil
}

// ExtendSubscriptionRenewalDateForAll https://developer.apple.com/documentation/appstoreserverapi/extend_subscription_renewal_dates_for_all_active_subscribers
func (a *StoreClient) ExtendSubscriptionRenewalDateForAll(ctx context.Context, body MassExtendRenewalDateRequest) (statusCode int, err error) {
	URL := HostProduction + PathExtendSubscriptionRenewalDateForAll
	if a.Token.Sandbox {
		URL = HostSandBox + PathExtendSubscriptionRenewalDateForAll
	}

	bodyBuf := new(bytes.Buffer)
	err = json.NewEncoder(bodyBuf).Encode(body)
	if err != nil {
		return 0, err
	}

	statusCode, _, err = a.Do(ctx, http.MethodPost, URL, bodyBuf)
	if err != nil {
		return statusCode, err
	}
	return statusCode, nil
}

// GetSubscriptionRenewalDataStatus https://developer.apple.com/documentation/appstoreserverapi/get_status_of_subscription_renewal_date_extensions
func (a *StoreClient) GetSubscriptionRenewalDataStatus(ctx context.Context, productId, requestIdentifier string) (statusCode int, rsp *MassExtendRenewalDateStatusResponse, err error) {
	URL := HostProduction + PathGetStatusOfSubscriptionRenewalDate
	if a.Token.Sandbox {
		URL = HostSandBox + PathGetStatusOfSubscriptionRenewalDate
	}
	URL = strings.Replace(URL, "{productId}", productId, -1)
	URL = strings.Replace(URL, "{requestIdentifier}", requestIdentifier, -1)

	statusCode, body, err := a.Do(ctx, http.MethodGet, URL, nil)
	if err != nil {
		return statusCode, nil, err
	}

	if statusCode != http.StatusOK {
		return statusCode, nil, fmt.Errorf("appstore api: %v return status code %v", URL, statusCode)
	}

	err = json.Unmarshal(body, &rsp)
	if err != nil {
		return statusCode, nil, err
	}

	return statusCode, rsp, nil
}

// GetAllNotificationHistory returns all the NotificationHistoryResponseItem using the paginationToken on behalf of you.
func (a *StoreClient) GetAllNotificationHistory(ctx context.Context, body NotificationHistoryRequest, duration time.Duration) (responses []NotificationHistoryResponseItem, err error) {
	paginationToken := ""
	for {
		rsp, err := a.GetNotificationHistory(ctx, body, paginationToken)
		if err != nil {
			return nil, err
		}

		responses = append(responses, rsp.NotificationHistory...)

		if rsp.HasMore {
			paginationToken = rsp.PaginationToken
		} else {
			break
		}

		time.Sleep(duration)
	}

	return responses, nil
}

// GetNotificationHistory https://developer.apple.com/documentation/appstoreserverapi/get_notification_history
// Note: Notification history is available starting on June 6, 2022. Use a startDate of June 6, 2022 or later in your request.
func (a *StoreClient) GetNotificationHistory(ctx context.Context, body NotificationHistoryRequest, paginationToken string) (rsp *NotificationHistoryResponses, err error) {
	baseURL := HostProduction + PathGetNotificationHistory
	if a.Token.Sandbox {
		baseURL = HostSandBox + PathGetNotificationHistory
	}

	URL := baseURL
	if paginationToken != "" {
		query := url.Values{}
		query.Set("paginationToken", paginationToken)
		URL += "?" + query.Encode()
	}

	bodyBuf := new(bytes.Buffer)
	err = json.NewEncoder(bodyBuf).Encode(body)
	if err != nil {
		return nil, err
	}

	statusCode, rspBody, err := a.Do(ctx, http.MethodPost, URL, bodyBuf)
	if err != nil {
		return nil, err
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("appstore api: %v return status code %v", URL, statusCode)
	}

	if err = json.Unmarshal(rspBody, &rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

// SendRequestTestNotification https://developer.apple.com/documentation/appstoreserverapi/request_a_test_notification
func (a *StoreClient) SendRequestTestNotification(ctx context.Context) (int, []byte, error) {
	URL := HostProduction + PathRequestTestNotification
	if a.Token.Sandbox {
		URL = HostSandBox + PathRequestTestNotification
	}

	return a.Do(ctx, http.MethodPost, URL, nil)
}

// GetTestNotificationStatus https://developer.apple.com/documentation/appstoreserverapi/get_test_notification_status
func (a *StoreClient) GetTestNotificationStatus(ctx context.Context, testNotificationToken string) (int, []byte, error) {
	URL := HostProduction + PathGetTestNotificationStatus
	if a.Token.Sandbox {
		URL = HostSandBox + PathGetTestNotificationStatus
	}
	URL = strings.Replace(URL, "{testNotificationToken}", testNotificationToken, -1)

	return a.Do(ctx, http.MethodGet, URL, nil)
}

// ParseSignedTransactions parse the jws singed transactions
// Per doc: https://datatracker.ietf.org/doc/html/rfc7515#section-4.1.6
func (a *StoreClient) ParseSignedTransactions(transactions []string) ([]*JWSTransaction, error) {
	result := make([]*JWSTransaction, 0)
	for _, v := range transactions {
		trans, err := a.ParseSignedTransaction(v)
		if err == nil && trans != nil {
			result = append(result, trans)
		}
	}

	return result, nil
}

// ParseJWSEncodeString parse the jws encode string, such as JWSTransaction and JWSRenewalInfoDecodedPayload
func (a *StoreClient) ParseJWSEncodeString(jwsEncode string) (interface{}, error) {
	// Split the JWS format string into its three parts
	parts := strings.Split(jwsEncode, ".")

	// Decode the payload part of the JWS format string
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}

	// Determine which struct to use based on the payload contents
	if strings.Contains(string(payload), "transactionId") {
		transaction := &JWSTransaction{}
		err = a.parseJWS(jwsEncode, transaction)
		return transaction, err
	} else if strings.Contains(string(payload), "renewalDate") {
		renewalInfo := &JWSRenewalInfoDecodedPayload{}
		err = a.parseJWS(jwsEncode, renewalInfo)
		return renewalInfo, err
	}

	return nil, nil
}

func (a *StoreClient) parseJWS(jwsEncode string, claims jwt.Claims) error {
	rootCertBytes, err := a.cert.extractCertByIndex(jwsEncode, 2)
	if err != nil {
		return err
	}
	rootCert, err := x509.ParseCertificate(rootCertBytes)
	if err != nil {
		return fmt.Errorf("appstore failed to parse root certificate")
	}

	intermediaCertBytes, err := a.cert.extractCertByIndex(jwsEncode, 1)
	if err != nil {
		return err
	}
	intermediaCert, err := x509.ParseCertificate(intermediaCertBytes)
	if err != nil {
		return fmt.Errorf("appstore failed to parse intermediate certificate")
	}

	leafCertBytes, err := a.cert.extractCertByIndex(jwsEncode, 0)
	if err != nil {
		return err
	}
	leafCert, err := x509.ParseCertificate(leafCertBytes)
	if err != nil {
		return fmt.Errorf("appstore failed to parse leaf certificate")
	}
	if err = a.cert.verifyCert(rootCert, intermediaCert, leafCert); err != nil {
		return err
	}

	pk, ok := leafCert.PublicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("appstore public key must be of type ecdsa.PublicKey")
	}

	_, err = jwt.ParseWithClaims(jwsEncode, claims, func(token *jwt.Token) (interface{}, error) {
		return pk, nil
	})
	return err
}

// ParseSignedTransaction parse one jws singed transaction for API like GetTransactionInfo
func (a *StoreClient) ParseSignedTransaction(transaction string) (*JWSTransaction, error) {
	tran := &JWSTransaction{}

	err := a.parseJWS(transaction, tran)
	if err != nil {
		return nil, err
	}

	return tran, nil
}

// Do Per doc: https://developer.apple.com/documentation/appstoreserverapi#topics
func (a *StoreClient) Do(ctx context.Context, method string, url string, body io.Reader) (int, []byte, error) {
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
	req = req.WithContext(ctx)

	resp, err := a.httpCli.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("appstore http client do err %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, nil, fmt.Errorf("appstore read http body err %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if rErr, ok := newAppStoreAPIError(bodyBytes, resp.Header); ok {
			return resp.StatusCode, bodyBytes, rErr
		}
	}

	return resp.StatusCode, bodyBytes, err
}
