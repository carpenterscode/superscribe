package superscribe

import (
	"time"
)

//go:generate mockgen -package=superscribe -destination=./mock.go -self_package=github.com/carpenterscode/superscribe github.com/carpenterscode/superscribe EventListener,Subscription

// ExpiringSubscriptions returns a list of App Store receipts for subscriptions nearing
// expiration for a specified current time.
type ExpiringSubscriptions func(time.Time) []string

// SubscriptionFetch returns the last known state of a subscription by original transaction ID,
// which can determine what changes have happened when compared to the latest receipt info.
type SubscriptionFetch func(string) (Subscription, error)

type Note interface {
	Type() NoteType
	Environment() Env

	ReceiptInfo

	AutoRenewProduct() string
	AutoRenewChangedAt() time.Time
	RefundedAt() time.Time
	StartedTrialAt() time.Time
}

type ReceiptInfo interface {
	Status() int
	AutoRenewStatus() bool
	CancelledAt() time.Time
	ExpiresAt() time.Time
	IsTrialPeriod() bool
	OriginalTransactionID() string
	PaidAt() time.Time
	ProductID() string
}

type EventListener interface {

	// Name describes the listener for identification in the logs
	Name() string

	// ChangedAutoRenewProduct indicates the next renewal period's product ID
	ChangedAutoRenewProduct(AutoRenewEvent) error

	// ChangedAutoRenewStatus indicates new on/off state
	ChangedAutoRenewStatus(AutoRenewEvent) error

	// Paid indicates a successful charge
	Paid(PayEvent) error

	// Refunded indicates App Store customer support issued a subscription refund of some sort
	Refunded(RefundEvent) error

	// StartedTrial indicates a subscription free trial began
	StartedTrial(StartTrialEvent) error
}

type User interface {
	UserID() string
	FacebookID() string
	SignedUpAt() time.Time

	FirstName() string
	LastName() string
	Email() string
	ImageURL() string

	AdvertisingID() string
	DeviceIP() string
	PremiumAccess() bool
	GetString(string) string
}

type Subscription interface {
	OriginalTransactionID() string
	ProductID() string

	AutoRenewStatus() bool
	IsTrialPeriod() bool
	ExpiresAt() time.Time

	Currency() string
	Price() float64

	User
}

type AutoRenewEvent interface {
	Subscription
	AutoRenewProduct() string
	AutoRenewChangedAt() time.Time
}

type PayEvent interface {
	Subscription
	PaidAt() time.Time
}

type RefundEvent interface {
	Subscription
	RefundedAt() time.Time
}

type StartTrialEvent interface {
	Subscription
	StartedTrialAt() time.Time
}
