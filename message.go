package sse

import (
    "bytes"
    "fmt"
)

type Message struct {
    Id,
    Data,
    Event string
    Retry int
}

func (m *Message) String() string {
    var buffer bytes.Buffer

    if len(m.Id) > 0 {
        buffer.WriteString(fmt.Sprintf("id: %s\n", m.Id))
    }

    if m.Retry > 0 {
        buffer.WriteString(fmt.Sprintf("retry: %d\n", m.Retry))
    }

    if len(m.Event) > 0 {
        buffer.WriteString(fmt.Sprintf("event: %s\n", m.Event))
    }

    if len(m.Data) > 0 {
        buffer.WriteString(fmt.Sprintf("data: %s\n", m.Data))
    }

    buffer.WriteString("\n")

    return buffer.String()
}
