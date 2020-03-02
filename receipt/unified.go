package receipt

import (
	"encoding/json"
	"log"
	"time"
)

type UnifiedEnv string

const (
	Sandbox    UnifiedEnv = "Sandbox"
	Production UnifiedEnv = "Production"
)

type latestReceiptInfo struct {
	CancellationDate      *Millistamp `json:"cancellation_date_ms,string,omitempty"`
	CancellationReason    int         `json:"cancellation_reason,string"`
	ExpiresDate           Millistamp  `json:"expires_date_ms,string"`
	IsInIntroOfferPeriod  bool        `json:"is_in_intro_offer_period,string"`
	IsTrialPeriod         bool        `json:"is_trial_period,string"`
	IsUpgraded            bool        `json:"is_upgraded,string,omitempty"`
	OriginalPurchaseDate  Millistamp  `json:"original_purchase_date_ms,string"`
	OriginalTransactionID string      `json:"original_transaction_id"`
	ProductID             string      `json:"product_id"`
	PurchaseDate          Millistamp  `json:"purchase_date_ms,string"`
	Quantity              int         `json:"quantity,string"`
	SubscriptionGroupID   string      `json:"subscription_group_identifier"`
	TransactionID         string      `json:"transaction_id"`
}

type pendingRenewalInfo struct {
	AutoRenewProductID     *string     `json:"auto_renew_product_id,omitempty"`
	AutoRenewStatus        int         `json:"auto_renew_status,string"`
	ExpirationIntent       *int        `json:"expiration_intent,string,omitempty"`
	GracePeriodExpiresDate *Millistamp `json:"grace_period_expires_date_ms,string,omitempty"`
	IsInBillingRetryPeriod *int        `json:"is_in_billing_retry_period,string,omitempty"`
	OriginalTransactionID  string      `json:"original_transaction_id"`
	PriceConsentStatus     *int        `json:"price_consent_status,string,omitempty"`
	ProductID              string      `json:"product_id"`
}

type Unified struct {
	Environment     UnifiedEnv           `json:"environment"`
	LatestReceipts  []latestReceiptInfo  `json:"latest_receipt_info"`
	PendingRenewals []pendingRenewalInfo `json:"pending_renewal_info"`
	Status          int                  `json:"status"`
}

func (u Unified) PendingRenewalInfo() (pendingRenewalInfo, bool) {
	var pending pendingRenewalInfo
	if len(u.PendingRenewals) > 0 {
		return u.PendingRenewals[0], true
	}
	return pending, false
}

func (u Unified) LatestReceiptInfo() (latestReceiptInfo, bool) {
	var latest latestReceiptInfo
	if len(u.LatestReceipts) > 0 {
		return u.LatestReceipts[0], true
	}
	return latest, false
}

func (u Unified) GracePeriodExpiresDate() (time.Time, bool) {
	if info, ok := u.PendingRenewalInfo(); !ok {
		return time.Time{}, false
	} else if info.GracePeriodExpiresDate == nil {
		return time.Time{}, false
	} else {
		return (*info.GracePeriodExpiresDate).Time(), true
	}
}

func ParseUnified(data []byte) (Unified, error) {
	var rcpt Unified
	if err := json.Unmarshal(data, &rcpt); err != nil {
		log.Println("err", err)
		return rcpt, err
	}
	return rcpt, nil
}
