package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/AdminXVII/go-libsse"
)

func main() {
	s := sse.NewServer(nil)

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.Handle("/events/", s)

	go func() {
		i := 0
		for {
			i++
			s.SendMessage(sse.Message{Data: strconv.Itoa(i)})
			time.Sleep(time.Second)
		}
	}()

	log.Println("Listening at :3000")
	http.ListenAndServe(":3000", nil)
}
