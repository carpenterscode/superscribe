# Superscriber

A consistent way to manage App Store subscriptions

## Overview

Have you tried using App Store Status Update Notifications for your app’s in-app subscriptions
and wondered why there are so few _RENEWAL_ notifications, or why neither _CANCEL_ nor
_DID_CHANGE_RENEWAL_PREF_ indicate when _auto_renew_status_ was switched on or off,
or why sometimes you get iOS 6 style receipts when your app never supported iOS 6 to begin with,
or… why–just why? Subscriptions may not be simple, but the App Store makes it feel unnecessarily
difficult.

If you don't believe me, check these out:

- [statusUpdateNotification not getting renewal notification type](https://forums.developer.apple.com/message/283579#283579) from Apple Developer Forums
- [Apple Subscription Notifications Are Almost Useless](https://www.revenuecat.com/2018/09/24/apple-subscription-notifications-are-almost-useless) from RevenueCat
- [Never Get RENEWAL Notification](https://stackoverflow.com/q/48049771/5477264) from StackOverflow
- [Server to Server Polling Auto Renewable Subscription](https://stackoverflow.com/q/50947948/5477264) from StackOverflow

_Superscriber_ provides a basic solution for the main App Store subscription use cases, making
it straightforward to reason about, and thus easier to extend for sophisticated or special cases.

## Get started

### Install

```sh
go get github.com/carpenterscode/superscriber
```

### Usage

_Superscriber_ provides a server that both

- scans for expiring subscriptions to check for state changes (like successful renewals) and
- listens for App Store notifications for a limited number of subscription events.

The server needs

1.  Your App Store shared secret
2.  A func you define to retrieve expiring subscriptions from the database, called during a _scan_
    operation
3.  A listener to update the database after the scan

```go
package main

import (
	"log"
	"os"
	"os/signal"
)

func matchWithDB(now time.Time, lookahead time.Duration, db sql.DB) []string {
	// TODO: Implement
}

func fetch(originalTransactionID string) (ss.SubscriptionValues, error) {
	// TODO: Implement
}

func main() {

	match := func(now time.Time) []string {
		return matchWithDB(now, lookahead, db)
	}

	interval := time.Duration(30) * time.Minute

	srv := ss.NewServer(addr, "secret", match, fetch, interval)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	srv.Start()
	log.Println("Started listening on", srv.Addr())

	<-quit
	log.Println("Draining…")
	srv.Stop()
	log.Println("Server shutdown")

	close(quit)
}
```

Optionally you can provide

- More listeners, such as for server-side conversion analytics, etc.

```go
srv.AddListener(AnalyticsListener{db, tracker})
```

- HTTP server request handlers, such as for `200 OK` responses to /healthz pings

```go
srv.HandleFunc("/healthz", func(writer http.ResponseWriter, req *http.Request) {
	writer.Write([]byte("OK"))
})
```

You cannot currently

- Modify the App Store Status Update Notification endpoint. It's currently hardcoded to
  `/superscriber`, but we can change that in the future.
- Use anything more sophisticated than a Go `time.Ticker`.

## Caveats

Currently, _Superscriber_ should only be run in a single instance setup. I personally run it on
production in a single-pod Kubernetes deployment, but we should figure out how to solve for
redundancy and performance by adding some kind of scalability.

## Future work

There's a lot of unfortunate complexity to subscription management, so the longer term goal is to
increase extensibility and robustness.

**Most important:** Let’s gather real use-cases and requirements to draft a prioritized roadmap.

- **Distinguish among first, first year's worth of, and remaining payments.** The _paid at_ event
  could be made more versatile and track 30% vs 15% App Store fee. Or to filter out renewals from
  first payments.
- **Make it easy for listeners to use the same single database query.** For instance, a server with
  listeners for multiple analytics services may need to get price and currency from the database,
  but for now, the easy solution involves making one query per listener, per event.
- **Track plan upgrade responses from customers.** For instance, moving all monthly subscriptions
  from 7.99/mo to 9.99/mo.
- **Publish some internal listeners.** Such as AppsFlyer and Amplitude.
- **Offer a scalable solution.** Subscriptions in the local database should only be scanned by a
  single process, but multiple instances of listeners should be able to coexist. The current 1:1
  model limits _Superscriber_ to one instance.
