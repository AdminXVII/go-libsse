package sse

import (
    "log"
)

type Options struct {
    // RetryInterval change EventSource default retry interval (milliseconds).
    RetryInterval int
    // Headers allow to set custom headers (useful for CORS support).
    Headers map[string]string
    // All usage logs end up in Logger
    Logger *log.Logger
    // Called when a new client appears. Return false if client should not be added
    InitClient func(client *Client, LastEventId string) bool
}

func (opt *Options) HasHeaders() bool {
    return opt.Headers != nil && len(opt.Headers) > 0
}
