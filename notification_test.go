package superscribe

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

const (
	currency     = "USD"
	price        = 9.99
	productID    = "year-premium"
	newProductID = "month-premium"
)

var autoRenewStatusChangedDate = time.Date(2019, time.June, 10, 21, 39, 47, 0, time.UTC)
var cancellationDate = time.Date(2019, time.March, 6, 17, 30, 17, 0, time.UTC)
var expiresDate = time.Date(2019, time.March, 13, 19, 11, 36, 0, time.UTC)
var gracePeriodExpiresDate = time.Date(2020, time.March, 8, 13, 00, 56, 0, time.UTC)
var purchaseDate = time.Date(2019, time.March, 6, 20, 11, 36, 0, time.UTC)
var originalPurchaseDate = time.Date(2019, time.March, 6, 20, 11, 36, 0, time.UTC)
var originalTransactionID = "123456789012345"

func dataFromFile(fileName string) []byte {
	if data, err := ioutil.ReadFile("testdata/" + fileName); err != nil {
		panic(err)
	} else {
		return data
	}
}

func notificationFromFile(fileName string) *notification {
	var body Notification
	if err := json.Unmarshal(dataFromFile(fileName), &body); err != nil {
		panic(fmt.Errorf("Should have unmarshalled JSON: %s", err.Error()))
	}
	return &notification{body}
}

func TestParseCancel(t *testing.T) {
	n := notificationFromFile("CANCEL.json")

	if n.Environment() != Prod {
		t.Error("Should have parsed environment: PROD")
	} else if n.AutoRenewStatus() {
		t.Error("Should have correct autorenew status: false")
	} else if n.Type() != Cancel {
		t.Error("Should have parsed notification type: CANCEL")
	} else if !n.RefundedAt().Equal(cancellationDate) {
		t.Error("Should have correctly parsed canceledAt", cancellationDate)
	} else if n.OriginalTransactionID() != originalTransactionID {
		t.Error("Should have parsed original transaction ID:", originalTransactionID)
	} else if n.IsTrialPeriod() {
		t.Error("Should have parsed not in trial period")
	}
}

func TestParseInitialBuy(t *testing.T) {
	n := notificationFromFile("INITIAL_BUY_to_trial.json")

	if n.Environment() != Prod {
		t.Error("Should have parsed environment: PROD")
	} else if !n.AutoRenewStatus() {
		t.Error("Should have correct autorenew status: true")
	} else if n.Type() != InitialBuy {
		t.Error("Should have parsed notification type: INITIAL_BUY")
	} else if n.body.LatestReceipt == "" {
		t.Error("Should have parsed latest_receipt field")
	} else if n.OriginalTransactionID() != originalTransactionID {
		t.Error("Should have parsed original transaction ID:", originalTransactionID)
	} else if !n.IsTrialPeriod() {
		t.Error("Should have parsed as in trial period")
	} else if !n.ExpiresAt().Equal(expiresDate) {
		t.Error("Should have parsed expires date as", expiresDate)
	} else if !n.PaidAt().Equal(purchaseDate) {
		t.Error("Should have parsed purchase date as", purchaseDate)
	}
}

func TestParseRenewal(t *testing.T) {
	n := notificationFromFile("RENEWAL.json")

	if n.Environment() != Prod {
		t.Error("Should have parsed environment: PROD")
	} else if !n.AutoRenewStatus() {
		t.Error("Should have correct autorenew status: true")
	} else if n.Type() != Renewal {
		t.Error("Should have parsed notification type: RENEWAL")
	} else if n.body.LatestReceipt == "" {
		t.Error("Should have parsed latest_receipt field")
	} else if n.OriginalTransactionID() != originalTransactionID {
		t.Error("Should have parsed original transaction ID:", originalTransactionID)
	} else if n.IsTrialPeriod() {
		t.Error("Should have parsed not in trial period")
	} else if !n.ExpiresAt().Equal(expiresDate) {
		t.Error("Should have parsed expires date as", expiresDate)
	} else if !n.PaidAt().Equal(purchaseDate) {
		t.Error("Should have parsed purchase date as", purchaseDate)
	}
}

func TestParseInteractiveRenewal(t *testing.T) {
	n := notificationFromFile("INTERACTIVE_RENEWAL.json")

	if n.Environment() != Prod {
		t.Error("Should have parsed environment: PROD")
	} else if !n.AutoRenewStatus() {
		t.Error("Should have correct autorenew status: true")
	} else if n.Type() != InteractiveRenewal {
		t.Error("Should have parsed notification type: INTERACTIVE_RENEWAL")
	} else if n.body.LatestReceipt == "" {
		t.Error("Should have parsed latest_receipt field")
	} else if n.OriginalTransactionID() != originalTransactionID {
		t.Error("Should have parsed original transaction ID:", originalTransactionID)
	} else if n.IsTrialPeriod() {
		t.Error("Should have parsed not in trial period")
	} else if !n.ExpiresAt().Equal(expiresDate) {
		t.Error("Should have parsed expires date", expiresDate)
	} else if !n.PaidAt().Equal(purchaseDate) {
		t.Error("Should have parsed purchase date", purchaseDate)
	}
}

func TestParseDidChangeRenewalPref(t *testing.T) {
	n := notificationFromFile("DID_CHANGE_RENEWAL_PREF.json")

	if n.Environment() != Prod {
		t.Error("Should have parsed environment: PROD")
	} else if !n.AutoRenewStatus() {
		t.Error("Should have correct autorenew status: true")
	} else if n.Type() != DidChangeRenewalPref {
		t.Error("Should have parsed notification type: DID_CHANGE_RENEWAL_PREF")
	} else if n.body.LatestReceipt == "" {
		t.Error("Should have parsed latest_receipt field")
	} else if n.OriginalTransactionID() != originalTransactionID {
		t.Error("Should have parsed original transaction ID:", originalTransactionID)
	} else if !n.IsTrialPeriod() {
		t.Error("Should have parsed as in trial period")
	} else if !n.ExpiresAt().Equal(expiresDate) {
		t.Error("Should have parsed expires date as", expiresDate)
	} else if !n.PaidAt().Equal(purchaseDate) {
		t.Error("Should have parsed purchase date as", purchaseDate)
	} else if n.AutoRenewProduct() != newProductID {
		t.Error("Should have parsed new product ID as", newProductID)
	} else if n.AutoRenewChangedAt().IsZero() {
		t.Error("Should have parsed auto renewed status changed as non-zero")
	}
}

func TestParseDidChangeRenewalStatus(t *testing.T) {
	n := notificationFromFile("DID_CHANGE_RENEWAL_STATUS_to_off.json")

	if n.Environment() != Prod {
		t.Error("Should have parsed environment: PROD")
	} else if n.AutoRenewStatus() {
		t.Error("Should have correct autorenew status: false")
	} else if n.Type() != DidChangeRenewalStatus {
		t.Error("Should have parsed notification type: DID_CHANGE_RENEWAL_STATUS")
	} else if n.body.LatestReceipt == "" {
		t.Error("Should have parsed latest_receipt field")
	} else if n.OriginalTransactionID() != originalTransactionID {
		t.Error("Should have parsed original transaction ID:", originalTransactionID)
	} else if !n.IsTrialPeriod() {
		t.Error("Should have parsed as in trial period")
	} else if !n.ExpiresAt().Equal(expiresDate) {
		t.Error("Should have parsed expires date as", expiresDate)
	} else if !n.PaidAt().Equal(purchaseDate) {
		t.Error("Should have parsed purchase date as", purchaseDate)
	} else if !n.AutoRenewChangedAt().Equal(autoRenewStatusChangedDate) {
		t.Error("Should have parsed auto renewed status changed as", autoRenewStatusChangedDate)
	}
}

func TestParseDidFailToRenew(t *testing.T) {
	n := notificationFromFile("DID_FAIL_TO_RENEW.json")

	if n.Environment() != Prod {
		t.Error("Should have parsed environment: PROD")
	} else if n.Type() != DidFailToRenew {
		t.Error("Should have parsed notification type: DID_FAIL_TO_RENEW")
	}

	if expiry, ok := n.GracePeriodEndsAt(); ok && expiry.Equal(gracePeriodExpiresDate) {
		t.Error("Should have parsed grace period expiration date as", expiry)
	}
}
