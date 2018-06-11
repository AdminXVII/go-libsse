package sse

import (
  "net/http"
  "fmt"
)

type Client struct {
    messages chan *Message
    response http.ResponseWriter
}

// NewClient creates a new http client which will wait for messages
func NewClient(res http.ResponseWriter) *Client {
    return &Client{messages: make(chan *Message, 100), response: res}
}

// SendMessage sends a message to client.
func (c *Client) SendMessage(message *Message) {
    c.messages <- message
}

// Listen makes the client wait for messages and emit the messages as SSE
func (c *Client) Listen() {
    flusher, ok := c.response.(http.Flusher)
    if !ok {
        http.Error(c.response, "Streaming unsupported.", http.StatusInternalServerError)
        return
    }
    
    c.response.WriteHeader(http.StatusOK)
    flusher.Flush()
    
    for msg := range c.messages {
        fmt.Fprintf(c.response, msg.String())
        flusher.Flush()
    }
}

// Close closes the client and exit the Listen function if applicable
func (c *Client) Close() {
    close(c.messages)
}