package sse

import "testing"

func TestEmptyMessage(t *testing.T) {
	msg := Message{}

	if msg.ToBuffer().String() != "\n" {
		t.Fatal("Message not empty.")
	}
}

func TestMessage(t *testing.T) {
	msg := Message{
		"123",
		"test",
		"myevent",
		10 * 1000,
	}

	if msg.ToBuffer().String() != "id: 123\nretry: 10000\nevent: myevent\ndata: test\n\n" {
		t.Fatal("Message not correct.")
	}
}
