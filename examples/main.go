package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	".."
)

func fetch(originalTransactionID string) appstore.SubscriptionValues {
	log.Println("original transaction ID", originalTransactionID)
	return nil
}

func match(now time.Time) []string {
	results := []string{"1234567890", "2345678901"}
	log.Println("Match()", now, results)
	return results
}

type stubListener struct{}

func (h stubListener) Paid(evt appstore.PayEvent) {
	log.Println("Paid", evt.PaidAt(), evt.ExpiresAt())
}

func (h stubListener) ChangedAutoRenewProduct(evt appstore.AutoRenewEvent) {
	log.Println("ChangedAutoRenewProduct", evt.AutoRenewChangedAt(), evt.AutoRenewProduct())
}

func (h stubListener) ChangedAutoRenewStatus(evt appstore.AutoRenewEvent) {
	log.Println("ChangedAutoRenewStatus", evt.AutoRenewChangedAt(), evt.AutoRenewOn())
}

func (h stubListener) Refunded(evt appstore.RefundEvent) {
	log.Println("Refund", evt.RefundedAt())
}

func (h stubListener) StartedTrial(evt appstore.StartTrialEvent) {
	log.Println("StartTrial", evt.StartedTrialAt())
}

func main() {
	srv := appstore.NewServer(
		":8080",
		"password",
		match,
		fetch,
		time.Second,
	)
	srv.AddListener(stubListener{})
	srv.AddListener(stubListener{})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	srv.Start()

	<-quit

	srv.Stop()

	close(quit)
}
