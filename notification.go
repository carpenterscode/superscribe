package superscribe

import (
	"time"

	"github.com/carpenterscode/superscribe/receipt"
)

type Notification struct {
	Env              Env      `json:"environment"`
	NotificationType NoteType `json:"notification_type"`
	Password         string   `json:"password"`

	CancellationDate   *receipt.AppleTime `json:"cancellation_date,omitempty"`
	WebOrderLineItemID string             `json:"web_order_line_item_id"`

	LatestReceipt            string                  `json:"latest_receipt,omitempty"`
	LatestReceiptInfo        receipt.ReceiptInfoBody `json:"latest_receipt_info,omitempty"`
	LatestExpiredReceipt     string                  `json:"latest_expired_receipt,omitempty"`
	LatestExpiredReceiptInfo receipt.ReceiptInfoBody `json:"latest_expired_receipt_info,omitempty"`

	AutoRenewStatus          bool              `json:"auto_renew_status,string"`
	AutoRenewStatusChangedAt receipt.AppleTime `json:"auto_renew_status_change_date"`
	AutoRenewAdamID          string            `json:"auto_renew_adam_id"`
	AutoRenewProductID       string            `json:"auto_renew_product_id"`
	ExpirationIntent         string            `json:"expiration_intent"`
}

type notification struct {
	body Notification
}

func (n notification) AutoRenewStatus() bool {
	return n.body.AutoRenewStatus
}

func (n notification) AutoRenewProduct() string {
	return n.body.AutoRenewProductID
}

func (n notification) AutoRenewChangedAt() time.Time {
	if n.body.NotificationType == DidChangeRenewalPref {
		return time.Now().UTC()
	}
	return n.body.AutoRenewStatusChangedAt.Time
}

func (n notification) CancelledAt() time.Time {
	if n.body.CancellationDate != nil {
		return n.body.CancellationDate.Time
	}
	return time.Time{}
}

func (n notification) Environment() Env {
	return n.body.Env
}

func (n notification) ExpiresAt() time.Time {
	if n.body.CancellationDate != nil {
		return n.body.LatestExpiredReceiptInfo.ExpiresDateFormatted.Time
	}
	return n.body.LatestReceiptInfo.ExpiresDateFormatted.Time
}

func (n notification) IsTrialPeriod() bool {
	if n.body.CancellationDate != nil {
		return n.body.LatestExpiredReceiptInfo.IsTrialPeriod
	}
	return n.body.LatestReceiptInfo.IsTrialPeriod
}

func (n notification) OriginalTransactionID() string {
	if n.body.CancellationDate != nil {
		return n.body.LatestExpiredReceiptInfo.OriginalTransactionID
	}
	return n.body.LatestReceiptInfo.OriginalTransactionID
}

func (n notification) OriginalPurchaseDate() time.Time {
	if n.body.CancellationDate != nil {
		return n.body.LatestExpiredReceiptInfo.OriginalPurchaseDate.Time
	}
	return n.body.LatestReceiptInfo.OriginalPurchaseDate.Time
}

func (n notification) PaidAt() time.Time {
	if n.body.CancellationDate != nil {
		return n.body.LatestExpiredReceiptInfo.PurchaseDate.Time
	}
	return n.body.LatestReceiptInfo.PurchaseDate.Time
}

func (n notification) ProductID() string {
	if n.body.CancellationDate != nil {
		return n.body.LatestExpiredReceiptInfo.ProductID
	}
	return n.body.LatestReceiptInfo.ProductID
}

func (n notification) RefundedAt() time.Time {
	if n.body.CancellationDate != nil {
		return (*(n.body.CancellationDate)).Time
	}
	return time.Time{}
}

func (n notification) StartedTrialAt() time.Time {
	if n.body.CancellationDate != nil {
		return n.body.LatestExpiredReceiptInfo.OriginalPurchaseDate.Time
	}
	return n.body.LatestReceiptInfo.OriginalPurchaseDate.Time
}

func (n notification) Status() int {
	return receipt.StatusValid // TODO: Update to use unified receipt in Fall 2019
}

func (n notification) Type() NoteType {
	return n.body.NotificationType
}
