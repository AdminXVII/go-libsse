package sse

import (
	"net/http"
	"sync"
)

// Server holds the info for a server
type Server struct {
	sync.RWMutex

	options     *Options
	lastEventId string
	clients     map[*client]bool // mimic a set
}

// NewServer creates a new SSE server.
func NewServer(options *Options) *Server {
	if options == nil {
		options = &Options{}
	}

	s := &Server{
		options: options,
		clients: make(map[*client]bool),
	}

	return s
}

// ServeHTTP is the basic handler for go's http package
func (s *Server) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		c := newClient(response, s.options.Headers)
		s.addClient(c)
		c.sendMessage(Message{Retry: s.options.RetryInterval})

		if s.options.InitMessages != nil {
			lastEventId := request.Header.Get("Last-Event-ID")
			c.intro = s.options.InitMessages(lastEventId, s.LastEventId())
		}

		c.listen()
		s.removeClient(c)
	} else if request.Method != "OPTIONS" {
		response.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// SendMessage broadcast a message to all clients
func (s *Server) SendMessage(message Message) {
	s.log("sending message")
	s.RLock()
	for c, open := range s.clients {
		if open {
			go c.sendMessage(message)
		}
	}
	s.RUnlock()
}

// LastEventId returns the event id of the last event to have been sent
func (s *Server) LastEventId() string {
	return s.lastEventId
}

// GetClientsCount outputs the current number of active http connections
func (s *Server) GetClientsCount() int {
	s.RLock()
	num := len(s.clients)
	s.RUnlock()

	return num
}

func (s *Server) addClient(client *client) {
	s.log("new client")
	s.Lock()
	s.clients[client] = true
	s.Unlock()
}

func (s *Server) removeClient(client *client) {
	s.log("removing client")
	s.Lock()
	delete(s.clients, client)
	s.Unlock()
}

func (s *Server) log(info string) {
	if s.options.Logger != nil {
		s.options.Logger.Print(info)
	}
}

// Restart closes all clients and allow new connections.
func (s *Server) Restart() {
	s.log("restarting server.")

	s.Lock()
	for client := range s.clients {
		s.removeClient(client)
	}
	s.Unlock()
}
