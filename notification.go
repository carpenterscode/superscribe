package superscriber

import (
	"time"
)

type env string

const (
	Sandbox env = "Sandbox"
	Prod    env = "PROD"
)

type nType string

const (
	Cancel               nType = "CANCEL"
	DidChangeRenewalPref nType = "DID_CHANGE_RENEWAL_PREF"
	InitialBuy           nType = "INITIAL_BUY"
	InteractiveRenewal   nType = "INTERACTIVE_RENEWAL"
	Renewal              nType = "RENEWAL"
)

type Notification struct {
	Environment      env    `json:"environment"`
	NotificationType nType  `json:"notification_type"`
	Password         string `json:"password"`

	CancellationDate   *AppleTime `json:"cancellation_date,omitempty"`
	WebOrderLineItemID string     `json:"web_order_line_item_id"`

	LatestReceipt            string          `json:"latest_receipt,omitempty"`
	LatestReceiptInfo        ios6ReceiptInfo `json:"latest_receipt_info,omitempty"`
	LatestExpiredReceipt     string          `json:"latest_expired_receipt,omitempty"`
	LatestExpiredReceiptInfo ios6ReceiptInfo `json:"latest_expired_receipt_info,omitempty"`

	AutoRenewStatus    bool   `json:"auto_renew_status,string"`
	AutoRenewAdamID    string `json:"auto_renew_adam_id"`
	AutoRenewProductID string `json:"auto_renew_product_id"`
	ExpirationIntent   string `json:"expiration_intent"`
}

func (n Notification) AutoRenewOn() bool {
	return n.AutoRenewStatus
}

func (n Notification) AutoRenewProduct() string {
	return n.AutoRenewProductID
}

func (n Notification) AutoRenewChangedAt() time.Time {
	return time.Now()
}

func (n Notification) ExpiresAt() time.Time {
	if n.CancellationDate != nil {
		return n.LatestExpiredReceiptInfo.ExpiresAt()
	}
	return n.LatestReceiptInfo.ExpiresAt()
}

func (n Notification) IsTrialPeriod() bool {
	if n.CancellationDate != nil {
		return n.LatestExpiredReceiptInfo.IsTrialPeriod()
	}
	return n.LatestReceiptInfo.IsTrialPeriod()
}

func (n Notification) OriginalTransactionID() string {
	if n.CancellationDate != nil {
		return n.LatestExpiredReceiptInfo.OriginalTransactionID()
	}
	return n.LatestReceiptInfo.OriginalTransactionID()
}

func (n Notification) PaidAt() time.Time {
	if n.CancellationDate != nil {
		return n.LatestExpiredReceiptInfo.PaidAt()
	}
	return n.LatestReceiptInfo.PaidAt()
}

func (n Notification) ProductID() string {
	if n.CancellationDate != nil {
		return n.LatestExpiredReceiptInfo.ProductID()
	}
	return n.LatestReceiptInfo.ProductID()
}

func (n Notification) RefundedAt() time.Time {
	if n.CancellationDate != nil {
		return (*(n.CancellationDate)).Time
	}
	return time.Now()
}

func (n Notification) StartedTrialAt() time.Time {
	if n.CancellationDate != nil {
		return n.LatestExpiredReceiptInfo.OriginalPurchaseDate.Time
	}
	return n.LatestReceiptInfo.OriginalPurchaseDate.Time
}
