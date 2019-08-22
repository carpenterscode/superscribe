package listener

import (
	"fmt"

	af "github.com/carpenterscode/appsflyer-go"
	ss "github.com/carpenterscode/superscribe"
)

const (
	CancelSubscription  af.EventName = "cancel_subscription"
	CancelTrial         af.EventName = "cancel_trial"
	ChangeRenewalPref   af.EventName = "change_renewal_pref"
	RestartSubscription af.EventName = "restart_subscription"
)

const (
	ParamExpirationDate af.EventParam = "expiration_date"
)

const appsflyerKey = "appsflyer_id"

type AppsFlyer struct {
	Tracker *af.Tracker
}

func (l AppsFlyer) setup(user ss.User, configure func(*af.Event)) error {

	if user.GetString(appsflyerKey) == "" {
		return fmt.Errorf("%s required for AppsFlyer event", appsflyerKey)
	}

	// Configure AppsFlyer event with account info
	afEvent := af.NewEvent(user.GetString(appsflyerKey), af.IOS)
	afEvent.SetAdvertisingID(user.AdvertisingID())
	afEvent.SetDeviceIP(user.DeviceIP())

	configure(afEvent)

	// Send event to AppsFlyer
	return l.Tracker.Send(afEvent)
}

func (l AppsFlyer) Name() string {
	return "AppsFlyer"
}

func (l AppsFlyer) ChangedAutoRenewProduct(evt ss.AutoRenewEvent) error {
	return l.setup(evt, func(afEvent *af.Event) {
		afEvent.SetEventTime(evt.AutoRenewChangedAt())
		afEvent.SetName(ChangeRenewalPref)
	})
}

func (l AppsFlyer) ChangedAutoRenewStatus(evt ss.AutoRenewEvent) error {
	return l.setup(evt, func(afEvent *af.Event) {
		afEvent.SetEventTime(evt.AutoRenewChangedAt())
		if evt.AutoRenewStatus() {
			afEvent.SetName(CancelSubscription)
		} else {
			afEvent.SetName(RestartSubscription)
		}
	})
}

func (l AppsFlyer) Paid(evt ss.PayEvent) error {
	return l.setup(evt, func(afEvent *af.Event) {
		afEvent.SetEventTime(evt.PaidAt())
		afEvent.SetName(af.Subscribe)
		afEvent.SetRevenue(evt.Price(), evt.Currency())
	})
}

func (l AppsFlyer) Refunded(evt ss.RefundEvent) error {
	return l.setup(evt, func(afEvent *af.Event) {
		if evt.IsTrialPeriod() {
			afEvent.SetName(CancelTrial)
		} else {
			afEvent.SetName(CancelSubscription)
		}
		afEvent.SetEventTime(evt.RefundedAt())
	})
}

func (l AppsFlyer) StartedTrial(evt ss.StartTrialEvent) error {
	return l.setup(evt, func(afEvent *af.Event) {
		afEvent.SetEventTime(evt.StartedTrialAt())
		afEvent.SetName(af.StartTrial)
		afEvent.SetPrice(evt.Price(), evt.Currency())
	})
}
