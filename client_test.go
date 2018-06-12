package sse

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Implement the CloseNotifier interface
type CloseRecorder struct {
	*httptest.ResponseRecorder

	closeNotify chan bool
}

// Dummy struct to see if code fail if flush is not available
// Adapted from httptest.ResponseRecorder
type NoFlushRecorder struct {
	Code        int
	HeaderMap   http.Header
	Body        *bytes.Buffer
	result      *http.Response // cache of Result's return value
	snapHeader  http.Header    // snapshot of HeaderMap at first Write
	wroteHeader bool
}

func NewNoFlushRecorder() *NoFlushRecorder {
	return &NoFlushRecorder{
		HeaderMap: make(http.Header),
		Body:      new(bytes.Buffer),
		Code:      200,
	}
}

// Header returns the response headers.
func (rw *NoFlushRecorder) Header() http.Header {
	m := rw.HeaderMap
	if m == nil {
		m = make(http.Header)
		rw.HeaderMap = m
	}
	return m
}

func (rw *NoFlushRecorder) writeHeader(b []byte, str string) {
	if rw.wroteHeader {
		return
	}
	if len(str) > 512 {
		str = str[:512]
	}

	m := rw.Header()

	_, hasType := m["Content-Type"]
	hasTE := m.Get("Transfer-Encoding") != ""
	if !hasType && !hasTE {
		if b == nil {
			b = []byte(str)
		}
		m.Set("Content-Type", http.DetectContentType(b))
	}

	rw.WriteHeader(200)
}

// WriteHeader sets rw.Code. After it is called, changing rw.Header
// will not affect rw.HeaderMap.
func (rw *NoFlushRecorder) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}
	rw.Code = code
	rw.wroteHeader = true
	if rw.HeaderMap == nil {
		rw.HeaderMap = make(http.Header)
	}
	rw.snapHeader = cloneHeader(rw.HeaderMap)
}

func cloneHeader(h http.Header) http.Header {
	h2 := make(http.Header, len(h))
	for k, vv := range h {
		vv2 := make([]string, len(vv))
		copy(vv2, vv)
		h2[k] = vv2
	}
	return h2
}

// Write always succeeds and writes to rw.Body, if not nil.
func (rw *NoFlushRecorder) Write(buf []byte) (int, error) {
	rw.writeHeader(buf, "")
	if rw.Body != nil {
		rw.Body.Write(buf)
	}
	return len(buf), nil
}

func NewCloseRecorder() *CloseRecorder {
	return &CloseRecorder{
		httptest.NewRecorder(),
		make(chan bool),
	}
}

func (r *CloseRecorder) CloseNotify() <-chan bool {
	return r.closeNotify
}

func NewBaseHeaders() map[string]string {
	baseHeaders := make(map[string]string)

	baseHeaders["Content-Type"] = "text/event-stream"
	baseHeaders["Cache-Control"] = "no-cache"
	baseHeaders["Connection"] = "keep-alive"
	return baseHeaders
}

func TestNewEmptyClientHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	newClient(rr, nil)

	headers := NewBaseHeaders()

	for header, value := range headers {
		if rr.Header().Get(header) != value {
			t.Errorf("returned unexpected header: got %v want %v",
				rr.Header().Get(header), value)
		}
	}
}

func TestNewClient(t *testing.T) {
	rr := httptest.NewRecorder()

	c := newClient(rr, nil)

	if rr != c.response {
		t.Error("Client does not save the response")
	}
}

func TestNewClientHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	headers := make(map[string]string)
	headers["custom-header"] = "keep"

	newClient(rr, headers)

	for header, value := range headers {
		if rr.Header().Get(header) != value {
			t.Errorf("returned unexpected header: got %v want %v",
				rr.Header().Get(header), value)
		}
	}
}

func TestClientResponse(t *testing.T) {
	rr := NewCloseRecorder()

	c := newClient(rr, nil)
	msg := Message{Data: "Hello"}

	go func() {
		c.listen()

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		if rr.Body.String() != msg.ToBuffer().String() {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), msg.ToBuffer().String())
		}
	}()

	c.sendMessage(msg)
	c.close()
}

func TestClientResponseOrder(t *testing.T) {
	msg := Message{Data: "Hello"}
	intro := Message{Data: "Yikes"}

	rr := NewCloseRecorder()
	c := newClient(rr, nil)
	c.intro = []Message{intro}

	go func() {
		c.listen()

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		buffer := new(bytes.Buffer)
		intro.ToBuffer().WriteTo(buffer)
		msg.ToBuffer().WriteTo(buffer)

		expected := buffer.String()
		if rr.Body.String() != expected {
			t.Errorf("handler returned unexpected body: got %v want %v",
				rr.Body.String(), expected)
		}
	}()

	c.sendMessage(msg)
	c.close()
}

func TestClientNoFlushResponse(t *testing.T) {
	msg := Message{Data: "Hello"}

	rr := NewNoFlushRecorder()
	c := newClient(rr, nil)

	go func() {
		c.listen()

		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusInternalServerError)
		}
	}()

	c.sendMessage(msg)
	c.close()
}
