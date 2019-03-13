package superscriber

import (
	"time"
)

// ExpiringSubscriptions returns a list of original transaction IDs for subscriptions nearing
// expiration for a specified current time.
type ExpiringSubscriptions func(time.Time) []string

// LastKnownSubscription returns the last known state of a subscription which can determine
// what changes have happened when compared to the latest receipt info.
type LastKnownSubscription func(string) (SubscriptionValues, error)

type SubscriptionValues interface {
	AutoRenewStatus() bool
	IsTrialPeriod() bool

	ExpiresAt() time.Time
	PaidAt() time.Time
}
