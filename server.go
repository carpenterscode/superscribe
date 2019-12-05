package superscribe

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/carpenterscode/superscribe/receipt"
)

type server struct {
	Match    ExpiringSubscriptions
	Listener *MultiEventListener
	Fetch    SubscriptionFetch
	Updater  SubscriptionUpdater
	mux      *http.ServeMux
	secret   string
	server   *http.Server
	Ticker   *time.Ticker
}

func (s server) Start() {
	go func() {
		now := time.Now()
		log.Println("Scan at", now)
		s.reviewSubscriptions(s.Match(now))
		for tick := range s.Ticker.C {
			log.Println("Scan at", tick)
			s.reviewSubscriptions(s.Match(tick))
		}
	}()

	go func() {
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
}

func (s server) Stop() {
	s.Ticker.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		log.Printf("Shutdown error %v\n", err)
		panic(err)
	}
}

func notificationHandler(w http.ResponseWriter, r *http.Request, listener EventListener,
	fetch SubscriptionFetch, updater SubscriptionUpdater) {

	data, bodyErr := ioutil.ReadAll(r.Body)
	if bodyErr != nil {
		log.Println("Should have read notification", bodyErr, r)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var body Notification
	if err := json.Unmarshal(data, &body); err != nil {
		log.Println("Should have unmarshaled notification", err, r)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	n := notification{body}

	if n.Environment() == Sandbox {
		log.Println("Received Sandbox notification")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if err := updater.UpdateWithNotification(n); err != nil {
		log.Println(n.OriginalTransactionID(), err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	sub, fetchErr := fetch(n.OriginalTransactionID())
	if fetchErr != nil {
		log.Println(fetchErr, n.OriginalTransactionID())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	evt := Event{}
	evt.SetNote(n)
	evt.SetRevenue(sub.Currency(), sub.Price())
	evt.SetUser(sub)

	var err error

	switch n.Type() {
	case Cancel:
		err = listener.Refunded(evt)

	case Renewal, InteractiveRenewal:
		err = listener.Paid(evt)

	case InitialBuy:
		if n.IsTrialPeriod() {
			err = listener.StartedTrial(evt)
		} else {
			err = listener.Paid(evt)
		}

	case DidChangeRenewalPref:
		err = listener.ChangedAutoRenewProduct(evt)

	case DidChangeRenewalStatus:
		err = listener.ChangedAutoRenewStatus(evt)

	}

	if err != nil {
		log.Println("Notification handler returns 500", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s server) reviewSubscriptions(receipts []string) {
	for _, receiptData := range receipts {
		resp, err := receipt.Validate(s.secret, receiptData)
		if err != nil {
			log.Println(err, receiptData)
			continue
		}

		if err := s.Updater.UpdateWithReceipt(resp); err != nil {
			log.Println(resp.OriginalTransactionID(), err)
			continue
		}

		sub, fetchErr := s.Fetch(resp.OriginalTransactionID())
		if fetchErr != nil {
			log.Println(fetchErr, resp.OriginalTransactionID())
			continue
		}

		// Check if expiration was pushed back before marking as paid
		if !sub.ExpiresAt().Before(resp.ExpiresAt()) {
			log.Println("Expiring has not renewed", sub.UserID())
			continue
		}

		evt := Event{}
		evt.SetReceiptInfo(resp)
		evt.SetRevenue(sub.Currency(), sub.Price())
		evt.SetUser(sub)

		if err := s.Listener.Paid(evt); err != nil {
			log.Println("Expiring Paid event error", err)
		}
	}
}

func (s server) AddListener(l EventListener) {
	s.Listener.Add(l)
}

func (s server) Addr() string {
	return s.server.Addr
}

func (s server) HandleFunc(pattern string, handlerFunc http.HandlerFunc) {
	s.mux.HandleFunc(pattern, handlerFunc)
}

func NewServer(addr, secret string, matcher ExpiringSubscriptions,
	fetch SubscriptionFetch, updater SubscriptionUpdater, interval time.Duration) *server {

	mux := http.NewServeMux()
	srv := server{
		Match:    matcher,
		Listener: NewMultiEventListener(),
		Fetch:    fetch,
		Updater:  updater,
		mux:      mux,
		secret:   secret,
		server:   &http.Server{Addr: addr, Handler: mux},
		Ticker:   time.NewTicker(interval),
	}

	mux.HandleFunc("/superscribe", func(w http.ResponseWriter, r *http.Request) {
		notificationHandler(w, r, srv.Listener, fetch, updater)
	})

	return &srv
}
