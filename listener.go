package superscribe

import (
	"log"
)

type listener struct {
	EventListener
}

type MultiEventListener struct {
	listeners []listener
}

func NewMultiEventListener() *MultiEventListener {
	return &MultiEventListener{make([]listener, 0)}
}

func (multi *MultiEventListener) Add(l EventListener) {
	multi.listeners = append(multi.listeners, listener{l})
}

func (multi *MultiEventListener) Name() string {
	return "Internal event bus"
}

func (multi MultiEventListener) ChangedAutoRenewProduct(evt AutoRenewEvent) error {
	for _, l := range multi.listeners {
		if err := l.ChangedAutoRenewProduct(evt); err != nil {
			log.Printf("%s listener ChangedAutoRenewProduct error: %v\n", l.Name(), err)
		}
	}
	return nil
}
func (multi MultiEventListener) ChangedAutoRenewStatus(evt AutoRenewEvent) error {
	for _, l := range multi.listeners {
		if err := l.ChangedAutoRenewStatus(evt); err != nil {
			log.Printf("%s listener ChangedAutoRenewStatus error: %v\n", l.Name(), err)
		}
	}
	return nil
}
func (multi MultiEventListener) Paid(evt PayEvent) error {
	for _, l := range multi.listeners {
		if err := l.Paid(evt); err != nil {
			log.Printf("%s listener Paid error: %v\n", l.Name(), err)
		}
	}
	return nil
}
func (multi MultiEventListener) Refunded(evt RefundEvent) error {
	for _, l := range multi.listeners {
		if err := l.Refunded(evt); err != nil {
			log.Printf("%s listener Refunded error: %v\n", l.Name(), err)
		}
	}
	return nil
}

func (multi MultiEventListener) StartedTrial(evt StartTrialEvent) error {
	for _, l := range multi.listeners {
		if err := l.StartedTrial(evt); err != nil {
			log.Printf("%s listener StartedTrial error: %v\n", l.Name(), err)
		}
	}
	return nil
}
