package sse

import (
	"bytes"
	"fmt"
)

// Message decribes what the server can send
type Message struct {
	Id,
	Data,
	Event string
	Retry int
}

// ToBuffer exports the message to a bytes buffer
func (m *Message) ToBuffer() (buffer *bytes.Buffer) {
	buffer = new(bytes.Buffer)

	if len(m.Id) > 0 {
		fmt.Fprintf(buffer, "id: %s\n", m.Id)
	}

	if m.Retry > 0 {
		fmt.Fprintf(buffer, "retry: %d\n", m.Retry)
	}

	if len(m.Event) > 0 {
		fmt.Fprintf(buffer, "event: %s\n", m.Event)
	}

	if len(m.Data) > 0 {
		fmt.Fprintf(buffer, "data: %s\n", m.Data)
	}

	buffer.WriteRune('\n')

	return
}
