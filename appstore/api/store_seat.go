//go:generate mockgen -destination=../mocks/seat.go -package=mocks github.com/awa/go-iap/appstore/api SeatManager

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"unicode/utf8"
)

// GetSeatInfo returns seat counts and active claim UUIDs for a subscription.
// This endpoint is currently Sandbox-only; configure StoreConfig.Sandbox accordingly.
// https://developer.apple.com/documentation/appstoreserverapi/get-seat-info
func (a *StoreClient) GetSeatInfo(ctx context.Context, originalTransactionId string) (*GetSeatInfoResponse, error) {
	if originalTransactionId == "" {
		return nil, fmt.Errorf("appstore GetSeatInfo: originalTransactionId is required")
	}
	URL := a.host + strings.ReplaceAll(PathGetSeatInfo, "{originalTransactionId}", url.PathEscape(originalTransactionId))
	status, body, err := a.Do(ctx, http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, newHTTPStatusError(status, URL)
	}
	var rsp *GetSeatInfoResponse
	if err := json.Unmarshal(body, &rsp); err != nil {
		return nil, fmt.Errorf("appstore GetSeatInfo: decode response: %w", err)
	}
	if rsp == nil {
		return nil, fmt.Errorf("appstore GetSeatInfo: expected a response object, got null")
	}
	return rsp, nil
}

// ClaimSeat claims a seat, optionally associating it with a customer's app transaction.
// Save and reuse body.RequestIdentifier when retrying to avoid claiming another seat.
// This endpoint is currently Sandbox-only; configure StoreConfig.Sandbox accordingly.
// https://developer.apple.com/documentation/appstoreserverapi/claim-seat
func (a *StoreClient) ClaimSeat(ctx context.Context, originalTransactionId string, body ClaimSeatRequest) (*ClaimSeatResponse, error) {
	if originalTransactionId == "" {
		return nil, fmt.Errorf("appstore ClaimSeat: originalTransactionId is required")
	}
	if body.RequestIdentifier == "" || !utf8.ValidString(body.RequestIdentifier) || utf8.RuneCountInString(body.RequestIdentifier) > 128 {
		return nil, fmt.Errorf("appstore ClaimSeat: requestIdentifier must contain 1 to 128 valid UTF-8 characters")
	}
	URL := a.host + strings.ReplaceAll(PathClaimSeat, "{originalTransactionId}", url.PathEscape(originalTransactionId))
	requestBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("appstore ClaimSeat: encode request: %w", err)
	}
	status, responseBody, err := a.Do(ctx, http.MethodPost, URL, bytes.NewReader(requestBody))
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, newHTTPStatusError(status, URL)
	}
	var rsp *ClaimSeatResponse
	if err := json.Unmarshal(responseBody, &rsp); err != nil {
		return nil, fmt.Errorf("appstore ClaimSeat: decode response: %w", err)
	}
	if rsp == nil {
		return nil, fmt.Errorf("appstore ClaimSeat: expected a response object, got null")
	}
	return rsp, nil
}

// UnclaimSeat releases a seat claim. The caller must revoke access in its own system.
// This endpoint is currently Sandbox-only; configure StoreConfig.Sandbox accordingly.
// https://developer.apple.com/documentation/appstoreserverapi/unclaim-seat
func (a *StoreClient) UnclaimSeat(ctx context.Context, originalTransactionId, seatClaimId string) (int, error) {
	if originalTransactionId == "" {
		return 0, fmt.Errorf("appstore UnclaimSeat: originalTransactionId is required")
	}
	if seatClaimId == "" {
		return 0, fmt.Errorf("appstore UnclaimSeat: seatClaimId is required")
	}
	URL := a.host + strings.ReplaceAll(PathUnclaimSeat, "{originalTransactionId}", url.PathEscape(originalTransactionId))
	URL = strings.ReplaceAll(URL, "{seatClaimId}", url.PathEscape(seatClaimId))
	status, _, err := a.Do(ctx, http.MethodDelete, URL, nil)
	if err != nil {
		return status, err
	}
	if status != http.StatusOK {
		return status, newHTTPStatusError(status, URL)
	}
	return status, nil
}

// AssociateSeat associates a customer with an existing seat claim; it doesn't grant access.
// This endpoint is currently Sandbox-only; configure StoreConfig.Sandbox accordingly.
// https://developer.apple.com/documentation/appstoreserverapi/associate-seat
func (a *StoreClient) AssociateSeat(ctx context.Context, originalTransactionId, seatClaimId string, body AssociateSeatRequest) (int, error) {
	if originalTransactionId == "" {
		return 0, fmt.Errorf("appstore AssociateSeat: originalTransactionId is required")
	}
	if seatClaimId == "" {
		return 0, fmt.Errorf("appstore AssociateSeat: seatClaimId is required")
	}
	if body.AppTransactionId == "" {
		return 0, fmt.Errorf("appstore AssociateSeat: appTransactionId is required")
	}
	URL := a.host + strings.ReplaceAll(PathAssociateSeat, "{originalTransactionId}", url.PathEscape(originalTransactionId))
	URL = strings.ReplaceAll(URL, "{seatClaimId}", url.PathEscape(seatClaimId))
	requestBody, err := json.Marshal(body)
	if err != nil {
		return 0, fmt.Errorf("appstore AssociateSeat: encode request: %w", err)
	}
	status, _, err := a.Do(ctx, http.MethodPost, URL, bytes.NewReader(requestBody))
	if err != nil {
		return status, err
	}
	if status != http.StatusOK {
		return status, newHTTPStatusError(status, URL)
	}
	return status, nil
}
