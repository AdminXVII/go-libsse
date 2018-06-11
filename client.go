package sse

import (
  "net/http"
)

type client struct {
    messages chan Message
    intro []Message
    response http.ResponseWriter
}

// newClient creates a new http client which will wait for messages
func newClient(res http.ResponseWriter, headers map[string]string) *client {
    responseHeaders := res.Header()

    for header, value := range headers {
        responseHeaders.Set(header, value)
    }
    responseHeaders.Set("Content-Type", "text/event-stream")
    responseHeaders.Set("Cache-Control", "no-cache")
    responseHeaders.Set("Connection", "keep-alive")
    
    return &client{messages: make(chan Message, 10), response: res}
}

// sendMessage sends a message to client.
func (c *client) sendMessage(message Message) {
    c.messages <- message
}

// listen makes the client wait for messages and emit the messages as SSE
func (c *client) listen() {
    flusher, ok := c.response.(http.Flusher)
    if !ok {
        http.Error(c.response, "Streaming unsupported.", http.StatusInternalServerError)
        return
    }
    
    c.response.WriteHeader(http.StatusOK)
    flusher.Flush()

    closeNotify := c.response.(http.CloseNotifier).CloseNotify()
    go func() {
        <-closeNotify
        c.close()
    }()
    
    for _, msg := range c.intro {
        msg.ToBuffer().WriteTo(c.response)
        flusher.Flush()
    }
    
    for msg := range c.messages {
        msg.ToBuffer().WriteTo(c.response)
        flusher.Flush()
    }
}

// close closes the client and exit the Listen function if applicable
func (c *client) close() {
    close(c.messages)
}