package superscriber

import (
	"time"
)

type Notification struct {
	Env              Env      `json:"environment"`
	NotificationType NoteType `json:"notification_type"`
	Password         string   `json:"password"`

	CancellationDate   *AppleTime `json:"cancellation_date,omitempty"`
	WebOrderLineItemID string     `json:"web_order_line_item_id"`

	LatestReceipt            string          `json:"latest_receipt,omitempty"`
	LatestReceiptInfo        ios6ReceiptInfo `json:"latest_receipt_info,omitempty"`
	LatestExpiredReceipt     string          `json:"latest_expired_receipt,omitempty"`
	LatestExpiredReceiptInfo ios6ReceiptInfo `json:"latest_expired_receipt_info,omitempty"`

	AutoRenewStatus          bool      `json:"auto_renew_status,string"`
	AutoRenewStatusChangedAt AppleTime `json:"auto_renew_status_change_date"`
	AutoRenewAdamID          string    `json:"auto_renew_adam_id"`
	AutoRenewProductID       string    `json:"auto_renew_product_id"`
	ExpirationIntent         string    `json:"expiration_intent"`
}

type notification struct {
	Notification
}

func (n notification) AutoRenewStatus() bool {
	return n.Notification.AutoRenewStatus
}

func (n notification) AutoRenewProduct() string {
	return n.Notification.AutoRenewProductID
}

func (n notification) AutoRenewChangedAt() time.Time {
	if n.Notification.NotificationType == DidChangeRenewalPref {
		return time.Now().UTC()
	}
	return n.Notification.AutoRenewStatusChangedAt.Time
}

func (n notification) CancelledAt() time.Time {
	if n.Notification.CancellationDate != nil {
		return n.Notification.CancellationDate.Time
	}
	return time.Time{}
}

func (n notification) Environment() Env {
	return n.Notification.Env
}

func (n notification) ExpiresAt() time.Time {
	if n.Notification.CancellationDate != nil {
		return n.Notification.LatestExpiredReceiptInfo.ExpiresAt()
	}
	return n.Notification.LatestReceiptInfo.ExpiresAt()
}

func (n notification) IsTrialPeriod() bool {
	if n.Notification.CancellationDate != nil {
		return n.Notification.LatestExpiredReceiptInfo.IsTrialPeriod()
	}
	return n.Notification.LatestReceiptInfo.IsTrialPeriod()
}

func (n notification) OriginalTransactionID() string {
	if n.Notification.CancellationDate != nil {
		return n.Notification.LatestExpiredReceiptInfo.OriginalTransactionID()
	}
	return n.Notification.LatestReceiptInfo.OriginalTransactionID()
}

func (n notification) PaidAt() time.Time {
	if n.Notification.CancellationDate != nil {
		return n.Notification.LatestExpiredReceiptInfo.PaidAt()
	}
	return n.Notification.LatestReceiptInfo.PaidAt()
}

func (n notification) ProductID() string {
	if n.Notification.CancellationDate != nil {
		return n.Notification.LatestExpiredReceiptInfo.ProductID()
	}
	return n.Notification.LatestReceiptInfo.ProductID()
}

func (n notification) RefundedAt() time.Time {
	if n.Notification.CancellationDate != nil {
		return (*(n.Notification.CancellationDate)).Time
	}
	return time.Time{}
}

func (n notification) StartedTrialAt() time.Time {
	if n.Notification.CancellationDate != nil {
		return n.Notification.LatestExpiredReceiptInfo.StartedTrialAt()
	}
	return n.Notification.LatestReceiptInfo.StartedTrialAt()
}

func (n notification) Status() int {
	return StatusValid // TODO: Update to use unified receipt in Fall 2019
}

func (n notification) Type() NoteType {
	return n.Notification.NotificationType
}
