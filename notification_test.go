package superscriber

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

var cancellationDate = time.Date(2019, time.March, 6, 17, 30, 17, 0, time.UTC)
var expiresDate = time.Date(2019, time.March, 13, 19, 11, 36, 0, time.UTC)
var purchaseDate = time.Date(2019, time.March, 6, 20, 11, 36, 0, time.UTC)
var originalPurchaseDate = time.Date(2019, time.March, 6, 20, 11, 36, 0, time.UTC)
var originalTransactionID = "123456789012345"

func notificationFromFile(fileName string) *Notification {
	var n Notification
	if data, err := ioutil.ReadFile("testdata/" + fileName); err != nil {
		panic(err)
	} else if err := json.Unmarshal(data, &n); err != nil {
		panic(fmt.Errorf("Should have unmarshalled JSON"))
	}
	return &n
}

func TestCancel(t *testing.T) {
	n := notificationFromFile("notification_cancel.json")

	if n.Environment != Prod {
		t.Error("Should have parsed environment: PROD")
	} else if n.AutoRenewOn() {
		t.Error("Should have correct autorenew status: false")
	} else if n.NotificationType != Cancel {
		t.Error("Should have parsed notification type: CANCEL")
	} else if !n.RefundedAt().Equal(cancellationDate) {
		t.Error("Should have correctly parsed canceledAt", cancellationDate)
	} else if n.OriginalTransactionID() != originalTransactionID {
		t.Error("Should have parsed original transaction ID:", originalTransactionID)
	} else if n.IsTrialPeriod() {
		t.Error("Should have parsed not in trial period")
	}
}

func TestInitialBuy(t *testing.T) {
	n := notificationFromFile("notification_initial_buy.json")

	if n.Environment != Prod {
		t.Error("Should have parsed environment: PROD")
	} else if !n.AutoRenewStatus {
		t.Error("Should have correct autorenew status: true")
	} else if n.NotificationType != InitialBuy {
		t.Error("Should have parsed notification type: INITIAL_BUY")
	} else if n.LatestReceipt == "" {
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

func TestRenewal(t *testing.T) {
	n := notificationFromFile("notification_renewal.json")

	if n.Environment != Prod {
		t.Error("Should have parsed environment: PROD")
	} else if !n.AutoRenewStatus {
		t.Error("Should have correct autorenew status: true")
	} else if n.NotificationType != Renewal {
		t.Error("Should have parsed notification type: RENEWAL")
	} else if n.LatestReceipt == "" {
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

func TestInteractiveRenewal(t *testing.T) {
	n := notificationFromFile("notification_interactive_renewal.json")

	if n.Environment != Prod {
		t.Error("Should have parsed environment: PROD")
	} else if !n.AutoRenewStatus {
		t.Error("Should have correct autorenew status: true")
	} else if n.NotificationType != InteractiveRenewal {
		t.Error("Should have parsed notification type: INTERACTIVE_RENEWAL")
	} else if n.LatestReceipt == "" {
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

func TestDidChangeRenewalPref(t *testing.T) {
	n := notificationFromFile("notification_did_change_renewal_pref.json")

	if n.Environment != Prod {
		t.Error("Should have parsed environment: PROD")
	} else if !n.AutoRenewStatus {
		t.Error("Should have correct autorenew status: true")
	} else if n.NotificationType != DidChangeRenewalPref {
		t.Error("Should have parsed notification type: DID_CHANGE_RENEWAL_PREF")
	} else if n.LatestReceipt == "" {
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
