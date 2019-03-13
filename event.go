package superscriber

import (
	"time"
)

type EventListener interface {

	// ChangedAutoRenewProduct indicates the next renewal period's product ID
	ChangedAutoRenewProduct(AutoRenewEvent)

	// ChangedAutoRenewStatus indicates new on/off state
	ChangedAutoRenewStatus(AutoRenewEvent)

	// Paid indicates a successful charge
	Paid(PayEvent)

	// Refunded indicates App Store customer support issued a subscription refund of some sort
	Refunded(RefundEvent)

	// StartedTrial indicates a subscription free trial began
	StartedTrial(StartTrialEvent)
}

type AutoRenewEvent interface {
	Event
	AutoRenewOn() bool
	AutoRenewProduct() string
	AutoRenewChangedAt() time.Time
}

type Event interface {
	OriginalTransactionID() string
	ProductID() string
}

type PayEvent interface {
	Event
	ExpiresAt() time.Time
	PaidAt() time.Time
}

type RefundEvent interface {
	Event
	RefundedAt() time.Time
	IsTrialPeriod() bool
}

type StartTrialEvent interface {
	Event
	StartedTrialAt() time.Time
}

type MultiEventListener struct {
	listeners []EventListener
}

func NewMultiEventListener() *MultiEventListener {
	return &MultiEventListener{make([]EventListener, 0)}
}

func (multi *MultiEventListener) Add(l EventListener) {
	multi.listeners = append(multi.listeners, l)
}

func (multi MultiEventListener) ChangedAutoRenewProduct(evt AutoRenewEvent) {
	for _, l := range multi.listeners {
		l.ChangedAutoRenewProduct(evt)
	}
}
func (multi MultiEventListener) ChangedAutoRenewStatus(evt AutoRenewEvent) {
	for _, l := range multi.listeners {
		l.ChangedAutoRenewStatus(evt)
	}
}
func (multi MultiEventListener) Paid(evt PayEvent) {
	for _, l := range multi.listeners {
		l.Paid(evt)
	}
}
func (multi MultiEventListener) Refunded(evt RefundEvent) {
	for _, l := range multi.listeners {
		l.Refunded(evt)
	}
}

func (multi MultiEventListener) StartedTrial(evt StartTrialEvent) {
	for _, l := range multi.listeners {
		l.StartedTrial(evt)
	}
}
