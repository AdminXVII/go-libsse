package main

import (
    "log"
    "net/http"
    "strconv"
    "time"

    "github.com/AdminXVII/go-sse"
)

func main() {
    s := sse.NewServer(nil)

    http.Handle("/", http.FileServer(http.Dir("./static")))
    http.Handle("/events/", s)

    go func () {
        for {
            s.SendMessage(&sse.Message{Data: time.Now().String()})
            time.Sleep(5 * time.Second)
        }
    }()

    go func () {
        i := 0
        for {
            i++
            s.SendMessage(&sse.Message{Data: strconv.Itoa(i)})
            time.Sleep(5 * time.Second)
        }
    }()

    log.Println("Listening at :3000")
    http.ListenAndServe(":3000", nil)
}
