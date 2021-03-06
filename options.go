package sse

import (
	"log"
)

// Options contains configuration for the server
type Options struct {
	// RetryInterval change EventSource default retry interval (milliseconds).
	RetryInterval int
	// Headers allow to set custom headers (useful for CORS support).
	Headers map[string]string
	// All usage logs end up in Logger
	Logger *log.Logger
	// Called when a new client appears. Return a set of messages to send before current messages
	InitMessages func(ClientLastEventId string, ServerLastEventId string) []Message
}
