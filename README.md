# go-sse [![Build Status](https://travis-ci.org/AdminXVII/go-sse.svg?branch=master)](https://travis-ci.org/AdminXVII/go-sse) [![GoDoc](https://godoc.org/github.com/AdminXVII/go-sse?status.svg)](http://godoc.org/github.com/AdminXVII/go-sse)

Server-Sent Events for Go

Inspired from [alexandrevincenzi's work](https://github.com/alexandrevicenzi/go-sse)

## About

[Server-sent events](http://www.html5rocks.com/en/tutorials/eventsource/basics/) is a method of continuously sending data from a server to the browser, rather than repeatedly requesting it, replacing the "long polling way".

`go-sse` is a small library to create a Server-Sent Events server in Go.

## Features

- Fully thread-safe
- Custom initialization after connection from browser (for example to patch the gap in last-event-id)
- Custom headers (useful for CORS)
- `Last-Event-ID` support (resend lost messages)
- [Follow SSE specification](https://html.spec.whatwg.org/multipage/comms.html#server-sent-events)

## Getting Started

Simple Go example that send messages to all clients.

```go
package main

import (
    "log"
    "net/http"
    "strconv"
    "time"

    "github.com/AdminXVII/go-sse"
)

func main() {
    // Create the server.
    s := sse.NewServer(nil)

    // Register with /events endpoint.
    http.Handle("/events/", s)

    // Dispatch messages to channel-1.
    go func () {
        for {
            s.SendMessage(sse.Message{Data: time.Now().String()})
            time.Sleep(5 * time.Second)
        }
    }()

    http.ListenAndServe(":3000", nil)
}
```

Connecting to our server from JavaScript:

```js
e1 = new EventSource('/events/');
e1.onmessage = function(event) {
    // do something...
};
```
