package api

// GetSeatInfoResponse https://developer.apple.com/documentation/appstoreserverapi/getseatinforesponse
type GetSeatInfoResponse struct {
	AvailableSeatCount   int32    `json:"availableSeatCount"`
	TotalSeatCount       int32    `json:"totalSeatCount"`
	AssignedSeatClaimIds []string `json:"assignedSeatClaimIds"` // UUID strings.
}

// ClaimSeatRequest https://developer.apple.com/documentation/appstoreserverapi/claimseatrequest
type ClaimSeatRequest struct {
	RequestIdentifier string `json:"requestIdentifier"` // Required, at most 128 characters; reuse when retrying a claim.
	AppTransactionId  string `json:"appTransactionId,omitempty"`
}

// ClaimSeatResponse https://developer.apple.com/documentation/appstoreserverapi/claimseatresponse
type ClaimSeatResponse struct {
	SeatClaimId string `json:"seatClaimId"` // UUID string.
}

// AssociateSeatRequest https://developer.apple.com/documentation/appstoreserverapi/associateseatrequest
type AssociateSeatRequest struct {
	AppTransactionId string `json:"appTransactionId"`
}

// InvitationManagement https://developer.apple.com/documentation/appstoreserverapi/invitationmanagement
type InvitationManagement string

const (
	InvitationManagementApple     InvitationManagement = "APPLE"
	InvitationManagementDeveloper InvitationManagement = "DEVELOPER"
)

// AdditionalTransactionReason https://developer.apple.com/documentation/appstoreserverapi/additionaltransactionreason
type AdditionalTransactionReason string

const AdditionalTransactionReasonSeatAddition AdditionalTransactionReason = "SEAT_ADDITION"

// Offer types are int32 values in both transaction and renewal payloads.
// https://developer.apple.com/documentation/appstoreserverapi/offertype
const (
	OfferTypeIntroductory int32 = 1
	OfferTypePromotional  int32 = 2
	OfferTypeCode         int32 = 3
	OfferTypeWinBack      int32 = 4
	OfferTypeRetention    int32 = 5
)

// TransactionMultiSeatInfo https://developer.apple.com/documentation/appstoreserverapi/transactionmultiseatinfo
type TransactionMultiSeatInfo struct {
	InvitationManagement InvitationManagement `json:"invitationManagement,omitempty"`
	PriceTiers           []MultiSeatPriceTier `json:"priceTiers,omitempty"`
	RevokedSeatQuantity  int32                `json:"revokedSeatQuantity,omitzero"`
}

// RenewalMultiSeatInfo https://developer.apple.com/documentation/appstoreserverapi/renewalmultiseatinfo
type RenewalMultiSeatInfo struct {
	RenewalPriceTiers []MultiSeatPriceTier `json:"renewalPriceTiers,omitempty"`
}

// MultiSeatPriceTier https://developer.apple.com/documentation/appstoreserverapi/multiseatpricetier
type MultiSeatPriceTier struct {
	Price    int64 `json:"price"` // Per-seat price in milliunits of the currency.
	Quantity int32 `json:"quantity"`
}

// AdvancedCommerceBundleItem https://developer.apple.com/documentation/appstoreserverapi/advancedcommercebundleitem
type AdvancedCommerceBundleItem struct {
	SKU         string `json:"SKU,omitempty"`
	Description string `json:"description,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
	Price       int64  `json:"price"` // Milliunits of the currency.
	TaxCode     string `json:"taxCode,omitempty"`
}

// AdvancedCommerceSeatPriceTier https://developer.apple.com/documentation/appstoreserverapi/advancedcommerceseatpricetier
type AdvancedCommerceSeatPriceTier struct {
	// The number of seats in this tier, not a cumulative upper bound.
	// Absent for the final tier, which has no upper limit.
	MaximumSeatQuantity *int32 `json:"maximumSeatQuantity,omitempty"`
	Price               int64  `json:"price"` // Per-seat price in milliunits of the currency.
}

// AdvancedCommercePartner https://developer.apple.com/documentation/appstoreserverapi/advancedcommercepartner
type AdvancedCommercePartner struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
