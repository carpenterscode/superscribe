package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	ss "github.com/carpenterscode/superscribe"
	"github.com/carpenterscode/superscribe/listener"
)

func fetch(receipt string) (ss.Subscription, error) {
	log.Println("receipt", receipt)
	return nil, nil
}

func match(now time.Time) []string {
	results := []string{"1234567890", "2345678901"}
	log.Println("Match()", now, results)
	return results
}

type updater struct{}

func (u updater) UpdateWithNotification(note ss.Note) error {
	return nil
}

func (u updater) UpdateWithReceipt(r receipt.Info) error {
	return nil
}

func main() {
	srv := ss.NewServer(
		":8080",
		"password",
		match,
		fetch,
		updater{},
		time.Second,
	)
	srv.AddListener(listener.AppsFlyer{})
	srv.AddListener(listener.Stub{})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	srv.Start()

	<-quit

	srv.Stop()

	close(quit)
}
