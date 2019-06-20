package superscriber

import (
	"fmt"
	"time"
)

type Event struct {

	// In-app subscription
	autoRenewProductID    string
	autoRenewStatus       bool
	currency              string
	isTrialPeriod         bool
	originalTransactionID string
	price                 float64
	productID             string

	autoRenewChangedAt time.Time
	cancelledAt        time.Time
	expiresAt          time.Time
	paidAt             time.Time
	refundedAt         time.Time
	startedTrialAt     time.Time

	// User data
	user User
}

func (evt *Event) SetNote(note Note) {
	evt.cancelledAt = note.CancelledAt()
	evt.expiresAt = note.ExpiresAt()
	evt.isTrialPeriod = note.IsTrialPeriod()
	evt.originalTransactionID = note.OriginalTransactionID()
	evt.paidAt = note.PaidAt()
	evt.productID = note.ProductID()

	evt.refundedAt = note.RefundedAt()
	evt.autoRenewProductID = note.AutoRenewProduct()
	evt.startedTrialAt = note.StartedTrialAt()
	evt.autoRenewStatus = note.AutoRenewStatus()
	evt.autoRenewChangedAt = note.AutoRenewChangedAt()
}

func (evt *Event) SetReceiptInfo(resp ReceiptInfo) {
	evt.cancelledAt = resp.CancelledAt()
	evt.expiresAt = resp.ExpiresAt()
	evt.isTrialPeriod = resp.IsTrialPeriod()
	evt.originalTransactionID = resp.OriginalTransactionID()
	evt.paidAt = resp.PaidAt()
	evt.productID = resp.ProductID()
	evt.autoRenewStatus = resp.AutoRenewStatus()
}

func (evt *Event) SetRevenue(currency string, price float64) {
	evt.currency = currency
	evt.price = price
}

func (evt *Event) SetStartedTrialAt(startedTrialAt time.Time) {
	evt.startedTrialAt = startedTrialAt
}

func (evt *Event) SetUser(user User) {
	evt.user = user
}

func (evt Event) OriginalTransactionID() string {
	return evt.originalTransactionID
}

func (evt Event) ProductID() string {
	return evt.productID
}

func (evt Event) AutoRenewStatus() bool {
	return evt.autoRenewStatus
}

func (evt Event) AutoRenewProduct() string {
	return evt.autoRenewProductID
}

func (evt Event) IsTrialPeriod() bool {
	return evt.isTrialPeriod
}

func (evt Event) ExpiresAt() time.Time {
	return evt.expiresAt
}

func (evt Event) Currency() string {
	return evt.currency
}

func (evt Event) Price() float64 {
	return evt.price
}

func (evt Event) AutoRenewChangedAt() time.Time {
	return evt.autoRenewChangedAt
}

func (evt Event) CancelledAt() time.Time {
	return evt.cancelledAt
}

func (evt Event) PaidAt() time.Time {
	return evt.paidAt
}

func (evt Event) RefundedAt() time.Time {
	return evt.refundedAt
}

func (evt Event) StartedTrialAt() time.Time {
	return evt.startedTrialAt
}

func (evt Event) Status() int {
	return StatusValid
}

func (evt Event) String() string {
	return fmt.Sprintf("\n%s: %v\n", "autoRenewProductID", evt.autoRenewProductID) +
		fmt.Sprintf("%s: %v\n", "autoRenewStatus", evt.autoRenewStatus) +
		fmt.Sprintf("%s: %v\n", "isTrialPeriod", evt.isTrialPeriod) +
		fmt.Sprintf("%s: %v\n", "originalTransactionID", evt.originalTransactionID) +
		fmt.Sprintf("%s: %v\n", "productID", evt.productID) +
		fmt.Sprintf("%s: %v\n", "autoRenewChangedAt", evt.autoRenewChangedAt) +
		fmt.Sprintf("%s: %v\n", "cancelledAt", evt.cancelledAt) +
		fmt.Sprintf("%s: %v\n", "expiresAt", evt.expiresAt) +
		fmt.Sprintf("%s: %v\n", "paidAt", evt.paidAt) +
		fmt.Sprintf("%s: %v\n", "refundedAt", evt.refundedAt) +
		fmt.Sprintf("%s: %v\n", "startedTrialAt", evt.startedTrialAt)
}

func (evt Event) UserID() string {
	return evt.user.UserID()
}

func (evt Event) SignedUpAt() time.Time {
	return evt.user.SignedUpAt()
}

func (evt Event) FirstName() string {
	return evt.user.FirstName()
}

func (evt Event) LastName() string {
	return evt.user.LastName()
}

func (evt Event) Email() string {
	return evt.user.Email()
}

func (evt Event) ImageURL() string {
	return evt.user.ImageURL()
}

func (evt Event) AdvertisingID() string {
	return evt.user.AdvertisingID()
}

func (evt Event) DeviceIP() string {
	return evt.user.DeviceIP()
}

func (evt Event) PremiumAccess() bool {
	return evt.user.PremiumAccess()
}

func (evt Event) FacebookID() string {
	return evt.user.FacebookID()
}

func (evt Event) GetString(key string) string {
	return evt.user.GetString(key)
}
