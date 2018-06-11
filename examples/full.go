package main

import (
    "log"
    "net/http"
    "strconv"
    "time"
    "os"

    "github.com/AdminXVII/go-libsse"
)

func main() {
    s := sse.NewServer(&sse.Options{
        // Increase default retry interval to 10s.
        RetryInterval: 10 * 1000,
        // CORS headers
        Headers: map[string]string {
            "Access-Control-Allow-Origin": "*",
            "Access-Control-Allow-Methods": "GET, OPTIONS",
            "Access-Control-Allow-Headers": "Keep-Alive,X-Requested-With,Cache-Control,Content-Type,Last-Event-ID",
        },
        // Print debug info
        Logger: log.New(os.Stdout,
          "go-sse: ",
          log.Ldate|log.Ltime|log.Lshortfile),
        // Add pertinent info first
        InitMessages: func(ClientLastEventId string, ServerLastEventId string) []sse.Message {
          return []sse.Message{sse.Message{Data: ServerLastEventId}, sse.Message{Id: "42", Data: "This is the answer to life, to the universe and to everything else"}}
        },
    })

    http.Handle("/", http.FileServer(http.Dir("./static")))
    http.Handle("/events/", s)

    go func () {
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
