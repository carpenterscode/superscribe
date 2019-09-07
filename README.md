# Superscribe

[![Build Status](https://travis-ci.org/carpenterscode/superscribe.svg?branch=master)](https://travis-ci.org/carpenterscode/superscribe)
[![GoDoc](https://godoc.org/github.com/carpenterscode/superscribe?status.svg)](https://godoc.org/github.com/carpenterscode/superscribe)

An easier way to handle App Store subscriptions

## Overview

Getting App Store Status Update Notifications (now Server-to-Server Notifications) for in-app
subscriptions can be tricky. _RENEWAL_ notifications may not occur as expected, _CANCEL_ does not
indicate when _auto_renew_status_ was switched on or off, and it can seem arbitrary when you
get older iOS 6 style of receipts or the new style.

Others have described these challenges too:

- [statusUpdateNotification not getting renewal notification type](https://forums.developer.apple.com/message/283579#283579) from Apple Developer Forums
- [Apple Subscription Notifications Are Almost Useless](https://www.revenuecat.com/2018/09/24/apple-subscription-notifications-are-almost-useless) from RevenueCat
- [Never Get RENEWAL Notification](https://stackoverflow.com/q/48049771/5477264) from StackOverflow
- [Server to Server Polling Auto Renewable Subscription](https://stackoverflow.com/q/50947948/5477264) from StackOverflow

_Superscribe_ intends to provide a basic “just works”, correct solution for the main
subscription use cases.

## Get started

### Install

```sh
go get github.com/carpenterscode/superscribe
```

### Configure

_Superscribe_ provides a server that both

- scans for expiring subscriptions to check for state changes (like successful renewals) and
- listens for App Store notifications for a limited number of subscription events.

The server needs

1.  Your App Store shared secret
2.  A func you define to retrieve expiring subscriptions from the database, called during a _scan_
    operation
3.  A listener to update the database after the scan

See how to connect listeners to Superscribe in [example/main.go](examples/main.go). This shows how
to use the included [AppsFlyer listener](listener/appsflyer.go) that attributes events using
server-to-server API.

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
  `/superscribe`, but we can change that in the future.
- Use anything more sophisticated than a Go `time.Ticker`.

### Usage

### Run automated tests

Generate mocks first

```sh
go generate
```

Test

```sh
go test ./... ./receipt
```

## Caveats

Currently, _Superscribe_ should only be run in a single instance setup. I personally run it on
production in a single-pod Kubernetes deployment, but we should figure out how to solve for
redundancy and performance by adding some kind of scalability.

## Future work

There's a lot of unfortunate complexity to subscription management, so the longer term goal is to
increase extensibility and robustness.

**Most important:** Let’s gather real use-cases and requirements to draft a prioritized roadmap.

- **Distinguish among first, first year's worth of, and remaining payments.** The _paid at_ event
  could be made more versatile and track 30% vs 15% App Store fee. Or to filter out renewals from
  first payments.
- **Track plan upgrade responses from customers.** For instance, moving all monthly subscriptions
  from 7.99/mo to 9.99/mo.
- **Offer a scalable solution.** Subscriptions in the local database should only be scanned by a
  single process, but multiple instances of listeners should be able to coexist. The current 1:1
  model limits _Superscribe_ to one instance.
