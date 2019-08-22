package listener

import (
	"log"

	ss "github.com/carpenterscode/superscribe"
)

type Stub struct{}

func (l Stub) Name() string {
	return "stub"
}

func (l Stub) Paid(evt ss.PayEvent) error {
	log.Println("Paid", evt.PaidAt(), evt.ExpiresAt())
	return nil
}

func (l Stub) ChangedAutoRenewProduct(evt ss.AutoRenewEvent) error {
	log.Println("ChangedAutoRenewProduct", evt.AutoRenewChangedAt(), evt.AutoRenewProduct())
	return nil
}

func (l Stub) ChangedAutoRenewStatus(evt ss.AutoRenewEvent) error {
	log.Println("ChangedAutoRenewStatus", evt.AutoRenewChangedAt(), evt.AutoRenewStatus())
	return nil
}

func (l Stub) Refunded(evt ss.RefundEvent) error {
	log.Println("Refund", evt.RefundedAt())
	return nil
}

func (l Stub) StartedTrial(evt ss.StartTrialEvent) error {
	log.Println("StartTrial", evt.StartedTrialAt())
	return nil
}
