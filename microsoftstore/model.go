package microsoftstore

import "time"

type UserIdentity struct {
	IdentityType         string `json:"identityType"`
	IdentityValue        string `json:"identityValue"`
	LocalTicketReference string `json:"localTicketReference"`
}

type ProductSkuId struct {
	ProductId string `json:"productId"`
	SkuId     string `json:"skuId"`
}

type ProductType string

const (
	Application         ProductType = "Application"
	Durable             ProductType = "Durable"
	Game                ProductType = "Game"
	UnmanagedConsumable ProductType = "UnmanagedConsumable"
)

type IAPRequest struct {
	Beneficiaries     []UserIdentity `json:"beneficiaries"`
	ContinuationToken string         `json:"continuationToken,omitempty"`
	MaxPageSize       int            `json:"maxPageSize,omitempty"`
	ModifiedAfter     *time.Time     `json:"modifiedAfter,omitempty"`
	ParentProductId   string         `json:"parentProductId,omitempty"`
	ProductSkuIds     []ProductSkuId `json:"productSkuIds,omitempty"`
	ProductTypes      []ProductType  `json:"productTypes"`
	ValidityType      string         `json:"validityType,omitempty"`
}

type IdentityContractV6 struct {
	IdentityType  string `json:"identityType"`  // Contains the value "pub".
	IdentityValue string `json:"identityValue"` // The string value of the publisherUserId from the specified Microsoft Store ID key.
}

// CollectionItemContractV6 represents an item in the user's collection.
type CollectionItemContractV6 struct {
	AcquiredDate         time.Time           `json:"acquiredDate"`               // The date on which the user acquired the item.
	CampaignId           *string             `json:"campaignId,omitempty"`       // The campaign ID that was provided at purchase time for this item.
	DevOfferId           *string             `json:"devOfferId,omitempty"`       // The offer ID from an in-app purchase.
	EndDate              time.Time           `json:"endDate"`                    // The end date of the item.
	FulfillmentData      []string            `json:"fulfillmentData,omitempty"`  // N/A
	InAppOfferToken      *string             `json:"inAppOfferToken,omitempty"`  // The developer-specified product ID string assigned to the item in Partner Center.
	ItemId               string              `json:"itemId"`                     // An ID that identifies this collection item from other items the user owns.
	LocalTicketReference string              `json:"localTicketReference"`       // The ID of the previously supplied localTicketReference in the request body.
	ModifiedDate         time.Time           `json:"modifiedDate"`               // The date this item was last modified.
	OrderId              *string             `json:"orderId,omitempty"`          // If present, the order ID of which this item was obtained.
	OrderLineItemId      *string             `json:"orderLineItemId,omitempty"`  // If present, the line item of the particular order for which this item was obtained.
	OwnershipType        string              `json:"ownershipType"`              // The string "OwnedByBeneficiary".
	ProductId            string              `json:"productId"`                  // The Store ID for the product in the Microsoft Store catalog.
	ProductType          string              `json:"productType"`                // One of the following product types: Application, Durable, UnmanagedConsumable.
	PurchasedCountry     *string             `json:"purchasedCountry,omitempty"` // N/A
	Purchaser            *IdentityContractV6 `json:"purchaser,omitempty"`        // Represents the identity of the purchaser of the item.
	Quantity             *int                `json:"quantity,omitempty"`         // The quantity of the item. Currently, this will always be 1.
	SkuId                string              `json:"skuId"`                      // The Store ID for the product's SKU in the Microsoft Store catalog.
	SkuType              string              `json:"skuType"`                    // Type of the SKU. Possible values include Trial, Full, and Rental.
	StartDate            time.Time           `json:"startDate"`                  // The date that the item starts being valid.
	Status               string              `json:"status"`                     // The status of the item. Possible values include Active, Expired, Revoked, and Banned.
	Tags                 []string            `json:"tags"`                       // N/A
	TransactionId        string              `json:"transactionId"`              // The transaction ID as a result of the purchase of this item.
}

type IAPResponse struct {
	ContinuationToken *string                    `json:"continuationToken,omitempty"` // Token to retrieve remaining products if there are multiple sets.
	Items             []CollectionItemContractV6 `json:"items,omitempty"`             // An array of products for the specified user.
}
