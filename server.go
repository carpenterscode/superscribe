package superscriber

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type server struct {
	Match    ExpiringSubscriptions
	Listener *MultiEventListener
	Fetch    LastKnownSubscription
	mux      *http.ServeMux
	secret   string
	server   http.Server
	Ticker   *time.Ticker
}

func (s server) Start() {
	go func() {
		s.reviewSubscriptions(s.Match(time.Now()))
		for tick := range s.Ticker.C {
			log.Println("TICK at", tick)
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

	ctx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	if err := s.server.Shutdown(ctx); err != nil {
		log.Printf("Shutdown error %v\n", err)
		panic(err)
	}
}

func notificationHandler(w http.ResponseWriter, r *http.Request, listener EventListener) {
	var n Notification
	decodeErr := json.NewDecoder(r.Body).Decode(&n)
	defer r.Body.Close()
	if decodeErr != nil {
		log.Println("Should have decoded notification", decodeErr, r)
		return
	}

	if n.Environment == Sandbox {
		log.Println("Received Sandbox notification")
		return
	}

	switch n.NotificationType {
	case Cancel:
		if n.CancellationDate != nil {
			listener.Refunded(n)
		} else {
			log.Println("CANCEL notification should have cancellation date")
		}

	case Renewal, InteractiveRenewal:
		listener.Paid(n)

	case InitialBuy:
		if n.IsTrialPeriod() {
			listener.StartedTrial(n)
		} else {
			listener.Paid(n)
		}

	case DidChangeRenewalPref:
		listener.ChangedAutoRenewProduct(n)
	}
}

func (s server) reviewSubscriptions(receipts []string) {
	log.Println("num receipts", len(receipts))
	for _, r := range receipts {
		resp, err := verifyReceipt(s.secret, r)
		if err != nil {
			log.Println(err)
			continue
		}

		sub, fetchErr := s.Fetch(resp.OriginalTransactionID())
		if fetchErr != nil {
			return
		}

		if sub.ExpiresAt().Before(resp.ExpiresAt()) {
			s.Listener.Paid(resp)
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

func NewServer(addr, secret string, matcher ExpiringSubscriptions, fetcher LastKnownSubscription,
	interval time.Duration) *server {

	mux := http.NewServeMux()
	srv := server{
		Match:    matcher,
		Listener: NewMultiEventListener(),
		Fetch:    fetcher,
		mux:      mux,
		secret:   secret,
		server:   http.Server{Addr: addr, Handler: mux},
		Ticker:   time.NewTicker(interval),
	}

	mux.HandleFunc("/superscriber", func(w http.ResponseWriter, r *http.Request) {
		notificationHandler(w, r, srv.Listener)
	})

	return &srv
}
