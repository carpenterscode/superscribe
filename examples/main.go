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

func main() {
	srv := ss.NewServer(
		":8080",
		"password",
		match,
		fetch,
		time.Second,
	)
	srv.AddListener(listener.AppsFlyer{}, true)
	srv.AddListener(listener.Stub{}, false)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	srv.Start()

	<-quit

	srv.Stop()

	close(quit)
}
