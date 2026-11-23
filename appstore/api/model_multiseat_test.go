package api

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestTransactionMultiSeatPayload(t *testing.T) {
	t.Parallel()
	var got JWSTransaction
	err := json.Unmarshal([]byte(`{
		"offerType":5,"additionalTransactionReason":"SEAT_ADDITION",
		"multiSeatInfo":{"invitationManagement":"DEVELOPER","priceTiers":[{"price":9007199254740993,"quantity":3}],"revokedSeatQuantity":1},
		"advancedCommerceInfo":{"items":[{"SKU":"subscription","bundleItem":{"SKU":"item","description":"description","displayName":"Item","price":9007199254740993,"taxCode":"tax"},"partners":[{"id":"partner-id","name":"Partner"}],"seatPriceTiers":[{"maximumSeatQuantity":5,"price":9007199254740993},{"price":0}],"refunds":[{"revokedSeatQuantity":2}]}]}
	}`), &got)
	if err != nil {
		t.Fatal(err)
	}
	if got.OfferType != OfferTypeRetention || got.AdditionalTransactionReason != AdditionalTransactionReasonSeatAddition {
		t.Fatalf("offer/reason = %d/%s", got.OfferType, got.AdditionalTransactionReason)
	}
	wantSeats := &TransactionMultiSeatInfo{
		InvitationManagement: InvitationManagementDeveloper,
		PriceTiers:           []MultiSeatPriceTier{{Price: 9007199254740993, Quantity: 3}},
		RevokedSeatQuantity:  1,
	}
	if !reflect.DeepEqual(got.MultiSeatInfo, wantSeats) {
		t.Errorf("multiSeatInfo = %+v, want %+v", got.MultiSeatInfo, wantSeats)
	}
	if len(got.AdvancedCommerceInfo.Items) != 1 {
		t.Fatalf("items = %+v", got.AdvancedCommerceInfo.Items)
	}
	item := got.AdvancedCommerceInfo.Items[0]
	wantBundle := &AdvancedCommerceBundleItem{SKU: "item", Description: "description", DisplayName: "Item", Price: 9007199254740993, TaxCode: "tax"}
	if !reflect.DeepEqual(item.BundleItem, wantBundle) {
		t.Errorf("bundleItem = %+v, want %+v", item.BundleItem, wantBundle)
	}
	if want := []AdvancedCommercePartner{{ID: "partner-id", Name: "Partner"}}; !reflect.DeepEqual(item.Partners, want) {
		t.Errorf("partners = %+v, want %+v", item.Partners, want)
	}
	wantTiers := []AdvancedCommerceSeatPriceTier{{MaximumSeatQuantity: new(int32(5)), Price: 9007199254740993}, {Price: 0}}
	if !reflect.DeepEqual(item.SeatPriceTiers, wantTiers) {
		t.Errorf("seatPriceTiers = %+v, want %+v", item.SeatPriceTiers, wantTiers)
	}
	if len(item.Refunds) != 1 || item.Refunds[0].RevokedSeatQuantity != 2 {
		t.Errorf("refunds = %+v", item.Refunds)
	}
}

func TestRenewalMultiSeatPayload(t *testing.T) {
	t.Parallel()
	var got JWSRenewalInfoDecodedPayload
	err := json.Unmarshal([]byte(`{
		"offerType":5,"renewalQuantity":7,
		"multiSeatInfo":{"renewalPriceTiers":[{"price":9007199254740993,"quantity":7}]},
		"advancedCommerceInfo":{"items":[{"SKU":"subscription","bundleItem":{"SKU":"item","price":1999},"partners":[{"id":"partner-id","name":"Partner"}],"seatRenewalPriceTiers":[{"maximumSeatQuantity":5,"price":1999},{"price":999}]}]}
	}`), &got)
	if err != nil {
		t.Fatal(err)
	}
	if got.OfferType != OfferTypeRetention || got.RenewalQuantity != 7 {
		t.Fatalf("offer/quantity = %d/%d", got.OfferType, got.RenewalQuantity)
	}
	wantSeats := &RenewalMultiSeatInfo{RenewalPriceTiers: []MultiSeatPriceTier{{Price: 9007199254740993, Quantity: 7}}}
	if !reflect.DeepEqual(got.MultiSeatInfo, wantSeats) {
		t.Errorf("multiSeatInfo = %+v, want %+v", got.MultiSeatInfo, wantSeats)
	}
	if len(got.AdvancedCommerceInfo.Items) != 1 {
		t.Fatalf("items = %+v", got.AdvancedCommerceInfo.Items)
	}
	item := got.AdvancedCommerceInfo.Items[0]
	if want := (&AdvancedCommerceBundleItem{SKU: "item", Price: 1999}); !reflect.DeepEqual(item.BundleItem, want) {
		t.Errorf("bundleItem = %+v, want %+v", item.BundleItem, want)
	}
	if want := []AdvancedCommercePartner{{ID: "partner-id", Name: "Partner"}}; !reflect.DeepEqual(item.Partners, want) {
		t.Errorf("partners = %+v, want %+v", item.Partners, want)
	}
	wantTiers := []AdvancedCommerceSeatPriceTier{{MaximumSeatQuantity: new(int32(5)), Price: 1999}, {Price: 999}}
	if !reflect.DeepEqual(item.SeatRenewalPriceTiers, wantTiers) {
		t.Errorf("seatRenewalPriceTiers = %+v, want %+v", item.SeatRenewalPriceTiers, wantTiers)
	}
}

func TestAdvancedCommerceUnlimitedTierJSON(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name string
		tier AdvancedCommerceSeatPriceTier
		want string
	}{
		{name: "unlimited free tier", tier: AdvancedCommerceSeatPriceTier{Price: 0}, want: `{"price":0}`},
		{name: "bounded tier", tier: AdvancedCommerceSeatPriceTier{MaximumSeatQuantity: new(int32(5)), Price: 1999}, want: `{"maximumSeatQuantity":5,"price":1999}`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.tier)
			if err != nil {
				t.Fatal(err)
			}
			if string(data) != tt.want {
				t.Errorf("JSON = %s, want %s", data, tt.want)
			}
		})
	}
}

func TestLegacyPayloadOmitsMultiSeatFields(t *testing.T) {
	t.Parallel()
	for _, payload := range []any{JWSTransaction{}, JWSRenewalInfoDecodedPayload{}, AdvancedCommerceTransactionItem{}, AdvancedCommerceRenewalItems{}} {
		data, err := json.Marshal(payload)
		if err != nil {
			t.Fatal(err)
		}
		var fields map[string]json.RawMessage
		if err := json.Unmarshal(data, &fields); err != nil {
			t.Fatal(err)
		}
		for _, key := range []string{"multiSeatInfo", "renewalQuantity", "additionalTransactionReason", "bundleItem", "partners", "seatPriceTiers", "seatRenewalPriceTiers"} {
			if _, exists := fields[key]; exists {
				t.Errorf("%T unexpectedly serialized %s", payload, key)
			}
		}
	}
}
