package sse

import (
    "net/http"
    "sync"
)

// Server holds the info for a server
type Server struct {
    sync.RWMutex
    
    options *Options
    shutdown chan bool
    clients map[*Client]bool // mimic a set
}

// NewServer creates a new SSE server.
func NewServer(options *Options) *Server {
    if options == nil {
        options = &Options{}
    }

    s := &Server{
        options: options,
        shutdown: make(chan bool),
        clients: make(map[*Client]bool),
    }

    return s
}

// ServeHTTP is the basic handler for go's http package
func (s *Server) ServeHTTP(response http.ResponseWriter, request *http.Request) {
    if request.Method == "GET" {
        headers := response.Header()
    
        if s.options.HasHeaders() {
            for header, value := range s.options.Headers {
                headers.Set(header, value)
            }
        }
        headers.Set("Content-Type", "text/event-stream")
        headers.Set("Cache-Control", "no-cache")
        headers.Set("Connection", "keep-alive")
        
        c := NewClient(response)
        
        if s.options.InitClient != nil {
            lastEventId := request.Header.Get("Last-Event-ID")
            s.options.InitClient(c, lastEventId)
        }
        
        s.addClient(c)

        closeNotify := response.(http.CloseNotifier).CloseNotify()
        go func() {
            <-closeNotify
            s.removeClient(c)
        }()

        c.SendMessage(&Message{Retry: s.options.RetryInterval})
        c.Listen()
    } else if request.Method != "OPTIONS" {
        response.WriteHeader(http.StatusMethodNotAllowed)
    }
}

// SendMessage broadcast a message to all clients
func (s *Server) SendMessage(message *Message) {
    s.options.Logger.Print("sending message")
    s.RLock()
    for c, open := range s.clients {
        if open {
            c.SendMessage(message)
        }
    }
    s.RUnlock()
}

// GetClientCount outputs the current number of active http connections
func (s *Server) GetClientCount() int {
    s.RLock()
    num := len(s.clients)
    s.RUnlock()
    
    return num
}

func (s *Server) addClient(client *Client) {
    s.options.Logger.Print("new client")
    s.Lock()
    s.clients[client] = true
    s.Unlock()
}

func (s *Server) removeClient(client *Client) {
    s.options.Logger.Print("removing client")
    s.Lock()
    delete(s.clients, client)
    s.Unlock()
    client.Close()
}


// Restart closes all clients and allow new connections.
func (s *Server) Restart() {
    s.options.Logger.Print("restarting server.")
    
    s.Lock()
    for client, _ := range s.clients {
        s.removeClient(client)
    }
    s.Unlock()
}