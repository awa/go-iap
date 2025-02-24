package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Error struct {
	// Only errorCode and errorMessage are returned by App Store Server API.
	errorCode    int
	errorMessage string

	// retryAfter is the number of seconds after which the client can retry the request.
	// This field is only set to the `Retry-After` header if you receive the HTTP 429 error, that informs you when you can next send a request.
	retryAfter int64
}

func newError(errorCode int, errorMessage string) *Error {
	return &Error{
		errorCode:    errorCode,
		errorMessage: errorMessage,
	}
}

type appStoreAPIErrorResp struct {
	ErrorCode    int    `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func newAppStoreAPIError(b []byte, hd http.Header) (*Error, bool) {
	if len(b) == 0 {
		return nil, false
	}
	var rErr appStoreAPIErrorResp
	if err := json.Unmarshal(b, &rErr); err != nil {
		return nil, false
	}
	if rErr.ErrorCode == 0 {
		return nil, false
	}
	if rErr.ErrorCode == 4290000 {
		retryAfter, err := strconv.ParseInt(hd.Get("Retry-After"), 10, 64)
		if err == nil {
			return &Error{errorCode: rErr.ErrorCode, errorMessage: rErr.ErrorMessage, retryAfter: retryAfter}, true
		}
	}
	return &Error{errorCode: rErr.ErrorCode, errorMessage: rErr.ErrorMessage}, true
}

func newErrorFromJSON(b []byte) (*Error, bool) {
	if len(b) == 0 {
		return nil, false
	}
	var rErr appStoreAPIErrorResp
	if err := json.Unmarshal(b, &rErr); err != nil {
		return nil, false
	}
	if rErr.ErrorCode == 0 {
		return nil, false
	}
	return &Error{errorCode: rErr.ErrorCode, errorMessage: rErr.ErrorMessage}, true
}

func (e *Error) Error() string {
	return fmt.Sprintf("errorCode: %d, errorMessage: %s", e.errorCode, e.errorMessage)
}

func (e *Error) As(target interface{}) bool {
	if targetErr, ok := target.(*Error); ok {
		*targetErr = *e
		return true
	}
	return false
}

func (e *Error) Is(target error) bool {
	if other, ok := target.(*Error); ok && other.errorCode == e.errorCode {
		return true
	}
	return false
}

func (e *Error) ErrorCode() int {
	return e.errorCode
}

func (e *Error) ErrorMessage() string {
	return e.errorMessage
}

func (e *Error) RetryAfter() int64 {
	return e.retryAfter
}

func (e *Error) Retryable() bool {
	// NOTE:
	// RateLimitExceededError[1] could also be considered as a retryable error.
	// But limits are enforced on an hourly basis[2], so you should handle exceeded rate limits gracefully instead of retrying immediately.
	// Refs:
	// [1] https://developer.apple.com/documentation/appstoreserverapi/ratelimitexceedederror
	// [2] https://developer.apple.com/documentation/appstoreserverapi/identifying_rate_limits
	switch e.errorCode {
	case 4040002, 4040004, 5000001, 4040006:
		return true
	default:
		return false
	}
}

// All Error lists in https://developer.apple.com/documentation/appstoreserverapi/error_codes.
var (
	// Retryable errors
	AccountNotFoundRetryableError               = newError(4040002, "Account not found. Please try again.")
	AppNotFoundRetryableError                   = newError(4040004, "App not found. Please try again.")
	GeneralInternalRetryableError               = newError(5000001, "An unknown error occurred. Please try again.")
	OriginalTransactionIdNotFoundRetryableError = newError(4040006, "Original transaction id not found. Please try again.")
	// Errors
	AccountNotFoundError                             = newError(4040001, "Account not found.")
	AppNotFoundError                                 = newError(4040003, "App not found.")
	FamilySharedSubscriptionExtensionIneligibleError = newError(4030007, "Subscriptions that users obtain through Family Sharing can't get a renewal date extension directly.")
	GeneralInternalError                             = newError(5000000, "An unknown error occurred.")
	GeneralBadRequestError                           = newError(4000000, "Bad request.")
	InvalidAppIdentifierError                        = newError(4000002, "Invalid request app identifier.")
	InvalidEmptyStorefrontCountryCodeListError       = newError(4000027, "Invalid request. If provided, the list of storefront country codes must not be empty.")
	InvalidExtendByDaysError                         = newError(4000009, "Invalid extend by days value.")
	InvalidExtendReasonCodeError                     = newError(4000010, "Invalid extend reason code.")
	InvalidOriginalTransactionIdError                = newError(4000008, "Invalid original transaction id.")
	InvalidRequestIdentifierError                    = newError(4000011, "Invalid request identifier.")
	InvalidRequestRevisionError                      = newError(4000005, "Invalid request revision.")
	InvalidRevokedError                              = newError(4000030, "Invalid request. The revoked parameter is invalid.")
	InvalidStatusError                               = newError(4000031, "Invalid request. The status parameter is invalid.")
	InvalidStorefrontCountryCodeError                = newError(4000028, "Invalid request. A storefront country code was invalid.")
	InvalidTransactionIdError                        = newError(4000006, "Invalid transaction id.")
	OriginalTransactionIdNotFoundError               = newError(4040005, "Original transaction id not found.")
	RateLimitExceededError                           = newError(4290000, "Rate limit exceeded.")
	StatusRequestNotFoundError                       = newError(4040009, "The server didn't find a subscription-renewal-date extension request for this requestIdentifier and productId combination.")
	SubscriptionExtensionIneligibleError             = newError(4030004, "Forbidden - subscription state ineligible for extension.")
	SubscriptionMaxExtensionError                    = newError(4030005, "Forbidden - subscription has reached maximum extension count.")
	TransactionIdNotFoundError                       = newError(4040010, "Transaction id not found.")
	// Notification test and history errors
	InvalidEndDateError                     = newError(4000016, "Invalid request. The end date is not a timestamp value represented in milliseconds.")
	InvalidNotificationTypeError            = newError(4000018, "Invalid request. The notification type or subtype is invalid.")
	InvalidPaginationTokenError             = newError(4000014, "Invalid request. The pagination token is invalid.")
	InvalidStartDateError                   = newError(4000015, "Invalid request. The start date is not a timestamp value represented in milliseconds.")
	InvalidTestNotificationTokenError       = newError(4000020, "Invalid request. The test notification token is invalid.")
	InvalidInAppOwnershipTypeError          = newError(4000026, "Invalid request. The in-app ownership type parameter is invalid.")
	InvalidProductIdError                   = newError(4000023, "Invalid request. The product id parameter is invalid.")
	InvalidProductTypeError                 = newError(4000022, "Invalid request. The product type parameter is invalid.")
	InvalidSortError                        = newError(4000021, "Invalid request. The sort parameter is invalid.")
	InvalidSubscriptionGroupIdentifierError = newError(4000024, "Invalid request. The subscription group identifier parameter is invalid.")
	MultipleFiltersSuppliedError            = newError(4000019, "Invalid request. Supply either a transaction id or a notification type, but not both.")
	PaginationTokenExpiredError             = newError(4000017, "Invalid request. The pagination token is expired.")
	ServerNotificationURLNotFoundError      = newError(4040007, "No App Store Server Notification URL found for provided app. Check that a URL is configured in App Store Connect for this environment.")
	StartDateAfterEndDateError              = newError(4000013, "Invalid request. The end date precedes the start date or the dates are the same.")
	StartDateTooFarInPastError              = newError(4000012, "Invalid request. The start date is earlier than the allowed start date.")
	TestNotificationNotFoundError           = newError(4040008, "Either the test notification token is expired or the notification and status are not yet available.")
	InvalidAccountTenureError               = newError(4000032, "Invalid request. The account tenure field is invalid.")
	InvalidAppAccountTokenError             = newError(4000033, "Invalid request. The app account token field must contain a valid UUID or an empty string.")
	InvalidConsumptionStatusError           = newError(4000034, "Invalid request. The consumption status field is invalid.")
	InvalidCustomerConsentedError           = newError(4000035, "Invalid request. The customer consented field is required and must indicate the customer consented")
	InvalidDeliveryStatusError              = newError(4000036, "Invalid request. The delivery status field is invalid")
	InvalidLifetimeDollarsPurchasedError    = newError(4000037, "Invalid request. The lifetime dollars purchased field is invalid")
	InvalidLifetimeDollarsRefundedError     = newError(4000038, "Invalid request. The lifetime dollars refunded field is invalid")
	InvalidPlatformError                    = newError(4000039, "Invalid request. The platform field is invalid")
	InvalidPlayTimeError                    = newError(4000040, "Invalid request. The playtime field is invalid")
	InvalidSampleContentProvidedError       = newError(4000041, "Invalid request. The sample content provided field is invalid")
	InvalidUserStatusError                  = newError(4000042, "Invalid request. The user status field is invalid")
	InvalidTransactionNotConsumableError    = newError(4000043, "Invalid request. The transaction id parameter must represent a consumable in-app purchase")
	InvalidTransactionTypeNotSupportedError = newError(4000047, "Invalid request. The transaction id doesn't represent a supported in-app purchase type")
	AppTransactionIdNotSupportedError       = newError(4000048, "Invalid request. Invalid request. App transactions aren't supported by this endpoint")
)
